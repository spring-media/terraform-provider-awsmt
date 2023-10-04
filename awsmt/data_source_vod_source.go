package awsmt

import (
	"context"

	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
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

func (d *dataSourceVodSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vod_source"
}

func (d *dataSourceVodSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = buildDatasourceSchema()
}

func (d *dataSourceVodSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*mediatailor.MediaTailor)
}

func (d *dataSourceVodSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data vodSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceLocationName := data.SourceLocationName
	vodSourceName := data.Name

	vodSource, err := d.client.DescribeVodSource(&mediatailor.DescribeVodSourceInput{SourceLocationName: sourceLocationName, VodSourceName: vodSourceName})
	if err != nil {
		resp.Diagnostics.AddError("Error while describing vod source", err.Error())
		return
	}

	data = readVodSourceToPlan(data, mediatailor.CreateVodSourceOutput(*vodSource))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
