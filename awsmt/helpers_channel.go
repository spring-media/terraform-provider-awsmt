package awsmt

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strings"
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

	if v, ok := d.GetOk("name"); ok {
		params.ChannelName = aws.String(v.(string))
	}

	if f := getFillerSlate(d); f != nil {
		params.FillerSlate = f
	}

	if o := getOutputs(d); o != nil {
		params.Outputs = o
	}

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

	if v, ok := d.GetOk("name"); ok {
		params.ChannelName = aws.String(v.(string))
	}

	params.FillerSlate = getFillerSlate(d)

	if o := getOutputs(d); o != nil {
		params.Outputs = o
	}

	return params
}

func createChannelPolicy(client *mediatailor.MediaTailor, d *schema.ResourceData) error {
	if v, ok := d.GetOk("policy"); ok {
		var putChannelPolicyParams = mediatailor.PutChannelPolicyInput{
			ChannelName: aws.String((d.Get("name")).(string)),
			Policy:      aws.String(v.(string)),
		}

		_, err := client.PutChannelPolicy(&putChannelPolicyParams)
		if err != nil {
			return fmt.Errorf("error while creating the policy: %v", err)
		}
	}
	return nil
}

func updateChannelPolicy(client *mediatailor.MediaTailor, d *schema.ResourceData, channelName *string) error {
	_, err := client.PutChannelPolicy(&mediatailor.PutChannelPolicyInput{ChannelName: channelName, Policy: aws.String(d.Get("policy").(string))})

	if err != nil && !strings.Contains(err.Error(), "NotFound") {
		return fmt.Errorf("error while getting the channel policy: %v", err)
	}
	return nil
}

func updatePolicy(client *mediatailor.MediaTailor, d *schema.ResourceData, channelName *string) error {
	if d.HasChange("policy") {
		_, newValue := d.GetChange("policy")
		if len(newValue.(string)) > 0 {
			err := updateChannelPolicy(client, d, channelName)
			if err != nil {
				return err
			}
		} else {
			err := deleteChannelPolicy(client, d, channelName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func deleteChannelPolicy(client *mediatailor.MediaTailor, d *schema.ResourceData, channelName *string) error {
	_, err := client.DeleteChannelPolicy(&mediatailor.DeleteChannelPolicyInput{ChannelName: channelName})
	if err != nil {
		return fmt.Errorf("error while deleting the policy: %v", err)
	}
	if err := d.Set("policy", ""); err != nil {
		return fmt.Errorf("error while unsetting the policy: %v", err)
	}
	return nil
}

func getResourceName(d *schema.ResourceData, fieldName string) (*string, error) {
	resourceName := d.Get(fieldName).(string)
	if len(resourceName) == 0 && len(d.Id()) > 0 {
		resourceArn, err := arn.Parse(d.Id())
		if err != nil {
			return nil, fmt.Errorf("error parsing the name from resource arn: %v", err)
		}
		arnSections := strings.Split(resourceArn.Resource, "/")
		resourceName = arnSections[len(arnSections)-1]
	}
	return &resourceName, nil
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

func flattenOutput(o *mediatailor.ResponseOutputItem) map[string]interface{} {
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
	return temp
}

func setOutputs(values *mediatailor.DescribeChannelOutput, d *schema.ResourceData) error {
	var outputs []map[string]interface{}
	for _, o := range values.Outputs {
		temp := flattenOutput(o)
		outputs = append(outputs, temp)
	}
	if err := d.Set("outputs", outputs); err != nil {
		return fmt.Errorf("error while setting the outputs: %w", err)
	}
	return nil
}

func setChannel(res *mediatailor.DescribeChannelOutput, d *schema.ResourceData) error {
	var errors []error

	errors = append(errors, d.Set("arn", res.Arn))
	errors = append(errors, d.Set("name", res.ChannelName))
	errors = append(errors, d.Set("channel_state", res.ChannelState))
	errors = append(errors, d.Set("creation_time", res.CreationTime.String()))
	errors = append(errors, setFillerState(res, d))
	errors = append(errors, d.Set("last_modified_time", res.LastModifiedTime.String()))
	errors = append(errors, setOutputs(res, d))
	errors = append(errors, d.Set("tags", res.Tags))
	errors = append(errors, d.Set("playback_mode", res.PlaybackMode))
	errors = append(errors, d.Set("tier", res.Tier))

	for _, e := range errors {
		if e != nil {
			return fmt.Errorf("the following error occured while setting the values: %w", e)
		}
	}
	return nil
}

func setChannelPolicy(res *mediatailor.GetChannelPolicyOutput, d *schema.ResourceData) error {
	if res.Policy != nil {
		if err := d.Set("policy", res.Policy); err != nil {
			return fmt.Errorf("error while setting the  the channel policy: %v", err)
		}
	}
	return nil
}

func getOutput(output interface{}) *mediatailor.RequestOutputItem {
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
	return &temp
}

func getOutputs(d *schema.ResourceData) []*mediatailor.RequestOutputItem {
	if v, ok := d.GetOk("outputs"); ok && v.([]interface{})[0] != nil {
		outputs := v.([]interface{})

		var res []*mediatailor.RequestOutputItem

		for _, output := range outputs {

			var temp = getOutput(output)

			res = append(res, temp)
		}
		return res
	}
	return nil
}

func startChannel(client *mediatailor.MediaTailor, channelName string) error {
	_, err := client.StartChannel(&mediatailor.StartChannelInput{
		ChannelName: aws.String(channelName),
	})
	if err != nil {
		return fmt.Errorf("error while starting the channel: %v", err)
	}
	return nil
}

func stopChannel(client *mediatailor.MediaTailor, channelName string) error {
	_, err := client.StopChannel(&mediatailor.StopChannelInput{
		ChannelName: aws.String(channelName),
	})
	if err != nil {
		return fmt.Errorf("error while stopping the channel: %v", err)
	}
	return nil
}

func checkStatusAndStartChannel(client *mediatailor.MediaTailor, d *schema.ResourceData) error {
	if v, ok := d.GetOk("channel_state"); ok && v != nil && v.(string) != "" {
		if v.(string) == "RUNNING" {
			if err := startChannel(client, d.Get("name").(string)); err != nil {
				return err
			}
		}
	}
	return nil
}
