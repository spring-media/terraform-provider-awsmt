package awsmt

import (
	"context"

	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var (
	_ datasource.DataSource              = &dataSourceLiveSource{}
	_ datasource.DataSourceWithConfigure = &dataSourceLiveSource{}
)

func DataSourceLiveSource() datasource.DataSource {
	return &dataSourceLiveSource{}
}

type dataSourceLiveSource struct {
	client *mediatailor.MediaTailor
}

func (d *dataSourceLiveSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_live_source"
}

func (d *dataSourceLiveSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = buildDatasourceSchema()
}

func (d *dataSourceLiveSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*mediatailor.MediaTailor)
}

func (d *dataSourceLiveSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data liveSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceLocationName := data.SourceLocationName
	liveSourceName := data.Name

	liveSource, err := d.client.DescribeLiveSource(&mediatailor.DescribeLiveSourceInput{SourceLocationName: sourceLocationName, LiveSourceName: liveSourceName})
	if err != nil {
		resp.Diagnostics.AddError("Error while describing live source", err.Error())
		return
	}

	data = readLiveSourceToPlan(data, mediatailor.CreateLiveSourceOutput(*liveSource))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
