# Data Source: awsmt_live_source

Use this resource to get information on a MediaTailor Live Source.


## Example Usage

```terraform
data "awsmt_live_source" "data_test" {
  source_location_name = awsmt_source_location.example.source_location_name
  live_source_name = awsmt_live_source.test.live_source_name
}
```

## Arguments Reference
The following arguments are supported:

* `source_location_name` - (Required) The name of the Source Location to which the Live Source refers.
* `live_source_name` - (Required) The name of the Live Source.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the channel.
* `creation_time` - The timestamp of when the channel was created.
* `last_modified_time` - The timestamp of when the channel was last modified.
* `tags` - Key-value mapping of resource tags.

### `http_package_configurations` - (Required)
A list of HTTP package configuration parameters for this Live Source.

* `path` - (Required) The relative path to the URL for this Live Source. This is combined with the http_configuration_url specified in the SourceLocation to form a valid URL.
* `source_group` - (Required) The name of the source group. This has to match one of the source groups specified in the channel.
* `type` - (Required) the streaming protocol for this package configuration. Can be Either 'HLS' or 'DASH'.
