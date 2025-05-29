resource "awsmt_playback_configuration" "r1" {
  ad_decision_server_url = "https://exampleurl.com/"
  avail_suppression = {
    fill_policy = "FULL_AVAIL_ONLY"
    mode        = "BEHIND_LIVE_EDGE"
    value       = "00:00:00"
  }
  bumper = {
    end_url   = "https://exampleurl.com/"
    start_url = "https://exampleurl.com/"
  }
  cdn_configuration = {
    ad_segment_url_prefix      = "https://exampleurl.com/"
    content_segment_url_prefix = "https://exampleurl.com/"
  }
  dash_configuration = {
    mpd_location         = "EMT_DEFAULT"
    origin_manifest_type = "MULTI_PERIOD"
  }
  live_pre_roll_configuration = {
    ad_decision_server_url = "https://exampleurl.com/"
    max_duration_seconds   = 2
  }
  manifest_processing_rules = {
    ad_marker_passthrough = {
      enabled = "false"
    }
  }
  name                                         = "example-playback-configuration-awsmt"
  personalization_threshold_seconds            = 2
  slate_ad_url                                 = "https://exampleurl.com/"
  tags                                         = { "Environment" : "dev" }
  video_content_source_url                     = "https://exampleurl.com/"
  transcode_profile_name                       = "profile_configured_in_your_account"
  log_configuration_percent_enabled            = 0
  log_configuration_enabled_logging_strategies = ["LEGACY_CLOUDWATCH", "VENDED_LOGS"]
}

data "awsmt_playback_configuration" "test" {
  name = awsmt_playback_configuration.r1.name
}

output "playback_configuration_out" {
  value = data.awsmt_playback_configuration.test
}