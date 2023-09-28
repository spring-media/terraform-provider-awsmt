package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	var data vodSourceModel
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

	data = readVodSourceToPlan(data, mediatailor.CreateVodSourceOutput(*vodSource))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
