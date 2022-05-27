package awsmt

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strings"
)

func resourceVodSource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVodSourceCreate,
		ReadContext:   resourceVodSourceRead,
		UpdateContext: resourceVodSourceUpdate,
		DeleteContext: resourceVodSourceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"arn":           &computedString,
			"creation_time": &computedString,
			"http_package_configurations": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path":         &requiredString,
						"source_group": &requiredString,
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"DASH", "HLS"}, false),
						},
					},
				},
			},
			"last_modified_time":   &computedString,
			"source_location_name": &requiredString,
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"vod_source_name": &requiredString,
		},
		CustomizeDiff: customdiff.Sequence(
			customdiff.ForceNewIfChange("vod_source_name", func(ctx context.Context, old, new, meta interface{}) bool { return old.(string) != new.(string) }),
			customdiff.ForceNewIfChange("source_location_name", func(ctx context.Context, old, new, meta interface{}) bool { return old.(string) != new.(string) }),
		),
	}
}

func resourceVodSourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	params := getCreateVodSourceInput(d)
	vodSource, err := client.CreateVodSource(&params)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while creating the vod source: %v", err))
	}
	d.SetId(aws.StringValue(vodSource.Arn))

	return resourceVodSourceRead(ctx, d, meta)
}

func resourceVodSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)
	resourceName := d.Get("vod_source_name").(string)
	if len(resourceName) == 0 && len(d.Id()) > 0 {
		resourceArn, err := arn.Parse(d.Id())
		if err != nil {
			return diag.FromErr(fmt.Errorf("error parsing the name from resource arn: %v", err))
		}
		arnSections := strings.Split(resourceArn.Resource, "/")
		resourceName = arnSections[len(arnSections)-1]
	}

	res, err := client.DescribeVodSource(&mediatailor.DescribeVodSourceInput{VodSourceName: aws.String(resourceName), SourceLocationName: aws.String(d.Get("source_location_name").(string))})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while retrieving the vod location: %v", err))
	}

	if err = setVodSource(res, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceVodSourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	if d.HasChange("tags") {
		oldValue, newValue := d.GetChange("tags")

		resourceName := d.Get("vod_source_name").(string)
		sourceLocationName := d.Get("source_location_name").(string)
		res, err := client.DescribeVodSource(&mediatailor.DescribeVodSourceInput{SourceLocationName: &sourceLocationName, VodSourceName: &resourceName})
		if err != nil {
			return diag.FromErr(err)
		}

		if err := updateTags(client, res.Arn, oldValue, newValue); err != nil {
			return diag.FromErr(err)
		}
	}

	var params = getUpdateVodSourceInput(d)
	vodSource, err := client.UpdateVodSource(&params)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while updating the vod source: %v", err))
	}
	d.SetId(aws.StringValue(vodSource.Arn))

	return resourceVodSourceRead(ctx, d, meta)
}

func resourceVodSourceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	_, err := client.DeleteVodSource(&mediatailor.DeleteVodSourceInput{VodSourceName: aws.String(d.Get("vod_source_name").(string)), SourceLocationName: aws.String(d.Get("source_location_name").(string))})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while deleting the resource: %v", err))
	}

	return nil
}
