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

func ResourceChannel() *schema.Resource {
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

	d.Set("arn", res.Arn)
	d.Set("channel_name", res.ChannelName)
	d.Set("channel_state", res.ChannelState)
	d.Set("creation_time", res.CreationTime.String())
	if res.FillerSlate != nil && res.FillerSlate != &(mediatailor.SlateSource{}) {
		temp := map[string]interface{}{}
		if res.FillerSlate.SourceLocationName != nil {
			temp["source_location_name"] = res.FillerSlate.SourceLocationName
		}
		if res.FillerSlate.VodSourceName != nil {
			temp["vod_source_name"] = res.FillerSlate.VodSourceName
		}
		d.Set("filler_slate", []interface{}{temp})
	}
	d.Set("last_modified_time", res.LastModifiedTime.String())

	var outputs []map[string]interface{}
	for _, o := range res.Outputs {
		temp := map[string]interface{}{}
		temp["manifest_name"] = o.ManifestName
		temp["playback_url"] = o.PlaybackUrl
		temp["source_group"] = o.SourceGroup

		if o.HlsPlaylistSettings != nil && o.HlsPlaylistSettings.ManifestWindowSeconds != nil {
			temp["hls_manifest_windows_seconds"] = o.HlsPlaylistSettings.ManifestWindowSeconds
		}

		if o.DashPlaylistSettings != nil && o.DashPlaylistSettings != &(mediatailor.DashPlaylistSettings{}) {
			if o.DashPlaylistSettings.ManifestWindowSeconds != nil {
				temp["dash_manifest_windows_seconds"] = o.DashPlaylistSettings.ManifestWindowSeconds
			}
			if o.DashPlaylistSettings.MinBufferTimeSeconds != nil {
				temp["dash_min_buffer_time_seconds"] = o.DashPlaylistSettings.MinBufferTimeSeconds
			}
			if o.DashPlaylistSettings.MinUpdatePeriodSeconds != nil {
				temp["dash_min_update_period_seconds"] = o.DashPlaylistSettings.MinUpdatePeriodSeconds
			}
			if o.DashPlaylistSettings.SuggestedPresentationDelaySeconds != nil {
				temp["dash_suggested_presentation_delay_seconds"] = o.DashPlaylistSettings.SuggestedPresentationDelaySeconds
			}
		}

		outputs = append(outputs, temp)
	}
	d.Set("outputs", outputs)
	d.Set("tags", res.Tags)
	d.Set("playback_mode", res.PlaybackMode)
	d.Set("tier", res.Tier)
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

func getFillerSlate(d *schema.ResourceData) *mediatailor.SlateSource {
	if v, ok := d.GetOk("filler_slate"); ok && v.([]interface{})[0] != nil {
		val := v.([]interface{})[0].(map[string]interface{})
		temp := mediatailor.SlateSource{}
		if str, ok := val["source_location_name"]; ok {
			temp.SourceLocationName = aws.String(str.(string))
		}
		if str, ok := val["vod_source_name"]; ok {
			temp.VodSourceName = aws.String(str.(string))
		}
		return &temp
	}
	return nil
}

func getCreateChannelInput(d *schema.ResourceData) mediatailor.CreateChannelInput {
	var params mediatailor.CreateChannelInput

	if v, ok := d.GetOk("channel_name"); ok {
		params.ChannelName = aws.String(v.(string))
	}

	params.FillerSlate = getFillerSlate(d)

	params.Outputs = getOutputs(d)

	if v, ok := d.GetOk("playback_mode"); ok {
		params.PlaybackMode = aws.String(v.(string))
	}

	outputMap := make(map[string]*string)
	if v, ok := d.GetOk("tags"); ok {
		val := v.(map[string]interface{})
		for k, value := range val {
			temp := value.(string)
			outputMap[k] = &temp
		}
	}
	params.Tags = outputMap

	if v, ok := d.GetOk("tier"); ok {
		params.Tier = aws.String(v.(string))
	}

	return params
}

func getUpdateChannelInput(d *schema.ResourceData) mediatailor.UpdateChannelInput {
	var params mediatailor.UpdateChannelInput

	if v, ok := d.GetOk("channel_name"); ok {
		params.ChannelName = aws.String(v.(string))
	}

	params.FillerSlate = getFillerSlate(d)

	params.Outputs = getOutputs(d)

	return params
}
