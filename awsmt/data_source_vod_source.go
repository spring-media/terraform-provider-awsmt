package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var (
	_ datasource.DataSource              = &dataSourceVodSource{}
	_ datasource.DataSourceWithConfigure = &dataSourceVodSource{}
)

func DataSourceVodSource() datasource.DataSource {
	return &dataSourceVodSource{}
}

type dataSourceVodSource struct {
	client *mediatailor.Client
}

func (d *dataSourceVodSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vod_source"
}

func (d *dataSourceVodSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                                   computedString,
			"source_location_name":                 requiredString,
			"http_package_configurations":          httpPackageConfigurationsDataSourceSchema,
			"creation_time":                        computedString,
			"tags":                                 computedMap,
			"last_modified_time":                   computedString,
			"arn":                                  computedString,
			"name":                                 requiredString,
			"ad_break_opportunities_offset_millis": computedInt64List,
		},
	}
}

func (d *dataSourceVodSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*mediatailor.Client)
}

func (d *dataSourceVodSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data vodSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceLocationName := data.SourceLocationName
	vodSourceName := data.Name

	vodSource, err := d.client.DescribeVodSource(ctx, &mediatailor.DescribeVodSourceInput{SourceLocationName: sourceLocationName, VodSourceName: vodSourceName})
	if err != nil {
		resp.Diagnostics.AddError("Error while describing vod source", err.Error())
		return
	}

	data = readVodSourceToState(data, *vodSource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
