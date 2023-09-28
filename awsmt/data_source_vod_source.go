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
	_ datasource.DataSource              = &dataSourceVodSource{}
	_ datasource.DataSourceWithConfigure = &dataSourceVodSource{}
)

func DataSourceVodSource() datasource.DataSource {
	return &dataSourceVodSource{}
}

type dataSourceVodSource struct {
	client *mediatailor.MediaTailor
}
type dataSourceVodSourceModel struct {
	ID                        types.String                       `tfsdk:"id"`
	Arn                       types.String                       `tfsdk:"arn"`
	CreationTime              types.String                       `tfsdk:"creation_time"`
	HttpPackageConfigurations []httpPackageConfigurationsDSModel `tfsdk:"http_package_configurations"`
	LastModifiedTime          types.String                       `tfsdk:"last_modified_time"`
	SourceLocationName        *string                            `tfsdk:"source_location_name"`
	Tags                      map[string]*string                 `tfsdk:"tags"`
	VodSourceName             *string                            `tfsdk:"vod_source_name"`
}

type httpPackageConfigurationsDSModel struct {
	Path        types.String `tfsdk:"path"`
	SourceGroup types.String `tfsdk:"source_group"`
	Type        types.String `tfsdk:"type"`
}

func (d *dataSourceVodSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vod_source"
}

func (d *dataSourceVodSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":            computedString,
			"arn":           computedString,
			"creation_time": computedString,
			"http_package_configurations": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"path":         computedString,
						"source_group": computedString,
						"type": schema.StringAttribute{
							Computed: true,
							Validators: []validator.String{
								stringvalidator.OneOf("HLS", "DASH"),
							},
						},
					},
				},
			},
			"last_modified_time":   computedString,
			"source_location_name": requiredString,
			"tags":                 computedMap,
			"vod_source_name":      requiredString,
		},
	}
}

func (d *dataSourceVodSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*mediatailor.MediaTailor)
}

func (d *dataSourceVodSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dataSourceVodSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceLocationName := data.SourceLocationName
	vodSourceName := data.VodSourceName

	vodSource, err := d.client.DescribeVodSource(&mediatailor.DescribeVodSourceInput{SourceLocationName: sourceLocationName, VodSourceName: vodSourceName})
	if err != nil {
		resp.Diagnostics.AddError("Error while describing vod source", err.Error())
		return
	}

	data.ID = types.StringValue(*vodSource.VodSourceName)
	if vodSource.Arn != nil && *vodSource.Arn != "" {
		data.Arn = types.StringValue(*vodSource.Arn)
	}
	if vodSource.CreationTime != nil {
		data.CreationTime = types.StringValue((aws.TimeValue(vodSource.CreationTime)).String())
	}

	if vodSource.HttpPackageConfigurations != nil && len(vodSource.HttpPackageConfigurations) > 0 {
		data.HttpPackageConfigurations = []httpPackageConfigurationsDSModel{}
		for _, httpPackageConfiguration := range vodSource.HttpPackageConfigurations {
			httpPackageConfigurations := httpPackageConfigurationsDSModel{}
			httpPackageConfigurations.Path = types.StringValue(*httpPackageConfiguration.Path)
			httpPackageConfigurations.SourceGroup = types.StringValue(*httpPackageConfiguration.SourceGroup)
			httpPackageConfigurations.Type = types.StringValue(*httpPackageConfiguration.Type)
			data.HttpPackageConfigurations = append(data.HttpPackageConfigurations, httpPackageConfigurations)
		}
	}

	if vodSource.LastModifiedTime != nil {
		data.LastModifiedTime = types.StringValue((aws.TimeValue(vodSource.LastModifiedTime)).String())
	}

	if vodSource.SourceLocationName != nil && *vodSource.SourceLocationName != "" {
		data.SourceLocationName = vodSource.SourceLocationName
	}

	if len(vodSource.Tags) > 0 {
		data.Tags = make(map[string]*string)
		for key, value := range vodSource.Tags {
			data.Tags[key] = value
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
