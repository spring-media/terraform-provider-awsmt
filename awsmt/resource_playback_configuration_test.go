package awsmt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccPlaybackConfigurationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `resource "awsmt_playback_configuration" "r1" {
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
						`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "id", "example-playback-configuration-awsmt"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "ad_decision_server_url", "https://exampleurl.com/"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "avail_supression.fill_policy", "FULL_AVAIL_ONLY"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "avail_supression.mode", "BEHIND_LIVE_EDGE"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "avail_supression.value", "00:00:00"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "bumper.end_url", "https://wxample.com/endbumper"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "bumper.start_url", "https://wxample.com/startbumper"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "cdn_configuration.ad_segment_url_prefix", "https://exampleurl.com/"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "dash_configuration.mpd_location", "DISABLED"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "dash_configuration.origin_manifest_type", "SINGLE_PERIOD"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "live_pre_roll_configuration.ad_decision_server_url", "https://exampleurl.com/"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "live_pre_roll_configuration.max_duration_seconds", "2"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "manifest_processing_rules.ad_marker_passthrough.enabled", "false"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "name", "example-playback-configuration-awsmt"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "personalization_threshold_seconds", "2"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "slate_ad_url", "https://exampleurl.com/"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "tags.Environment", "dev"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "video_content_source_url", "https://exampleurl.com/"),
				),
			},
			// ImportState testing
			{
				ResourceName: "awsmt_playback_configuration.r1",
				ImportState:  true,
			},
			// Update and Read testing
			{
				Config: `resource "awsmt_playback_configuration" "r1" {
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
  							personalization_threshold_seconds = 3
							slate_ad_url = "https://exampleurl.com/"
  							tags = {"Environment": "dev", "Name": "example-playback-configuration-awsmt"}
 	 						video_content_source_url = "https://exampleurl.com/"
						}

						data "awsmt_playback_configuration" "test"{
  							name = awsmt_playback_configuration.r1.name
						}

						output "playback_configuration_out" {
  							value = data.awsmt_playback_configuration.test
						}
						`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "name", "example-playback-configuration-awsmt"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "personalization_threshold_seconds", "3"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "log_configuration_percent_enabled", "0"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "tags.Name", "example-playback-configuration-awsmt"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
