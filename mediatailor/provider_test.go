package mediatailor

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

//var testAccProviders map[string]*schema.Provider
//var testAccProvider *schema.Provider

var ProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	//testAccProvider = Provider()
	//testAccProviders = map[string]*schema.Provider{
	//	"mediatailor": testAccProvider,
	//}
	ProviderFactories = map[string]func() (*schema.Provider, error){
		"mediatailor": func() (*schema.Provider, error) { return Provider(), nil }, //nolint:unparam
	}
}
