package awsmt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccSourceLocationResource(t *testing.T) {
	name := "test_source_location"
	base_url := "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com"
	base_url2 := "https://example.com/"
	k1 := "Environment"
	v1 := "dev"
	k2 := "Testing"
	v2 := "pass"
	k3 := "Environment"
	v3 := "prod"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: basicSourceLocation(name, base_url, k1, v1, k2, v2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "id", "test_source_location"),
					resource.TestMatchResourceAttr("awsmt_source_location.test_source_location", "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:sourceLocation\/.*$`)),
					resource.TestMatchResourceAttr("awsmt_source_location.test_source_location", "creation_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "default_segment_delivery_configuration.base_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com"),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "http_configuration.base_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com"),
					resource.TestMatchResourceAttr("awsmt_source_location.test_source_location", "last_modified_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "segment_delivery_configurations.0.base_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com"),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "name", "test_source_location"),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "tags.Testing", "pass"),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "tags.Environment", "dev"),
				),
			},
			// ImportState testing
			{
				ResourceName: "awsmt_source_location.test_source_location",
				ImportState:  true,
			},
			// Update and Read testing
			{
				Config: basicSourceLocationWithAccessConfig(name, base_url2, k3, v3, k2, v2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "id", "test_source_location"),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "access_configuration.access_type", "S3_SIGV4"),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "name", "test_source_location"),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "http_configuration.base_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com"),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "default_segment_delivery_configuration.base_url", "https://example.com/"),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "segment_delivery_configurations.0.base_url", "https://example.com/"),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "tags.Testing", "pass"),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "tags.Environment", "prod")),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
func TestAccSourceLocationDeleteVodResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: basicSourceLocationWithVodSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "id", "test_source_location"),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "name", "test_source_location"),
				),
			},
		},
	})
}

func basicSourceLocation(name, base_url, k1, v1, k2, v2 string) string {
	return fmt.Sprintf(`resource "awsmt_source_location" "test_source_location"{
  							name = "%[1]s"
  							http_configuration = {
    							base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com"
  							}
  							default_segment_delivery_configuration = {
    							base_url = "%[2]s"
  							}
							segment_delivery_configurations = [{
    							base_url = "%[2]s"
								name = "default"
  							}]
							tags = {
   		 						"%[3]s": "%[4]s",
								"%[5]s": "%[6]s"
							}
						}
						data "awsmt_source_location" "read" {
  							name = awsmt_source_location.test_source_location.name
						}
						output "awsmt_source_location" {
  							value = data.awsmt_source_location.read
						}
						`, name, base_url, k1, v1, k2, v2)

}

func basicSourceLocationWithAccessConfig(name, base_url, k1, v1, k2, v2 string) string {
	return fmt.Sprintf(`resource "awsmt_source_location" "test_source_location"{
  							name = "%[1]s"
  							http_configuration = {
    							base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com"
  							}
							access_configuration = {
								access_type = "S3_SIGV4"
							}
  							default_segment_delivery_configuration = {
    							base_url = "%[2]s"
  							}
							segment_delivery_configurations = [{
    							base_url = "%[2]s"
								name = "default"
  							}]
							tags = {
   		 						"%[3]s": "%[4]s",
								"%[5]s": "%[6]s"
							}
						}
						data "awsmt_source_location" "read" {
  							name = awsmt_source_location.test_source_location.name
						}
						output "awsmt_source_location" {
  							value = data.awsmt_source_location.read
						}
						`, name, base_url, k1, v1, k2, v2)
}

func basicSourceLocationWithVodSource() string {
	return `resource "awsmt_source_location" "test_source_location"{
  							name = "test_source_location"
  							http_configuration = {
    							base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com"
  							}
  							default_segment_delivery_configuration = {
    							base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"
  							}
							segment_delivery_configurations = [{
    							base_url = "https://example.com/"
								name = "default"
  							}]
							tags = {"Environment": "dev"}
						}
						data "awsmt_source_location" "read" {
  							name = awsmt_source_location.test_source_location.name
						}
						output "awsmt_source_location" {
  							value = data.awsmt_source_location.read
						}
						resource "awsmt_vod_source" "test" {
  					http_package_configurations = [{
						path = "/test"
						source_group = "default"
    					type = "HLS"
  					}]
  					source_location_name = awsmt_source_location.test_source_location.name
  					name = "vod_source_example"
				}
`
}
