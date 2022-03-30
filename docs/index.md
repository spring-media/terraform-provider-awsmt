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
      version = "1.10.0"
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
  version = "1.10.0"
  region = "eu-central-1"
}
```

## Argument Reference

The AWSMT Provider supports the following argument:

* `region` - (Optional, type string).<br/> AWS region code, defaults to `eu-central-1`. 
You can learn more about aws regions and the available codes [here](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html).
