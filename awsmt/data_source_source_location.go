package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
type dataSourceSourceLocationlModel struct {
	ID                                  types.String                                            `tfsdk:"id"`
	AccessConfiguration                 *sourceLocationAccessConfigurationModel                 `tfsdk:"access_configuration"`
	Arn                                 types.String                                            `tfsdk:"arn"`
	CreationTime                        types.String                                            `tfsdk:"creation_time"`
	DefaultSegmentDeliveryConfiguration *sourceLocationDefaultSegmentDeliveryConfigurationModel `tfsdk:"default_segment_delivery_configuration"`
	HttpConfiguration                   *sourceLocationHttpConfigurationModel                   `tfsdk:"http_configuration"`
	LastModifiedTime                    types.String                                            `tfsdk:"last_modified_time"`
	SegmentDeliveryConfigurations       *[]sourceLocationSegmentDeliveryConfigurationsModel     `tfsdk:"segment_delivery_configuration"`
	SourceLocationName                  *string                                                 `tfsdk:"name"`
	Tags                                map[string]*string                                      `tfsdk:"tags"`
}

type sourceLocationAccessConfigurationModel struct {
	AccessType                             *string                                                    `tfsdk:"access_type"`
	SecretsManagerAccessTokenConfiguration *sourceLocationSecretsManagerAccessTokenConfigurationModel `tfsdk:"smatc"`
	// SMATC is short for Secret Manager Access Token Configuration
}

type sourceLocationSecretsManagerAccessTokenConfigurationModel struct {
	HeaderName      *string `tfsdk:"header_name"`
	SecretArn       *string `tfsdk:"secret_arn"`
	SecretStringKey *string `tfsdk:"secret_string_key"`
}

type sourceLocationDefaultSegmentDeliveryConfigurationModel struct {
	BaseUrl *string `tfsdk:"base_url"`
}

type sourceLocationHttpConfigurationModel struct {
	BaseUrl *string `tfsdk:"base_url"`
}

type sourceLocationSegmentDeliveryConfigurationsModel struct {
	BaseUrl *string `tfsdk:"base_url"`
	Name    *string `tfsdk:"name"`
}

func (d *dataSourceSourceLocation) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_location"
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
							Validators: []validator.String{
								stringvalidator.OneOf("S3_SIGV4", "SECRETS_MANAGER_ACCESS_TOKEN"),
							},
							CustomType: types.StringType,
						},
						"smatc": schema.MapNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"header_name": schema.StringAttribute{
										Computed:   true,
										CustomType: types.StringType,
									},
									"secret_arn": schema.StringAttribute{
										Computed:   true,
										CustomType: types.StringType,
									},
									"secret_string_key": schema.StringAttribute{
										Computed:   true,
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
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"base_url": schema.StringAttribute{
						Computed:   true,
						CustomType: types.StringType,
					},
				},
			},
			"http_configuration": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"base_url": schema.StringAttribute{
						Computed: true,
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
							Computed:   true,
							CustomType: types.StringType,
						},
						"name": schema.StringAttribute{
							Computed:   true,
							CustomType: types.StringType,
						},
					},
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"tags": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
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
	var data dataSourceSourceLocationlModel
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

	data.Arn = types.StringValue(aws.StringValue(sourceLocation.Arn))
	data.SourceLocationName = sourceLocation.SourceLocationName
	data.CreationTime = types.StringValue((aws.TimeValue(sourceLocation.CreationTime)).String())
	data.LastModifiedTime = types.StringValue((aws.TimeValue(sourceLocation.LastModifiedTime)).String())

	data.AccessConfiguration = &sourceLocationAccessConfigurationModel{}
	if sourceLocation.AccessConfiguration != nil {
		if sourceLocation.AccessConfiguration.AccessType != nil {
			data.AccessConfiguration.AccessType = sourceLocation.AccessConfiguration.AccessType
		}
	}
	data.AccessConfiguration.SecretsManagerAccessTokenConfiguration = &sourceLocationSecretsManagerAccessTokenConfigurationModel{}
	if sourceLocation.AccessConfiguration != nil && sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration != nil {
		if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName != nil {
			data.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName = sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName
		}
		if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn != nil {
			data.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn = sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn
		}
		if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey != nil {
			data.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey = sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey
		}
	}

	data.DefaultSegmentDeliveryConfiguration = &sourceLocationDefaultSegmentDeliveryConfigurationModel{}
	if sourceLocation.DefaultSegmentDeliveryConfiguration != nil {
		if sourceLocation.DefaultSegmentDeliveryConfiguration.BaseUrl != nil {
			data.DefaultSegmentDeliveryConfiguration.BaseUrl = sourceLocation.DefaultSegmentDeliveryConfiguration.BaseUrl
		}
	}

	data.HttpConfiguration = &sourceLocationHttpConfigurationModel{}
	if sourceLocation.HttpConfiguration != nil {
		if sourceLocation.HttpConfiguration.BaseUrl != nil {
			data.HttpConfiguration.BaseUrl = sourceLocation.HttpConfiguration.BaseUrl
		}
	}

	data.SegmentDeliveryConfigurations = &[]sourceLocationSegmentDeliveryConfigurationsModel{}

	data.ID = types.StringValue(aws.StringValue(sourceLocation.SourceLocationName))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
