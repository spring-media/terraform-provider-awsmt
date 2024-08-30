package awsmt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccSourceLocationResource(t *testing.T) {
	resourceName := "awsmt_source_location.test_source_location"
	name := "test_source_location"
	baseUrl := "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com"
	baseUrl2 := "https://example.com/"
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
				Config: basicSourceLocation(name, baseUrl, k1, v1, k2, v2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "test_source_location"),
					resource.TestMatchResourceAttr(resourceName, "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:sourceLocation\/.*$`)),
					resource.TestMatchResourceAttr(resourceName, "creation_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr(resourceName, "default_segment_delivery_configuration.base_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com"),
					resource.TestCheckResourceAttr(resourceName, "http_configuration.base_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com"),
					resource.TestMatchResourceAttr(resourceName, "last_modified_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr(resourceName, "segment_delivery_configurations.0.base_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com"),
					resource.TestCheckResourceAttr(resourceName, "name", "test_source_location"),
					resource.TestCheckResourceAttr(resourceName, "tags.Testing", "pass"),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "dev"),
				),
			},
			// ImportState testing
			{
				ResourceName: resourceName,
				ImportState:  true,
			},
			// Update and Read testing
			{
				Config: basicSourceLocationWithAccessConfig(name, baseUrl2, k3, v3, k2, v2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "test_source_location"),
					resource.TestCheckResourceAttr(resourceName, "access_configuration.access_type", "S3_SIGV4"),
					resource.TestCheckResourceAttr(resourceName, "name", "test_source_location"),
					resource.TestCheckResourceAttr(resourceName, "http_configuration.base_url", "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com"),
					resource.TestCheckResourceAttr(resourceName, "default_segment_delivery_configuration.base_url", "https://example.com/"),
					resource.TestCheckResourceAttr(resourceName, "segment_delivery_configurations.0.base_url", "https://example.com/"),
					resource.TestCheckResourceAttr(resourceName, "tags.Testing", "pass"),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "prod")),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
func TestAccSourceLocationDeleteVodResource(t *testing.T) {
	resourceName := "awsmt_source_location.test_source_location"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: basicSourceLocationWithVodSource(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", "test_source_location"),
					resource.TestCheckResourceAttr(resourceName, "name", "test_source_location"),
				),
			},
		},
	})
}

func basicSourceLocation(name, baseUrl, k1, v1, k2, v2 string) string {
	return fmt.Sprintf(`
		resource "awsmt_source_location" "test_source_location"{
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
		`, name, baseUrl, k1, v1, k2, v2)

}

func basicSourceLocationWithAccessConfig(name, baseUrl, k1, v1, k2, v2 string) string {
	return fmt.Sprintf(`
		resource "awsmt_source_location" "test_source_location"{
			name = "%[1]s"
			http_configuration = {
				baseUrl = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com"
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
		`, name, baseUrl, k1, v1, k2, v2)
}

func basicSourceLocationWithVodSource() string {
	return `
		resource "awsmt_source_location" "test_source_location"{
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
