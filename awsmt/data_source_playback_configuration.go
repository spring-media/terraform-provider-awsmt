package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &dataSourcePlaybackConfiguration{}
	_ datasource.DataSourceWithConfigure = &dataSourcePlaybackConfiguration{}
)

func DataSourcePlaybackConfiguration() datasource.DataSource {
	return &dataSourcePlaybackConfiguration{}
}

type dataSourcePlaybackConfiguration struct {
	client *mediatailor.Client
}

func (d *dataSourcePlaybackConfiguration) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_playback_configuration"
}

func (d *dataSourcePlaybackConfiguration) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                     computedString,
			"ad_decision_server_url": computedString,
			"avail_suppression": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"fill_policy": computedString,
					"mode":        computedString,
					"value":       computedString,
				},
			},
			"bumper": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"end_url":   computedString,
					"start_url": computedString,
				},
			},
			"cdn_configuration": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"ad_segment_url_prefix":      computedString,
					"content_segment_url_prefix": computedString,
				},
			},
			"configuration_aliases": schema.ListAttribute{
				Computed: true,
				ElementType: types.MapType{
					ElemType: types.MapType{
						ElemType: types.StringType,
					},
				},
			},
			"dash_configuration": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"manifest_endpoint_prefix": computedString,
					"mpd_location":             computedString,
					"origin_manifest_type":     computedString,
				},
			},
			"hls_configuration_manifest_endpoint_prefix": computedString,
			"live_pre_roll_configuration": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"ad_decision_server_url": computedString,
					"max_duration_seconds":   computedInt64,
				},
			},
			"log_configuration_percent_enabled": computedInt64,
			"manifest_processing_rules": schema.SingleNestedAttribute{
				Computed: true,

				Attributes: map[string]schema.Attribute{
					"ad_marker_passthrough": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"enabled": computedBool,
						},
					},
				},
			},
			"name":                                   requiredString,
			"personalization_threshold_seconds":      computedInt64,
			"playback_configuration_arn":             computedString,
			"playback_endpoint_prefix":               computedString,
			"session_initialization_endpoint_prefix": computedString,
			"slate_ad_url":                           computedString,
			"tags":                                   computedMap,
			"transcode_profile_name":                 computedString,
			"video_content_source_url":               computedString,
		},
	}
}

func (d *dataSourcePlaybackConfiguration) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(clients).v2
}

func (d *dataSourcePlaybackConfiguration) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data playbackConfigurationModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name

	playbackConfiguration, err := d.client.GetPlaybackConfiguration(context.TODO(), &mediatailor.GetPlaybackConfigurationInput{Name: name})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while retrieving the playback configuration "+err.Error(),
			err.Error(),
		)
		return
	}

	m := putPlaybackConfigurationModelbuilder{model: &data, output: mediatailor.PutPlaybackConfigurationOutput(*playbackConfiguration), isResource: false}

	resp.Diagnostics.Append(resp.State.Set(ctx, m.getModel())...)
}
