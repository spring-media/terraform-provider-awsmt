# Data Source: awsmt_vod_source

Use this resource to get information on a MediaTailor VOD Source.


## Example Usage

```terraform
data "awsmt_vod_source" "data_test" {
  source_location_name = awsmt_source_location.example.source_location_name
  name = awsmt_vod_source.test.name
}
```

## Arguments Reference
The following arguments are supported:

* `source_location_name` - (Required) The name of the Source Location to which the VOD source refers.
* `name` - (Required) The name of the VOD Source.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the channel.
* `creation_time` - The timestamp of when the channel was created.
* `http_package_configurations` - A list of HTTP package configuration parameters for this VOD source.
  * `path` - The relative path to the URL for this VOD source. This is combined with the http_configuration_url specified in the SourceLocation to form a valid URL.
  * `source_group` - The name of the source group. This has to match one of the source groups specified in the channel.
  * `type` - the streaming protocol for this package configuration. Can be Either 'HLS' or 'DASH'.
* `last_modified_time` - The timestamp of when the channel was last modified.
* `tags` - Key-value mapping of resource tags.
