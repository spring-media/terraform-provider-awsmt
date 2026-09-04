package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type PlaybackConfigurationModel struct {
	ID                              types.String                          `tfsdk:"id"`
	AdConditioningConfiguration     *AdConditioningConfigurationModel     `tfsdk:"ad_conditioning_configuration"`
	AdDecisionServerConfiguration   *AdDecisionServerConfigurationModel   `tfsdk:"ad_decision_server_configuration"`
	AdDecisionServerUrl             *string                               `tfsdk:"ad_decision_server_url"`
	AvailSuppression                *AvailSuppressionModel                `tfsdk:"avail_suppression"`
	Bumper                          *BumperModel                          `tfsdk:"bumper"`
	CdnConfiguration                *CdnConfigurationModel                `tfsdk:"cdn_configuration"`
	ConfigurationAliases            map[string]map[string]string          `tfsdk:"configuration_aliases"`
	DashConfiguration               *DashConfigurationModel               `tfsdk:"dash_configuration"`
	FunctionMapping                 map[string]string                     `tfsdk:"function_mapping"`
	// @ADR
	// Context: The Provider Framework does not allow computed blocks
	// Decision: We decided to flatten the Log Configuration and the HLS Configuration blocks into the resource.
	// Consequences: The schema of the object differs from that of the SDK.
	HlsConfigurationManifestEndpointPrefix                types.String                        `tfsdk:"hls_configuration_manifest_endpoint_prefix"`
	InsertionMode                                         *string                             `tfsdk:"insertion_mode"`
	LogConfigurationPercentEnabled                        types.Int64                         `tfsdk:"log_configuration_percent_enabled"`
	LogConfigurationEnabledLoggingStrategies              types.List                          `tfsdk:"log_configuration_enabled_logging_strategies"`
	LogConfigurationAdsInteractionLog                     *AdsInteractionLogModel             `tfsdk:"log_configuration_ads_interaction_log"`
	LogConfigurationManifestServiceInteractionLog         *ManifestServiceInteractionLogModel `tfsdk:"log_configuration_manifest_service_interaction_log"`
	LivePreRollConfiguration                              *LivePreRollConfigurationModel      `tfsdk:"live_pre_roll_configuration"`
	ManifestProcessingRules                               *ManifestProcessingRulesModel       `tfsdk:"manifest_processing_rules"`
	Name                                                  *string                             `tfsdk:"name"`
	PersonalizationThresholdSeconds                       *int32                              `tfsdk:"personalization_threshold_seconds"`
	PlaybackConfigurationArn                              types.String                        `tfsdk:"playback_configuration_arn"`
	PlaybackEndpointPrefix                                types.String                        `tfsdk:"playback_endpoint_prefix"`
	SessionInitializationEndpointPrefix                   types.String                        `tfsdk:"session_initialization_endpoint_prefix"`
	SlateAdUrl                                            *string                             `tfsdk:"slate_ad_url"`
	Tags                                                  map[string]string                   `tfsdk:"tags"`
	TranscodeProfileName                                  *string                             `tfsdk:"transcode_profile_name"`
	VideoContentSourceUrl                                 *string                             `tfsdk:"video_content_source_url"`
}

type AdConditioningConfigurationModel struct {
	StreamingMediaFileConditioning *string `tfsdk:"streaming_media_file_conditioning"`
}

type AdDecisionServerConfigurationModel struct {
	HttpRequest *HttpRequestModel `tfsdk:"http_request"`
}

type HttpRequestModel struct {
	Body            *string           `tfsdk:"body"`
	CompressRequest *string           `tfsdk:"compress_request"`
	Headers         map[string]string `tfsdk:"headers"`
	Method          *string           `tfsdk:"method"`
}

type AdsInteractionLogModel struct {
	ExcludeEventTypes     []string `tfsdk:"exclude_event_types"`
	PublishOptInEventTypes []string `tfsdk:"publish_opt_in_event_types"`
}

type ManifestServiceInteractionLogModel struct {
	ExcludeEventTypes     []string `tfsdk:"exclude_event_types"`
	PublishOptInEventTypes []string `tfsdk:"publish_opt_in_event_types"`
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
