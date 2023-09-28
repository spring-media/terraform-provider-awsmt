package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go/service/mediatailor"
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
	client *mediatailor.MediaTailor
}
type dataSourcePlaybackConfigurationModel struct {
	ID                                  types.String                   `tfsdk:"id"`
	AdDecisionServerUrl                 *string                        `tfsdk:"ad_decision_server_url"`
	AvailSupression                     *availSupressionModel          `tfsdk:"avail_supression"`
	Bumper                              *bumperModel                   `tfsdk:"bumper"`
	CdnConfiguration                    *cdnConfigurationModel         `tfsdk:"cdn_configuration"`
	ConfigurationAliases                map[string]map[string]*string  `tfsdk:"configuration_aliases"`
	DashConfiguration                   *dashConfigurationModel        `tfsdk:"dash_configuration"`
	HlsConfiguration                    *hlsConfigurationModel         `tfsdk:"hls_configuration"`
	LivePreRollConfiguration            *livePreRollConfigurationModel `tfsdk:"live_pre_roll_configuration"`
	LogConfiguration                    *logConfigurationModel         `tfsdk:"log_configuration"`
	ManifestProcessingRules             *manifestProcessingRulesModel  `tfsdk:"manifest_processing_rules"`
	Name                                *string                        `tfsdk:"name"`
	PersonalizationThresholdSeconds     *int64                         `tfsdk:"personalization_threshold_seconds"`
	PlaybackConfigurationArn            *string                        `tfsdk:"playback_configuration_arn"`
	PlaybackEndpointPrefix              *string                        `tfsdk:"playback_endpoint_prefix"`
	SessionInitializationEndpointPrefix *string                        `tfsdk:"session_initialization_endpoint_prefix"`
	SlateAdUrl                          *string                        `tfsdk:"slate_ad_url"`
	Tags                                map[string]*string             `tfsdk:"tags"`
	TranscodeProfileName                *string                        `tfsdk:"transcode_profile_name"`
	VideoContentSourceUrl               *string                        `tfsdk:"video_content_source_url"`
}

type availSupressionModel struct {
	FillPolicy *string `tfsdk:"fill_policy"`
	Mode       *string `tfsdk:"mode"`
	Value      *string `tfsdk:"value"`
}

type bumperModel struct {
	EndUrl   *string `tfsdk:"end_url"`
	StartUrl *string `tfsdk:"start_url"`
}

type cdnConfigurationModel struct {
	AdSegmentUrlPrefix      *string `tfsdk:"ad_segment_url_prefix"`
	ContentSegmentUrlPrefix *string `tfsdk:"content_segment_url_prefix"`
}

type dashConfigurationModel struct {
	ManifestEndpointPrefix *string `tfsdk:"manifest_endpoint_prefix"`
	MpdLocation            *string `tfsdk:"mpd_location"`
	OriginManifestType     *string `tfsdk:"origin_manifest_type"`
}

type hlsConfigurationModel struct {
	ManifestEndpointPrefix *string `tfsdk:"manifest_endpoint_prefix"`
}

type livePreRollConfigurationModel struct {
	AdDecisionServerUrl *string `tfsdk:"ad_decision_server_url"`
	MaxDurationSeconds  *int64  `tfsdk:"max_duration_seconds"`
}

type logConfigurationModel struct {
	PercentEnabled *int64 `tfsdk:"percent_enabled"`
}

type manifestProcessingRulesModel struct {
	AdMarkerPassthrough *adMarkerPassthroughModel `tfsdk:"ad_marker_passthrough"`
}
type adMarkerPassthroughModel struct {
	Enabled *bool `tfsdk:"enabled"`
}

func (d *dataSourcePlaybackConfiguration) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_playback_configuration"
}

func (d *dataSourcePlaybackConfiguration) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                     computedString,
			"ad_decision_server_url": computedString,
			"avail_supression": schema.SingleNestedAttribute{
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
			"hls_configuration": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"manifest_endpoint_prefix": computedString,
				},
			},
			"live_pre_roll_configuration": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"ad_decision_server_url": computedString,
					"max_duration_seconds":   computedInt64,
				},
			},
			"log_configuration": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"percent_enabled": computedInt64,
				},
			},
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

	d.client = req.ProviderData.(*mediatailor.MediaTailor)
}

func (d *dataSourcePlaybackConfiguration) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dataSourcePlaybackConfigurationModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name

	playbackConfiguration, err := d.client.GetPlaybackConfiguration(&mediatailor.GetPlaybackConfigurationInput{Name: name})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while retrieving the playback configuration "+err.Error(),
			err.Error(),
		)
		return
	}

	data = readPlaybackConfigToData(data, *playbackConfiguration)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
