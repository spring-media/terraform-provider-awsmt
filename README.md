# ott-tfprovider-awsmt
A Terraform provider for AWS MediaTailor

## Building the Provider

1. Open `Makefile` and make sure that the value of `OS_ARCH` matches your system's architecture;
2. Run `make`.

## Provider Setup

By default, the provider sends requests to the `eu-central-1` aws region. You can override this default value by setting a region variable in the Terraform provider configuration.
For example, in `main.tf`:
```
provider "mediatailor" {
    region = "us-west-1"
}
```

## Querying Configurations

An example of how to query configurations from aws can be found in `./examples/main.tf`. Make sure that `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` are exported as environmental variables.
You can query a single configuration by specifying the `name` of the configuration, or all the configurations if you do not specify anything.

Run `terraform init` and then `terraform apply` inside the `./examples` directory to get a result.

## Testing

### Acceptance Testing
Run `make testacc` to execute the acceptance tests.

### Unit Testing
1. Navigate to `./mediatailor`;
2. Run `go test`.