package awsmt

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"
	"strings"
)

var (
	_ resource.Resource                = &resourceVodSource{}
	_ resource.ResourceWithConfigure   = &resourceVodSource{}
	_ resource.ResourceWithImportState = &resourceVodSource{}
)

func ResourceVodSource() resource.Resource {
	return &resourceVodSource{}
}

type resourceVodSource struct {
	client *mediatailor.MediaTailor
}

func (r *resourceVodSource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vod_source"
}

func (r *resourceVodSource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                   computedString,
			"source_location_name": requiredString,
			"http_package_configurations": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"path":         requiredString,
						"source_group": requiredString,
						"type": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("HLS", "DASH"),
							},
						},
					},
				},
			},
			"creation_time":      computedString,
			"tags":               optionalMap,
			"last_modified_time": computedString,
			"arn":                computedString,
			"name":               requiredString,
			"ad_break_opportunities_offset_millis": schema.ListAttribute{
				Optional:    true,
				ElementType: types.Int64Type,
			},
		},
	}
}

func (r *resourceVodSource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(clients).v1
}

func (r *resourceVodSource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vodSourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := vodSourceInput(plan)

	vodSource, err := r.client.CreateVodSource(&input)
	if err != nil {
		resp.Diagnostics.AddError("Error while creating vod source", err.Error())
		return
	}

	plan = readVodSourceToPlan(plan, *vodSource)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceVodSource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vodSourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var sourceLocationName, name string

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("source_location_name"), &sourceLocationName)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("name"), &name)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := &mediatailor.DescribeVodSourceInput{}
	input.VodSourceName = &name
	input.SourceLocationName = &sourceLocationName

	vodSource, err := r.client.DescribeVodSource(input)
	if err != nil {
		resp.Diagnostics.AddError("Error while describing vod source", "Could not describe the vod source: "+*input.SourceLocationName+" and "+*input.VodSourceName+": "+err.Error())
		return
	}

	state = readVodSourceToState(state, *vodSource)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceVodSource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan vodSourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &mediatailor.DescribeVodSourceInput{}
	input.VodSourceName = plan.Name
	input.SourceLocationName = plan.SourceLocationName

	vodSource, err := r.client.DescribeVodSource(input)
	if err != nil {
		resp.Diagnostics.AddError("Error while describing Vod source", err.Error())
		return
	}

	oldTags := vodSource.Tags
	newTags := plan.Tags

	// Check if tags are different
	if !reflect.DeepEqual(oldTags, newTags) {
		err = updatesTags(r.client, oldTags, newTags, *vodSource.Arn)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error while updating vod source tags"+err.Error(),
				err.Error(),
			)
		}
	}

	updateInput := vodSourceUpdateInput(plan)
	updatedVodSource, err := r.client.UpdateVodSource(&updateInput)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while updating vod source "+err.Error(),
			err.Error(),
		)
	}

	plan = readVodSourceToPlan(plan, mediatailor.CreateVodSourceOutput(*updatedVodSource))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceVodSource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state vodSourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &mediatailor.DeleteVodSourceInput{}
	input.VodSourceName = state.Name
	input.SourceLocationName = state.SourceLocationName

	_, err := r.client.DeleteVodSource(input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while deleting vod source "+err.Error(),
			err.Error(),
		)
	}
}

func (r *resourceVodSource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: source_location_name, name. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("source_location_name"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
}
