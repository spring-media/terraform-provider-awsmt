package awsmt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccFunctionCustomOutput(t *testing.T) {
	resourceName := "awsmt_function.test"
	functionId := "test_acc_function_custom_output"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: customOutputFunctionConfig(functionId, "first description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "function_id", functionId),
					resource.TestCheckResourceAttr(resourceName, "function_type", "CUSTOM_OUTPUT"),
					resource.TestCheckResourceAttr(resourceName, "description", "first description"),
					resource.TestCheckResourceAttr(resourceName, "custom_output_configuration.runtime", "JSONATA"),
					resource.TestCheckResourceAttr(resourceName, "custom_output_configuration.output.player_params.device_type", "mobile"),
				),
			},
			{
				Config: customOutputFunctionConfig(functionId, "updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "function_id", functionId),
					resource.TestCheckResourceAttr(resourceName, "description", "updated description"),
				),
			},
			{
				ResourceName: resourceName,
				ImportState:  true,
			},
		},
	})
}

func TestAccFunctionHttpRequest(t *testing.T) {
	resourceName := "awsmt_function.test"
	functionId := "test_acc_function_http_request"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: httpRequestFunctionConfig(functionId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "function_id", functionId),
					resource.TestCheckResourceAttr(resourceName, "function_type", "HTTP_REQUEST"),
					resource.TestCheckResourceAttr(resourceName, "http_request_configuration.method_type", "GET"),
					resource.TestCheckResourceAttr(resourceName, "http_request_configuration.request_timeout_milliseconds", "1000"),
					resource.TestCheckResourceAttr(resourceName, "http_request_configuration.runtime", "JSONATA"),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
				),
			},
		},
	})
}

func customOutputFunctionConfig(functionId, description string) string {
	return fmt.Sprintf(`
		resource "awsmt_function" "test" {
			function_id   = "%[1]s"
			function_type = "CUSTOM_OUTPUT"
			description   = "%[2]s"
			custom_output_configuration = {
				runtime = "JSONATA"
				output = {
					"player_params.device_type" = "mobile"
				}
			}
		}
	`, functionId, description)
}

func httpRequestFunctionConfig(functionId string) string {
	return fmt.Sprintf(`
		resource "awsmt_function" "test" {
			function_id   = "%[1]s"
			function_type = "HTTP_REQUEST"
			description   = "HTTP request function"
			http_request_configuration = {
				method_type                  = "GET"
				request_timeout_milliseconds = 1000
				runtime                      = "JSONATA"
				url                          = "https://example.com/api"
				output = {
					"player_params.result" = "response.body.value"
				}
			}
		}
	`, functionId)
}
