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
	_ resource.Resource                = &resourceSourceLocation{}
	_ resource.ResourceWithConfigure   = &resourceSourceLocation{}
	_ resource.ResourceWithImportState = &resourceSourceLocation{}
)

func ResourceSourceLocation() resource.Resource {
	return &resourceSourceLocation{}
}

type resourceSourceLocation struct {
	client *mediatailor.MediaTailor
}

type resourceSourceLocationModel struct {
	ID                                  types.String                              `tfsdk:"id"`
	AccessConfiguration                 accessConfigurationRModel                 `tfsdk:"access_configuration"`
	Arn                                 types.String                              `tfsdk:"arn"`
	CreationTime                        types.String                              `tfsdk:"creation_time"`
	DefaultSegmentDeliveryConfiguration defaultSegmentDeliveryConfigurationRModel `tfsdk:"default_segment_delivery_configuration"`
	HttpConfiguration                   httpConfigurationRModel                   `tfsdk:"http_configuration"`
	LastModifiedTime                    types.String                              `tfsdk:"last_modified_time"`
	SegmentDeliveryConfigurations       []segmentDeliveryConfigurationsRModel     `tfsdk:"segment_delivery_configuration"`
	SourceLocationName                  types.String                              `tfsdk:"source_location_name"`
	Tags                                map[string]*string                        `tfsdk:"tags"`
}

type accessConfigurationRModel struct {
	AccessType                             types.String                                 `tfsdk:"access_type"`
	SecretsManagerAccessTokenConfiguration secretsManagerAccessTokenConfigurationRModel `tfsdk:"secrets_manager_access_token_configuration"`
}

type secretsManagerAccessTokenConfigurationRModel struct {
	HeaderName      types.String `tfsdk:"header_name"`
	SecretArn       types.String `tfsdk:"secret_arn"`
	SecretStringKey types.String `tfsdk:"secret_string_key"`
}

type defaultSegmentDeliveryConfigurationRModel struct {
	BaseUrl types.String `tfsdk:"dsdc_base_url"`
}

type httpConfigurationRModel struct {
	BaseUrl types.String `tfsdk:"hc_base_url"`
}

type segmentDeliveryConfigurationsRModel struct {
	BaseUrl types.String `tfsdk:"sdc_base_url"`
	Name    types.String `tfsdk:"name"`
}

func (r *resourceSourceLocation) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_location"
}

func (r *resourceSourceLocation) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"access_configuration": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"access_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("S3_SIGV4"),
						},
					},
					"secrets_manager_access_token_configuration": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"header_name": schema.StringAttribute{
								Optional: true,
							},
							"secret_arn": schema.StringAttribute{
								Optional: true,
							},
							"secret_string_key": schema.StringAttribute{
								Optional: true,
							},
						},
					},
				},
			},
			"arn": schema.StringAttribute{
				Computed: true,
			},
			"creation_time": schema.StringAttribute{
				Computed: true,
			},
			"default_segment_delivery_configuration": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"dsdc_base_url": schema.StringAttribute{
						Optional: true,
					},
				},
			},
			"http_configuration": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"hc_base_url": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"last_modified_time": schema.StringAttribute{
				Computed: true,
			},
			"segment_delivery_configurations": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"sdc_base_url": schema.StringAttribute{
							Optional: true,
						},
						"name": schema.StringAttribute{
							Optional: true,
						},
					},
				},
			},
			"sourceLocationName": schema.StringAttribute{
				Required: true,
			},
			"tags": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *resourceSourceLocation) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*mediatailor.MediaTailor)
}

func (r *resourceSourceLocation) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceSourceLocationModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := sourceLocationInput(plan)

	// Create Source Location
	sourceLocation, err := r.client.CreateSourceLocation(&params)
	if err != nil {
		resp.Diagnostics.AddError("Error while creating source location", err.Error())
		return
	}

	plan = readSourceLocationToPlan(plan, *sourceLocation)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSourceLocation) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceSourceLocationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceLocationName := aws.String(state.SourceLocationName.String())

	// Describe Source Location

	sourceLocation, err := r.client.DescribeSourceLocation(&mediatailor.DescribeSourceLocationInput{SourceLocationName: sourceLocationName})
	if err != nil {
		resp.Diagnostics.AddError("Error while describing source location", "Could not describe the source location: "+state.SourceLocationName.ValueString()+": "+err.Error())
		return
	}

	state = readSourceLocationToState(state, *sourceLocation)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSourceLocation) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resourceSourceLocationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.SourceLocationName.String()

	sourceLocation, err := r.client.DescribeSourceLocation(&mediatailor.DescribeSourceLocationInput{SourceLocationName: aws.String(name)})
	if err != nil {
		resp.Diagnostics.AddError("Error while describing source location", "Could not describe the source location: "+name+": "+err.Error())
		return
	}

	oldTags := sourceLocation.Tags
	newTags := plan.Tags

	// Check if tags are different
	if !reflect.DeepEqual(oldTags, newTags) {
		err = updatesTags(r.client, oldTags, newTags, *sourceLocation.Arn)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error while updating playback configuration tags"+err.Error(),
				err.Error(),
			)
		}
	}

	if !reflect.DeepEqual(sourceLocation.AccessConfiguration, plan.AccessConfiguration) {
		// delete source location
		name := aws.String(plan.SourceLocationName.String())
		err = deleteSourceLocation(r.client, name)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error while deleting source location "+err.Error(),
				err.Error(),
			)
			return
		}

		// create source location
		params := sourceLocationInput(plan)
		sourceLocation, err := r.client.CreateSourceLocation(&params)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error while creating new source location with new access configuration "+err.Error(),
				err.Error(),
			)
			return
		}

		plan = readSourceLocationToPlan(plan, *sourceLocation)
	}

	params := updateSourceLocationInput(plan)

	sourceLocationUpdated, err := r.client.UpdateSourceLocation(&params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating source location "+err.Error(),
			err.Error(),
		)
		return
	}

	plan = readSourceLocationToPlanUpdate(plan, *sourceLocationUpdated)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSourceLocation) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceSourceLocationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := aws.String(state.SourceLocationName.String())

	vodSourcesList, err := r.client.ListVodSources(&mediatailor.ListVodSourcesInput{SourceLocationName: name})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving vod sources "+err.Error(),
			err.Error(),
		)
		return
	}
	for _, vodSource := range vodSourcesList.Items {
		if _, err := r.client.DeleteVodSource(&mediatailor.DeleteVodSourceInput{VodSourceName: vodSource.VodSourceName, SourceLocationName: name}); err != nil {
			resp.Diagnostics.AddError(
				"Error deleting vod sources "+err.Error(),
				err.Error(),
			)
			return
		}
	}
	liveSourcesList, err := r.client.ListLiveSources(&mediatailor.ListLiveSourcesInput{SourceLocationName: name})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving live sources "+err.Error(),
			err.Error(),
		)
		return
	}
	for _, liveSource := range liveSourcesList.Items {
		if _, err := r.client.DeleteLiveSource(&mediatailor.DeleteLiveSourceInput{LiveSourceName: liveSource.LiveSourceName, SourceLocationName: name}); err != nil {
			resp.Diagnostics.AddError(
				"Error deleting live sources "+err.Error(),
				err.Error(),
			)
			return
		}
	}
	_, err = r.client.DeleteSourceLocation(&mediatailor.DeleteSourceLocationInput{SourceLocationName: name})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting resource "+err.Error(),
			err.Error(),
		)
		return
	}
}

func (r *resourceSourceLocation) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("arn"), req, resp)
}
