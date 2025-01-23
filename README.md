# terraform-provider-awsmt

A Terraform provider for AWS MediaTailor

[![GitHub Actions](https://github.com/spring-media/ott-tfprovider-awsmt/workflows/CI/badge.svg?branch=main)](https://github.com/spring-media/ott-tfprovider-awsmt/actions?workflow=CI)

## Documentation

You can find the documentation about the provider, its resources and its data sources [here](https://registry.terraform.io/providers/spring-media/awsmt/latest/docs).

## Building the Provider

Run `make`.

To use a local version of the provider, create a ~/.terraformrc file with the following content:

```
provider_installation {
    dev_overrides {
      "spring-media/awsmt" = "/Users/<USERNAME>/.terraform.d/plugins/github.com/spring-media/terraform-provider-awsmt/0.0.1/<SYSTEM_ARCHITECTURE>"
    }
    direct {}
}
```

## Prerequisites

Make sure to be logged in. To learn more about log in methods, please refer to the [official documentation](https://registry.terraform.io/providers/spring-media/awsmt/latest/docs).

## Testing

> **NOTE:** AWS credentials for AWS must be provided through environment variables `AWS_ACCESS_KEY_ID`,
> `AWS_SECRET_ACCESS_KEY` and `AWS_SESSION_TOKEN`.

The tests require the following environment variables to be defined:

```bash
export AWS_REGION=eu-central-1
export AWS_PROFILE=as-nmt-ott-test
export AWS_ACCOUNT_ID=985600762523
```

Run `make clean sweep test` to execute both acceptance and unit tests.
Run `make sweep` to delete resources that might not have been automatically destroyed after the tests were run.
