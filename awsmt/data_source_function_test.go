package awsmt

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccFunctionDataSourceBasic(t *testing.T) {
	resourceName := "data.awsmt_function.test"
	functionId := "test_acc_ds_function"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: functionDataSourceConfig(functionId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "function_id", functionId),
					resource.TestCheckResourceAttr(resourceName, "function_type", "CUSTOM_OUTPUT"),
					resource.TestCheckResourceAttr(resourceName, "description", "test data source function"),
					resource.TestCheckResourceAttr(resourceName, "custom_output_configuration.runtime", "JSONATA"),
					resource.TestCheckResourceAttr(resourceName, "custom_output_configuration.output.player_params.device_type", "mobile"),
					resource.TestMatchResourceAttr(resourceName, "arn", regexp.MustCompile(`^arn:aws:mediatailor:.*`)),
				),
			},
		},
	})
}

func TestAccFunctionDataSourceErrors(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      functionDataSourceErrorConfig(),
				ExpectError: regexp.MustCompile("Error reading function"),
			},
		},
	})
}

func functionDataSourceConfig(functionId string) string {
	return fmt.Sprintf(`
resource "awsmt_function" "test" {
  function_id   = "%[1]s"
  function_type = "CUSTOM_OUTPUT"
  description   = "test data source function"
  custom_output_configuration = {
    runtime = "JSONATA"
    output = {
      "player_params.device_type" = "mobile"
    }
  }
}

data "awsmt_function" "test" {
  function_id = awsmt_function.test.function_id
}
`, functionId)
}

func functionDataSourceErrorConfig() string {
	return `
data "awsmt_function" "test" {
  function_id = "nonexistent_function_id"
}
`
}
