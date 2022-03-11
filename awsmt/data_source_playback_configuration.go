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

func flattenPlaybackConfiguration(configuration *mediatailor.PlaybackConfiguration) map[string]interface{} {
	if configuration != nil {
		output := make(map[string]interface{})

		output["ad_decision_server_url"] = configuration.AdDecisionServerUrl
		output["avail_suppression"] = []interface{}{map[string]interface{}{
			"mode":  configuration.AvailSuppression.Mode,
			"value": configuration.AvailSuppression.Value,
		}}
		output["bumper"] = []interface{}{map[string]interface{}{
			"end_url":   configuration.Bumper.EndUrl,
			"start_url": configuration.Bumper.StartUrl,
		}}
		output["cdn_configuration"] = []interface{}{map[string]interface{}{
			"ad_segment_url_prefix":      configuration.CdnConfiguration.AdSegmentUrlPrefix,
			"content_segment_url_prefix": configuration.CdnConfiguration.ContentSegmentUrlPrefix,
		}}
		if configuration.ConfigurationAliases != nil {
			output["configuration_aliases"] = configuration.ConfigurationAliases
		}
		output["dash_configuration"] = []interface{}{map[string]interface{}{
			"manifest_endpoint_prefix": configuration.DashConfiguration.ManifestEndpointPrefix,
			"mpd_location":             configuration.DashConfiguration.MpdLocation,
			"origin_manifest_type":     configuration.DashConfiguration.OriginManifestType,
		}}
		output["hls_configuration"] = []interface{}{map[string]interface{}{
			"manifest_endpoint_prefix": configuration.HlsConfiguration.ManifestEndpointPrefix,
		}}
		if configuration.LivePreRollConfiguration.MaxDurationSeconds != nil && *configuration.LivePreRollConfiguration.MaxDurationSeconds > 0 {
			output["live_pre_roll_configuration"] = []interface{}{map[string]interface{}{
				"ad_decision_server_url": configuration.LivePreRollConfiguration.AdDecisionServerUrl,
				"max_duration_seconds":   configuration.LivePreRollConfiguration.MaxDurationSeconds,
			}}
		}
		if configuration.LogConfiguration != nil {
			output["log_configuration"] = []interface{}{map[string]interface{}{
				"percent_enabled": configuration.LogConfiguration.PercentEnabled,
			}}
		} else {
			output["log_configuration"] = []interface{}{map[string]interface{}{
				"percent_enabled": 0,
			}}
		}

		if configuration.ManifestProcessingRules.AdMarkerPassthrough.Enabled != nil && (*configuration.ManifestProcessingRules.AdMarkerPassthrough.Enabled) {
			output["manifest_processing_rules"] = []interface{}{map[string]interface{}{
				"ad_marker_passthrough": []interface{}{map[string]interface{}{
					"enabled": configuration.ManifestProcessingRules.AdMarkerPassthrough.Enabled,
				}},
			}}
		}
		output["name"] = configuration.Name
		output["personalization_threshold_seconds"] = configuration.PersonalizationThresholdSeconds
		output["playback_configuration_arn"] = configuration.PlaybackConfigurationArn
		output["playback_endpoint_prefix"] = configuration.PlaybackEndpointPrefix
		output["session_initialization_endpoint_prefix"] = configuration.SessionInitializationEndpointPrefix
		output["slate_ad_url"] = configuration.SlateAdUrl
		output["tags"] = configuration.Tags
		output["transcode_profile_name"] = configuration.TranscodeProfileName
		output["video_content_source_url"] = configuration.VideoContentSourceUrl
		return output
	}
	return map[string]interface{}{}
}
