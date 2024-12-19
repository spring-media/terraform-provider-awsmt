package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-mediatailor/awsmt/models"
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
	client *mediatailor.Client
}

func (r *resourceVodSource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vod_source"
}

func (r *resourceVodSource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                          computedStringWithStateForUnknown,
			"source_location_name":        requiredString,
			"http_package_configurations": httpPackageConfigurationsResourceSchema,
			"creation_time":               computedStringWithStateForUnknown,
			"tags":                        optionalMap,
			"last_modified_time":          computedString,
			"arn":                         computedStringWithStateForUnknown,
			"name":                        requiredStringWithRequiresReplace,
			"ad_break_opportunities_offset_millis": schema.ListAttribute{
				Computed:    true,
				Optional:    true,
				ElementType: types.Int64Type,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *resourceVodSource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*mediatailor.Client)
}

func (r *resourceVodSource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.VodSourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	vodSource, err := r.client.CreateVodSource(ctx, getCreateVodSourceInput(plan))
	if err != nil {
		resp.Diagnostics.AddError("Error while creating vod source", err.Error())
		return
	}

	plan = readVodSourceToPlan(plan, *vodSource)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceVodSource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.VodSourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var sourceLocationName, name string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("source_location_name"), &sourceLocationName)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("name"), &name)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &mediatailor.DescribeVodSourceInput{
		VodSourceName:      &name,
		SourceLocationName: &sourceLocationName,
	}

	vodSource, err := r.client.DescribeVodSource(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error while describing vod source", "Could not describe the vod source: "+*input.SourceLocationName+":"+*input.VodSourceName+". "+err.Error())
		return
	}

	state = readVodSourceToState(state, *vodSource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceVodSource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.VodSourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &mediatailor.DescribeVodSourceInput{
		VodSourceName:      plan.Name,
		SourceLocationName: plan.SourceLocationName,
	}

	vodSource, err := r.client.DescribeVodSource(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error while describing Vod source", err.Error())
		return
	}

	// Update tags
	err = UpdatesTags(r.client, vodSource.Tags, plan.Tags, *vodSource.Arn)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while updating vod source tags"+err.Error(),
			err.Error(),
		)
	}

	updateInput := getUpdateVodSourceInput(plan)
	updatedVodSource, err := r.client.UpdateVodSource(ctx, &updateInput)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while updating vod source "+err.Error(),
			err.Error(),
		)
	}

	plan = readVodSourceToPlan(plan, mediatailor.CreateVodSourceOutput(*updatedVodSource))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceVodSource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.VodSourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &mediatailor.DeleteVodSourceInput{
		VodSourceName:      state.Name,
		SourceLocationName: state.SourceLocationName,
	}

	_, err := r.client.DeleteVodSource(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while deleting vod source "+err.Error(),
			err.Error(),
		)
	}
}

func (r *resourceVodSource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importStateForContentSources(ctx, req, resp)
}
