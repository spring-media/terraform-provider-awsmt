# Resource: awsmt_playback_configuration

Use this resource to manage a MediaTailor playback configuration.

## Example Usage

You can specify the arguments inside a resource block like this:

```terraform
resource "awsmt_playback_configuration" "conf" {
  ad_decision_server_url = "https://exampleurl.com/"
  ad_conditioning_configuration = {
    streaming_media_file_conditioning = "TRANSCODE"
  }
  ad_decision_server_configuration = {
    http_request = {
      method           = "POST"
      compress_request = "GZIP"
      headers = {
        "Content-Type" = "application/json"
      }
      body = "{\"key\": \"value\"}"
    }
  }
  avail_suppression = {
   mode = "OFF"
  }
  bumper = {}
  cdn_configuration = {
    ad_segment_url_prefix      = "test"
    content_segment_url_prefix = "test"
  }
  dash_configuration = {
    mpd_location         = "EMT_DEFAULT"
    origin_manifest_type = "MULTI_PERIOD"
  }
  function_mapping = {
    "PRE_ADS_REQUEST" = "my-function-id"
  }
  insertion_mode = "STITCHED_ONLY"
  name = "test-playback-configuration-awsmt"
  manifest_processing_rules = {
    ad_marker_passthrough = {
      enabled = "false"
    }
  }
  log_configuration_ads_interaction_log = {
    exclude_event_types      = ["BEACON_FIRED"]
    publish_opt_in_event_types = ["RAW_ADS_RESPONSE"]
  }
  log_configuration_manifest_service_interaction_log = {
    exclude_event_types = ["TRACKING_RESPONSE"]
  }
  slate_ad_url             = "https://exampleurl.com/"
  transcode_profile_name    = "profile_configured_in_your_account"
  video_content_source_url = "https://exampleurl.com/"
}
```

## Arguments Reference

