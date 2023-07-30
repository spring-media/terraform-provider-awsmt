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
	_ datasource.DataSource              = &dataSourceLiveSource{}
	_ datasource.DataSourceWithConfigure = &dataSourceLiveSource{}
)

func DataSourceLiveSource() datasource.DataSource {
	return &dataSourceLiveSource{}
}

type dataSourceLiveSource struct {
	client *mediatailor.MediaTailor
}
type dataSourceLiveSourceModel struct {
	ID                          types.String                         `tfsdk:"id"`
	Arn                         types.String                         `tfsdk:"arn"`
	CreationTime                types.String                         `tfsdk:"creation_time"`
	HttpPackageConfigurationsLS []httpPackageConfigurationsLSDSModel `tfsdk:"http_package_configuration"`
	LastModifiedTime            types.String                         `tfsdk:"last_modified_time"`
	LiveSourceName              types.String                         `tfsdk:"live_source_name"`
	SourceLocationName          types.String                         `tfsdk:"source_location_name"`
	Tags                        map[string]*string                   `tfsdk:"tags"`
}

type httpPackageConfigurationsLSDSModel struct {
	Path        types.String `tfsdk:"path"`
	SourceGroup types.String `tfsdk:"source_group"`
	Type        types.String `tfsdk:"type"`
}

func (d *dataSourceLiveSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_location"
}

func (d *dataSourceLiveSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"arn": schema.StringAttribute{
				Computed: true,
			},
			"creation_time": schema.StringAttribute{
				Computed: true,
			},
			"http_package_configuration": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"path": schema.StringAttribute{
							Computed: true,
						},
						"source_group": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
							Validators: []validator.String{
								stringvalidator.OneOf("HLS", "DASH"),
							},
						},
					},
				},
			},
			"last_modified_time": schema.StringAttribute{
				Computed: true,
			},
			"live_source_name": schema.StringAttribute{
				Required: true,
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

func (d *dataSourceLiveSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*mediatailor.MediaTailor)
}

func (d *dataSourceLiveSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dataSourceLiveSourceModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sourceLocationName := aws.String(data.SourceLocationName.String())
	liveSourceName := aws.String(data.LiveSourceName.String())

	liveSource, err := d.client.DescribeLiveSource(&mediatailor.DescribeLiveSourceInput{SourceLocationName: sourceLocationName, LiveSourceName: liveSourceName})
	if err != nil {
		resp.Diagnostics.AddError("Error while describing live source", err.Error())
		return
	}

	if liveSource.Arn != nil {
		data.Arn = types.StringValue(*liveSource.Arn)
	}

	if liveSource.CreationTime != nil {
		data.CreationTime = types.StringValue((aws.TimeValue(liveSource.CreationTime)).String())
	}

	if liveSource.HttpPackageConfigurations != nil && len(liveSource.HttpPackageConfigurations) > 0 {
		for _, httpPackageConfiguration := range liveSource.HttpPackageConfigurations {
			httpPackageConfigurations := httpPackageConfigurationsLSDSModel{}
			httpPackageConfigurations.Path = types.StringValue(*httpPackageConfiguration.Path)
			httpPackageConfigurations.SourceGroup = types.StringValue(*httpPackageConfiguration.SourceGroup)
			httpPackageConfigurations.Type = types.StringValue(*httpPackageConfiguration.Type)
			data.HttpPackageConfigurationsLS = append(data.HttpPackageConfigurationsLS, httpPackageConfigurations)
		}
	}

	if liveSource.LastModifiedTime != nil {
		data.LastModifiedTime = types.StringValue((aws.TimeValue(liveSource.LastModifiedTime)).String())
	}

	if liveSource.SourceLocationName != nil && *liveSource.SourceLocationName != "" {
		data.SourceLocationName = types.StringValue(*liveSource.SourceLocationName)
	}

	if liveSource.Tags != nil && len(liveSource.Tags) > 0 {
		data.Tags = make(map[string]*string)
		for key, value := range liveSource.Tags {
			data.Tags[key] = value
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
