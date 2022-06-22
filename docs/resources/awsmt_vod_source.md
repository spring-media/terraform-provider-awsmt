# Resource: awsmt_vod_source

Use this resource to manage a MediaTailor VOD Source.


## Example Usage

```terraform
resource "awsmt_vod_source" "example" {
  http_package_configurations {
    path = "/"
    source_group = "default"
    type = "HLS"
  }
  source_location_name = "existing_source_location"
  name = "vod_source_example"
}
```

## Arguments Reference
The following arguments are supported:

* `http_package_configurations` - (Required) A list of HTTP package configuration parameters for this VOD source.
  * `path` - (Required) The relative path to the URL for this VOD source. This is combined with the http_configuration_url specified in the SourceLocation to form a valid URL.
  * `source_group` - (Required) The name of the source group. This has to match one of the source groups specified in the channel.
  * `type` - (Required) the streaming protocol for this package configuration. Can be Either 'HLS' or 'DASH'.
* `source_location_name` - (Required) The name of the Source Location to which the VOD source refers.
* `tags` - (Optional) Key-value mapping of resource tags.
* `name` - (Required) The name of the VOD Source.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the channel.
* `creation_time` - The timestamp of when the channel was created.
* `last_modified_time` - The timestamp of when the channel was last modified.

## Import

VOD Sources can be imported using their ARN as identifier. For example:

```sh
  $ terraform import awsmt_vod_source.example arn:aws:mediatailor:us-east-1:000000000000:vodSource/sourceLocationName/VodSourceName
```