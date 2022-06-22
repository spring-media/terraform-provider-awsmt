# AWSMT Provider

The AWSMT Provider allows you to interact with AWS Elemental MediaTailor
using Terraform. You need to export your `AW_ACCESS_KEY_ID` and
`AWS_SECRET_ACCESS_KEY` as environmental variables in order to use this provider.

## Configuration

Example configuration (using Terraform 0.13 or newer): 
```
terraform {
  required_providers {
    awsmt = {
      version = "1.17.0"
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
  version = "1.17.0"
  region = "eu-central-1"
}
```

## Argument Reference

The AWSMT Provider supports the following argument:

* `region` - (Optional).<br/> AWS region code, defaults to `eu-central-1`. 
You can learn more about aws regions and the available codes [here](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html).

* `profile` - (Optional).<br/> AWS configuration profile.
  You can find the profile(s) name in '~/.aws/config' (Mac & Linux) or '%USERPROFILE%\.aws\config' (Windows). SSO login will be used if the profile name is specified or if an environmental variable called 'AWS_PROFILE' is found. Please note that the value of the environmental variable is ignored if an explicit declaration is found.