package awsmt

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"awsmt": providerserver.NewProtocol6WithError(New()),
	}
)

/* func TestMain(m *testing.M) {
	resource.TestMain(m)
} */
