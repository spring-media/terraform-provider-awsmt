package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePlaybackConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePlaybackConfigurationRead,
		// schema based on https://docs.aws.amazon.com/sdk-for-go/api/service/mediatailor/#GetPlaybackConfigurationOutput
		// with types found on https://sourcegraph.com/github.com/aws/aws-sdk-go/-/docs/service/mediatailor
		Schema: map[string]*schema.Schema{
			"name":                   &requiredString,
			"ad_decision_server_url": &computedString,
			"avail_suppression": createComputedList(map[string]*schema.Schema{
				"mode":  &computedString,
				"value": &computedString,
			}),
			"bumper": createComputedList(map[string]*schema.Schema{
				"end_url":   &computedString,
				"start_url": &computedString,
			}),
			"cdn_configuration": createComputedList(map[string]*schema.Schema{
				"ad_segment_url_prefix":      &computedString,
				"content_segment_url_prefix": &computedString,
			}),
			"configuration_aliases": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type:     schema.TypeMap,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			"dash_configuration": createComputedList(map[string]*schema.Schema{
				"manifest_endpoint_prefix": &computedString,
				"mpd_location":             &computedString,
				"origin_manifest_type":     &computedString,
			}),
			"hls_configuration": createComputedList(map[string]*schema.Schema{
				"manifest_endpoint_prefix": &computedString,
			}),
			"live_pre_roll_configuration": createComputedList(map[string]*schema.Schema{
				"ad_decision_server_url": &computedString,
				"max_duration_seconds":   &computedInt,
			}),
			"log_configuration": createComputedList(map[string]*schema.Schema{
				"percent_enabled": &computedInt,
			}),
			"manifest_processing_rules": createComputedList(map[string]*schema.Schema{
				"ad_marker_passthrough": createComputedList(map[string]*schema.Schema{
					"enabled": &computedBool,
				}),
			}),
			"personalization_threshold_seconds":      &computedInt,
			"playback_configuration_arn":             &computedString,
			"playback_endpoint_prefix":               &computedString,
			"session_initialization_endpoint_prefix": &computedString,
			"slate_ad_url":                           &computedString,
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"transcode_profile_name":   &computedString,
			"video_content_source_url": &computedString,
		},
	}
}

func dataSourcePlaybackConfigurationRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mediatailor.MediaTailor)
	var diags diag.Diagnostics

	name := d.Get("name").(string)

	res, err := getSinglePlaybackConfiguration(client, name)
	if err != nil {
		return diag.FromErr(err)
	}
	output := flattenPlaybackConfiguration(res)
	returnPlaybackConfiguration(d, output, diags)
	d.SetId(*res.PlaybackConfigurationArn)

	return diags
}
