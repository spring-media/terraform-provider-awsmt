package awsmt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
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
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "id", "test_source_location"),
					resource.TestMatchResourceAttr("awsmt_source_location.test_source_location", "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:sourceLocation\/.*$`)),
					resource.TestMatchResourceAttr("awsmt_source_location.test_source_location", "creation_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "default_segment_delivery_configuration.base_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "http_configuration.base_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com"),
					resource.TestMatchResourceAttr("awsmt_source_location.test_source_location", "last_modified_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "segment_delivery_configurations.0.base_url", "https://example.com/"),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "name", "test_source_location"),
				),
			},
			// ImportState testing
			{
				ResourceName: "awsmt_source_location.test_source_location",
				ImportState:  true,
			},
			// Update and Read testing
			{
				Config: `resource "awsmt_source_location" "test_source_location"{
							access_configuration = {
								access_type = "S3_SIGV4"
							}
  							name = "test_source_location"
  							http_configuration = {
    							base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/"
  							}
  							default_segment_delivery_configuration = {
    							base_url = "https://example.com/"
  							}
							segment_delivery_configurations = [{
    							base_url = "https://example.com/"
								name = "default"
  							}]
							tags = {"Environment": "dev", "Name": "test_source_location"}
						}
						data "awsmt_source_location" "read" {
  							name = awsmt_source_location.test_source_location.name
						}
						output "awsmt_source_location" {
  							value = data.awsmt_source_location.read
						}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "id", "test_source_location"),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "access_configuration.access_type", "S3_SIGV4"),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "name", "test_source_location"),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "http_configuration.base_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/"),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "default_segment_delivery_configuration.base_url", "https://example.com/"),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "tags.Name", "test_source_location"),
				),
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
				Config: `resource "awsmt_source_location" "test_source_location"{
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
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "id", "test_source_location"),
					resource.TestCheckResourceAttr("awsmt_source_location.test_source_location", "name", "test_source_location"),
				),
			},
		},
	})
}
