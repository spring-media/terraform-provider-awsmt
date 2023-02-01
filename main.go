package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"terraform-provider-mediatailor/awsmt"
)

func main() {
	err := providerserver.Serve(context.Background(), awsmt.New, providerserver.ServeOpts{

		Address: "registry.terraform.io/spring-media/awsmt",
	})
	if err != nil {
		return
	}
}
