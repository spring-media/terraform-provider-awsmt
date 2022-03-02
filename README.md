# terraform-provider-awsmt
A Terraform provider for AWS MediaTailor


[![GitHub Actions](https://github.com/spring-media/ott-tfprovider-awsmt/workflows/CI/badge.svg?branch=main)](https://github.com/spring-media/ott-tfprovider-awsmt/actions?workflow=CI)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=spring-media_ott-tfprovider-awsmt&metric=alert_status&token=06d658832169745b96bb3266679443282e48ace4)](https://sonarcloud.io/summary/new_code?id=spring-media_ott-tfprovider-awsmt)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=spring-media_ott-tfprovider-awsmt&metric=coverage&token=06d658832169745b96bb3266679443282e48ace4)](https://sonarcloud.io/summary/new_code?id=spring-media_ott-tfprovider-awsmt)

## Documentation

You can find the documentation about the provider, its data sources, and its resources [here](https://github.com/spring-media/terraform-provider-awsmt/tree/main/docs).

## Building the Provider

to use a local version of the provider, run `make` and create a ~/.terraformrc file with the following content:
```
provider_installation {
    dev_overrides {
      "spring-media/awsmt" = "/Users/<USERNAME>/.terraform.d/plugins/github.com/spring-media/terraform-provider-awsmt/0.0.1/<SYSTEM_ARCHITECTURE>"
    }
    direct {}
}
```

## Querying Configurations

An example of how to query configurations from aws can be found in `./examples/main.tf`. 
Make sure that `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` are exported as environmental variables.

You can query a single configuration by specifying the `name` of the configuration. 

Run `terraform init` and then `terraform apply` inside the `./examples` directory to get a result.

## Testing

Run `make test` to execute both acceptance and unit tests.
