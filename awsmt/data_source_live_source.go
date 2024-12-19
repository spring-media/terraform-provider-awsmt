package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"terraform-provider-mediatailor/awsmt/models"
)

var (
	_ datasource.DataSource              = &dataSourceLiveSource{}
	_ datasource.DataSourceWithConfigure = &dataSourceLiveSource{}
)

func DataSourceLiveSource() datasource.DataSource {
	return &dataSourceLiveSource{}
}

type dataSourceLiveSource struct {
	client *mediatailor.Client
}

func (d *dataSourceLiveSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_live_source"
}

func (d *dataSourceLiveSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                          computedString,
			"arn":                         computedString,
			"creation_time":               computedString,
			"http_package_configurations": httpPackageConfigurationsDataSourceSchema,
			"last_modified_time":          computedString,
			"name":                        requiredString,
			"source_location_name":        requiredString,
			"tags":                        computedMap,
		},
	}
}

func (d *dataSourceLiveSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*mediatailor.Client)
}

func (d *dataSourceLiveSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.LiveSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceLocationName := data.SourceLocationName
	liveSourceName := data.Name

	liveSource, err := d.client.DescribeLiveSource(ctx, &mediatailor.DescribeLiveSourceInput{SourceLocationName: sourceLocationName, LiveSourceName: liveSourceName})
	if err != nil {
		resp.Diagnostics.AddError("Error while describing live source", err.Error())
		return
	}

	data = readLiveSource(data, mediatailor.CreateLiveSourceOutput(*liveSource))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
