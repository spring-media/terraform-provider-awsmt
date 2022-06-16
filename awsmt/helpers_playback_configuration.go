package awsmt

import (
	"github.com/aws/aws-sdk-go/aws"
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

type CreateInput struct {
	d     *schema.ResourceData
	input *mediatailor.PutPlaybackConfigurationInput
}

func (i CreateInput) getAdDecisionServerUrlInput() {
	if v, ok := i.d.GetOk("ad_decision_server_url"); ok {
		val := v.(string)
		i.input.AdDecisionServerUrl = &val
	}
}

func (i CreateInput) getAvailSuppressionInput() {
	if v, ok := i.d.GetOk("avail_suppression"); ok && v.([]interface{})[0] != nil {
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
		i.input.AvailSuppression = &output
	}
}

func (i CreateInput) getManifestProcessingRulesInput() {
	if v, ok := i.d.GetOk("manifest_processing_rules"); ok && v.([]interface{})[0] != nil {
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

		i.input.ManifestProcessingRules = &output
	}
}

func (i CreateInput) getBumperInput() {
	if v, ok := i.d.GetOk("bumper"); ok && v.([]interface{})[0] != nil {
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
		i.input.Bumper = &output
	}
}

func (i CreateInput) getConfigurationAliasesInput() {
	if v, ok := i.d.GetOk("configuration_aliases"); ok {
		val := v.(map[string]map[string]*string)
		i.input.ConfigurationAliases = val
	}
}

func (i CreateInput) getCDNConfigurationInput() {
	if v, ok := i.d.GetOk("cdn_configuration"); ok && v.([]interface{})[0] != nil {
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
		i.input.CdnConfiguration = &output
	}
}

func (i CreateInput) getDashConfigurationInput() {
	if v, ok := i.d.GetOk("dash_configuration"); ok && v.([]interface{})[0] != nil {
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
		i.input.DashConfiguration = &output
	}
}

func (i CreateInput) getLivePreRollConfigurationInput() {
	if v, ok := i.d.GetOk("live_pre_roll_configuration"); ok && v.([]interface{})[0] != nil {
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
		i.input.LivePreRollConfiguration = &output
	}
}

func (i CreateInput) getTagsInput() {
	outputMap := make(map[string]*string)
	if v, ok := i.d.GetOk("tags"); ok {
		val := v.(map[string]interface{})
		for k, value := range val {
			temp := value.(string)
			outputMap[k] = &temp
		}
	}
	i.input.Tags = outputMap
}

func (i CreateInput) getNameInput() {
	if v, ok := i.d.GetOk("name"); ok {
		val := v.(string)
		i.input.Name = &val
	}
}

func (i CreateInput) getPersonalizationThresholdSecondsInput() {
	if v, ok := i.d.GetOk("personalization_threshold_seconds"); ok {
		val := int64(v.(int))
		i.input.PersonalizationThresholdSeconds = &val
	}
}

func (i CreateInput) getSlateAdUrlInput() {
	if v, ok := i.d.GetOk("slate_ad_url"); ok {
		val := v.(string)
		i.input.SlateAdUrl = &val
	}
}

func (i CreateInput) getTranscodeProfileNameInput() {
	if v, ok := i.d.GetOk("transcode_profile_name"); ok {
		val := v.(string)
		i.input.TranscodeProfileName = &val
	}
}

func (i CreateInput) getVideoContentSourceUrlInput() {
	if v, ok := i.d.GetOk("video_content_source_url"); ok {
		val := v.(string)
		i.input.VideoContentSourceUrl = &val
	}
}

func getPlaybackConfigurationInput(d *schema.ResourceData) mediatailor.PutPlaybackConfigurationInput {
	input := mediatailor.PutPlaybackConfigurationInput{}
	i := CreateInput{
		d:     d,
		input: &input,
	}
	i.getAdDecisionServerUrlInput()
	i.getAvailSuppressionInput()
	i.getBumperInput()
	i.getConfigurationAliasesInput()
	i.getCDNConfigurationInput()
	i.getDashConfigurationInput()
	i.getLivePreRollConfigurationInput()
	i.getManifestProcessingRulesInput()
	i.getNameInput()
	i.getPersonalizationThresholdSecondsInput()
	i.getSlateAdUrlInput()
	i.getTagsInput()
	i.getTranscodeProfileNameInput()
	i.getVideoContentSourceUrlInput()
	return input
}

func returnPlaybackConfiguration(d *schema.ResourceData, values map[string]interface{}, diags diag.Diagnostics) diag.Diagnostics {
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

func makeBaseList(fields map[string]*schema.Schema) *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}
}

func createOptionalList(fields map[string]*schema.Schema) *schema.Schema {
	s := makeBaseList(fields)
	s.Optional = true
	s.MaxItems = 1
	return s
}

func createRequiredList(fields map[string]*schema.Schema) *schema.Schema {
	s := makeBaseList(fields)
	s.Required = true
	s.MaxItems = 1
	return s
}

func createComputedList(fields map[string]*schema.Schema) *schema.Schema {
	s := makeBaseList(fields)
	s.Computed = true
	return s
}

func deleteTags(client *mediatailor.MediaTailor, resourceArn string, removedTags []string) error {
	if len(removedTags) != 0 {

		var removedValuesPointer []*string
		for i := range removedTags {
			removedValuesPointer = append(removedValuesPointer, &removedTags[i])
		}

		untagInput := mediatailor.UntagResourceInput{ResourceArn: aws.String(resourceArn), TagKeys: removedValuesPointer}
		_, err := client.UntagResource(&untagInput)
		if err != nil {
			return err
		}
	}
	return nil
}
