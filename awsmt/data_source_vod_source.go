package awsmt

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVodSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVodSourceRead,
		Schema: map[string]*schema.Schema{
			"arn":           &computedString,
			"creation_time": &computedString,
			"http_package_configurations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path":         &computedString,
						"source_group": &computedString,
						"type":         &computedString,
					},
				},
			},
			"last_modified_time":   &computedString,
			"source_location_name": &requiredString,
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vod_source_name": &requiredString,
		},
	}
}

func dataSourceVodSourceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)
	resourceName := d.Get("vod_source_name").(string)
	sourceLocationName := d.Get("source_location_name").(string)

	input := &mediatailor.DescribeVodSourceInput{SourceLocationName: &(sourceLocationName), VodSourceName: aws.String(resourceName)}

	res, err := client.DescribeVodSource(input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while reading the vod source: %v", err))
	}

	if err = setVodSource(res, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
