terraform {
  required_providers {
    awsmt = {
      version = "1.4.0"
      source  = "spring-media/awsmt"
      // to use a local version of the provider,
      // run `make` and create a ~/.terraformrc file with the following content:
      //provider_installation {
      //  dev_overrides {
      //      "spring-media/awsmt" = "/Users/<USERNAME>/.terraform.d/plugins/github.com/spring-media/terraform-provider-awsmt/0.0.1/<SYSTEM_ARCHITECTURE>"
      //  }
      //  direct {}
      //}
    }
  }
}

#data "awsmt_playback_configuration" "c1" {
#  name="broadcast-staging-live-stream"
#}

resource "awsmt_playback_configuration" "r1" {
  ad_decision_server_url = "https://exampleurl.com/"
  cdn_configuration {
    ad_segment_url_prefix = "test"
    content_segment_url_prefix = "test"
  }
  dash_configuration {
    mpd_location = "EMT_DEFAULT"
    origin_manifest_type = "MULTI_PERIOD"
  }
  name = "test-playback-configuration-awsmt"
  slate_ad_url = "https://exampleurl.com/"
  tags = {"Environment": "dev"}
  video_content_source_url = "https://exampleurl.com/"
}

output "out" {
  value = resource.awsmt_playback_configuration.r1
}
