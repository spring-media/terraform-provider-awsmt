package awsmt

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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
	client *mediatailor.Client
}

func (r *resourceLiveSource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_live_source"
}

func (r *resourceLiveSource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                          computedStringWithStateForUnknown,
			"arn":                         computedStringWithStateForUnknown,
			"creation_time":               computedStringWithStateForUnknown,
			"http_package_configurations": httpPackageConfigurationsResourceSchema,
			"last_modified_time":          computedString,
			"source_location_name":        requiredString,
			"tags":                        optionalMap,
			"name":                        requiredStringWithRequiresReplace,
		},
	}
}

func (r *resourceLiveSource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(clients).v2
}

func (r *resourceLiveSource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan liveSourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	liveSource, err := r.client.CreateLiveSource(ctx, getCreateLiveSourceInput(plan))
	if err != nil {
		resp.Diagnostics.AddError("Error while creating live source", err.Error())
		return
	}

	plan = readLiveSource(plan, *liveSource)

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceLiveSource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state liveSourceModel

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

	input := &mediatailor.DescribeLiveSourceInput{
		LiveSourceName:     &name,
		SourceLocationName: &sourceLocationName,
	}

	liveSource, err := r.client.DescribeLiveSource(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error while describing live source", err.Error())
		return
	}

	state = readLiveSource(state, mediatailor.CreateLiveSourceOutput(*liveSource))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceLiveSource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan liveSourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &mediatailor.DescribeLiveSourceInput{
		LiveSourceName:     plan.Name,
		SourceLocationName: plan.SourceLocationName,
	}

	liveSource, err := r.client.DescribeLiveSource(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error while describing live source", err.Error())
		return
	}

	// Update tags
	err = V2UpdatesTags(r.client, liveSource.Tags, plan.Tags, *liveSource.Arn)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while updating live source tags"+err.Error(),
			err.Error(),
		)
	}

	updateInput := getUpdateLiveSourceInput(plan)
	updatedLiveSource, err := r.client.UpdateLiveSource(ctx, &updateInput)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while updating live source "+err.Error(),
			err.Error(),
		)
	}

	plan = readLiveSource(plan, mediatailor.CreateLiveSourceOutput(*updatedLiveSource))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceLiveSource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state liveSourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &mediatailor.DeleteLiveSourceInput{
		LiveSourceName:     state.Name,
		SourceLocationName: state.SourceLocationName,
	}

	_, err := r.client.DeleteLiveSource(ctx, params)
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
			fmt.Sprintf("Expected import identifier with format: source_location_name, name. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("source_location_name"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)

}
