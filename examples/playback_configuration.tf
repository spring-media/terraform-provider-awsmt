resource "awsmt_playback_configuration" "r1" {
  ad_decision_server_url = "https://exampleurl.com/"
  avail_suppression = {
    mode = "AFTER_LIVE_EDGE"
    fill_policy = "FULL_AVAIL_ONLY"
  }
  cdn_configuration = {
    ad_segment_url_prefix = "https://exampleurl.com/"
  }
  dash_configuration = {
    mpd_location         = "DISABLED",
    origin_manifest_type = "SINGLE_PERIOD"
  }
  manifest_processing_rules = {
    ad_marker_passthrough = {
      enabled = "false"
    }
  }
  name                              = "example-playback-configuration-awsmt"
  personalization_threshold_seconds = 2
  tags                              = { "Environment" : "dev" }
  video_content_source_url          = "https://exampleurl.com/"
}

data "awsmt_playback_configuration" "test" {
  name = awsmt_playback_configuration.r1.name
}

output "playback_configuration_out" {
  value = data.awsmt_playback_configuration.test
}