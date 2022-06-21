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

func resourceLiveSource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLiveSourceCreate,
		ReadContext:   resourceLiveSourceRead,
		UpdateContext: resourceLiveSourceUpdate,
		DeleteContext: resourceLiveSourceDelete,
		Schema: map[string]*schema.Schema{
			"arn":           &computedString,
			"creation_time": &computedString,
			"http_package_configurations": createRequiredList(map[string]*schema.Schema{
				"path":         &requiredString,
				"source_group": &requiredString,
				"type": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"DASH", "HLS"}, false),
				},
			}),
			"last_modified_time":   &computedString,
			"name":                 &requiredString,
			"source_location_name": &requiredString,
			"tags":                 &optionalTags,
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: customdiff.Sequence(
			customdiff.ForceNewIfChange("name", func(ctx context.Context, old, new, meta interface{}) bool { return old.(string) != new.(string) }),
			customdiff.ForceNewIfChange("source_location_name", func(ctx context.Context, old, new, meta interface{}) bool { return old.(string) != new.(string) }),
		),
	}
}

func resourceLiveSourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	params := getCreateLiveSourceInput(d)
	liveSource, err := client.CreateLiveSource(&params)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while creating the live source: %v", err))
	}
	d.SetId(aws.StringValue(liveSource.Arn))

	return resourceLiveSourceRead(ctx, d, meta)
}

func resourceLiveSourceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)
	liveSourceName := d.Get("name").(string)
	sourceLocationName := d.Get("source_location_name").(string)

	if len(liveSourceName) == 0 && len(d.Id()) > 0 {
		resourceArn, err := arn.Parse(d.Id())
		if err != nil {
			return diag.FromErr(fmt.Errorf("error parsing the name from resource arn: %v", err))
		}
		arnSections := strings.Split(resourceArn.Resource, "/")
		liveSourceName = arnSections[len(arnSections)-1]
		sourceLocationName = arnSections[len(arnSections)-2]
	}

	input := &mediatailor.DescribeLiveSourceInput{SourceLocationName: &(sourceLocationName), LiveSourceName: aws.String(liveSourceName)}

	res, err := client.DescribeLiveSource(input)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while reading the live source: %v", err))
	}

	if err = setLiveSource(res, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceLiveSourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	if d.HasChange("tags") {
		oldValue, newValue := d.GetChange("tags")

		resourceName := d.Get("name").(string)
		sourceLocationName := d.Get("source_location_name").(string)
		res, err := client.DescribeLiveSource(&mediatailor.DescribeLiveSourceInput{SourceLocationName: &sourceLocationName, LiveSourceName: &resourceName})
		if err != nil {
			return diag.FromErr(err)
		}

		if err := updateTags(client, res.Arn, oldValue, newValue); err != nil {
			return diag.FromErr(err)
		}
	}

	var params = getUpdateLiveSourceInput(d)
	liveSource, err := client.UpdateLiveSource(&params)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while updating the live source: %v", err))
	}
	d.SetId(aws.StringValue(liveSource.Arn))

	return resourceLiveSourceRead(ctx, d, meta)
}

func resourceLiveSourceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	_, err := client.DeleteLiveSource(&mediatailor.DeleteLiveSourceInput{LiveSourceName: aws.String(d.Get("name").(string)), SourceLocationName: aws.String(d.Get("source_location_name").(string))})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while deleting the resource: %v", err))
	}

	return nil
}
