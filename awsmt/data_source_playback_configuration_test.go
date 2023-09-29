package awsmt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccPlaybackConfigurationDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: playbackConfigDS(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.awsmt_playback_configuration.test", "id", "example-playback-configuration-awsmt"),
					resource.TestCheckResourceAttr("data.awsmt_playback_configuration.test", "ad_decision_server_url", "https://exampleurl.com/"),
					resource.TestCheckResourceAttr("data.awsmt_playback_configuration.test", "avail_supression.fill_policy", "FULL_AVAIL_ONLY"),
					resource.TestCheckResourceAttr("data.awsmt_playback_configuration.test", "avail_supression.mode", "BEHIND_LIVE_EDGE"),
					resource.TestCheckResourceAttr("data.awsmt_playback_configuration.test", "avail_supression.value", "00:00:00"),
					resource.TestCheckResourceAttr("data.awsmt_playback_configuration.test", "bumper.end_url", "https://wxample.com/endbumper"),
					resource.TestCheckResourceAttr("data.awsmt_playback_configuration.test", "bumper.start_url", "https://wxample.com/startbumper"),
					resource.TestCheckResourceAttr("data.awsmt_playback_configuration.test", "cdn_configuration.ad_segment_url_prefix", "https://exampleurl.com/"),
					resource.TestCheckResourceAttr("data.awsmt_playback_configuration.test", "cdn_configuration.content_segment_url_prefix", "https://exampleurl.com/"),
					resource.TestCheckResourceAttr("data.awsmt_playback_configuration.test", "dash_configuration.mpd_location", "DISABLED"),
					resource.TestCheckResourceAttr("data.awsmt_playback_configuration.test", "dash_configuration.origin_manifest_type", "SINGLE_PERIOD"),
					resource.TestCheckResourceAttr("data.awsmt_playback_configuration.test", "live_pre_roll_configuration.ad_decision_server_url", "https://exampleurl.com/"),
					resource.TestCheckResourceAttr("data.awsmt_playback_configuration.test", "live_pre_roll_configuration.max_duration_seconds", "2"),
					resource.TestCheckResourceAttr("data.awsmt_playback_configuration.test", "manifest_processing_rules.ad_marker_passthrough.enabled", "false"),
					resource.TestCheckResourceAttr("data.awsmt_playback_configuration.test", "name", "example-playback-configuration-awsmt"),
					resource.TestCheckResourceAttr("data.awsmt_playback_configuration.test", "personalization_threshold_seconds", "2"),
					resource.TestCheckResourceAttr("data.awsmt_playback_configuration.test", "slate_ad_url", "https://exampleurl.com/"),
					resource.TestCheckResourceAttr("data.awsmt_playback_configuration.test", "tags.Environment", "dev"),
					resource.TestCheckResourceAttr("data.awsmt_playback_configuration.test", "video_content_source_url", "https://exampleurl.com/"),
					resource.TestCheckResourceAttr("data.awsmt_playback_configuration.test", "log_configuration_percent_enabled", "0"),
				),
			},
		},
	})
}

func TestAccPlaybackConfigurationDataSourceErrors(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      plabackConfigDSError(),
				ExpectError: regexp.MustCompile("Error while retrieving the playback configuration "),
			},
		},
	})
}

func playbackConfigDS() string {
	return `resource "awsmt_playback_configuration" "r1" {
  							ad_decision_server_url = "https://exampleurl.com/"
  							avail_supression = {
								fill_policy = "FULL_AVAIL_ONLY"
    							mode = "BEHIND_LIVE_EDGE"
								value = "00:00:00"
							}
							bumper = {
								end_url = "https://wxample.com/endbumper"
    							start_url = "https://wxample.com/startbumper"
  							}
  							cdn_configuration = {
    							ad_segment_url_prefix = "https://exampleurl.com/"
								content_segment_url_prefix = "https://exampleurl.com/"
  							}
  							dash_configuration = {
    							mpd_location = "DISABLED",
    							origin_manifest_type = "SINGLE_PERIOD"
  							}
							live_pre_roll_configuration = {
								ad_decision_server_url = "https://exampleurl.com/"
								max_duration_seconds = 2
							}
  							manifest_processing_rules = {
								ad_marker_passthrough = {
      								enabled = "false"
    							}
  							}
  							name = "example-playback-configuration-awsmt"
  							personalization_threshold_seconds = 2
							slate_ad_url = "https://exampleurl.com/"
  							tags = {"Environment": "dev"}
 	 						video_content_source_url = "https://exampleurl.com/"
						}

						data "awsmt_playback_configuration" "test"{
  							name = awsmt_playback_configuration.r1.name
						}

						output "playback_configuration_out" {
  							value = data.awsmt_playback_configuration.test
						}
						`
}

func plabackConfigDSError() string {
	return `resource "awsmt_playback_configuration" "r1" {
  							ad_decision_server_url = "https://exampleurl.com/"
  							avail_supression = {
								fill_policy = "FULL_AVAIL_ONLY"
    							mode = "BEHIND_LIVE_EDGE"
								value = "00:00:00"
							}
							bumper = {
								end_url = "https://wxample.com/endbumper"
    							start_url = "https://wxample.com/startbumper"
  							}
  							cdn_configuration = {
    							ad_segment_url_prefix = "https://exampleurl.com/"
								content_segment_url_prefix = "https://exampleurl.com/"
  							}
  							dash_configuration = {
    							mpd_location = "DISABLED",
    							origin_manifest_type = "SINGLE_PERIOD"
  							}
							live_pre_roll_configuration = {
								ad_decision_server_url = "https://exampleurl.com/"
								max_duration_seconds = 2
							}
  							manifest_processing_rules = {
								ad_marker_passthrough = {
      								enabled = "false"
    							}
  							}
  							name = "example-playback-configuration-awsmt"
  							personalization_threshold_seconds = 2
							slate_ad_url = "https://exampleurl.com/"
  							tags = {"Environment": "dev"}
 	 						video_content_source_url = "https://exampleurl.com/"
						}

						data "awsmt_playback_configuration" "test"{
  							name = "testingErrors"
						}

						output "playback_configuration_out" {
  							value = data.awsmt_playback_configuration.test
						}
						`
}
