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
		// @ADR
		// Context: Source locations have an option to set up a Secrets Manager Access Token Configuration to authenticate
		// to the S3 bucket that contains the media files. It requires an AWS MKN Key and a Secrets Manager Secret but
		// we didn't manage to properly set them up.
		// Decision: Since we cannot test the SMATC without the key and the secret, we did not include it in the resource.
		// Consequences: We cannot use S£ buckets that require authentication.
		Schema: map[string]*schema.Schema{
			"arn":           &computedString,
			"creation_time": &computedString,
			"default_segment_delivery_configuration_url": &optionalString,
			"http_configuration_url":                     &requiredString,
			"last_modified_time":                         &computedString,
			"segment_delivery_configurations": createOptionalList(
				map[string]*schema.Schema{
					"base_url": &optionalString,
					"name":     &optionalString,
				},
			),
			"name": &requiredString,
			"tags": &optionalTags,
		},
		CustomizeDiff: customdiff.Sequence(
			customdiff.ForceNewIfChange("name", func(ctx context.Context, old, new, meta interface{}) bool { return old.(string) != new.(string) }),
		),
	}
}

func resourceSourceLocationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	var params = getCreateSourceLocationInput(d)

	sourceLocation, err := client.CreateSourceLocation(&params)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while creating the source location: %v", err))
	}
	d.SetId(aws.StringValue(sourceLocation.Arn))

	return resourceSourceLocationRead(ctx, d, meta)
}

func resourceSourceLocationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	resourceName := d.Get("name").(string)
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

	if err = setSourceLocation(res, d); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceSourceLocationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	if d.HasChange("tags") {
		oldValue, newValue := d.GetChange("tags")

		resourceName := d.Get("name").(string)
		res, err := client.DescribeSourceLocation(&mediatailor.DescribeSourceLocationInput{SourceLocationName: &resourceName})
		if err != nil {
			return diag.FromErr(err)
		}

		if err := updateTags(client, res.Arn, oldValue, newValue); err != nil {
			return diag.FromErr(err)
		}
	}

	var params = getUpdateSourceLocationInput(d)
	sourceLocation, err := client.UpdateSourceLocation(&params)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while updating the source location: %v", err))
	}
	d.SetId(aws.StringValue(sourceLocation.Arn))

	return resourceSourceLocationRead(ctx, d, meta)
}

func resourceSourceLocationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)
	sourceLocationName := aws.String(d.Get("name").(string))
	vodSourcesList, err := client.ListVodSources(&mediatailor.ListVodSourcesInput{SourceLocationName: sourceLocationName})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while retrieving Vod Sources list: %v", err))
	}
	for _, vodSource := range vodSourcesList.Items {
		_, err = client.DeleteVodSource(&mediatailor.DeleteVodSourceInput{VodSourceName: vodSource.VodSourceName, SourceLocationName: sourceLocationName})
		if err != nil {
			return diag.FromErr(fmt.Errorf("error while deleting the resource: %v", err))
		}
	}

	liveSourcesList, err := client.ListLiveSources(&mediatailor.ListLiveSourcesInput{SourceLocationName: sourceLocationName})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while retrieving Live Sources list: %v", err))
	}
	for _, liveSource := range liveSourcesList.Items {
		_, err := client.DeleteLiveSource(&mediatailor.DeleteLiveSourceInput{LiveSourceName: liveSource.LiveSourceName, SourceLocationName: sourceLocationName})
		if err != nil {
			return diag.FromErr(fmt.Errorf("error while deleting the resource: %v", err))
		}
	}

	_, err = client.DeleteSourceLocation(&mediatailor.DeleteSourceLocationInput{SourceLocationName: sourceLocationName})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while deleting the resource: %v", err))
	}

	return nil
}
