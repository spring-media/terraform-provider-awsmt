package awsmt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"regexp"
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
					resource.TestCheckResourceAttr("data.awsmt_playback_configuration.c1", "name", "test-acc-configuration"),
				),
			},
		},
	})
}

func TestAccPlaybackConfigurationDataSourceError(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPlaybackConfigurationDataSource2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.awsmt_playback_configuration.c2", "name", ""),
				),
				ExpectError: regexp.MustCompile("`name` parameter required"),
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if a, b := os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"); a == "" || b == "" {
		t.Fatal("AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY must both be set for acceptance tests")
	}
}

func testAccPlaybackConfigurationDataSource1() string {
	return `
data "awsmt_playback_configuration" "c1" {
  name = "test-acc-configuration"
}
output "out1" {
  value = data.awsmt_playback_configuration.c1
}
`
}
func testAccPlaybackConfigurationDataSource2() string {
	return `
data "awsmt_playback_configuration" "c2" {
  name = ""
}
output "out2" {
  value = data.awsmt_playback_configuration.c2
}
`
}
