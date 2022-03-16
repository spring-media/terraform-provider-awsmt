terraform {
  required_providers {
    awsmt = {
      version = "~> 1.8.0"
      source  = "spring-media/awsmt"
    }
  }
}

#data "awsmt_playback_configuration" "c1" {
#  name = "replay-live-stream"
#}

resource "awsmt_playback_configuration" "r1" {
  ad_decision_server_url = "https://exampleurl.com/"
  avail_suppression {
   mode = "OFF"
  }
  bumper {}
  cdn_configuration {
    ad_segment_url_prefix = "https://exampleurl.com/"
  }
  dash_configuration {
    mpd_location = "DISABLED"
    origin_manifest_type = "SINGLE_PERIOD"
  }
  live_pre_roll_configuration {}
  manifest_processing_rules {
    ad_marker_passthrough{
      enabled = "false"
    }
  }
  name = "example-playback-configuration-awsmt"
  personalization_threshold_seconds = 2
  tags = {"Environment": "dev"}
  video_content_source_url = "https://exampleurl.com/"
}

output "out" {
  value = resource.awsmt_playback_configuration.r1
  # value = data.awsmt_playback_configuration.c1
}
