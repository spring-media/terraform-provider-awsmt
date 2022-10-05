# Resource: awsmt_source_location

Use this resource to manage a MediaTailor Source Location.

~> **NOTE:** The source location data source currently does not support the use of access configuration using Amazon Secrets Manager Access Token.

## Example Usage

```terraform
resource "awsmt_source_location" "example" {
  default_segment_delivery_configuration_url = "https://example.com"
  http_configuration_url                     = "https://example.com"
  segment_delivery_configurations {
    base_url = "https://example.com",
    name =     "example"
  }
  name = "example"
  tags = {
    "key": "value"
  }
}
```

## Arguments Reference

The following arguments are supported:

- `name` - (Required) The name of the source location.
- `default_segment_delivery_configuration_url` - (Optional) The hostname of the server that will be used to serve segments.
- `http_configuration_url` - (Optional) The base URL for the source location host server.
- `segment_delivery_configurations` – (Optional List) A list of the segment delivery configurations associated with this resource.
  - `base_url` - (Optional) The base URL of the host or path of the segment delivery server that you're using to serve segments.
  - `name` - (Optional) A unique identifier used to distinguish between multiple segment delivery configurations in a source location.
- `tags` - (Optional) Key-value mapping of resource tags.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `arn` - The ARN of the channel.
- `creation_time` - The timestamp of when the channel was created.
- `last_modified_time` - The timestamp of when the channel was last modified.

## Import

Source Locations can be imported using their ARN as identifier. For example:

```
  $ terraform import awsmt_source_location.example arn:aws:mediatailor:us-east-1:000000000000:source-location/example
```
