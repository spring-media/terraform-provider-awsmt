package awsmt

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccVodSourceResourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "awsmt_vod_source" "test" {
  					http_package_configurations = [{
						path = "/"
						source_group = "default"
    					type = "HLS"
  					}]
  					source_location_name = awsmt_source_location.test_source_location.source_location_name
  					vod_source_name = "vod_source_example"
					tags = {"Environment": "dev"}
				}

				data "awsmt_vod_source" "data_test" {
  					source_location_name = awsmt_source_location.test_source_location.source_location_name
  					vod_source_name = awsmt_vod_source.test.vod_source_name
				}

				output "vod_source_out" {
  					value = data.awsmt_vod_source.data_test
				}
				resource "awsmt_source_location" "test_source_location"{
  					source_location_name = "test_source_location"
  					http_configuration = {
    					hc_base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/"
  					}
  					default_segment_delivery_configuration = {
    					dsdc_base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"
  					}
				}
				data "awsmt_source_location" "test" {
  					source_location_name = awsmt_source_location.test_source_location.source_location_name
				}
				output "awsmt_source_location" {
  					value = data.awsmt_source_location.test
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_vod_source.test", "id", "test_source_location,vod_source_example"),
					resource.TestMatchResourceAttr("awsmt_vod_source.test", "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:vodSource\/.*$`)),
					resource.TestMatchResourceAttr("awsmt_vod_source.test", "creation_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr("awsmt_vod_source.test", "http_package_configurations.0.path", "/"),
					resource.TestCheckResourceAttr("awsmt_vod_source.test", "http_package_configurations.0.source_group", "default"),
					resource.TestCheckResourceAttr("awsmt_vod_source.test", "http_package_configurations.0.type", "HLS"),
					resource.TestMatchResourceAttr("awsmt_vod_source.test", "last_modified_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr("awsmt_vod_source.test", "vod_source_name", "vod_source_example"),
					resource.TestCheckResourceAttr("awsmt_vod_source.test", "source_location_name", "test_source_location"),
					resource.TestCheckResourceAttr("awsmt_vod_source.test", "tags.Environment", "dev"),
				),
			},
			// ImportState testing
			/*{
				ResourceName:      "awsmt_vod_source.test",
				ImportState:       true,
			},*/
			{
				Config: `
				resource "awsmt_vod_source" "test" {
  					http_package_configurations = [{
						path = "/test"
						source_group = "default"
    					type = "HLS"
  					}]
  					source_location_name = awsmt_source_location.test_source_location.source_location_name
  					vod_source_name = "vod_source_example"
					tags = {"Environment": "dev", "Testing": "pass"}
				}

				data "awsmt_vod_source" "data_test" {
  					source_location_name = awsmt_source_location.test_source_location.source_location_name
  					vod_source_name = awsmt_vod_source.test.vod_source_name
				}

				output "vod_source_out" {
  					value = data.awsmt_vod_source.data_test
				}
				resource "awsmt_source_location" "test_source_location"{
  					source_location_name = "test_source_location"
  					http_configuration = {
    					hc_base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/"
  					}
  					default_segment_delivery_configuration = {
    					dsdc_base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"
  					}
				}
				data "awsmt_source_location" "test" {
  					source_location_name = awsmt_source_location.test_source_location.source_location_name
				}
				output "awsmt_source_location" {
  					value = data.awsmt_source_location.test
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("awsmt_vod_source.test", "id", "test_source_location,vod_source_example"),
					resource.TestMatchResourceAttr("awsmt_vod_source.test", "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:vodSource\/.*$`)),
					resource.TestMatchResourceAttr("awsmt_vod_source.test", "creation_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr("awsmt_vod_source.test", "http_package_configurations.0.path", "/test"),
					resource.TestCheckResourceAttr("awsmt_vod_source.test", "http_package_configurations.0.source_group", "default"),
					resource.TestCheckResourceAttr("awsmt_vod_source.test", "http_package_configurations.0.type", "HLS"),
					resource.TestMatchResourceAttr("awsmt_vod_source.test", "last_modified_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr("awsmt_vod_source.test", "vod_source_name", "vod_source_example"),
					resource.TestCheckResourceAttr("awsmt_vod_source.test", "source_location_name", "test_source_location"),
					resource.TestCheckResourceAttr("awsmt_vod_source.test", "tags.Environment", "dev"),
					resource.TestCheckResourceAttr("awsmt_vod_source.test", "tags.Testing", "pass"),
				),
			},
		},
	})
}
