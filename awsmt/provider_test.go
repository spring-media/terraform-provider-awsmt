package awsmt

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

//var testAccProviders map[string]*schema.Provider
//var testAccProvider *schema.Provider

var ProviderFactories map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider

func init() {
	//testAccProvider = Provider()
	//testAccProviders = map[string]*schema.Provider{
	//	"awsmt": testAccProvider,
	//}
	testAccProvider = Provider()
	ProviderFactories = map[string]func() (*schema.Provider, error){
		"awsmt": func() (*schema.Provider, error) { return testAccProvider, nil }, //nolint:unparam
	}
}
