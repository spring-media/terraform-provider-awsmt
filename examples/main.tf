terraform {
  required_providers {
    mediatailor = {
      version = "1.0.6"
      source  = "spring-media/awsmt"
      // to use a local version of the provider,
      // run `make` and set the source to:
      // source = "github.com/spring-media/terraform-provider-awsmt"
    }
  }
}

data "mediatailor_configuration" "c1" {
  max_results = 5
}

output "out" {
  value = data.mediatailor_configuration.c1
}
