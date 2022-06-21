package awsmt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccChannelDataSourceBasic(t *testing.T) {
	dataSourceName := "data.awsmt_channel.test"
	rName := "basic_channel"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccChannelDataSourceBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceName, "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:channel\/.*$`)),
					resource.TestMatchResourceAttr(dataSourceName, "creation_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \+\d{4} \w+$`)),
					resource.TestMatchResourceAttr(dataSourceName, "last_modified_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr(dataSourceName, "outputs.0.hls_manifest_windows_seconds", "30"),
					resource.TestCheckResourceAttr(dataSourceName, "name", rName),
				),
			},
		},
	})
}

func testAccChannelDataSourceBasic(rName string) string {
	return fmt.Sprintf(`
resource "awsmt_channel" "test" {
  name = "%[1]s"
  outputs {
    manifest_name                = "default"
    source_group                 = "default"
    hls_manifest_windows_seconds = 30
  }
  playback_mode = "LOOP"
  tier = "BASIC"
}

data "awsmt_channel" "test" {
  name = awsmt_channel.test.name
}
`, rName)
}
