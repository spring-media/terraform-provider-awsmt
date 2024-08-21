package awsmt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccPlaybackConfigurationDataSourceBasic(t *testing.T) {
	resourceName := "data.awsmt_playback_configuration.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: playbackConfigDS(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "example-playback-configuration-awsmt"),
					resource.TestCheckResourceAttr(resourceName, "ad_decision_server_url", "https://exampleurl.com/"),
					resource.TestCheckResourceAttr(resourceName, "avail_suppression.fill_policy", "FULL_AVAIL_ONLY"),
					resource.TestCheckResourceAttr(resourceName, "avail_suppression.mode", "BEHIND_LIVE_EDGE"),
					resource.TestCheckResourceAttr(resourceName, "avail_suppression.value", "00:00:00"),
					resource.TestCheckResourceAttr(resourceName, "bumper.end_url", "https://wxample.com/endbumper"),
					resource.TestCheckResourceAttr(resourceName, "bumper.start_url", "https://wxample.com/startbumper"),
					resource.TestCheckResourceAttr(resourceName, "cdn_configuration.ad_segment_url_prefix", "https://exampleurl.com/"),
					resource.TestCheckResourceAttr(resourceName, "cdn_configuration.content_segment_url_prefix", "https://exampleurl.com/"),
					resource.TestCheckResourceAttr(resourceName, "configuration_aliases.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "configuration_aliases.player_params.foo.player_params.bar", "player_params.buzz"),
					resource.TestCheckResourceAttr(resourceName, "dash_configuration.mpd_location", "DISABLED"),
					resource.TestCheckResourceAttr(resourceName, "dash_configuration.origin_manifest_type", "SINGLE_PERIOD"),
					resource.TestCheckResourceAttr(resourceName, "live_pre_roll_configuration.ad_decision_server_url", "https://exampleurl.com/"),
					resource.TestCheckResourceAttr(resourceName, "live_pre_roll_configuration.max_duration_seconds", "2"),
					resource.TestCheckResourceAttr(resourceName, "name", "example-playback-configuration-awsmt"),
					resource.TestCheckResourceAttr(resourceName, "personalization_threshold_seconds", "2"),
					resource.TestCheckResourceAttr(resourceName, "slate_ad_url", "https://exampleurl.com/"),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "dev"),
					resource.TestCheckResourceAttr(resourceName, "video_content_source_url", "https://exampleurl.com/"),
					resource.TestCheckResourceAttr(resourceName, "log_configuration_percent_enabled", "0"),
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
				Config:      playbackConfigDSError(),
				ExpectError: regexp.MustCompile("Error while retrieving the playback configuration "),
			},
		},
	})
}

func playbackConfigDS() string {
	return `resource "awsmt_playback_configuration" "r1" {
  							ad_decision_server_url = "https://exampleurl.com/"
  							avail_suppression = {
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
							configuration_aliases = {
								"player_params.foo" = {
									"player_params.bar" = "player_params.buzz"
								}
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

func playbackConfigDSError() string {
	return `resource "awsmt_playback_configuration" "r1" {
  							ad_decision_server_url = "https://exampleurl.com/"
  							avail_suppression = {
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
