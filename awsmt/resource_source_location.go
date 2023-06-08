package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"
)

var (
	_ resource.Resource              = &resourceSourceLocation{}
	_ resource.ResourceWithConfigure = &resourceSourceLocation{}
)

func ResourceSourceLocation() resource.Resource {
	return &resourceSourceLocation{}
}

type resourceSourceLocation struct {
	client *mediatailor.MediaTailor
}

type resourceSourceLocationModel struct {
	ID                                  types.String                                                    `tfsdk:"id"`
	AccessConfiguration                 *resourceSourceLocationAccessConfigurationModel                 `tfsdk:"access_configuration"`
	Arn                                 types.String                                                    `tfsdk:"arn"`
	CreationTime                        types.String                                                    `tfsdk:"creation_time"`
	DefaultSegmentDeliveryConfiguration *resourceSourceLocationDefaultSegmentDeliveryConfigurationModel `tfsdk:"default_segment_delivery_configuration"`
	HttpConfiguration                   resourceSourceLocationHttpConfigurationModel                    `tfsdk:"http_configuration"`
	LastModifiedTime                    types.String                                                    `tfsdk:"last_modified_time"`
	SegmentDeliveryConfigurations       *[]resourceSourceLocationSegmentDeliveryConfigurationsModel     `tfsdk:"segment_delivery_configuration"`
	SourceLocationName                  *string                                                         `tfsdk:"name"`
	Tags                                map[string]*string                                              `tfsdk:"tags"`
}

type resourceSourceLocationAccessConfigurationModel struct {
	AccessType                             *string                                                           `tfsdk:"access_type"`
	SecretsManagerAccessTokenConfiguration resourceSourceLocationSecretsManagerAccessTokenConfigurationModel `tfsdk:"smatc"`
	// SMATC is short for Secret Manager Access Token Configuration
}

type resourceSourceLocationSecretsManagerAccessTokenConfigurationModel struct {
	HeaderName      *string `tfsdk:"header_name"`
	SecretArn       *string `tfsdk:"secret_arn"`
	SecretStringKey *string `tfsdk:"secret_string_key"`
}

type resourceSourceLocationDefaultSegmentDeliveryConfigurationModel struct {
	BaseUrl *string `tfsdk:"base_url"`
}

type resourceSourceLocationHttpConfigurationModel struct {
	BaseUrl *string `tfsdk:"base_url"`
}

