package mediatailor

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"os"
	"testing"
)

func TestAccConfigurationDataSource_basic(t *testing.T) {
	//dataSourceName := "data.mediatailor_configuration.c1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigurationDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.mediatailor_configuration.c1", "name", "staging-live-stream"),
				),
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("AWS_ACCESS_KEY_ID"); v == "" {
		t.Fatal("AWS_ACCESS_KEY_ID must be set for acceptance tests")
	}
	if v := os.Getenv("AWS_SECRET_ACCESS_KEY"); v == "" {
		t.Fatal("AWS_SECRET_ACCESS_KEY must be set for acceptance tests")
	}
}

func testAccConfigurationDataSource_basic() string {
	return `
data "mediatailor_configuration" "c1" {
  name = "staging-live-stream"
}

output "out" {
  value = data.mediatailor_configuration.c1
}
`
}
