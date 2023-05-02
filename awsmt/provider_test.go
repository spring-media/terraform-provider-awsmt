package awsmt

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func sharedClientForRegion(region string) (interface{}, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return nil, errors.New("unable to initialize provider in the specified region")
	}
	c := mediatailor.New(sess)
	return c, nil
}

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"awsmt": providerserver.NewProtocol6WithError(New()),
	}
)

/* func TestMain(m *testing.M) {
	resource.TestMain(m)
} */
