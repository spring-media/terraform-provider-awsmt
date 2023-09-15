package awsmt

import (
	"context"
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
type dataSourceSourceLocationModel struct {
	ID                                  types.String                                `tfsdk:"id"`
	AccessConfiguration                 *accessConfigurationDSModel                 `tfsdk:"access_configuration"`
	Arn                                 types.String                                `tfsdk:"arn"`
	CreationTime                        types.String                                `tfsdk:"creation_time"`
	DefaultSegmentDeliveryConfiguration *defaultSegmentDeliveryConfigurationDSModel `tfsdk:"default_segment_delivery_configuration"`
	HttpConfiguration                   *httpConfigurationDSModel                   `tfsdk:"http_configuration"`
	LastModifiedTime                    types.String                                `tfsdk:"last_modified_time"`
	SegmentDeliveryConfigurations       []segmentDeliveryConfigurationsDSModel      `tfsdk:"segment_delivery_configurations"`
	SourceLocationName                  *string                                     `tfsdk:"source_location_name"`
	Tags                                map[string]*string                          `tfsdk:"tags"`
}

type accessConfigurationDSModel struct {
	AccessType                             types.String                                   `tfsdk:"access_type"`
	SecretsManagerAccessTokenConfiguration *secretsManagerAccessTokenConfigurationDSModel `tfsdk:"secrets_manager_access_token_configuration"`
}

type secretsManagerAccessTokenConfigurationDSModel struct {
	HeaderName      types.String `tfsdk:"header_name"`
	SecretArn       types.String `tfsdk:"secret_arn"`
	SecretStringKey types.String `tfsdk:"secret_string_key"`
}

type defaultSegmentDeliveryConfigurationDSModel struct {
	BaseUrl types.String `tfsdk:"dsdc_base_url"`
}

type httpConfigurationDSModel struct {
	BaseUrl types.String `tfsdk:"hc_base_url"`
}

type segmentDeliveryConfigurationsDSModel struct {
	BaseUrl types.String `tfsdk:"sdc_base_url"`
	SDCName types.String `tfsdk:"sdc_name"`
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
			"access_configuration": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"access_type": schema.StringAttribute{
						Computed: true,
						Validators: []validator.String{
							stringvalidator.OneOf("S3_SIGV4", "SECRETS_MANAGER_ACCESS_TOKEN"),
						},
					},
					"secrets_manager_access_token_configuration": schema.SingleNestedAttribute{
						Computed: true,
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
			"arn": schema.StringAttribute{
				Computed: true,
			},
			"creation_time": schema.StringAttribute{
				Computed: true,
			},
			"default_segment_delivery_configuration": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"dsdc_base_url": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"http_configuration": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"hc_base_url": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"last_modified_time": schema.StringAttribute{
				Computed: true,
			},
			"segment_delivery_configurations": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"sdc_base_url": schema.StringAttribute{
							Computed: true,
						},
						"sdc_name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
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

	sourceLocationName := data.SourceLocationName

	sourceLocation, err := d.client.DescribeSourceLocation(&mediatailor.DescribeSourceLocationInput{SourceLocationName: sourceLocationName})
	if err != nil {
		resp.Diagnostics.AddError("Error while describing source location", err.Error())
		return
	}

	data = readSourceLocationToData(data, *sourceLocation)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
