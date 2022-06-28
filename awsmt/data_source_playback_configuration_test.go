package awsmt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccPlaybackConfigurationDataSourceBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPlaybackConfigurationDataSource1(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.awsmt_playback_configuration.c1", "name", "testacc_example_playback"),
				),
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if a, b, c := os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), os.Getenv("AWS_PROFILE"); (a == "" || b == "") || c == "" {
		t.Fatal("Either AWS_PROFILE or both AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY must be set for acceptance tests")
	}
}

func testAccPlaybackConfigurationDataSource1() string {
	return `
resource "awsmt_playback_configuration" "test"{
  ad_decision_server_url = "https://exampleurl.com/"
  name= "testacc_example_playback"
  dash_configuration {
    mpd_location = "EMT_DEFAULT"
    origin_manifest_type = "MULTI_PERIOD"
  }
  video_content_source_url = "https://exampleurl.com"
}
data "awsmt_playback_configuration" "c1" {
  name = awsmt_playback_configuration.test.name
}
`
}
