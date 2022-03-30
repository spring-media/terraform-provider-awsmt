package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
			"ad_decision_server_url": &optionalString,
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
			"dash_configuration": createOptionalList(map[string]*schema.Schema{
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
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"transcode_profile_name":   &optionalString,
			"video_content_source_url": &requiredString,
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
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
	if !d.HasChanges("ad_decision_server_url", "avail_suppression", "bumper", "cdn_configuration", "configuration_aliases", "dash_configuration", "live_pre_roll_configuration", "log_configuration", "manifest_processing_rules", "name", "personalization_threshold_seconds", "slate_ad_url", "tags", "transcode_profile_name", "video_content_source_url") {
		resourcePlaybackConfigurationRead(ctx, d, m)
		return diags
	}
	client := m.(*mediatailor.MediaTailor)

	if d.HasChange("name") {
		oldValue, newValue := d.GetChange("name")
		oldName := oldValue.(string)
		deletePlaybackConfiguration(client, oldName)
		d.SetId(newValue.(string))
	}

	if d.HasChange("tags") {
		oldValue, newValue := d.GetChange("tags")
		for k := range oldValue.(map[string]interface{}) {
			if _, ok := (newValue.(map[string]interface{}))[k]; !ok {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Warning,
					Summary:  "Tag removal detected, but not supported.",
					Detail:   "This provider does not support tag removal. For more information about the issue, visit this link: https://github.com/aws/aws-sdk-go/issues/4337\nThe tag(s) will only be removed from the terraform state.",
				})
				break
			}
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
	returnPlaybackConfigurationResource(d, output, diags)
	return diags
}

func resourcePlaybackConfigurationDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mediatailor.MediaTailor)
	var diags diag.Diagnostics
	deletePlaybackConfiguration(client, d.Get("name").(string))
	d.SetId("")
	return diags
}