All the descriptions for the fields are from the [official AWS documentation](https://docs.aws.amazon.com/sdk-for-go/api/service/mediatailor/#MediaTailor.PutPlaybackConfiguration).

The following arguments are supported:

- `ad_conditioning_configuration` - (Optional) The setting that indicates what conditioning MediaTailor will perform on ads that the ADS returns.
  - `streaming_media_file_conditioning` - (Required) Indicates what transcoding action MediaTailor takes when it first receives ads from the ADS. Valid values: `TRANSCODE`, `NONE`.
- `ad_decision_server_configuration` - (Optional) The configuration for customizing HTTP requests to the ad decision server (ADS).
  - `http_request` - (Optional) The HTTP request configuration parameters for the ADS.
    - `body` - (Optional) The request body content to send with HTTP requests to the ADS. Only eligible for POST requests.
    - `compress_request` - (Optional) The compression method to apply to requests. Valid values: `NONE`, `GZIP`. Only eligible for POST requests.
    - `headers` - (Optional) Custom HTTP headers to include in requests to the ADS as key-value pairs. Only eligible for POST requests.
    - `method` - (Optional) The HTTP method to use when making requests to the ADS. Valid values: `GET`, `POST`.
- `ad_decision_server_url` - (Required) The URL for the ad decision server (ADS).
- `avail_suppression` - (Optional) The configuration for avail suppression, also known as ad suppression.
  - `fill_policy` - Defines the policy to apply to the avail suppression mode. Can be either full (BEHIND_LIVE_EDGE mode) or partial (AFTER_LIVE_EDGE).
  - `mode` - The ad suppression mode. Can either be "OFF", "BEHIND_LIVE_EDGE" "AFTER_LIVE_EDGE".
  - `value` - Time value in HH:MM:SS format after which MediaTailor will not fill any ad breaks.
- `bumper` - (Optional) The configuration for bumpers.
  - `end_url` - The URL for the end bumper asset.
  - `start_url` - The URL for the start bumper asset.
- `cdn_configuration` - (Optional) The configuration for using a content delivery network (CDN) for content and ad segment management.
  - `ad_segment_url_prefix` - A non-default CDN to serve ads segments.
  - `content_segment_url_prefix` - A CDN to cache content segments.
- `configuration_aliases` - (Optional) The player parameters and aliases used as dynamic variables during session initialization.
- `dash_configuration` - (Optional) The configuration for DASH content.
  - `mpd_location` - Controls whether MediaTailor includes the Location tag in Dash manifest files. Can either be "DISABLED" or "EMT_DEFAULT.
  - `origin_manifest_type` - Controls whether MediaTailor handles manifest files as single-period or multi-period manifest files. Can either be "SINGLE_PERIOD" or "MULTI_PERIOD".
- `function_mapping` - (Optional) A map of lifecycle hook event names to function identifiers. Valid keys are `PRE_SESSION_INITIALIZATION` and `PRE_ADS_REQUEST`. Values are function IDs.
- `insertion_mode` - (Optional) Controls whether players can use stitched or guided ad insertion. Valid values: `STITCHED_ONLY`, `PLAYER_SELECT`.
- `live_pre_roll_configuration` - (Optional) The configuration for pre-roll ad insertion.
  - `ad_decision_server_url` - The URL for the ad decision server (ADS) for pre-roll ads.
  - `max_duration_seconds` - The maximum allowed duration for the pre-roll ad avail.
- `log_configuration_percent_enabled` - (Optional) The percentage of session logs that MediaTailor sends to your Cloudwatch Logs account.
- `log_configuration_enabled_logging_strategies` - (Optional) The method used for collecting logs from AWS Elemental MediaTailor. Allowed values are "LEGACY_CLOUDWATCH" or "VENDED_LOGS", or both.
- `log_configuration_ads_interaction_log` - (Optional) Settings for customizing what events are included in logs for interactions with the ADS.
  - `exclude_event_types` - (Optional) List of event types that MediaTailor won't emit in the logs. See [ADS log event types](https://docs.aws.amazon.com/mediatailor/latest/ug/ads-log-format.html) for valid values.
  - `publish_opt_in_event_types` - (Optional) List of event types that MediaTailor will emit in the logs (not emitted by default). Valid values: `RAW_ADS_RESPONSE`, `RAW_ADS_REQUEST`, `PRE_ADS_REQUEST_HOOK_SUMMARY`, `PRE_ADS_REQUEST_FUNCTION_COMPLETED`.
- `log_configuration_manifest_service_interaction_log` - (Optional) Settings for customizing what events are included in logs for interactions with the origin server.
  - `exclude_event_types` - (Optional) List of event types that MediaTailor won't emit in the logs. See [manifest service log event types](https://docs.aws.amazon.com/mediatailor/latest/ug/log-types.html) for valid values.
  - `publish_opt_in_event_types` - (Optional) List of event types that MediaTailor will emit in the logs (not emitted by default). Valid values: `PRE_SESSION_INIT_HOOK_SUMMARY`, `PRE_SESSION_INIT_FUNCTION_COMPLETED`.
- `manifest_processing_rules` – (Optional) The configuration for manifest processing rules
  - `ad_marker_passthrough` – For HLS, when set to true, MediaTailor passes through EXT-X-CUE-IN, EXT-X-CUE-OUT, and EXT-X-SPLICEPOINT-SCTE35 ad markers from the origin manifest to the MediaTailor personalized manifest.
    - `enabled` - Enables ad marker passthrough for your configuration.
- `personalization_threshold_seconds` - (Optional) Defines the maximum duration of underfilled ad time (in seconds) allowed in an ad break.
- `slate_ad_url` - (Optional) The URL for a high-quality video asset to transcode and use to fill in time that's not used by ads.
- `tags` - (Optional) Key-value mapping of resource tags.
- `transcode_profile_name` - (Optional) The name that is used to associate this playback configuration with a custom transcode profile.
- `video_content_source_url` - (Required) The URL prefix for the parent manifest for the stream, minus the asset ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `dash_configuration` - The configuration for DASH content.
  - `manifest_endpoint_prefix` - URL generated by MediaTailor to initiate a playback session.
- `hls_configuration_manifest_endpoint_prefix` - URL generated by MediaTailor to initiate a playback session on devices that support Apple HLS.
- `playback_configuration_arn` - The Amazon Resource Name (ARN) for the playback configuration.
- `playback_endpoint_prefix` - The URL that the player accesses to get a manifest from AWS Elemental MediaTailor.
- `session_initialization_endpoint_prefix` - The URL that the player uses to initialize a session that uses client-side reporting.

## Import

`awsmt_playback_configuration` resources can be imported using their name as identifier. For example:

```sh
  $ terraform import awsmt_playback_configuration.example broadcast-live-stream
```
