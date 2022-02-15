package awsmt

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func init() {
	resource.AddTestSweepers("test_playback_configuration", &resource.Sweeper{
		Name: "test_playback_configuration",
		F: func(region string) error {
			client, err := sharedClientForRegion(region)
			if err != nil {
				return fmt.Errorf("Error getting client: %s", err)
			}
			conn := client.(*mediatailor.MediaTailor)
			name := "test-playback-configuration-awsmt"
			_, err = conn.DeletePlaybackConfiguration(&mediatailor.DeletePlaybackConfigurationInput{Name: &name})
			if err != nil {
				return err
			}
			return nil
		},
	})
}

func TestAccPlaybackConfigurationResourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPlaybackConfigurationResource(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_playback_configuration.r1", "name", "test-playback-configuration-awsmt"),
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
    mpd_location = "EMT_DEFAULT"
    origin_manifest_type = "MULTI_PERIOD"
  }
  name = "test-playback-configuration-awsmt"
  slate_ad_url = "https://exampleurl.com/"
  tags = {"Environment": "dev"}
  video_content_source_url = "https://exampleurl.com/"
}

`
}
