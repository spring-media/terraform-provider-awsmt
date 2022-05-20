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
			"access_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// SMATC is short for Secret Manager Access Token Configuration
						"smatc_header_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"smatc_secret_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"smatc_secret_string_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default_segment_delivery_configuration_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"http_configuration_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"segment_delivery_configurations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"base_url": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"source_location_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceSourceLocationRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mediatailor.MediaTailor)

	name := d.Get("source_location_name").(string)
	if name == "" {
		return diag.Errorf("`source_location_name` parameter required")
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
