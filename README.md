# ott-tfprovider-awsmt
A Terraform provider for AWS MediaTailor


[![GitHub Actions](https://github.com/spring-media/ott-tfprovider-awsmt/workflows/CI/badge.svg?branch=main)](https://github.com/spring-media/ott-tfprovider-awsmt/actions?workflow=CI)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=spring-media_ott-tfprovider-awsmt&metric=alert_status&token=06d658832169745b96bb3266679443282e48ace4)](https://sonarcloud.io/summary/new_code?id=spring-media_ott-tfprovider-awsmt)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=spring-media_ott-tfprovider-awsmt&metric=coverage&token=06d658832169745b96bb3266679443282e48ace4)](https://sonarcloud.io/summary/new_code?id=spring-media_ott-tfprovider-awsmt)

## Building the Provider

Run `make`.

## Provider Setup

By default, the provider sends requests to the `eu-central-1` aws region. You can override this default value by setting a region variable in the Terraform provider configuration.
For example, in `main.tf`:
```
provider "mediatailor" {
    region = "us-west-1"
}
```

## Querying Configurations

An example of how to query configurations from aws can be found in `./examples/main.tf`. 
Make sure that `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` are exported as environmental variables.

You can query a single configuration by specifying the `name` of the configuration, or all the 
configurations if you do not specify anything. Here are all the available parameters

Name | Type | Notes
---|---|---|
`name` | string | N/A
`max_results` | int | Ignored if `name` is specified
`next_token` | string | Ignored if `name` is specified

Run `terraform init` and then `terraform apply` inside the `./examples` directory to get a result.

## Testing

Run `make test` to execute both acceptance and unit tests.
