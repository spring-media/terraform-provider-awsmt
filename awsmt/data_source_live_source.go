package awsmt

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

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
			"name":                 requiredString,
			"source_location_name": requiredString,
			"tags":                 computedMap,
		},
	}
}

func (d *dataSourceLiveSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(clients).v1
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
