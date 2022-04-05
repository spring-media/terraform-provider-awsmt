# Resource: awsmt_playback_configuration

Use this resource to manage a MediaTailor playback configuration.

## Example Usage
You can specify the arguments inside a resource block like this:
```
resource "awsmt_playback_configuration" "conf" {
  ad_decision_server_url = "https://exampleurl.com/"
  avail_suppression {
   mode = "OFF"
  }
  bumper {}
  cdn_configuration {
    ad_segment_url_prefix = "test"
    content_segment_url_prefix = "test"
  }
  dash_configuration {
    mpd_location = "EMT_DEFAULT"
    origin_manifest_type = "MULTI_PERIOD"
  }
  name = "test-playback-configuration-awsmt"
  live_pre_roll_configuration {}
  manifest_processing_rules {
    ad_marker_passthrough{
      enabled = "false"
    }
  }
  slate_ad_url = "https://exampleurl.com/"
  tags = {}
  transcode_profile_name = "profile_configured_in_your_account"
  video_content_source_url = "https://exampleurl.com/"
}
```

## Arguments Reference
All the descriptions for the fields are from the [official AWS documentation](https://docs.aws.amazon.com/sdk-for-go/api/service/mediatailor/#MediaTailor.PutPlaybackConfiguration).

The following arguments are required:

* `name` - (required, string) </br>
  The identifier for the playback configuration.
* `ad_decision_server_url` - (required, type string). <br/>
  The URL for the ad decision server (ADS). This includes the specification
  of static parameters and placeholders for dynamic parameters. AWS Elemental
  MediaTailor substitutes player-specific and session-specific parameters as
  needed when calling the ADS. Alternately, for testing you can provide a static
  VAST URL. The maximum length is 25,000 characters.
* `dash_configuration` - (required, structure) (see [below for nested schema](#dash_conf))<br/>
  The configuration for DASH content.
* `video_content_source_url` - (required, type string)<br/>
  The URL prefix for the parent manifest for the stream, minus the asset ID.
  The maximum length is 512 characters.


The following arguments are optional


* `avail_suppression` - (optional, structure) (see [below for nested schema](#avail_suppression))<br/>
  The configuration for avail suppression, also known as ad suppression.
* `bumper` - (required, structure) (see [below for nested schema](#bumper))<br/>
  The configuration for bumpers. Bumpers are short audio or video clips that play at the start or before the end of an ad break.
* `cdn_configuration` - (optional, structure) (see [below for nested schema](#cdn_conf))<br/>
  The configuration for using a content delivery network (CDN), like Amazon
  CloudFront, for content and ad segment management.
* `configuration_aliases` - (map)<br/>
  The player parameters and aliases used as dynamic variables during session initialization.
* `live_pre_roll_configuration` - (optional, structure) (see [below for nested schema](#live_pre_roll_configuration))<br/>
  The configuration for pre-roll ad insertion.
* `manifest_processing_rules` - (optional, structure) (see [below for nested schema](#manifest_processing_rules))<br/>
  The configuration for manifest processing rules. Manifest processing rules
  enable customization of the personalized manifests created by MediaTailor.
* `personalization_threshold_seconds` - (optional, integer)<br/>
  Defines the maximum duration of underfilled ad time (in seconds) allowed
  in an ad break. If the duration of underfilled ad time exceeds the personalization
  threshold, then the personalization of the ad break is abandoned and the
  underlying content is shown. This feature applies to ad replacement in live 
  and VOD streams, rather than ad insertion, because it relies on an underlying
  content stream.
* `slate_ad_url` - (optional, string)<br/>
  The URL for a high-quality video asset to transcode and use to fill in time
  that's not used by ads. AWS Elemental MediaTailor shows the slate to fill
  in gaps in media content. Configuring the slate is optional for non-VPAID
  playback configurations. For VPAID, the slate is required because MediaTailor
  provides it in the slots designated for dynamic ad content. The slate must
  be a high-quality asset that contains both audio and video.
* `tags` - (optional, map)<br/>
  The tags assigned to the playback configuration.
* `transcode_profile_name` - (optional, string)<br/>
  The name that is used to associate this playback configuration with a custom
  transcode profile. This overrides the dynamic transcoding defaults of MediaTailor.
  Use this only if you have already set up custom profiles with the help of
  AWS Support.

## Attributes Reference

* `hls_configuration` - (structure) (see [below for nested schema](#attribute_hls_conf))<br/>
  The configuration for HLS content.
* `log_configuration` - (structure) (see [below for nested schema](#log_configuration))<br/>
  The Amazon CloudWatch log settings for a playback configuration.
* `playback_configuration_arn` - (string)<br/>
  The Amazon Resource Name (ARN) for the playback configuration.
* `playback_endpoint_prefix` - (string)<br/>
  The URL that the player accesses to get a manifest from AWS Elemental MediaTailor.
  This session will use server-side reporting.
* `session_initialization_endpoint_prefix` - (string)<br/>
  The URL that the player uses to initialize a session that uses client-side
  reporting.

## Nested Schemas

<a id="avail_suppression"></a>
### `avail_suppression`

* `mode` - (required, string)<br/>
  Sets the ad suppression mode. By default, ad suppression is off and all ad
  breaks are filled with ads or slate. When Mode is set to BEHIND_LIVE_EDGE,
  ad suppression is active and MediaTailor won't fill ad breaks on or behind
  the ad suppression Value time in the manifest lookback window.
* `value` - (optional, string)<br/>
  A live edge offset time in HH:MM:SS. MediaTailor won't fill ad breaks on
  or behind this time in the manifest lookback window. If Value is set to 00:00:00,
  it is in sync with the live edge, and MediaTailor won't fill any ad breaks
  on or behind the live edge. If you set a Value time, MediaTailor won't fill
  any ad breaks on or behind this time in the manifest lookback window. For
  example, if you set 00:45:00, then MediaTailor will fill ad breaks that occur
  within 45 minutes behind the live edge, but won't fill ad breaks on or behind
  45 minutes behind the live edge.

<a id="bumper"></a>
### `bumper`

* `end_url` - (optional, string)<br/>
  The URL for the end bumper asset.
* `start_url` - (optional, string)<br/>
  The URL for the start bumper asset.

<a id="cdn_conf"></a>
### `cdn_configuration`

* `ad_segment_url_prefix` - (optional, string)<br/>
  A non-default content delivery network (CDN) to serve ad segments. By default,
  AWS Elemental MediaTailor uses Amazon CloudFront with default cache settings
  as its CDN for ad segments. To set up an alternate CDN, create a rule in
  your CDN for the origin ads.mediatailor.`region`.amazonaws.com. Then specify
  the rule's name in this AdSegmentUrlPrefix. When AWS Elemental MediaTailor
  serves a manifest, it reports your CDN as the source for ad segments.
* `content_segment_url_prefix` - (optional, string) <br/>
  A content delivery network (CDN) to cache content segments, so that content
  requests don’t always have to go to the origin server. First, create a
  rule in your CDN for the content segment origin server. Then specify the
  rule's name in this ContentSegmentUrlPrefix. When AWS Elemental MediaTailor
  serves a manifest, it reports your CDN as the source for content segments.

<a id="dash_conf"></a>
### `dash_configuration`

* `mpd_location` - (optional, string) <br/>
  The setting that controls whether MediaTailor includes the Location tag in
  DASH manifests. MediaTailor populates the Location tag with the URL for manifest
  update requests, to be used by players that don't support sticky redirects.
  Disable this if you have CDN routing rules set up for accessing MediaTailor
  manifests, and you are either using client-side reporting or your players
  support sticky HTTP redirects. Valid values are DISABLED and EMT_DEFAULT.
  The EMT_DEFAULT setting enables the inclusion of the tag and is the default
  value.
* `origin_manifest_type` - (optional, string) <br/>
  The setting that controls whether MediaTailor handles manifests from the
  origin server as multi-period manifests or single-period manifests. If your
  origin server produces single-period manifests, set this to SINGLE_PERIOD.
  The default setting is MULTI_PERIOD. For multi-period manifests, omit this
  setting or set it to MULTI_PERIOD.

<a id="attribute_hls_conf"></a>
### `hls_configuration`

* `manifest_endpoint_prefix` - (attribute, type string)<br/>
  The URL that is used to initiate a playback session for devices that support
  Apple HLS. The session uses server-side reporting.

<a id="live_pre_roll_configuration"></a>
### `live_preroll_configuration`

* `ad_decision_server_url` - (optional, string)<br/>
  The URL for the ad decision server (ADS) for pre-roll ads. This includes
  the specification of static parameters and placeholders for dynamic parameters.
  AWS Elemental MediaTailor substitutes player-specific and session-specific
  parameters as needed when calling the ADS. Alternately, for testing, you
  can provide a static VAST URL. The maximum length is 25,000 characters.
* `max_duration_server` - (optional, string)<br/>
  The maximum allowed duration for the pre-roll ad avail. AWS Elemental MediaTailor
  won't play pre-roll ads to exceed this duration, regardless of the total
  duration of ads that the ADS returns.

<a id="log_configuration"></a>
### `log_configuration`

* `percent_enabled` - (attribute, string)<br/>
  The percentage of session logs that MediaTailor sends to your Cloudwatch
  Logs account. For example, if your playback configuration has 1000 sessions
  and percentEnabled is set to 60, MediaTailor sends logs for 600 of the sessions
  to CloudWatch Logs. MediaTailor decides at random which of the playback configuration
  sessions to send logs for. If you want to view logs for a specific session,
  you can use the debug log mode.

<a id="manifest_processing_rules"></a>
### `manifest_processing_rules`

* `ad_marker_passthrough` - (required, structure) (see [below for nested schema](#ad_marker_passthrough))<br/>
  For HLS, when set to true, MediaTailor passes through EXT-X-CUE-IN, EXT-X-CUE-OUT,
  and EXT-X-SPLICEPOINT-SCTE35 ad markers from the origin manifest to the MediaTailor
  personalized manifest.

<a id="ad_marker_passthrough"></a>
### `ad_marker_passthrough`

* `enabled` - (optional, bool)<br/>
  Enables ad marker passthrough for your configuration.

## Import

`awsmt_playback_configuration` resources can be imported using their name as identifier. For example:
```
  $ terraform import awsmt_playback_configuration.example broadcast-live-stream
```