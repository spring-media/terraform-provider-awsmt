package awsmt

import (
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getSinglePlaybackConfiguration(c *mediatailor.MediaTailor, name string) (*mediatailor.PlaybackConfiguration, error) {
	output, err := c.GetPlaybackConfiguration(&mediatailor.GetPlaybackConfigurationInput{Name: &name})
	if err != nil {
		return nil, err
	}
	return (*mediatailor.PlaybackConfiguration)(output), nil
}

func flattenPlaybackConfiguration(c *mediatailor.PlaybackConfiguration) map[string]interface{} {
	if c == nil {
		return map[string]interface{}{}
	}
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

func getPlaybackConfigurationInput(d *schema.ResourceData) mediatailor.PutPlaybackConfigurationInput {
	input := mediatailor.PutPlaybackConfigurationInput{}
	if v, ok := d.GetOk("ad_decision_server_url"); ok {
		val := v.(string)
		input.AdDecisionServerUrl = &val
	}
	if v, ok := d.GetOk("avail_suppression"); ok && v.([]interface{})[0] != nil {
		val := v.([]interface{})[0].(map[string]interface{})
		output := mediatailor.AvailSuppression{}
		if str, ok := val["mode"]; ok {
			converted := str.(string)
			output.Mode = &converted
		}
		if str, ok := val["value"]; ok {
			converted := str.(string)
			output.Value = &converted
		}
		input.AvailSuppression = &output
	}
	if v, ok := d.GetOk("bumper"); ok && v.([]interface{})[0] != nil {
		val := v.([]interface{})[0].(map[string]interface{})
		output := mediatailor.Bumper{}
		if str, ok := val["end_url"]; ok {
			converted := str.(string)
			output.EndUrl = &converted
		}
		if str, ok := val["start_url"]; ok {
			converted := str.(string)
			output.StartUrl = &converted
		}
		input.Bumper = &output
	}
	if v, ok := d.GetOk("configuration_aliases"); ok {
		val := v.(map[string]map[string]*string)
		input.ConfigurationAliases = val
	}
	if v, ok := d.GetOk("cdn_configuration"); ok && v.([]interface{})[0] != nil {
		val := v.([]interface{})[0].(map[string]interface{})
		output := mediatailor.CdnConfiguration{}
		if str, ok := val["ad_segment_url_prefix"]; ok {
			converted := str.(string)
			output.AdSegmentUrlPrefix = &converted
		}
		if str, ok := val["content_segment_url_prefix"]; ok {
			converted := str.(string)
			output.ContentSegmentUrlPrefix = &converted
		}
		input.CdnConfiguration = &output
	}
	if v, ok := d.GetOk("dash_configuration"); ok && v.([]interface{})[0] != nil {
		val := v.([]interface{})[0].(map[string]interface{})
		output := mediatailor.DashConfigurationForPut{}
		if str, ok := val["mpd_location"]; ok {
			converted := str.(string)
			output.MpdLocation = &converted
		}
		if str, ok := val["origin_manifest_type"]; ok {
			converted := str.(string)
			output.OriginManifestType = &converted
		}
		input.DashConfiguration = &output
	}
	if v, ok := d.GetOk("live_pre_roll_configuration"); ok && v.([]interface{})[0] != nil {
		val := v.([]interface{})[0].(map[string]interface{})
		output := mediatailor.LivePreRollConfiguration{}
		if str, ok := val["ad_decision_server_url"]; ok {
			converted := str.(string)
			output.AdDecisionServerUrl = &converted
		}
		if integer, ok := val["max_duration_seconds"]; ok {
			converted := int64(integer.(int))
			output.MaxDurationSeconds = &converted
		}
		input.LivePreRollConfiguration = &output
	}
	if v, ok := d.GetOk("manifest_processing_rules"); ok && v.([]interface{})[0] != nil {
		val := v.([]interface{})[0].(map[string]interface{})
		output := mediatailor.ManifestProcessingRules{}
		if v2, ok := val["ad_marker_passthrough"]; ok {
			output2 := mediatailor.AdMarkerPassthrough{}
			val2 := v2.([]interface{})[0].(map[string]interface{})
			if boolean, ok := val2["enabled"]; ok {
				converted := boolean.(bool)
				output2.Enabled = &converted
			}
			output.AdMarkerPassthrough = &output2
		}

		input.ManifestProcessingRules = &output
	}

	if v, ok := d.GetOk("name"); ok {
		val := v.(string)
		input.Name = &val
	}
	if v, ok := d.GetOk("personalization_threshold_seconds"); ok {
		val := int64(v.(int))
		input.PersonalizationThresholdSeconds = &val
	}
	if v, ok := d.GetOk("slate_ad_url"); ok {
		val := v.(string)
		input.SlateAdUrl = &val
	}
	outputMap := make(map[string]*string)
	if v, ok := d.GetOk("tags"); ok {
		val := v.(map[string]interface{})
		for k, value := range val {
			temp := value.(string)
			outputMap[k] = &temp
		}
		input.Tags = outputMap
	} else {
		input.Tags = outputMap
	}
	if v, ok := d.GetOk("transcode_profile_name"); ok {
		val := v.(string)
		input.TranscodeProfileName = &val
	}
	if v, ok := d.GetOk("video_content_source_url"); ok {
		val := v.(string)
		input.VideoContentSourceUrl = &val
	}
	return input
}

func returnPlaybackConfigurationResource(d *schema.ResourceData, values map[string]interface{}, diags diag.Diagnostics) diag.Diagnostics {
	for k := range values {
		setSingleValue(d, values, diags, k)
	}
	return diags
}

func setSingleValue(d *schema.ResourceData, values map[string]interface{}, diags diag.Diagnostics, name string) diag.Diagnostics {
	err := d.Set(name, values[name])
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func deletePlaybackConfiguration(client *mediatailor.MediaTailor, name string) diag.Diagnostics {
	_, err := client.DeletePlaybackConfiguration(&mediatailor.DeletePlaybackConfigurationInput{Name: &name})
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
