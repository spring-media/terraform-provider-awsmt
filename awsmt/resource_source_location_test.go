package awsmt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccSourceLocationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: `resource "awsmt_source_location" "test_source_location"{
  name = "test_source_location"
  http_configuration = {
    base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/"

  }
}
data "awsmt_source_location" "read" {
  name = awsmt_source_location.test_source_location.name
}
output "awsmt_source_location" {
  value = data.awsmt_source_location.read
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "name", "test_source_location"),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "http_configuration.base_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/"),
				),
			},
			// ImportState testing
			/*
				{
					ResourceName: "awsmt_source_location.test_source_location",
					ImportState:  true,
				}, */
			// Update and Read testing
			{
				Config: `resource "awsmt_source_location" "test_source_location"{
  name = "test_source_location"
  http_configuration = {
    base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/"
  }
  default_segment_delivery_configuration = {
    base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"
  }
}
data "awsmt_source_location" "read" {
  name = awsmt_source_location.test_source_location.name
}
output "awsmt_source_location" {
  value = data.awsmt_source_location.read
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "name", "test_source_location"),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "http_configuration.base_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/"),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "default_segment_delivery_configuration.base_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
