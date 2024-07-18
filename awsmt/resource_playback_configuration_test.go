package awsmt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccPlaybackConfigurationResource(t *testing.T) {
	name := "example-playback-configuration-awsmt"
	ad_url := "https://exampleurl.com/"
	ad_url2 := "https://exampleurl2.com/"
	bumper_e := "https://wxample.com/endbumper"
	bumper_e2 := "https://wxample.com/endbumper2"
	bumper_s := "https://wxample.com/startbumper"
	bumper_s2 := "https://wxample.com/startbumper2"
	cdn_url := "https://exampleurl.com/"
	cdn_url2 := "https://exampleurl2.com/"
	max_d := "2"
	max_d2 := "3"
	p_s := "2"
	p_s2 := "3"
	k1 := "Environment"
	v1 := "dev"
	k2 := "Testing"
	v2 := "pass"
	k3 := "Environment"
	v3 := "prod"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: basicPlaybackConfiguration(name, ad_url, bumper_e, bumper_s, cdn_url, max_d, p_s, k1, v1, k2, v2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "id", "example-playback-configuration-awsmt"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "ad_decision_server_url", "https://exampleurl.com/"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "avail_suppression.fill_policy", "FULL_AVAIL_ONLY"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "avail_suppression.mode", "BEHIND_LIVE_EDGE"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "avail_suppression.value", "00:00:00"),
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
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "tags.Testing", "pass"),
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
				Config: basicPlaybackConfiguration(name, ad_url2, bumper_e2, bumper_s2, cdn_url2, max_d2, p_s2, k3, v3, k2, v2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "name", "example-playback-configuration-awsmt"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "personalization_threshold_seconds", "3"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "log_configuration_percent_enabled", "0"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "ad_decision_server_url", "https://exampleurl2.com/"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "bumper.end_url", "https://wxample.com/endbumper2"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "bumper.start_url", "https://wxample.com/startbumper2"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "cdn_configuration.ad_segment_url_prefix", "https://exampleurl2.com/"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "tags.Environment", "prod"),
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "tags.Testing", "pass"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func basicPlaybackConfiguration(name, ad_url, bumper_e, bumper_s, cdn_url, max_d, p_s, k1, v1, k2, v2 string) string {
	return fmt.Sprintf(`resource "awsmt_playback_configuration" "r1" {
  							ad_decision_server_url = "%[2]s"
  							avail_suppression = {
								fill_policy = "FULL_AVAIL_ONLY"
    							mode = "BEHIND_LIVE_EDGE"
								value = "00:00:00"
							}
							bumper = {
								end_url = "%[3]s"
    							start_url = "%[4]s"
  							}
  							cdn_configuration = {
    							ad_segment_url_prefix = "%[5]s"
  							}
  							dash_configuration = {
    							mpd_location = "DISABLED",
    							origin_manifest_type = "SINGLE_PERIOD"
  							}
							live_pre_roll_configuration = {
								ad_decision_server_url = "%[2]s"
								max_duration_seconds = "%[6]s"
							}
  							manifest_processing_rules = {
								ad_marker_passthrough = {
      								enabled = "false"
    							}
  							}
  							name = "%[1]s"
  							personalization_threshold_seconds = "%[7]s"
							slate_ad_url = "https://exampleurl.com/"
  							tags = {
   		 						"%[8]s": "%[9]s",
								"%[10]s": "%[11]s"
							}
 	 						video_content_source_url = "%[2]s"
						}

						data "awsmt_playback_configuration" "test"{
  							name = awsmt_playback_configuration.r1.name
						}

						output "playback_configuration_out" {
  							value = data.awsmt_playback_configuration.test
						}
						`, name, ad_url, bumper_e, bumper_s, cdn_url, max_d, p_s, k1, v1, k2, v2)

}
