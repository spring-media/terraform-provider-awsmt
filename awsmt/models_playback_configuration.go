package awsmt

import "github.com/hashicorp/terraform-plugin-framework/types"

type playbackConfigurationModel struct {
	ID                   types.String                  `tfsdk:"id"`
	AdDecisionServerUrl  *string                       `tfsdk:"ad_decision_server_url"`
	AvailSupression      *availSupressionModel         `tfsdk:"avail_supression"`
	Bumper               *bumperModel                  `tfsdk:"bumper"`
	CdnConfiguration     *cdnConfigurationModel        `tfsdk:"cdn_configuration"`
	ConfigurationAliases map[string]map[string]*string `tfsdk:"configuration_aliases"`
	DashConfiguration    *dashConfigurationModel       `tfsdk:"dash_configuration"`
	// @ADR
	// Context: The Provider Framework does not allow computed blocks
	// Decision: We decided to flatten the Log Configuration and the HLS Configuration blocks into the resource.
	// Consequences: The schema of the object differs from that of the SDK.
	HlsConfigurationManifestEndpointPrefix types.String                   `tfsdk:"hls_configuration_manifest_endpoint_prefix"`
	LogConfigurationPercentEnabled         types.Int64                    `tfsdk:"log_configuration_percent_enabled"`
	LivePreRollConfiguration               *livePreRollConfigurationModel `tfsdk:"live_pre_roll_configuration"`
	ManifestProcessingRules                *manifestProcessingRulesModel  `tfsdk:"manifest_processing_rules"`
	Name                                   *string                        `tfsdk:"name"`
	PersonalizationThresholdSeconds        *int64                         `tfsdk:"personalization_threshold_seconds"`
	PlaybackConfigurationArn               types.String                   `tfsdk:"playback_configuration_arn"`
	PlaybackEndpointPrefix                 types.String                   `tfsdk:"playback_endpoint_prefix"`
	SessionInitializationEndpointPrefix    types.String                   `tfsdk:"session_initialization_endpoint_prefix"`
	SlateAdUrl                             *string                        `tfsdk:"slate_ad_url"`
	Tags                                   map[string]*string             `tfsdk:"tags"`
	TranscodeProfileName                   *string                        `tfsdk:"transcode_profile_name"`
	VideoContentSourceUrl                  *string                        `tfsdk:"video_content_source_url"`
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
	ManifestEndpointPrefix types.String `tfsdk:"manifest_endpoint_prefix"`
	MpdLocation            *string      `tfsdk:"mpd_location"`
	OriginManifestType     *string      `tfsdk:"origin_manifest_type"`
}

type livePreRollConfigurationModel struct {
	AdDecisionServerUrl *string `tfsdk:"ad_decision_server_url"`
	MaxDurationSeconds  *int64  `tfsdk:"max_duration_seconds"`
}

type manifestProcessingRulesModel struct {
	AdMarkerPassthrough *adMarkerPassthroughModel `tfsdk:"ad_marker_passthrough"`
}
type adMarkerPassthroughModel struct {
	Enabled *bool `tfsdk:"enabled"`
}
