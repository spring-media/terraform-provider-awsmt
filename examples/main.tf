terraform {
  required_providers {
    mediatailor = {
      version = "0.1"
      source  = "test/ott/mediatailor"
    }
  }
}

data "mediatailor_configuration" "c1" {
  //name = "staging-live-stream"
}

output "out" {
  value = data.mediatailor_configuration.c1
}