# Data Source: awsmt_source_location

This data source provides information about a MediaTailor Source Location.

## Example Usage

```terraform
data "awsmt_source_location" "example" {
  source_location_name = "example"
}
```

## Arguments Reference

The following arguments are supported:

* `source_location_name` - (Required) The name of the source location.

## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `arn` - The ARN of the channel.
* `default_segment_delivery_configuration_url` - The hostname of the server that will be used to serve segments.
* `creation_time` - The timestamp of when the channel was created.
* `http_configuration_url` - The base URL for the source location host server.
* `last_modified_time` - The timestamp of when the channel was last modified.
* `tags` - Key-value mapping of resource tags.


### `access_configuration`
Access configuration parameters. Configures the type of authentication used to access content from your source location.

* `access_type` - The type of authentication used to access content from HttpConfiguration::BaseUrl on your source location. Accepted values: "S3_SIGV4" and "SECRETS_MANAGER_ACCESS_TOKEN". [Read More](https://docs.aws.amazon.com/sdk-for-go/api/service/mediatailor/#AccessConfiguration).
* `smatc_header_name` - Part of Secrets Manager Access Token Configuration. The name of the HTTP header used to supply the access token in requests to the source location.
* `smatc_secret_arn` - Part of Secrets Manager Access Token Configuration. The Amazon Resource Name (ARN) of the AWS Secrets Manager secret that contains the access token.
* `smatc_secret_string_key` - Part of Secrets Manager Access Token Configuration. The AWS Secrets Manager SecretString key associated with the access token.

### `segment_delivery_configurations`
(List) A list of the segment delivery configurations associated with this resource.

* `base_url` - (Optional) The base URL of the host or path of the segment delivery server that you're using to serve segments.
* `name` - (Optional) A unique identifier used to distinguish between multiple segment delivery configurations in a source location.