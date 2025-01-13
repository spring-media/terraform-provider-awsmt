package models

import (
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ChannelModel struct {
	ID           types.String `tfsdk:"id"`
	Arn          types.String `tfsdk:"arn"`
	Name         *string      `tfsdk:"name"`
	ChannelState *string      `tfsdk:"channel_state"`
	CreationTime types.String `tfsdk:"creation_time"`
	// @ADR
	// Context: Managing the enablement and disablement of logs requires a configuration structure in the SDK.
	// Decision: As the only log type available for channels is AS_RUN, we simplified the configuration by
	// converting this option into a boolean for the provider.
	// Consequences: The process for enabling and disabling logs differs slightly from the SDK's approach.
	EnableAsRunLogs  types.Bool           `tfsdk:"enable_as_run_logs"`
	FillerSlate      *FillerSlateModel    `tfsdk:"filler_slate"`
	LastModifiedTime types.String         `tfsdk:"last_modified_time"`
	Outputs          []OutputsModel       `tfsdk:"outputs"`
	PlaybackMode     *string              `tfsdk:"playback_mode"`
	Policy           jsontypes.Normalized `tfsdk:"policy"`
	Tags             map[string]string    `tfsdk:"tags"`
	Tier             *string              `tfsdk:"tier"`
}

type FillerSlateModel struct {
	SourceLocationName *string `tfsdk:"source_location_name"`
	VodSourceName      *string `tfsdk:"vod_source_name"`
}

type OutputsModel struct {
	DashPlaylistSettings *DashPlaylistSettingsModel `tfsdk:"dash_playlist_settings"`
	HlsPlaylistSettings  *HlsPlaylistSettingsModel  `tfsdk:"hls_playlist_settings"`
	ManifestName         *string                    `tfsdk:"manifest_name"`
	PlaybackUrl          types.String               `tfsdk:"playback_url"`
	SourceGroup          *string                    `tfsdk:"source_group"`
}

type DashPlaylistSettingsModel struct {
	ManifestWindowSeconds             *int64 `tfsdk:"manifest_window_seconds"`
	MinBufferTimeSeconds              *int64 `tfsdk:"min_buffer_time_seconds"`
	MinUpdatePeriodSeconds            *int64 `tfsdk:"min_update_period_seconds"`
	SuggestedPresentationDelaySeconds *int64 `tfsdk:"suggested_presentation_delay_seconds"`
}
type HlsPlaylistSettingsModel struct {
	AdMarkupType          []*string `tfsdk:"ad_markup_type"`
	ManifestWindowSeconds *int64    `tfsdk:"manifest_window_seconds"`
}
