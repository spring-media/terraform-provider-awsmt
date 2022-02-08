# AWSMT Provider

The AWSMT Provider allows you to interact with AWS Elemental Media Tailor
using Terraform. You need to export your `AW_ACCESS_KEY_ID` and
`AWS_SECRET_ACCESS_KEY` as environmental variables in order to use this provider.

## Configuration

Example configuration (using Terraform 0.13 or newer): 
```
terraform {
  required_providers {
    awsmt = {
      version = "1.1.0"
      source  = "spring-media/awsmt"
     }
  }
}

provider "awsmt" {
  region = "eu-central-1"
}
```

## Arguments

The AWSMT Provider supports the following argument:

* `region` - (optional, type string). AWS region code, defaults to `eu-central-1`. 
You can learn more about aws regions and the available codes [here](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-regions-availability-zones.html).
