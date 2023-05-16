package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &dataSourceSourceLocation{}
	_ datasource.DataSourceWithConfigure = &dataSourceSourceLocation{}
)

func DataSourceSourceLocation() datasource.DataSource {
	return &dataSourceSourceLocation{}
}

type dataSourceSourceLocation struct {
	client *mediatailor.MediaTailor
}

type dataSourceSourceLocationModel struct {
	ID                                  types.String                                           `tfsdk:"id"`
	AccessConfiguration                 sourceLocationAccessConfigurationModel                 `tfsdk:"access_configuration"`
	Arn                                 types.String                                           `tfsdk:"arn"`
	CreationTime                        types.String                                           `tfsdk:"creation_time"`
	DefaultSegmentDeliveryConfiguration sourceLocationDefaultSegmentDeliveryConfigurationModel `tfsdk:"default_segment_delivery_configuration"`
	HttpConfiguration                   sourceLocationHttpConfigurationModel                   `tfsdk:"http_configuration"`
	LastModifiedTime                    types.String                                           `tfsdk:"last_modified_time"`
	SegmentDeliveryConfigurations       []sourceLocationSegmentDeliveryConfigurationsModel     `tfsdk:"segment_delivery_configuration"`
	SourceLocationName                  *string                                                `tfsdk:"source_location_name"`
	tags                                map[string]*string                                     `tfsdk:"tags"`
}

type sourceLocationAccessConfigurationModel struct {
	AccessType                             types.String                                              `tfsdk:"access_type"`
	SecretsManagerAccessTokenConfiguration sourceLocationSecretsManagerAccessTokenConfigurationModel `tfsdk:"smatc"`
	// SMATC is short for Secret Manager Access Token Configuration
}

type sourceLocationSecretsManagerAccessTokenConfigurationModel struct {
	HeaderName      types.String `tfsdk:"header_name"`
	SecretArn       types.String `tfsdk:"secret_arn"`
	SecretStringKey types.String `tfsdk:"secret_string_key"`
}

type sourceLocationDefaultSegmentDeliveryConfigurationModel struct {
	BaseUrl types.String `tfsdk:"base_url"`
}

type sourceLocationHttpConfigurationModel struct {
	BaseUrl types.String `tfsdk:"base_url"`
}

type sourceLocationSegmentDeliveryConfigurationsModel struct {
	BaseUrl types.String `tfsdk:"base_url"`
	Name    types.String `tfsdk:"name"`
}

func (d *dataSourceSourceLocation) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_source"
}

func (d *dataSourceSourceLocation) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"access_configuration": schema.MapNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"access_type": schema.StringAttribute{
							Computed: true,
						},
						"smatc": schema.MapNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"header_name": schema.StringAttribute{
										Computed: true,
									},
									"secret_arn": schema.StringAttribute{
										Computed: true,
									},
									"secret_string_key": schema.StringAttribute{
										Computed: true,
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
			"default_segment_delivery_configuration": schema.MapNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"base_url": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"http_configuration": schema.MapNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"base_url": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"last_modified_time": schema.StringAttribute{
				Computed: true,
			},
			"segment_delivery_configuration": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"base_url": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"source_location_name": schema.StringAttribute{
				Required: true,
			},
			"tags": schema.MapAttribute{
				Computed: true,
			},
		},
	}
}

func (d *dataSourceSourceLocation) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*mediatailor.MediaTailor)
}

func (d *dataSourceSourceLocation) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dataSourceSourceLocationModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.SourceLocationName

	sourceLocation, err := d.client.DescribeSourceLocation(&mediatailor.DescribeSourceLocationInput{SourceLocationName: name})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while retrieving the source location "+err.Error(),
			err.Error(),
		)
		return
	}

	data.AccessConfiguration.AccessType = types.StringValue(*sourceLocation.AccessConfiguration.AccessType)
	data.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName = types.StringValue(*sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName)
	data.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn = types.StringValue(*sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn)
	data.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey = types.StringValue(*sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey)
	data.Arn = types.StringValue(*sourceLocation.Arn)
	data.CreationTime = types.StringValue(aws.TimeValue(sourceLocation.CreationTime).String())
	data.DefaultSegmentDeliveryConfiguration.BaseUrl = types.StringValue(*sourceLocation.DefaultSegmentDeliveryConfiguration.BaseUrl)
	data.HttpConfiguration.BaseUrl = types.StringValue(*sourceLocation.HttpConfiguration.BaseUrl)
	data.LastModifiedTime = types.StringValue(aws.TimeValue(sourceLocation.LastModifiedTime).String())

	var segmentDeliveryConfigurations []sourceLocationSegmentDeliveryConfigurationsModel
	if data.SegmentDeliveryConfigurations != nil && len(data.SegmentDeliveryConfigurations) > 0 {
		for _, o := range data.SegmentDeliveryConfigurations {
			segmentDeliveryConfiguration := sourceLocationSegmentDeliveryConfigurationsModel{}
			segmentDeliveryConfiguration.BaseUrl = o.BaseUrl

			segmentDeliveryConfigurations = append(segmentDeliveryConfigurations, segmentDeliveryConfiguration)
		}
	}
	data.SegmentDeliveryConfigurations = segmentDeliveryConfigurations
	data.tags = sourceLocation.Tags

	data.ID = types.StringValue(aws.StringValue(sourceLocation.SourceLocationName))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
