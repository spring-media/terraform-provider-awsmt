# awsmt_playback_configuration (Resource)

Use this resource to create a playback configuration.

## Usage

```
resource "awsmt_playback_configuration" "conf" {
  ad_decision_server_url = "https://exampleurl.com/"
  cdn_configuration {
    ad_segment_url_prefix = "test"
    content_segment_url_prefix = "test"
  }
  dash_configuration {
    mpd_location = "test"
    origin_manifest_type = "MULTI_PERIOD"
  }
  name = "test-playback-configuration-awsmt"
  slate_ad_url = "https://exampleurl.com/"
  tags = {}
  transcode_profile_name = "test"
  video_content_source_url = "https://exampleurl.com/"
}
```

## Schema
All the descriptions for the fields are from the [official AWS documentation](https://docs.aws.amazon.com/sdk-for-go/api/service/mediatailor/#MediaTailor.PutPlaybackConfiguration) or this [SourceGraph Page](https://sourcegraph.com/github.com/aws/aws-sdk-go/-/docs/service/mediatailor#PutPlaybackConfigurationInput) .

* `ad_decision_server_url` - (optional, type string). <br/>
  The URL for the ad decision server (ADS). This includes the specification
  of static parameters and placeholders for dynamic parameters. AWS Elemental
  MediaTailor substitutes player-specific and session-specific parameters as
  needed when calling the ADS. Alternately, for testing you can provide a static
  VAST URL. The maximum length is 25,000 characters.
* `cdn_configuration` - (optional, type list of object) (see [below for nested schema](#cdn_conf))<br/>
  The configuration for using a content delivery network (CDN), like Amazon
  CloudFront, for content and ad segment management.
* `dash_configuration` - (optional, type list of object) (see [below for nested schema](#dash_conf))<br/>


<a id="cdn_conf"></a>
### Nested Schema for `cdn_configuration`

* `ad_segment_url_prefix` - (optional, type string)<br/>
  A non-default content delivery network (CDN) to serve ad segments. By default,
  AWS Elemental MediaTailor uses Amazon CloudFront with default cache settings
  as its CDN for ad segments. To set up an alternate CDN, create a rule in
  your CDN for the origin ads.mediatailor.`region`.amazonaws.com. Then specify
  the rule's name in this AdSegmentUrlPrefix. When AWS Elemental MediaTailor
  serves a manifest, it reports your CDN as the source for ad segments.
* `content_segment_url_prefix` - (optional, type string) <br/>
  A content delivery network (CDN) to cache content segments, so that content
  requests don’t always have to go to the origin server. First, create a
  rule in your CDN for the content segment origin server. Then specify the
  rule's name in this ContentSegmentUrlPrefix. When AWS Elemental MediaTailor
  serves a manifest, it reports your CDN as the source for content segments.
* `name` - (required, type string) </br>
  The identifier for the playback configuration.
* `slate_ad_url` - (optional, type string)<br/>
  The URL for a high-quality video asset to transcode and use to fill in time
  that's not used by ads. AWS Elemental MediaTailor shows the slate to fill
  in gaps in media content. Configuring the slate is optional for non-VPAID
  playback configurations. For VPAID, the slate is required because MediaTailor
  provides it in the slots designated for dynamic ad content. The slate must
  be a high-quality asset that contains both audio and video.
* - `tags` - (optional, type map of string)<br/>
  The tags assigned to the playback configuration.
*`transcode_profile_name` - (optional, type string)<br/>
  The name that is used to associate this playback configuration with a custom
  transcode profile. This overrides the dynamic transcoding defaults of MediaTailor.
  Use this only if you have already set up custom profiles with the help of
  AWS Support.
*`video_content_source_url` - (optional, type string)<br/>
  The URL prefix for the parent manifest for the stream, minus the asset ID.
  The maximum length is 512 characters.


<a id="dash_conf"></a>
### Nested Schema for `dash_configuration`

* `mpd_location` - (optional, type string) <br/>
  The setting that controls whether MediaTailor includes the Location tag in
  DASH manifests. MediaTailor populates the Location tag with the URL for manifest
  update requests, to be used by players that don't support sticky redirects.
  Disable this if you have CDN routing rules set up for accessing MediaTailor
  manifests, and you are either using client-side reporting or your players
  support sticky HTTP redirects. Valid values are DISABLED and EMT_DEFAULT.
  The EMT_DEFAULT setting enables the inclusion of the tag and is the default
  value.
* `origin_manifest_type` - (optional, type string, enum `SINGLE_PERIOD` | `MULTI_PERIOD`) <br/>
  The setting that controls whether MediaTailor handles manifests from the
  origin server as multi-period manifests or single-period manifests. If your
  origin server produces single-period manifests, set this to SINGLE_PERIOD.
  The default setting is MULTI_PERIOD. For multi-period manifests, omit this
  setting or set it to MULTI_PERIOD.