type resourceSourceLocationSegmentDeliveryConfigurationsModel struct {
	BaseUrl *string `tfsdk:"base_url"`
	Name    *string `tfsdk:"name"`
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
			"access_configuration": schema.MapNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"access_type": schema.StringAttribute{
							Optional: true,
							Validators: []validator.String{
								stringvalidator.OneOf("S3_SIGV4", "SECRETS_MANAGER_ACCESS_TOKEN"),
							},
							CustomType: types.StringType,
						},
						"smatc": schema.MapNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"header_name": schema.StringAttribute{
										Optional:   true,
										CustomType: types.StringType,
									},
									"secret_arn": schema.StringAttribute{
										Optional:   true,
										CustomType: types.StringType,
									},
									"secret_string_key": schema.StringAttribute{
										Optional:   true,
										CustomType: types.StringType,
									},
								},
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
					"base_url": schema.StringAttribute{
						Optional:   true,
						CustomType: types.StringType,
					},
				},
			},
			"http_configuration": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"base_url": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"last_modified_time": schema.StringAttribute{
				Computed: true,
			},
			"segment_delivery_configuration": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"base_url": schema.StringAttribute{
							Optional:   true,
							CustomType: types.StringType,
						},
						"name": schema.StringAttribute{
							Optional:   true,
							CustomType: types.StringType,
						},
					},
				},
			},
			"name": schema.StringAttribute{
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

	input := mediatailor.CreateSourceLocationInput{}
	if plan.AccessConfiguration != nil {
		input.AccessConfiguration = &mediatailor.AccessConfiguration{}
		if plan.AccessConfiguration.AccessType != nil {
			input.AccessConfiguration.AccessType = plan.AccessConfiguration.AccessType
		}
		if plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName != nil {
			input.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName = plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName
		}
		if plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn != nil {
			input.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn = plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn
		}
		if plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey != nil {
			input.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey = plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey
		}
	}
	if plan.DefaultSegmentDeliveryConfiguration != nil {
		input.DefaultSegmentDeliveryConfiguration = &mediatailor.DefaultSegmentDeliveryConfiguration{}
		if plan.DefaultSegmentDeliveryConfiguration.BaseUrl != nil {
			input.DefaultSegmentDeliveryConfiguration.BaseUrl = plan.DefaultSegmentDeliveryConfiguration.BaseUrl
		}
	}
	if plan.SegmentDeliveryConfigurations != nil {
		for _, v := range *plan.SegmentDeliveryConfigurations {
			if v.BaseUrl != nil && v.Name != nil {
				input.SegmentDeliveryConfigurations = append(input.SegmentDeliveryConfigurations, &mediatailor.SegmentDeliveryConfiguration{
					BaseUrl: aws.String(*v.BaseUrl),
					Name:    aws.String(*v.Name),
				})
			}
		}
	}
	if plan.Tags != nil && len(plan.Tags) > 0 {
		input.Tags = plan.Tags
	}
	input.SourceLocationName = plan.SourceLocationName
	input.HttpConfiguration = &mediatailor.HttpConfiguration{
		BaseUrl: plan.HttpConfiguration.BaseUrl,
	}

	sourceLocation, err := r.client.CreateSourceLocation(&input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while creating source location "+err.Error(),
			err.Error(),
		)
		return
	}

	plan.Arn = types.StringValue(*sourceLocation.Arn)
	plan.CreationTime = types.StringValue((aws.TimeValue(sourceLocation.CreationTime)).String())
	plan.LastModifiedTime = types.StringValue((aws.TimeValue(sourceLocation.LastModifiedTime)).String())

	plan.ID = types.StringValue(*sourceLocation.SourceLocationName)

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

	name := state.SourceLocationName
	// Get refreshed order value from AWS
	sourceLocation, err := r.client.DescribeSourceLocation(&mediatailor.DescribeSourceLocationInput{
		SourceLocationName: name,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving source location "+err.Error(),
			err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.Arn = types.StringValue(*sourceLocation.Arn)
	if sourceLocation.AccessConfiguration != nil {
		state.AccessConfiguration.AccessType = sourceLocation.AccessConfiguration.AccessType
		if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration != nil {
			if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName != nil {
				state.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName = sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName
			}
			if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn != nil {
				state.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn = sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn
			}
			if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey != nil {
				state.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey = sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey
			}
		}
	}
	if sourceLocation.DefaultSegmentDeliveryConfiguration != nil {
		if sourceLocation.DefaultSegmentDeliveryConfiguration.BaseUrl != nil {
			state.DefaultSegmentDeliveryConfiguration.BaseUrl = sourceLocation.DefaultSegmentDeliveryConfiguration.BaseUrl
		}
	}
	if sourceLocation.SegmentDeliveryConfigurations != nil {
		for _, v := range sourceLocation.SegmentDeliveryConfigurations {
			*state.SegmentDeliveryConfigurations = append(*state.SegmentDeliveryConfigurations, resourceSourceLocationSegmentDeliveryConfigurationsModel{
				BaseUrl: v.BaseUrl,
				Name:    v.Name,
			})
		}
	}
	if sourceLocation.Tags != nil {
		state.Tags = sourceLocation.Tags
	}
	state.HttpConfiguration.BaseUrl = sourceLocation.HttpConfiguration.BaseUrl
	state.CreationTime = types.StringValue((aws.TimeValue(sourceLocation.CreationTime)).String())
	state.LastModifiedTime = types.StringValue((aws.TimeValue(sourceLocation.LastModifiedTime)).String())

	state.ID = types.StringValue(*sourceLocation.SourceLocationName)
	// Set refreshed state
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

	name := plan.SourceLocationName
	sourceLocation, err := r.client.DescribeSourceLocation(&mediatailor.DescribeSourceLocationInput{
		SourceLocationName: name,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving source location "+err.Error(),
			err.Error(),
		)
		return
	}

	oldTags := sourceLocation.Tags
	newTags := plan.Tags
	if !reflect.DeepEqual(oldTags, newTags) {
		if oldTags != nil && len(oldTags) > 0 {
			var removeTags []*string
			for k := range oldTags {
				removeTags = append(removeTags, aws.String(k))
			}
			_, err := r.client.UntagResource(&mediatailor.UntagResourceInput{ResourceArn: sourceLocation.Arn, TagKeys: removeTags})
			if err != nil {
				resp.Diagnostics.AddError(
					"Error removing tags from source location "+err.Error(),
					err.Error(),
				)
				return
			}
		}
		if newTags != nil {
			_, err := r.client.TagResource(&mediatailor.TagResourceInput{ResourceArn: sourceLocation.Arn, Tags: newTags})
			if err != nil {
				resp.Diagnostics.AddError(
					"Error adding tags to source location "+err.Error(),
					err.Error(),
				)
				return
			}
		}

		plan.Tags = newTags
	}

	var params mediatailor.UpdateSourceLocationInput

	if plan.AccessConfiguration != nil {
		params.AccessConfiguration = &mediatailor.AccessConfiguration{}
		if plan.AccessConfiguration != nil || sourceLocation.AccessConfiguration != nil {
			if plan.AccessConfiguration.AccessType != nil || sourceLocation.AccessConfiguration.AccessType != nil {
				if plan.AccessConfiguration.AccessType != sourceLocation.AccessConfiguration.AccessType {
					params.AccessConfiguration.AccessType = plan.AccessConfiguration.AccessType
				}
			}
			if plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName != nil || sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName != nil {
				if plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName != sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName {
					params.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName = plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName
				}
			}
			if plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn != nil || sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn != nil {
				if plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn != sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn {
					params.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn = plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn
				}
			}
			if plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey != nil || sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey != nil {
				if plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey != sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey {
					params.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey = plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey
				}
			}
		}
	}
	if plan.DefaultSegmentDeliveryConfiguration != nil {
		if plan.DefaultSegmentDeliveryConfiguration.BaseUrl != nil || sourceLocation.DefaultSegmentDeliveryConfiguration.BaseUrl != nil {
			params.DefaultSegmentDeliveryConfiguration = &mediatailor.DefaultSegmentDeliveryConfiguration{}
			params.DefaultSegmentDeliveryConfiguration.BaseUrl = plan.DefaultSegmentDeliveryConfiguration.BaseUrl
		}
	}
	if plan.HttpConfiguration.BaseUrl != nil || sourceLocation.HttpConfiguration.BaseUrl != nil {
		params.HttpConfiguration = &mediatailor.HttpConfiguration{}
		params.HttpConfiguration.BaseUrl = plan.HttpConfiguration.BaseUrl
	} else {
		params.HttpConfiguration.BaseUrl = sourceLocation.HttpConfiguration.BaseUrl
	}
	if plan.SegmentDeliveryConfigurations != nil || sourceLocation.SegmentDeliveryConfigurations != nil {
		params.SegmentDeliveryConfigurations = []*mediatailor.SegmentDeliveryConfiguration{}
		if !reflect.DeepEqual(plan.SegmentDeliveryConfigurations, sourceLocation.SegmentDeliveryConfigurations) {
			for i, v := range *plan.SegmentDeliveryConfigurations {
				params.SegmentDeliveryConfigurations[i].BaseUrl = v.BaseUrl
				params.SegmentDeliveryConfigurations[i].Name = v.Name
			}
		}
	}
	params.SourceLocationName = name

	sourceLocationUpdated, err := r.client.UpdateSourceLocation(&params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating source location "+err.Error(),
			err.Error(),
		)
		return
	}

	plan.Arn = types.StringValue(*sourceLocationUpdated.Arn)
	plan.CreationTime = types.StringValue((aws.TimeValue(sourceLocationUpdated.CreationTime)).String())
	plan.LastModifiedTime = types.StringValue((aws.TimeValue(sourceLocationUpdated.LastModifiedTime)).String())
	if sourceLocationUpdated.AccessConfiguration != nil {
		plan.AccessConfiguration.AccessType = sourceLocationUpdated.AccessConfiguration.AccessType
		if sourceLocationUpdated.AccessConfiguration.SecretsManagerAccessTokenConfiguration != nil {
			if sourceLocationUpdated.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName != nil {
				plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName = sourceLocationUpdated.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName
			}
			if sourceLocationUpdated.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn != nil {
				plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn = sourceLocationUpdated.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn
			}
			if sourceLocationUpdated.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey != nil {
				plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey = sourceLocationUpdated.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey
			}
		}
	}
	if sourceLocationUpdated.DefaultSegmentDeliveryConfiguration != nil {
		if sourceLocationUpdated.DefaultSegmentDeliveryConfiguration.BaseUrl != nil {
			plan.DefaultSegmentDeliveryConfiguration.BaseUrl = sourceLocationUpdated.DefaultSegmentDeliveryConfiguration.BaseUrl
		}
	}
	if sourceLocationUpdated.SegmentDeliveryConfigurations != nil {
		for _, v := range sourceLocationUpdated.SegmentDeliveryConfigurations {
			*plan.SegmentDeliveryConfigurations = append(*plan.SegmentDeliveryConfigurations, resourceSourceLocationSegmentDeliveryConfigurationsModel{
				BaseUrl: v.BaseUrl,
				Name:    v.Name,
			})
		}
	}
	plan.HttpConfiguration.BaseUrl = sourceLocationUpdated.HttpConfiguration.BaseUrl
	plan.ID = types.StringValue(*sourceLocationUpdated.SourceLocationName)

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

	name := state.SourceLocationName

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
