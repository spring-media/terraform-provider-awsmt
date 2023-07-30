package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"
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

type resourceLiveSourceModel struct {
	ID                        types.String                        `tfsdk:"id"`
	Arn                       types.String                        `tfsdk:"arn"`
	CreationTime              types.String                        `tfsdk:"creation_time"`
	HttpPackageConfigurations []httpPackageConfigurationsLSRModel `tfsdk:"http_package_configurations"`
	LastModifiedTime          types.String                        `tfsdk:"last_modified_time"`
	LiveSourceName            types.String                        `tfsdk:"live_source_name"`
	SourceLocationName        types.String                        `tfsdk:"source_location_name"`
	Tags                      map[string]*string                  `tfsdk:"tags"`
}

type httpPackageConfigurationsLSRModel struct {
	Path        types.String `tfsdk:"path"`
	SourceGroup types.String `tfsdk:"source_group"`
	Type        types.String `tfsdk:"type"`
}

func (r *resourceLiveSource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_location"
}

func (r *resourceLiveSource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"arn": schema.StringAttribute{
				Computed: true,
			},
			"creation_time": schema.StringAttribute{
				Computed: true,
			},
			"http_package_configuration": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"path": schema.StringAttribute{
							Required: true,
						},
						"source_group": schema.StringAttribute{
							Required: true,
						},
						"type": schema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("HLS", "DASH"),
							},
						},
					},
				},
			},
			"last_modified_time": schema.StringAttribute{
				Computed: true,
			},
			"live_source_name": schema.StringAttribute{
				Required: true,
			},
			"source_location_name": schema.StringAttribute{
				Required: true,
			},
			"tags": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
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
	var plan resourceLiveSourceModel

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
	var state resourceLiveSourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &mediatailor.DescribeLiveSourceInput{}
	input.LiveSourceName = aws.String(state.LiveSourceName.String())
	input.SourceLocationName = aws.String(state.SourceLocationName.String())

	liveSource, err := r.client.DescribeLiveSource(input)
	if err != nil {
		resp.Diagnostics.AddError("Error while describing live source", err.Error())
		return
	}

	state = readLiveSourceToState(state, *liveSource)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceLiveSource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resourceLiveSourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &mediatailor.DescribeLiveSourceInput{}
	input.LiveSourceName = aws.String(plan.LiveSourceName.String())
	input.SourceLocationName = aws.String(plan.SourceLocationName.String())

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

	plan = readUpdatedLiveSourceToPlan(plan, *updatedLiveSource)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceLiveSource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceLiveSourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &mediatailor.DeleteLiveSourceInput{}
	input.LiveSourceName = aws.String(state.LiveSourceName.String())
	input.SourceLocationName = aws.String(state.SourceLocationName.String())

	_, err := r.client.DeleteLiveSource(input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while deleting live source "+err.Error(),
			err.Error(),
		)
	}
}

func (r *resourceLiveSource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("arn"), req, resp)
}
