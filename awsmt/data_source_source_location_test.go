package awsmt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccSourceLocationDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "awsmt_source_location" "test_source_location"{
  name = "test_source_location"
  http_configuration = {
    base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/"
  }
  default_segment_delivery_configuration = {
    base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"
  }
}
data "awsmt_source_location" "test" {
  name = awsmt_source_location.test_source_location.name
}
output "awsmt_source_location" {
  value = data.awsmt_source_location.test
}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.awsmt_source_location.test", "name", "test_source_location"),
					resource.TestCheckResourceAttr("data.awsmt_source_location.test", "http_configuration.base_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/"),
					resource.TestCheckResourceAttr("data.awsmt_source_location.test", "default_segment_delivery_configuration.base_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"),
				),
			},
		},
	})
}
