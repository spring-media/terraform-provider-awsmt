package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func resourcePlaybackConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePlaybackConfigurationCreate,
		UpdateContext: resourcePlaybackConfigurationUpdate,
		ReadContext:   resourcePlaybackConfigurationRead,
		DeleteContext: resourcePlaybackConfigurationDelete,
		// schema based on: https://docs.aws.amazon.com/mediatailor/latest/apireference/playbackconfiguration.html#playbackconfiguration-prop-putplaybackconfigurationrequest-personalizationthresholdseconds
		// and https://sourcegraph.com/github.com/aws/aws-sdk-go/-/docs/service/mediatailor#PutPlaybackConfigurationInput
		Schema: map[string]*schema.Schema{
			"ad_decision_server_url": &requiredString,
			"avail_suppression": createOptionalList(map[string]*schema.Schema{
				"mode":  &optionalString,
				"value": &optionalString,
			}),
			"bumper": createOptionalList(map[string]*schema.Schema{
				"end_url":   &optionalString,
				"start_url": &optionalString,
			}),
			"cdn_configuration": createOptionalList(map[string]*schema.Schema{
				"ad_segment_url_prefix":      &optionalString,
				"content_segment_url_prefix": &optionalString,
			}),
			"configuration_aliases": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type:     schema.TypeMap,
					Optional: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			"dash_configuration": createRequiredList(map[string]*schema.Schema{
				"manifest_endpoint_prefix": &computedString,
				"mpd_location":             &optionalString,
				"origin_manifest_type":     &optionalString,
			}),
			"hls_configuration": createComputedList(map[string]*schema.Schema{
				"manifest_endpoint_prefix": &computedString,
			}),
			"live_pre_roll_configuration": createOptionalList(map[string]*schema.Schema{
				"ad_decision_server_url": &optionalString,
				"max_duration_seconds":   &optionalInt,
			}),
			"log_configuration": createComputedList(map[string]*schema.Schema{
				"percent_enabled": &computedInt,
			}),
			"manifest_processing_rules": createOptionalList(map[string]*schema.Schema{
				"ad_marker_passthrough": createOptionalList(map[string]*schema.Schema{
					"enabled": &optionalBool,
				}),
			}),
			"name":                                   &requiredString,
			"personalization_threshold_seconds":      &optionalInt,
			"playback_configuration_arn":             &computedString,
			"playback_endpoint_prefix":               &computedString,
			"session_initialization_endpoint_prefix": &computedString,
			"slate_ad_url":                           &optionalString,
			"tags":                                   &optionalTags,
			"transcode_profile_name":                 &optionalString,
			"video_content_source_url":               &requiredString,
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		CustomizeDiff: customdiff.Sequence(
			customdiff.ForceNewIfChange("name", func(ctx context.Context, old, new, meta interface{}) bool { return old.(string) != new.(string) }),
		),
	}
}

func resourcePlaybackConfigurationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mediatailor.MediaTailor)
	var diags diag.Diagnostics

	input := getPlaybackConfigurationInput(d)

	_, err := client.PutPlaybackConfiguration(&input)
	if err != nil {
		return diag.FromErr(err)
	}
	resourcePlaybackConfigurationRead(ctx, d, m)
	d.SetId(*input.Name)
	return diags
}

func resourcePlaybackConfigurationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*mediatailor.MediaTailor)

	// @ADR
	// Context: Updating tags using the PutPlaybackConfiguration method does not allow to remove them.
	// Decision: We decided to check for removed tags and remove them using the UntagResource method, while we still use
	// the PutPlaybackConfiguration method to add and update tags. We use this approach for every resource in the provider.
	// Consequences: The Update function logic is now more complicated, but tag removal is supported.
	if d.HasChange("tags") {
		oldValue, newValue := d.GetChange("tags")
		var removedTags []string
		for k := range oldValue.(map[string]interface{}) {
			if _, ok := (newValue.(map[string]interface{}))[k]; !ok {
				removedTags = append(removedTags, k)
			}
		}
		resourceName := d.Get("name").(string)
		res, err := client.GetPlaybackConfiguration(&mediatailor.GetPlaybackConfigurationInput{Name: &resourceName})
		if err != nil {
			return diag.FromErr(err)
		}
		err = deleteTags(client, aws.StringValue(res.PlaybackConfigurationArn), removedTags)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	input := getPlaybackConfigurationInput(d)
	_, err := client.PutPlaybackConfiguration(&input)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("last_updated", time.Now().Format(time.RFC850)); err != nil {
		return diag.FromErr(err)
	}

	resourcePlaybackConfigurationRead(ctx, d, m)
	return diags
}

func resourcePlaybackConfigurationRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mediatailor.MediaTailor)
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	if len(name) == 0 && len(d.Id()) > 0 {
		name = d.Id()
	}
	res, err := client.GetPlaybackConfiguration(&mediatailor.GetPlaybackConfigurationInput{Name: &name})
	if err != nil {
		return diag.FromErr(err)
	}

	output := flattenPlaybackConfiguration((*mediatailor.PlaybackConfiguration)(res))
	returnPlaybackConfiguration(d, output, diags)
	return diags
}

func resourcePlaybackConfigurationDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mediatailor.MediaTailor)
	var diags diag.Diagnostics
	deletePlaybackConfiguration(client, d.Get("name").(string))
	d.SetId("")
	return diags
}
