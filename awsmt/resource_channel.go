package awsmt

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"channel_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"channel_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filler_slate": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_location_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"vod_source_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"last_modified_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
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
						"manifest_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"playback_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_group": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"playback_mode": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"LINEAR", "LOOP"}, false),
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tier": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"BASIC", "STANDARD"}, false),
			},
		},
		CustomizeDiff: customdiff.Sequence(
			customdiff.ForceNewIfChange("channel_name", func(ctx context.Context, old, new, meta interface{}) bool { return old.(string) != new.(string) }),
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
	d.SetId(aws.StringValue(channel.Arn))

	return resourceChannelRead(ctx, d, meta)
}

func resourceChannelRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	resourceName := d.Get("channel_name").(string)
	if len(resourceName) == 0 && len(d.Id()) > 0 {
		resourceArn, err := arn.Parse(d.Id())
		if err != nil {
			return diag.FromErr(fmt.Errorf("error parsing the name from resource arn: %v", err))
		}
		arnSections := strings.Split(resourceArn.Resource, "/")
		resourceName = arnSections[len(arnSections)-1]
	}
	res, err := client.DescribeChannel(&mediatailor.DescribeChannelInput{ChannelName: aws.String(resourceName)})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while retrieving the channel: %v", err))
	}

	err = returnChannel(res, d)
	if err != nil {
		diag.FromErr(err)
	}

	return nil
}

func resourceChannelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	if d.HasChange("tags") {
		oldValue, newValue := d.GetChange("tags")
		var removedTags []string
		for k := range oldValue.(map[string]interface{}) {
			if _, ok := (newValue.(map[string]interface{}))[k]; !ok {
				removedTags = append(removedTags, k)
			}
		}
		resourceName := d.Get("channel_name").(string)
		res, err := client.DescribeChannel(&mediatailor.DescribeChannelInput{ChannelName: &resourceName})
		if err != nil {
			return diag.FromErr(err)
		}
		err = deleteTags(client, aws.StringValue(res.Arn), removedTags)
		if err != nil {
			return diag.FromErr(err)
		}
		if newValue != nil {
			var newTags = make(map[string]*string)
			for k, v := range newValue.(map[string]interface{}) {
				val := v.(string)
				newTags[k] = &val
			}
			tagInput := mediatailor.TagResourceInput{ResourceArn: res.Arn, Tags: newTags}
			_, err := client.TagResource(&tagInput)
			if err != nil {
				return diag.FromErr(err)
			}
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

func resourceChannelDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	_, err := client.DeleteChannel(&mediatailor.DeleteChannelInput{ChannelName: aws.String(d.Get("channel_name").(string))})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while deleting the resource: %v", err))
	}

	return nil
}

func getOutputs(d *schema.ResourceData) []*mediatailor.RequestOutputItem {
	if v, ok := d.GetOk("outputs"); ok && v.([]interface{})[0] != nil {
		outputs := v.([]interface{})

		var res []*mediatailor.RequestOutputItem

		for _, output := range outputs {
			current := output.(map[string]interface{})
			temp := mediatailor.RequestOutputItem{}

			if str, ok := current["manifest_name"]; ok {
				temp.ManifestName = aws.String(str.(string))
			}
			if str, ok := current["source_group"]; ok {
				temp.SourceGroup = aws.String(str.(string))
			}

			if num, ok := current["hls_manifest_windows_seconds"]; ok && num.(int) != 0 {
				tempHls := mediatailor.HlsPlaylistSettings{}
				tempHls.ManifestWindowSeconds = aws.Int64(int64(num.(int)))
				temp.HlsPlaylistSettings = &tempHls
			}

			tempDash := mediatailor.DashPlaylistSettings{}
			if num, ok := current["dash_manifest_windows_seconds"]; ok && num.(int) != 0 {
				tempDash.ManifestWindowSeconds = aws.Int64(int64(num.(int)))
			}
			if num, ok := current["dash_min_buffer_time_seconds"]; ok && num.(int) != 0 {
				tempDash.MinBufferTimeSeconds = aws.Int64(int64(num.(int)))
			}
			if num, ok := current["dash_min_update_period_seconds"]; ok && num.(int) != 0 {
				tempDash.MinUpdatePeriodSeconds = aws.Int64(int64(num.(int)))
			}
			if num, ok := current["dash_suggested_presentation_delay_seconds"]; ok && num.(int) != 0 {
				tempDash.SuggestedPresentationDelaySeconds = aws.Int64(int64(num.(int)))
			}
			if tempDash != (mediatailor.DashPlaylistSettings{}) {
				temp.DashPlaylistSettings = &tempDash
			}

			res = append(res, &temp)
		}
		return res
	}
	return nil
}
