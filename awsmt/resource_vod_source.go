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

type resourceVodSourceModel struct {
	ID                        types.String                        `tfsdk:"id"`
	Arn                       types.String                        `tfsdk:"arn"`
	CreationTime              types.String                        `tfsdk:"creation_time"`
	HttpPackageConfigurations []httpPackageConfigurationsVSRModel `tfsdk:"http_package_configurations"`
	LastModifiedTime          types.String                        `tfsdk:"last_modified_time"`
	SourceLocationName        types.String                        `tfsdk:"source_location_name"`
	Tags                      map[string]*string                  `tfsdk:"tags"`
	VodSourceName             types.String                        `tfsdk:"vod_source_name"`
}

type httpPackageConfigurationsVSRModel struct {
	Path        types.String `tfsdk:"path"`
	SourceGroup types.String `tfsdk:"source_group"`
	Type        types.String `tfsdk:"type"`
}

func (r *resourceVodSource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_location"
}

func (r *resourceVodSource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"source_location_name": schema.StringAttribute{
				Required: true,
			},
			"tags": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"vod_source_name": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *resourceVodSource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*mediatailor.MediaTailor)
}

func (r *resourceVodSource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceVodSourceModel

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
	var state resourceVodSourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &mediatailor.DescribeVodSourceInput{}
	input.VodSourceName = aws.String(state.VodSourceName.String())
	input.SourceLocationName = aws.String(state.SourceLocationName.String())

	vodSource, err := r.client.DescribeVodSource(input)
	if err != nil {
		resp.Diagnostics.AddError("Error while describing vod source", err.Error())
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
	var plan resourceVodSourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &mediatailor.DescribeVodSourceInput{}
	input.VodSourceName = aws.String(plan.VodSourceName.String())
	input.SourceLocationName = aws.String(plan.SourceLocationName.String())

	VodSource, err := r.client.DescribeVodSource(input)
	if err != nil {
		resp.Diagnostics.AddError("Error while describing Vod source", err.Error())
		return
	}

	oldTags := VodSource.Tags
	newTags := plan.Tags

	// Check if tags are different
	if !reflect.DeepEqual(oldTags, newTags) {
		err = updatesTags(r.client, oldTags, newTags, *VodSource.Arn)
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

	plan = readUpdatedVodSourceToPlan(plan, *updatedVodSource)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceVodSource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceVodSourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := &mediatailor.DeleteVodSourceInput{}
	input.VodSourceName = aws.String(state.VodSourceName.String())
	input.SourceLocationName = aws.String(state.SourceLocationName.String())

	_, err := r.client.DeleteVodSource(input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while deleting vod source "+err.Error(),
			err.Error(),
		)
	}
}

func (r *resourceVodSource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("arn"), req, resp)
}
