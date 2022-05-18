package awsmt

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

func setFillerState(values *mediatailor.DescribeChannelOutput, d *schema.ResourceData) error {
	if values.FillerSlate != nil && values.FillerSlate != &(mediatailor.SlateSource{}) {
		temp := map[string]interface{}{}
		if values.FillerSlate.SourceLocationName != nil {
			temp["source_location_name"] = values.FillerSlate.SourceLocationName
		}
		if values.FillerSlate.VodSourceName != nil {
			temp["vod_source_name"] = values.FillerSlate.VodSourceName
		}
		if err := d.Set("filler_slate", []interface{}{temp}); err != nil {
			return fmt.Errorf("error while setting the filler slate: %w", err)
		}
	}
	return nil
}

func setOutputs(values *mediatailor.DescribeChannelOutput, d *schema.ResourceData) error {
	var outputs []map[string]interface{}
	for _, o := range values.Outputs {
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
	if err := d.Set("outputs", outputs); err != nil {
		return fmt.Errorf("error while setting the outputs: %w", err)
	}
	return nil
}
