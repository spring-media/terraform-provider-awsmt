package awsmt

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
						"smatc_header": {
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
			"segment_delivery_configuration": {
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
				Optional: true,
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
	return nil
}

func resourceSourceLocationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}

func resourceSourceLocationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceSourceLocationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
