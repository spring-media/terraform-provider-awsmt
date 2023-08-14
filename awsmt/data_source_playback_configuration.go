package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
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
			"id": schema.StringAttribute{
				Computed: true,
			},
			"ad_decision_server_url": schema.StringAttribute{
				Computed: true,
			},
			"avail_supression": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"fill_policy": schema.StringAttribute{
						Computed:   true,
						CustomType: types.StringType,
					},
					"mode": schema.StringAttribute{
						Computed:   true,
						CustomType: types.StringType,
					},
					"value": schema.StringAttribute{
						Computed:   true,
						CustomType: types.StringType,
					},
				},
			},
			"bumper": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"end_url": schema.StringAttribute{
						Computed: true,
					},
					"start_url": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"cdn_configuration": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"ad_segment_url_prefix": schema.StringAttribute{
						Computed: true,
					},
					"content_segment_url_prefix": schema.StringAttribute{
						Computed: true,
					},
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
					"manifest_endpoint_prefix": schema.StringAttribute{
						Computed: true,
					},
					"mpd_location": schema.StringAttribute{
						Computed: true,
					},
					"origin_manifest_type": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"hls_configuration": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"manifest_endpoint_prefix": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"live_pre_roll_configuration": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"ad_decision_server_url": schema.StringAttribute{
						Computed: true,
					},
					"max_duration_seconds": schema.Int64Attribute{
						Computed: true,
					},
				},
			},
			"log_configuration": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"percent_enabled": schema.Int64Attribute{
						Computed: true,
					},
				},
			},
			"manifest_processing_rules": schema.SingleNestedAttribute{
				Computed: true,

				Attributes: map[string]schema.Attribute{
					"ad_marker_passthrough": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Computed: true,
							},
						},
					},
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"personalization_threshold_seconds": schema.Int64Attribute{
				Computed: true,
			},
			"playback_configuration_arn": schema.StringAttribute{
				Computed: true,
			},
			"playback_endpoint_prefix": schema.StringAttribute{
				Computed: true,
			},
			"session_initialization_endpoint_prefix": schema.StringAttribute{
				Computed: true,
			},
			"slate_ad_url": schema.StringAttribute{
				Computed: true,
			},
			"tags": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"transcode_profile_name": schema.StringAttribute{
				Computed: true,
			},
			"video_content_source_url": schema.StringAttribute{
				Computed: true,
			},
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

	data.AdDecisionServerUrl = playbackConfiguration.AdDecisionServerUrl

	// AVAIL SUPRESSION
	if playbackConfiguration.AvailSuppression != nil {
		data.AvailSupression = &availSupressionModel{}
		if playbackConfiguration.AvailSuppression.Mode != nil {
			data.AvailSupression.Mode = playbackConfiguration.AvailSuppression.Mode
		}
		if playbackConfiguration.AvailSuppression.Value != nil {
			data.AvailSupression.Value = playbackConfiguration.AvailSuppression.Value
		}
		if playbackConfiguration.AvailSuppression.FillPolicy != nil {
			data.AvailSupression.FillPolicy = playbackConfiguration.AvailSuppression.FillPolicy
		}

	}
	// BUMPER
	if playbackConfiguration.Bumper != nil {
		data.Bumper = &bumperModel{}
		if playbackConfiguration.Bumper.EndUrl != nil {
			data.Bumper.EndUrl = playbackConfiguration.Bumper.EndUrl
		}
		if playbackConfiguration.Bumper.StartUrl != nil {
			data.Bumper.StartUrl = playbackConfiguration.Bumper.StartUrl
		}
	}
	// CDN CONFIGURATION
	if playbackConfiguration.CdnConfiguration != nil {
		data.CdnConfiguration = &cdnConfigurationModel{}
		if playbackConfiguration.CdnConfiguration.AdSegmentUrlPrefix != nil {
			data.CdnConfiguration.AdSegmentUrlPrefix = playbackConfiguration.CdnConfiguration.AdSegmentUrlPrefix
		}
		if playbackConfiguration.CdnConfiguration.ContentSegmentUrlPrefix != nil {
			data.CdnConfiguration.ContentSegmentUrlPrefix = playbackConfiguration.CdnConfiguration.ContentSegmentUrlPrefix
		}
	}

	if playbackConfiguration.ConfigurationAliases != nil {
		data.ConfigurationAliases = playbackConfiguration.ConfigurationAliases
	}
	// DASH CONFIGURATION
	if playbackConfiguration.DashConfiguration != nil {
		data.DashConfiguration = &dashConfigurationModel{}
		if playbackConfiguration.DashConfiguration.ManifestEndpointPrefix != nil {
			data.DashConfiguration.ManifestEndpointPrefix = playbackConfiguration.DashConfiguration.ManifestEndpointPrefix
		}
		if playbackConfiguration.DashConfiguration.MpdLocation != nil {
			data.DashConfiguration.MpdLocation = playbackConfiguration.DashConfiguration.MpdLocation
		}
		if playbackConfiguration.DashConfiguration.OriginManifestType != nil {
			data.DashConfiguration.OriginManifestType = playbackConfiguration.DashConfiguration.OriginManifestType
		}
	}
	// HLS CONFIGURATION
	if playbackConfiguration.HlsConfiguration != nil {
		data.HlsConfiguration = &hlsConfigurationModel{}
		if playbackConfiguration.HlsConfiguration.ManifestEndpointPrefix != nil {
			data.HlsConfiguration.ManifestEndpointPrefix = playbackConfiguration.HlsConfiguration.ManifestEndpointPrefix
		}
	}
	// LIVE PRE ROLL CONFIGURATION
	if playbackConfiguration.LivePreRollConfiguration != nil {
		data.LivePreRollConfiguration = &livePreRollConfigurationModel{}
		if playbackConfiguration.LivePreRollConfiguration.AdDecisionServerUrl != nil {
			data.LivePreRollConfiguration.AdDecisionServerUrl = playbackConfiguration.LivePreRollConfiguration.AdDecisionServerUrl
		}
		if playbackConfiguration.LivePreRollConfiguration.MaxDurationSeconds != nil {
			data.LivePreRollConfiguration.MaxDurationSeconds = playbackConfiguration.LivePreRollConfiguration.MaxDurationSeconds
		}
	}
	// LOG CONFIGURATION
	if playbackConfiguration.LogConfiguration != nil {
		data.LogConfiguration = &logConfigurationModel{}
		if playbackConfiguration.LogConfiguration.PercentEnabled != nil {
			data.LogConfiguration.PercentEnabled = playbackConfiguration.LogConfiguration.PercentEnabled
		}
	} else {
		data.LogConfiguration = &logConfigurationModel{}
		data.LogConfiguration.PercentEnabled = aws.Int64(0)
	}
	// MANIFEST PROCESSING RULES
	if playbackConfiguration.ManifestProcessingRules != nil {
		data.ManifestProcessingRules = &manifestProcessingRulesModel{}
		if playbackConfiguration.ManifestProcessingRules.AdMarkerPassthrough != nil {
			data.ManifestProcessingRules.AdMarkerPassthrough = &adMarkerPassthroughModel{}
			if playbackConfiguration.ManifestProcessingRules.AdMarkerPassthrough.Enabled != nil {
				data.ManifestProcessingRules.AdMarkerPassthrough.Enabled = playbackConfiguration.ManifestProcessingRules.AdMarkerPassthrough.Enabled
			}
		}
	}

	data.Name = playbackConfiguration.Name
	data.PersonalizationThresholdSeconds = playbackConfiguration.PersonalizationThresholdSeconds
	data.PlaybackConfigurationArn = playbackConfiguration.PlaybackConfigurationArn
	data.PlaybackEndpointPrefix = playbackConfiguration.PlaybackEndpointPrefix
	data.SessionInitializationEndpointPrefix = playbackConfiguration.SessionInitializationEndpointPrefix
	data.SlateAdUrl = playbackConfiguration.SlateAdUrl
	data.Tags = playbackConfiguration.Tags
	data.TranscodeProfileName = playbackConfiguration.TranscodeProfileName
	data.VideoContentSourceUrl = playbackConfiguration.VideoContentSourceUrl
	data.ID = types.StringValue(*playbackConfiguration.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
