package awsmt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccSourceLocationDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: sourceLocationDS(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.awsmt_source_location.read", "id", "test_source_location"),
					resource.TestMatchResourceAttr("data.awsmt_source_location.read", "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:sourceLocation\/.*$`)),
					resource.TestMatchResourceAttr("data.awsmt_source_location.read", "creation_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr("data.awsmt_source_location.read", "default_segment_delivery_configuration.base_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"),
					resource.TestCheckResourceAttr("data.awsmt_source_location.read", "http_configuration.base_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com"),
					resource.TestMatchResourceAttr("data.awsmt_source_location.read", "last_modified_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr("data.awsmt_source_location.read", "segment_delivery_configurations.0.base_url", "https://example.com/"),
					resource.TestCheckResourceAttr("data.awsmt_source_location.read", "name", "test_source_location"),
				),
			},
		},
	})
}

func TestAccSourceLocationDataSourceErrors(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      sourceLocationDSError(),
				ExpectError: regexp.MustCompile("Error while describing source location"),
			},
		},
	})
}

func sourceLocationDS() string {
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
`
}

func sourceLocationDSError() string {
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
  							name = "testing_errors"
						}
						output "awsmt_source_location" {
  							value = data.awsmt_source_location.read
						}
`
}
