# Data Source: awsmt_source_location

This data source provides information about a MediaTailor Source Location.

~> **NOTE:** The source location data source currently does not support the use of access configuration using Amazon Secrets Manager Access Token Configuration.

## Example Usage

```terraform
data "awsmt_source_location" "example" {
  name = "example"
}
```

## Arguments Reference

The following arguments are supported:

- `name` - (Required) The name of the source location.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `access_configuration` - (Optional) The access configuration for the source location.
  - `access_type` - (Required) The type of authentication used to access content from HttpConfiguration::BaseUrl on your source location. Valid values are `SECRETS_MANAGER_ACCESS_TOKEN` and `S3_SIGV$`.
  - `smatc` - (Optional) Part of Secrets Manager Access Token Configuration. The Amazon Resource Name (ARN) of the AWS Secrets Manager secret that contains the access token.
    - `header_name` - (Optional) Part of Secrets Manager Access Token Configuration. The name of the HTTP header used to supply the access token in requests to the source location.
    - `secret_arn` - (Optional) Part of Secrets Manager Access Token Configuration. The Amazon Resource Name (ARN) of the AWS Secrets Manager secret that contains the access token.
    - `secret_string_key` - (Optional) Part of Secrets Manager Access Token Configuration. The AWS Secrets Manager SecretString key associated with the access token.
- `arn` - The ARN of the channel.
- `creation_time` - The timestamp of when the channel was created.
- `default_segment_delivery_configuration` - The default segment delivery configuration settings.
  - `base_url` - The hostname of the server that will be used to serve segments.
- `http_configuration` - The HTTP configuration for the source location.
  - `base_url` - The base URL for the source location host server.
- `last_modified_time` - The timestamp of when the channel was last modified.
- `segment_delivery_configurations` – (List) A list of the segment delivery configurations associated with this resource.
  - `base_url` - The base URL of the host or path of the segment delivery server that you're using to serve segments.
  - `name` - A unique identifier used to distinguish between multiple segment delivery configurations in a source location.
- `tags` - Key-value mapping of resource tags.
