# Data Source: awsmt_source_location

This data source provides information about a MediaTailor Source Location.

~> **NOTE:** The source location data source currently does not support the use of access configuration using Amazon Secrets Manager Access Token.

## Example Usage

```terraform
data "awsmt_source_location" "example" {
  name = "example"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the source location.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the channel.
* `default_segment_delivery_configuration_url` - The hostname of the server that will be used to serve segments.
* `creation_time` - The timestamp of when the channel was created.
* `http_configuration_url` - The base URL for the source location host server.
* `last_modified_time` - The timestamp of when the channel was last modified.
* `tags` - Key-value mapping of resource tags.

### `segment_delivery_configurations`
(List) A list of the segment delivery configurations associated with this resource.

* `base_url` - (Optional) The base URL of the host or path of the segment delivery server that you're using to serve segments.
* `name` - (Optional) A unique identifier used to distinguish between multiple segment delivery configurations in a source location.