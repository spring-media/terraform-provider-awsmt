package mediatailor

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccConfigurationDataSourceBasic(t *testing.T) {
	//dataSourceName := "data.mediatailor_configuration.c1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigurationDataSourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mediatailor_configuration.c1", "name", "staging-live-stream"),
					resource.TestCheckResourceAttr("data.mediatailor_configuration.c2", "max_results", "2"),
				),
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if a, b := os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"); a == "" || b == "" {
		t.Fatal("AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY must both be set for acceptance tests")
	}
}

func testAccConfigurationDataSourceBasic() string {
	return `
data "mediatailor_configuration" "c1" {
  name = "staging-live-stream"
}

data "mediatailor_configuration" "c2" {
  max_results = 2
}

output "out1" {
  value = data.mediatailor_configuration.c1
}

output "out2" {
  value = data.mediatailor_configuration.c2
}
`
}
