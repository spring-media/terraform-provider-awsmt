# AWSMT Provider

The AWSMT Provider allows you to interact with AWS Elemental MediaTailor
using Terraform.

# Authentication

This provider offers 3 authentication options, and tries to authenticate you in the following order:

1. Using SSO, using the `profile` from the provider configuration;
2. Using SSO, using an environmental variable called `AWS_PROFILE`;
3. Using the `AW_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` environmental variables.

## Configuration

Example configuration (using Terraform 0.13 or newer):

```
terraform {
  required_providers {
    awsmt = {
      version = "~> 1.17"
      source  = "spring-media/awsmt"
     }
  }
}

provider "awsmt" {
  region = "eu-central-1"
}
```

Terraform 0.12 or earlier:

```
provider "awsmt" {
  version = "–> 1.17"
  region = "eu-central-1"
}
```

## Argument Reference

The AWSMT Provider supports the following argument:

- `region` - (Optional) AWS region code, defaults to `eu-central-1`.
  You can learn more about aws regions and the available codes [here](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html).

- `profile` - (Optional) AWS configuration profile.
  You can find the profile(s) name in `~/.aws/config` (Mac & Linux) or `%USERPROFILE%\.aws\config` (Windows).
