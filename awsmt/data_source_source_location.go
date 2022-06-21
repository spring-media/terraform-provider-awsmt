package awsmt

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSourceLocation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSourceLocationRead,
		Schema: map[string]*schema.Schema{
			"arn":           &computedString,
			"creation_time": &computedString,
			"default_segment_delivery_configuration_url": &computedString,
			"http_configuration_url":                     &computedString,
			"last_modified_time":                         &computedString,
			"segment_delivery_configurations": createComputedList(map[string]*schema.Schema{
				"base_url": &computedString,
				"name":     &computedString,
			}),
			"name": &requiredString,
			"tags": &computedTags,
		},
	}
}

func dataSourceSourceLocationRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mediatailor.MediaTailor)

	name := d.Get("name").(string)
	if name == "" {
		return diag.Errorf("`name` parameter required")
	}
	res, err := client.DescribeSourceLocation(&mediatailor.DescribeSourceLocationInput{SourceLocationName: &name})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while retrieving the source location: %w", err))
	}

	d.SetId(aws.StringValue(res.SourceLocationName))

	err = setSourceLocation(res, d)
	if err != nil {
		diag.FromErr(err)
	}
	return nil
}
