package awsmt

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"
)

//var testAccProviders map[string]*schema.Provider
//var testAccProvider *schema.Provider

var ProviderFactories map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func sharedClientForRegion(region string) (interface{}, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return nil, errors.New("unable to initialize provider in the specified region")
	}
	c := mediatailor.New(sess)
	return c, nil
}

func init() {
	testAccProvider = Provider()
	ProviderFactories = map[string]func() (*schema.Provider, error){
		"awsmt": func() (*schema.Provider, error) { return testAccProvider, nil }, //nolint:unparam
	}
}
