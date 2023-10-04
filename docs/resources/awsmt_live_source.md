# Resource: awsmt_live_source

Use this resource to manage a MediaTailor Live Source.

## Example Usage

```terraform
resource "awsmt_live_source" "example" {
  http_package_configurations {
    path         = "/"
    source_group = "default"
    type         = "HLS"
  }
  source_location_name = "existing_source_location"
  name                 = "live_source_example"
}
```

## Arguments Reference

The following arguments are supported:

- `http_package_configurations` - (Required) A list of HTTP package configuration parameters for this Live source.
  - `path` - (Required) The relative path to the URL for this Live Source. This is combined with the http_configuration_url specified in the SourceLocation to form a valid URL.
  - `source_group` - (Required) The name of the source group. This has to match one of the source groups specified in the channel.
  - `type` - (Required) the streaming protocol for this package configuration. Can be Either 'HLS' or 'DASH'.
- `name` - (Required) The name of the Live Source.
- `source_location_name` - (Required) The name of the Source Location to which the Live Source refers.
- `tags` - (Optional) Key-value mapping of resource tags.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `arn` - The ARN of the channel.
- `creation_time` - The timestamp of when the channel was created.
- `last_modified_time` - The timestamp of when the channel was last modified.

## Import

Live Sources can be imported using their Name and SourceLocationName in one string as identifier. For example:

```shell
  $ terraform import awsmt_live_source.example name:example-live-source,example-source-location
```
