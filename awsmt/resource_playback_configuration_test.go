package awsmt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccPlaybackConfigurationResourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPlaybackConfigurationResource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("resource.awsmt_playback_configuration.c1", "name", "test-playback-configuration-awsmt"),
					//resource.TestCheckResourceAttr("resource.awsmt_playback_configuration.c1", "name", "broadcast-prod-live-stream"),
				),
			},
		},
	})
}

func testAccPlaybackConfigurationResource() string {
	return `
resource "awsmt_playback_configuration" "r1" {
  ad_decision_server_url = "https://exampleurl.com/"
  cdn_configuration {
    ad_segment_url_prefix = "test"
    content_segment_url_prefix = "test"
  }
  dash_configuration {
    mpd_location = "test"
    origin_manifest_type = "MULTI_PERIOD"
  }
  name = "test-playback-configuration-awsmt"
  slate_ad_url = "https://exampleurl.com/"
  tags = {}
  transcode_profile_name = "test"
  video_content_source_url = "https://exampleurl.com/"
}

output "out" {
  value = resource.awsmt_playback_configuration.r1
}
`
}
