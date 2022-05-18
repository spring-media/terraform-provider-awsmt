package awsmt

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"awsmt_playback_configuration": resourcePlaybackConfiguration(),
			"awsmt_channel":                ResourceChannel(),
		},
		DataSourcesMap:       map[string]*schema.Resource{"awsmt_playback_configuration": dataSourcePlaybackConfiguration()},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	region, ok := d.GetOk("region")
	if !ok {
		region = "eu-central-1"
	}
	sess, err := session.NewSession(&aws.Config{Region: aws.String(region.(string))})
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Unable to initialize session in region '%s'", region),
			Detail:   fmt.Sprintf("Unable to create a new session for the specified region: %s", err),
		})
		return nil, diags
	}
	c := mediatailor.New(sess)
	return c, diags
}
