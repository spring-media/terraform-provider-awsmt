package awsmt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccLiveSourceResourceBasic(t *testing.T) {
	terraformResourceName := "awsmt_live_source.live_source_acc_test"
	name := "live_source_example"
	path := "/"
	path2 := "/test"
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
			{
				Config: basicLiveSourceWithSourceLocation(name, path, k1, v1, k2, v2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(terraformResourceName, "id", "test_source_location,live_source_example"),
					resource.TestMatchResourceAttr(terraformResourceName, "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:liveSource\/.*$`)),
					resource.TestMatchResourceAttr(terraformResourceName, "creation_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr(terraformResourceName, "http_package_configurations.0.path", "/"),
					resource.TestCheckResourceAttr(terraformResourceName, "http_package_configurations.0.source_group", "default"),
					resource.TestCheckResourceAttr(terraformResourceName, "http_package_configurations.0.type", "HLS"),
					resource.TestMatchResourceAttr(terraformResourceName, "last_modified_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr(terraformResourceName, "name", "live_source_example"),
					resource.TestCheckResourceAttr(terraformResourceName, "source_location_name", "test_source_location"),
					resource.TestCheckResourceAttr(terraformResourceName, "tags.Environment", "dev"),
					resource.TestCheckResourceAttr(terraformResourceName, "tags.Testing", "pass"),
				),
			},
			{
				Config: basicLiveSourceWithSourceLocation(name, path2, k3, v3, k2, v2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(terraformResourceName, "id", "test_source_location,live_source_example"),
					resource.TestMatchResourceAttr(terraformResourceName, "arn", regexp.MustCompile(`^arn:aws:mediatailor:[\w-]+:\d+:liveSource\/.*$`)),
					resource.TestMatchResourceAttr(terraformResourceName, "creation_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr(terraformResourceName, "http_package_configurations.0.path", "/test"),
					resource.TestCheckResourceAttr(terraformResourceName, "http_package_configurations.0.source_group", "default"),
					resource.TestCheckResourceAttr(terraformResourceName, "http_package_configurations.0.type", "HLS"),
					resource.TestMatchResourceAttr(terraformResourceName, "last_modified_time", regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}(\.\d{1,3})? \+\d{4} \w+$`)),
					resource.TestCheckResourceAttr(terraformResourceName, "name", "live_source_example"),
					resource.TestCheckResourceAttr(terraformResourceName, "source_location_name", "test_source_location"),
					resource.TestCheckResourceAttr(terraformResourceName, "tags.Environment", "prod"),
					resource.TestCheckResourceAttr(terraformResourceName, "tags.Testing", "pass"),
				),
			},
		},
	})
}

func TestAccLiveSourceResourceCreationFailure(t *testing.T) {
	name := "live_source_failing"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      liveSourceWithoutSourceLocation(name),
				ExpectError: regexp.MustCompile(`.*The specified source-location doesn't exist.*`),
			},
		},
	})
}

func basicLiveSourceWithSourceLocation(name, path, k1, v1, k2, v2 string) string {
	return fmt.Sprintf(`
		resource "awsmt_live_source" "live_source_acc_test" {
			http_package_configurations = [{
				path = "%[2]s"
				source_group = "default"
				type = "HLS"
			}]
			source_location_name = awsmt_source_location.test_source_location.name
			name = "%[1]s"
			tags = {
				"%[3]s": "%[4]s",
				"%[5]s": "%[6]s"
			}
		}
		data "awsmt_live_source" "data_test" {
			source_location_name = awsmt_source_location.test_source_location.name
			name = awsmt_live_source.live_source_acc_test.name
		}

		output "live_source_out" {
			value = data.awsmt_live_source.data_test
		}
		resource "awsmt_source_location" "test_source_location" {
			name = "test_source_location"
			http_configuration = {
				base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/"
			}
			default_segment_delivery_configuration = {
				base_url = "https://ott-mediatailor-test.s3.eu-central-1.amazonaws.com/test-img.jpeg"
			}
		}
		`, name, path, k1, v1, k2, v2)
}

func liveSourceWithoutSourceLocation(name string) string {
	return fmt.Sprintf(`
		resource "awsmt_live_source" "live_source_acc_test" {
			http_package_configurations = [{
				path = "/"
				source_group = "default"
				type = "HLS"
			}]
			source_location_name = "unexisting_source_location"
			name = "%[1]s"
		}`, name)
}
