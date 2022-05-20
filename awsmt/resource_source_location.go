package awsmt

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
)

func resourceSourceLocation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSourceLocationCreate,
		ReadContext:   resourceSourceLocationRead,
		UpdateContext: resourceSourceLocationUpdate,
		DeleteContext: resourceSourceLocationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"access_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_type": {
							Type:     schema.TypeString,
							Optional: true,
						},
						// SMATC is short for Secret Manager Access Token Configuration
						"smatc_header_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"smatc_secret_arn": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"smatc_secret_string_key": {
							Type:     schema.TypeString,
							Optional: true,
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
				Optional: true,
			},
			"http_configuration_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"segment_delivery_configurations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"base_url": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeInt,
							Optional: true,
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
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSourceLocationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	var params = getCreateSourceLocationInput(d)

	sourceLocation, err := client.CreateSourceLocation(&params)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while creating the channel: %v", err))
	}
	d.SetId(aws.StringValue(sourceLocation.Arn))

	return resourceSourceLocationRead(ctx, d, meta)
}

func resourceSourceLocationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	resourceName := d.Get("source_location_name").(string)
	if len(resourceName) == 0 && len(d.Id()) > 0 {
		resourceArn, err := arn.Parse(d.Id())
		if err != nil {
			return diag.FromErr(fmt.Errorf("error parsing the name from resource arn: %v", err))
		}
		arnSections := strings.Split(resourceArn.Resource, "/")
		resourceName = arnSections[len(arnSections)-1]
	}
	res, err := client.DescribeSourceLocation(&mediatailor.DescribeSourceLocationInput{SourceLocationName: aws.String(resourceName)})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while retrieving the source location: %v", err))
	}

	if err = setResourceLocation(res, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSourceLocationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	if d.HasChange("tags") {
		oldValue, newValue := d.GetChange("tags")

		resourceName := d.Get("source_location_name").(string)
		res, err := client.DescribeSourceLocation(&mediatailor.DescribeSourceLocationInput{SourceLocationName: &resourceName})
		if err != nil {
			return diag.FromErr(err)
		}

		if err := updateTags(client, res.Arn, oldValue, newValue); err != nil {
			return diag.FromErr(err)
		}
	}

	var params = getUpdateChannelInput(d)
	channel, err := client.UpdateChannel(&params)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while creating the channel: %v", err))
	}
	d.SetId(aws.StringValue(channel.Arn))

	return resourceChannelRead(ctx, d, meta)
}

func resourceSourceLocationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
