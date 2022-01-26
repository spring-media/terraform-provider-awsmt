terraform {
  required_providers {
    mediatailor = {
      version = "0.0.1"
      source  = "github.com/spring-media/ott-tfprovider-awsmt"
    }
  }
}

data "mediatailor_configuration" "c1" {
  max_results = 1
  //name = "staging-live-stream"
}

output "out" {
  value = data.mediatailor_configuration.c1
}
