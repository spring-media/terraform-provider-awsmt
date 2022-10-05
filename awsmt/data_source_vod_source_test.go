package awsmt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccVodSourceDataSourceBasic(t *testing.T) {
	dataSourceName := "data.awsmt_vod_source.test"
	sourceLocationName := "vod_basic_sl"
	vodSourceName := "vod_source_data_source_test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVodSourceDataSourceBasic(sourceLocationName, vodSourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceName, "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:vodSource\/.*$`)),
					resource.TestMatchResourceAttr(dataSourceName, "creation_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{3})? \+\d{4} \w+$`)),
					resource.TestMatchResourceAttr(dataSourceName, "last_modified_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr(dataSourceName, "source_location_name", sourceLocationName),
					resource.TestCheckResourceAttr(dataSourceName, "name", vodSourceName),
				),
			},
		},
	})
}

func testAccVodSourceDataSourceBasic(sourceLocationName, vodSourceName string) string {
	return fmt.Sprintf(`
resource "awsmt_source_location" "example"{
  name = "%[1]s"
  http_configuration_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"
}

resource "awsmt_vod_source" "test" {
  http_package_configurations {
    path = "/"
    source_group = "default"
    type = "HLS"
  }
  source_location_name = awsmt_source_location.example.name
  name = "%[2]s"
}

data "awsmt_vod_source" "test" {
  source_location_name = awsmt_source_location.example.name
  name = awsmt_vod_source.test.name
}
`, sourceLocationName, vodSourceName)
}
