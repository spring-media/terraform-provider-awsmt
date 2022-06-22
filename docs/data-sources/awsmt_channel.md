# Resource: awsmt_channel

Provides information about an Elemental MediaTailor Channel.

## Example Usage

```terraform
data "aws_media_tailor_channel" "example" {
  name = "example-channel"
}
```

## Argument Reference
The following arguments are supported:

* `name` - (Required) The name of the channel.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the channel.
* `channel_state` - Returns whether the channel is running or not.
* `creation_time` - The timestamp of when the channel was created.
* `filler_slate` – The slate used to fill gaps between programs in the schedule. You must configure filler slate if your channel uses the LINEAR PlaybackMode.
  * `source_location_name` - The name of the source location where the slate VOD source is stored.
  * `vod_source_name` - The slate VOD source name. The VOD source must already exist in a source location before it can be used for slate.
* `last_modified_time` - The timestamp of when the channel was last modified.
* `outputs` – The channel's output properties.
  * `dash_manifest_windows_seconds` - The total duration (in seconds) of each dash manifest.
  * `dash_min_buffer_time_seconds` - Minimum amount of content (measured in seconds) that a player must keep available in the buffer.
  * `dash_min_update_period_seconds` - Minimum amount of time (in seconds) that the player should wait before requesting updates to the manifest.
  * `dash_suggested_presentation_delay_seconds` - Amount of time (in seconds) that the player should be from the live point at the end of the manifest.
  * `hls_manifest_windows_seconds` - The total duration (in seconds) of each hls manifest.
  * `manifest_name` - The name of the manifest for the channel. The name appears in the PlaybackUrl.
  * `playback_url` - The URL used for playback by content players.
* `playback_mode` - The type of playback mode for this channel. Can be either LINEAR or LOOP.
* `policy` - The IAM policy for the channel.
* `source_group` - A string used to match which HttpPackageConfiguration is used for each VodSource.
* `tags` - Key-value mapping of resource tags. If configured with a provider [`default_tags` configuration block](/docs/providers/aws/index.html#default_tags-configuration-block) present, tags with matching keys will overwrite those defined at the provider-level.
* `tier` - The tier for this channel. STANDARD tier channels can contain live programs.
