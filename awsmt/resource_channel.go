package awsmt

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
	"strings"
)

func resourceChannel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChannelCreate,
		ReadContext:   resourceChannelRead,
		UpdateContext: resourceChannelUpdate,
		DeleteContext: resourceChannelDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"arn":           &computedString,
			"name":          &requiredString,
			"channel_state": &computedString,
			"creation_time": &computedString,
			"filler_slate": createOptionalList(map[string]*schema.Schema{
				"source_location_name": &optionalString,
				"vod_source_name":      &optionalString,
			}),
			"last_modified_time": &computedString,
			"outputs": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dash_manifest_windows_seconds": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(30, 3600),
						},
						"dash_min_buffer_time_seconds": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(2, 60),
						},
						"dash_min_update_period_seconds": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(2, 60),
						},
						"dash_suggested_presentation_delay_seconds": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(2, 60),
						},
						"hls_manifest_windows_seconds": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(30, 3600),
						},
						"manifest_name": &requiredString,
						"playback_url":  &computedString,
						"source_group":  &requiredString,
					},
				},
			},
			"playback_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"LINEAR", "LOOP"}, false),
			},
			// @ADR
			// In the context of developing a specific resource for channel policies,
			// facing concerns with the fact that such resource does not have an own ARN
			// we decided for incorporating the resource in the channel
			// and neglected continuing to develop a standalone resource
			// to be able to manage channel policies,
			// accepting the downsides it comes with: mainly the fact that the CRUD functions for the channel resource now
			// have to perform more than 1 API calls, increasing the chances of error, and the fact that the policy requires
			// the developer to specify the ARN for the channel it refers to, ignoring the fact that the arn is not known
			// while declaring the resource and needs to be built by the developer knowing the account ID and resource name,
			// because we did not want to create a different resource with no arn.
			"policy": {
				Type:     schema.TypeString,
				Optional: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					re := regexp.MustCompile(`\s?|\r?|\n?`)
					return re.ReplaceAllString(old, "") == re.ReplaceAllString(new, "")
				},
			},
			"tags": &optionalTags,
			"tier": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"BASIC", "STANDARD"}, false),
			},
		},
		CustomizeDiff: customdiff.Sequence(
			customdiff.ForceNewIfChange("name", func(ctx context.Context, old, new, meta interface{}) bool { return old.(string) != new.(string) }),
		),
	}
}

func resourceChannelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	var params = getCreateChannelInput(d)

	channel, err := client.CreateChannel(&params)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while creating the channel: %v", err))
	}

	if err := createChannelPolicy(client, d); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(aws.StringValue(channel.Arn))

	return resourceChannelRead(ctx, d, meta)
}

func resourceChannelRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)
	var resourceName *string
	resourceName, err := getResourceName(d, "name")
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := client.DescribeChannel(&mediatailor.DescribeChannelInput{ChannelName: resourceName})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while retrieving the channel: %v", err))
	}
	err = setChannel(res, d)
	if err != nil {
		diag.FromErr(err)
	}

	policy, err := client.GetChannelPolicy(&mediatailor.GetChannelPolicyInput{ChannelName: resourceName})
	if err != nil && !strings.Contains(err.Error(), "NotFound") {
		return diag.FromErr(fmt.Errorf("error while getting the channel policy: %v", err))
	}
	if err := setChannelPolicy(policy, d); err != nil {
		diag.FromErr(err)
	}

	return nil
}

func resourceChannelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	resourceName := d.Get("name").(string)

	if d.HasChange("tags") {
		oldValue, newValue := d.GetChange("tags")
		res, err := client.DescribeChannel(&mediatailor.DescribeChannelInput{ChannelName: &resourceName})
		if err != nil {
			return diag.FromErr(err)
		}
		if err := updateTags(client, res.Arn, oldValue, newValue); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("policy") {
		_, newValue := d.GetChange("policy")
		if len(newValue.(string)) > 0 {
			err := updateChannelPolicy(client, d, &resourceName)
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			err := deleteChannelPolicy(client, d, &resourceName)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	var params = getUpdateChannelInput(d)
	channel, err := client.UpdateChannel(&params)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while updating the channel: %v", err))
	}
	d.SetId(aws.StringValue(channel.Arn))

	return resourceChannelRead(ctx, d, meta)
}

func resourceChannelDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	_, err := client.DeleteChannelPolicy(&mediatailor.DeleteChannelPolicyInput{ChannelName: aws.String(d.Get("name").(string))})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while deleting the channel policy: %v", err))
	}

	_, err = client.DeleteChannel(&mediatailor.DeleteChannelInput{ChannelName: aws.String(d.Get("name").(string))})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while deleting the resource: %v", err))
	}

	return nil
}
