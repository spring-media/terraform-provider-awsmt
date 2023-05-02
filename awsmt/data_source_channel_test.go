package awsmt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccChannelDataSourceBasic(t *testing.T) {
	dataSourceName := "data.awsmt_channel.testing"
	resourceName := "basic_channel"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccChannelDataSourceBasic(resourceName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceName, "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:channel\/.*$`)),
					resource.TestMatchResourceAttr(dataSourceName, "creation_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestMatchResourceAttr(dataSourceName, "last_modified_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr(dataSourceName, "outputs.0.hls_manifest_windows_seconds", "30"),
					resource.TestCheckResourceAttr(dataSourceName, "name", resourceName),
				),
			},
		},
	})
}

func testAccChannelDataSourceBasic(resourceName string) string {
	return fmt.Sprintf(`
resource "awsmt_channel" "testing"  {
  name = "%[1]s"
  outputs = [{
    manifest_name                = "default"
    source_group                 = "default"
    hls_playlist_settings = [{manifest_windows_seconds = 30}]
  }]
  playback_mode = "LOOP"
  tier = "BASIC"
}

data "awsmt_channel" "read" {
  name = "sample-1"
}
`, resourceName)
}
