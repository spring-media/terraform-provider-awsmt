terraform {
  required_providers {
    awsmt = {
      version = "1.0.9"
      source  = "spring-media/awsmt"
      // to use a local version of the provider,
      // run `make` and set the source to:
      // source = "github.com/spring-media/terraform-provider-awsmt"
    }
  }
}

data "awsmt_playback_configuration" "c1" {
  name="broadcast-staging-live-stream"
}

output "out" {
  value = data.awsmt_playback_configuration.c1
}
