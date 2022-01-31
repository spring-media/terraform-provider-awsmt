terraform {
  required_providers {
    mediatailor = {
      version = "1.0.6"
      source  = "github.com/spring-media/ott-tfprovider-awsmt"
    }
  }
}

data "mediatailor_configuration" "c1" {
  name = "staging-live-stream"
}

output "out" {
  value = data.mediatailor_configuration.c1
}
