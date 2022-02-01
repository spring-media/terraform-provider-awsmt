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
				Config: testAccPlaybackConfigurationDataSourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.aws_playback_configuration.c1", "name", "broadcast-prod-live-stream"),
					resource.TestCheckResourceAttr("data.aws_playback_configuration.c2", "max_results", "2"),
				),
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if a, b := os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"); a == "" || b == "" {
		t.Fatal("AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY must both be set for acceptance tests")
	}
}

func testAccPlaybackConfigurationDataSourceBasic() string {
	return `
data "aws_playback_configuration" "c1" {
  name = "broadcast-prod-live-stream"
}

data "aws_playback_configuration" "c2" {
  max_results = 2
}

output "out1" {
  value = data.aws_playback_configuration.c1
}

output "out2" {
  value = data.aws_playback_configuration.c2
}
`
}
