# Resource: awsmt_source_location

Use this resource to manage a MediaTailor Source Location.

~> **WARNING:** Deleting a Source Location also deletes all the Vod Sources and Live Sources connected to it.

## Example Usage

```terraform
resource "awsmt_source_location" "example" {
  access_configuration {
    access_type = "SECRETS_MANAGER_ACCESS_TOKEN"
    smatc_header_name =       "auth"
    smatc_secret_arn =        "arn:aws:secretsmanager:us-east-1:000000000000:secret/example"
    smatc_secret_string_key = "example"
  }
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
- `access_configuration` - (Optional) The access configuration for the source location.
  - `access_type` - (Required) The type of authentication used to access content from HttpConfiguration::BaseUrl on your source location. Valid values are `SECRETS_MANAGER_ACCESS_TOKEN` and `S3_SIGV$`.
  - `smatc_header_name` - (Optional) Part of Secrets Manager Access Token Configuration. The name of the HTTP header used to supply the access token in requests to the source location.
  - `smatc_secret_arn` - (Optional) Part of Secrets Manager Access Token Configuration. The Amazon Resource Name (ARN) of the AWS Secrets Manager secret that contains the access token.
  - `smatc_secret_string_key` - (Optional) Part of Secrets Manager Access Token Configuration. The AWS Secrets Manager SecretString key associated with the access token.
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
