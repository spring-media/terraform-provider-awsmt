package awsmt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccSourceLocationDataSourceBasic(t *testing.T) {
	dataSourceName := "data.awsmt_source_location.test"
	rName := "basic_source_location"
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceLocationDataSourceBasic(rName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceName, "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:sourceLocation\/.*$`)),
					resource.TestMatchResourceAttr(dataSourceName, "creation_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \+\d{4} \w+$`)),
					resource.TestMatchResourceAttr(dataSourceName, "last_modified_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr(dataSourceName, "source_location_name", rName),
				),
			},
		},
	})
}

func testAccSourceLocationDataSourceBasic(rName string) string {
	return fmt.Sprintf(`
resource "awsmt_source_location" "test_data_source"{
  access_configuration {
    access_type = "S3_SIGV4"
  }
  default_segment_delivery_configuration_url = "https://www.example.com"
  http_configuration_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"
  source_location_name = "%[1]s"
  segment_delivery_configurations {
    base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"
    name =     "example"
  }
}

data "awsmt_source_location" "test" {
  source_location_name = awsmt_source_location.test_data_source.source_location_name
}
`, rName)
}
