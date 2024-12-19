package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type PlaybackConfigurationModel struct {
	ID                   types.String                 `tfsdk:"id"`
	AdDecisionServerUrl  *string                      `tfsdk:"ad_decision_server_url"`
	AvailSuppression     *AvailSuppressionModel       `tfsdk:"avail_suppression"`
	Bumper               *BumperModel                 `tfsdk:"bumper"`
	CdnConfiguration     *CdnConfigurationModel       `tfsdk:"cdn_configuration"`
	ConfigurationAliases map[string]map[string]string `tfsdk:"configuration_aliases"`
	DashConfiguration    *DashConfigurationModel      `tfsdk:"dash_configuration"`
	// @ADR
	// Context: The Provider Framework does not allow computed blocks
	// Decision: We decided to flatten the Log Configuration and the HLS Configuration blocks into the resource.
	// Consequences: The schema of the object differs from that of the SDK.
	HlsConfigurationManifestEndpointPrefix types.String                   `tfsdk:"hls_configuration_manifest_endpoint_prefix"`
	LogConfigurationPercentEnabled         types.Int64                    `tfsdk:"log_configuration_percent_enabled"`
	LivePreRollConfiguration               *LivePreRollConfigurationModel `tfsdk:"live_pre_roll_configuration"`
	ManifestProcessingRules                *ManifestProcessingRulesModel  `tfsdk:"manifest_processing_rules"`
	Name                                   *string                        `tfsdk:"name"`
	PersonalizationThresholdSeconds        *int32                         `tfsdk:"personalization_threshold_seconds"`
	PlaybackConfigurationArn               types.String                   `tfsdk:"playback_configuration_arn"`
	PlaybackEndpointPrefix                 types.String                   `tfsdk:"playback_endpoint_prefix"`
	SessionInitializationEndpointPrefix    types.String                   `tfsdk:"session_initialization_endpoint_prefix"`
	SlateAdUrl                             *string                        `tfsdk:"slate_ad_url"`
	Tags                                   map[string]string              `tfsdk:"tags"`
	TranscodeProfileName                   *string                        `tfsdk:"transcode_profile_name"`
	VideoContentSourceUrl                  *string                        `tfsdk:"video_content_source_url"`
}

type AvailSuppressionModel struct {
	FillPolicy *string `tfsdk:"fill_policy"`
	Mode       *string `tfsdk:"mode"`
	Value      *string `tfsdk:"value"`
}

type BumperModel struct {
	EndUrl   *string `tfsdk:"end_url"`
	StartUrl *string `tfsdk:"start_url"`
}

type CdnConfigurationModel struct {
	AdSegmentUrlPrefix      *string `tfsdk:"ad_segment_url_prefix"`
	ContentSegmentUrlPrefix *string `tfsdk:"content_segment_url_prefix"`
}

type DashConfigurationModel struct {
	ManifestEndpointPrefix types.String `tfsdk:"manifest_endpoint_prefix"`
	MpdLocation            *string      `tfsdk:"mpd_location"`
	OriginManifestType     *string      `tfsdk:"origin_manifest_type"`
}

type LivePreRollConfigurationModel struct {
	AdDecisionServerUrl *string `tfsdk:"ad_decision_server_url"`
	MaxDurationSeconds  *int32  `tfsdk:"max_duration_seconds"`
}

type ManifestProcessingRulesModel struct {
	AdMarkerPassthrough *AdMarkerPassthroughModel `tfsdk:"ad_marker_passthrough"`
}
type AdMarkerPassthroughModel struct {
	Enabled bool `tfsdk:"enabled"`
}
