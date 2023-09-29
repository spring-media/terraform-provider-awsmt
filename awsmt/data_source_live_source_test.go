package awsmt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccLiveSourceDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: liveSourceDS(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.awsmt_live_source.data_test", "id", "test_source_location,live_source_example"),
					resource.TestMatchResourceAttr("data.awsmt_live_source.data_test", "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:liveSource\/.*$`)),
					resource.TestMatchResourceAttr("data.awsmt_live_source.data_test", "creation_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr("data.awsmt_live_source.data_test", "http_package_configurations.0.path", "/"),
					resource.TestCheckResourceAttr("data.awsmt_live_source.data_test", "http_package_configurations.0.source_group", "default"),
					resource.TestCheckResourceAttr("data.awsmt_live_source.data_test", "http_package_configurations.0.type", "HLS"),
					resource.TestMatchResourceAttr("data.awsmt_live_source.data_test", "last_modified_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr("data.awsmt_live_source.data_test", "name", "live_source_example"),
					resource.TestCheckResourceAttr("data.awsmt_live_source.data_test", "source_location_name", "test_source_location"),
					resource.TestCheckResourceAttr("data.awsmt_live_source.data_test", "tags.Environment", "dev"),
				),
			},
		},
	})
}

func TestAccLiveSourceDataSourceErrors(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      liveSourceError(),
				ExpectError: regexp.MustCompile("Error while describing live source"),
			},
		},
	})
}

func liveSourceDS() string {
	return `
				resource "awsmt_live_source" "test" {
  					http_package_configurations = [{
    					path = "/"
    					source_group = "default"
    					type = "HLS"
  					}]
  					source_location_name = awsmt_source_location.test_source_location.name
  					name = "live_source_example"
					tags = {"Environment": "dev"}
				}

				data "awsmt_live_source" "data_test" {
  					source_location_name = awsmt_source_location.test_source_location.name
  					name = awsmt_live_source.test.name
				}

				output "live_source_out" {
  					value = data.awsmt_live_source.data_test
				}
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
				`
}

func liveSourceError() string {
	return `
				resource "awsmt_live_source" "test" {
  					http_package_configurations = [{
    					path = "/"
    					source_group = "default"
    					type = "HLS"
  					}]
  					source_location_name = awsmt_source_location.test_source_location.name
  					name = "live_source_example"
					tags = {"Environment": "dev"}
				}

				data "awsmt_live_source" "data_test" {
  					source_location_name = awsmt_source_location.test_source_location.name
  					name = "testingError"
				}

				output "live_source_out" {
  					value = data.awsmt_live_source.data_test
				}
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
				`
}
