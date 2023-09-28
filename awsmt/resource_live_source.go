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
	"reflect"
	"strings"
)

var (
	_ resource.Resource                = &resourceLiveSource{}
	_ resource.ResourceWithConfigure   = &resourceLiveSource{}
	_ resource.ResourceWithImportState = &resourceLiveSource{}
)

func ResourceLiveSource() resource.Resource {
	return &resourceLiveSource{}
}

type resourceLiveSource struct {
	client *mediatailor.MediaTailor
}

func (r *resourceLiveSource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_live_source"
}

func (r *resourceLiveSource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":            computedString,
			"arn":           computedString,
			"creation_time": computedString,
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
			"last_modified_time":   computedString,
			"live_source_name":     requiredString,
			"source_location_name": requiredString,
			"tags":                 optionalMap,
		},
	}
}

func (r *resourceLiveSource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*mediatailor.MediaTailor)
}

func (r *resourceLiveSource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan liveSourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := liveSourceInput(plan)

	liveSource, err := r.client.CreateLiveSource(&input)
	if err != nil {
		resp.Diagnostics.AddError("Error while creating live source", err.Error())
		return
	}

	plan = readLiveSourceToPlan(plan, *liveSource)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceLiveSource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state liveSourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var sourceLocationName, liveSourceName string

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("source_location_name"), &sourceLocationName)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("live_source_name"), &liveSourceName)...)

	if resp.Diagnostics.HasError() {
		return
	}

	input := &mediatailor.DescribeLiveSourceInput{}
	input.LiveSourceName = &liveSourceName
	input.SourceLocationName = &sourceLocationName

	liveSource, err := r.client.DescribeLiveSource(input)
	if err != nil {
		resp.Diagnostics.AddError("Error while describing live source", err.Error())
		return
	}

	state = readLiveSourceToPlan(state, mediatailor.CreateLiveSourceOutput(*liveSource))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceLiveSource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan liveSourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &mediatailor.DescribeLiveSourceInput{}
	input.LiveSourceName = plan.LiveSourceName
	input.SourceLocationName = plan.SourceLocationName

	liveSource, err := r.client.DescribeLiveSource(input)
	if err != nil {
		resp.Diagnostics.AddError("Error while describing live source", err.Error())
		return
	}

	oldTags := liveSource.Tags
	newTags := plan.Tags

	// Check if tags are different
	if !reflect.DeepEqual(oldTags, newTags) {
		err = updatesTags(r.client, oldTags, newTags, *liveSource.Arn)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error while updating live source tags"+err.Error(),
				err.Error(),
			)
		}
	}

	updateInput := liveSourceUpdateInput(plan)
	updatedLiveSource, err := r.client.UpdateLiveSource(&updateInput)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while updating live source "+err.Error(),
			err.Error(),
		)
	}

	plan = readLiveSourceToPlan(plan, mediatailor.CreateLiveSourceOutput(*updatedLiveSource))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceLiveSource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state liveSourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &mediatailor.DeleteLiveSourceInput{}
	params.LiveSourceName = state.LiveSourceName
	params.SourceLocationName = state.SourceLocationName

	_, err := r.client.DeleteLiveSource(params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while deleting live source "+err.Error(),
			err.Error(),
		)
	}
}

func (r *resourceLiveSource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: source_location_name, live_source_name. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("source_location_name"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("live_source_name"), idParts[1])...)

}
