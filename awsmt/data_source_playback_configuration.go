package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePlaybackConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePlaybackConfigurationRead,
		// schema based on https://docs.aws.amazon.com/sdk-for-go/api/service/mediatailor/#GetPlaybackConfigurationOutput
		// with types found on https://sourcegraph.com/github.com/aws/aws-sdk-go/-/docs/service/mediatailor
		Schema: map[string]*schema.Schema{
			"name": &requiredString,
			"configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ad_decision_server_url": &computedString,
						"avail_suppression": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									//OFF | BEHIND_LIVE_EDGE
									"mode":  &computedString,
									"value": &computedString,
								},
							},
						},
						"bumper": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"end_url":   &computedString,
									"start_url": &computedString,
								},
							},
						},
						"cdn_configuration": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ad_segment_url_prefix":      &computedString,
									"content_segment_url_prefix": &computedString,
								},
							},
						},
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
						"dash_configuration": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"manifest_endpoint_prefix": &computedString,
									"mpd_location":             &computedString,
									"origin_manifest_type":     &computedString,
								},
							},
						},
						"hls_configuration": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"manifest_endpoint_prefix": &computedString,
								},
							},
						},
						"live_pre_roll_configuration": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ad_decision_server_url": &computedString,
									"max_duration_seconds":   &computedInt,
								},
							},
						},
						"log_configuration": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"percent_enabled": &computedInt,
								},
							},
						},
						"manifest_processing_rules": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ad_marker_passthrough": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": &computedBool,
											},
										},
									},
								},
							},
						},
						"name":                                   &computedString,
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
				},
			},
		},
	}
}

func dataSourcePlaybackConfigurationRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mediatailor.MediaTailor)
	var diags diag.Diagnostics
	var output []interface{}

	name := d.Get("name").(string)

	if name != "" {
		res, err := getSinglePlaybackConfiguration(client, name)
		if err != nil {
			return diag.FromErr(err)
		}
		output = []interface{}{flattenPlaybackConfiguration(res)}
		if err := d.Set("configuration", output); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(uuid.New().String())
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "`name` parameter required",
			Detail:   "You need to specify a `name` parameter in your configuration",
		})
	}
	return diags
}

func getSinglePlaybackConfiguration(c *mediatailor.MediaTailor, name string) (*mediatailor.PlaybackConfiguration, error) {
	output, err := c.GetPlaybackConfiguration(&mediatailor.GetPlaybackConfigurationInput{Name: &name})
	if err != nil {
		return nil, err
	}
	return (*mediatailor.PlaybackConfiguration)(output), nil
}

func flattenPlaybackConfiguration(c *mediatailor.PlaybackConfiguration) map[string]interface{} {
	if c != nil {
		output := make(map[string]interface{})

		output["ad_decision_server_url"] = c.AdDecisionServerUrl
		if !(*c.AvailSuppression.Mode == "OFF" && c.AvailSuppression.Value == nil) {
			output["avail_suppression"] = []interface{}{map[string]interface{}{
				"mode":  c.AvailSuppression.Mode,
				"value": c.AvailSuppression.Value,
			}}
		}
		if c.Bumper.EndUrl != nil || c.Bumper.StartUrl != nil {
			output["bumper"] = []interface{}{map[string]interface{}{
				"end_url":   c.Bumper.EndUrl,
				"start_url": c.Bumper.StartUrl,
			}}
		}
		if c.CdnConfiguration.AdSegmentUrlPrefix != nil || c.CdnConfiguration.ContentSegmentUrlPrefix != nil {
			output["cdn_configuration"] = []interface{}{map[string]interface{}{
				"ad_segment_url_prefix":      c.CdnConfiguration.AdSegmentUrlPrefix,
				"content_segment_url_prefix": c.CdnConfiguration.ContentSegmentUrlPrefix,
			}}
		}
		if c.ConfigurationAliases != nil {
			output["configuration_aliases"] = c.ConfigurationAliases
		}
		output["dash_configuration"] = []interface{}{map[string]interface{}{
			"manifest_endpoint_prefix": c.DashConfiguration.ManifestEndpointPrefix,
			"mpd_location":             c.DashConfiguration.MpdLocation,
			"origin_manifest_type":     c.DashConfiguration.OriginManifestType,
		}}
		output["hls_configuration"] = []interface{}{map[string]interface{}{
			"manifest_endpoint_prefix": c.HlsConfiguration.ManifestEndpointPrefix,
		}}
		if !(c.LivePreRollConfiguration.MaxDurationSeconds == nil && c.LivePreRollConfiguration.AdDecisionServerUrl == nil) {
			output["live_pre_roll_configuration"] = []interface{}{map[string]interface{}{
				"ad_decision_server_url": c.LivePreRollConfiguration.AdDecisionServerUrl,
				"max_duration_seconds":   c.LivePreRollConfiguration.MaxDurationSeconds,
			}}
		}
		if c.LogConfiguration != nil {
			output["log_configuration"] = []interface{}{map[string]interface{}{
				"percent_enabled": c.LogConfiguration.PercentEnabled,
			}}
		} else {
			output["log_configuration"] = []interface{}{map[string]interface{}{
				"percent_enabled": 0,
			}}
		}
		if *c.ManifestProcessingRules.AdMarkerPassthrough.Enabled {
			output["manifest_processing_rules"] = []interface{}{map[string]interface{}{
				"ad_marker_passthrough": []interface{}{map[string]interface{}{
					"enabled": c.ManifestProcessingRules.AdMarkerPassthrough.Enabled,
				}},
			}}
		}
		output["name"] = c.Name
		output["personalization_threshold_seconds"] = c.PersonalizationThresholdSeconds
		output["playback_configuration_arn"] = c.PlaybackConfigurationArn
		output["playback_endpoint_prefix"] = c.PlaybackEndpointPrefix
		output["session_initialization_endpoint_prefix"] = c.SessionInitializationEndpointPrefix
		output["slate_ad_url"] = c.SlateAdUrl
		output["tags"] = c.Tags
		output["transcode_profile_name"] = c.TranscodeProfileName
		output["video_content_source_url"] = c.VideoContentSourceUrl
		return output
	}
	return map[string]interface{}{}
}
