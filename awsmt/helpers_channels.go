package awsmt

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

// POLICY

func createChannelPolicy(ctx context.Context, req resource.CreateRequest, client *mediatailor.MediaTailor) error {

	var plan resourceChannelModel

	_ = req.Plan.Get(ctx, &plan)

	policy := plan.Policy.ValueString()

	var putChannelPolicyParams = mediatailor.PutChannelPolicyInput{
		ChannelName: plan.Name,
		Policy:      aws.String(policy),
	}

	_, err := client.PutChannelPolicy(&putChannelPolicyParams)
	if err != nil {
		return err
	}
	return err
}

func setChannelPolicy(res *mediatailor.GetChannelPolicyOutput) error {
	var state resourceChannelModel

	if res.Policy != nil {
		state.Policy = types.StringValue(aws.StringValue(res.Policy))
	}
	return nil
}

func updatePolicy(client *mediatailor.MediaTailor, channelName *string, oldPolicy string, newPolicy string) error {

	if oldPolicy != newPolicy {
		if len(newPolicy) > 0 {
			err := updateChannelPolicy(client, newPolicy, channelName)
			if err != nil {
				return fmt.Errorf("error while updating the policy: %v", err)
			}
		} else {
			err := deleteChannelPolicy(client, channelName)
			if err != nil {
				return fmt.Errorf("error while deleting the policy: %v", err)
			}
		}
	}
	return nil
}

func updateChannelPolicy(client *mediatailor.MediaTailor, newPolicy string, channelName *string) error {
	_, err := client.PutChannelPolicy(&mediatailor.PutChannelPolicyInput{ChannelName: channelName, Policy: &newPolicy})
	if err != nil && !strings.Contains(err.Error(), "NotFound") {
		return fmt.Errorf("error while updating the policy: %v", err)
	}
	return nil
}

func deleteChannelPolicy(client *mediatailor.MediaTailor, channelName *string) error {
	_, err := client.DeleteChannelPolicy(&mediatailor.DeleteChannelPolicyInput{ChannelName: channelName})
	if err != nil {
		return fmt.Errorf("error while deleting the policy: %v", err)
	}
	return nil
}

// CHANNEL STATE

func startChannel(client *mediatailor.MediaTailor, channelName *string) error {

	_, err := client.StartChannel(&mediatailor.StartChannelInput{
		ChannelName: channelName,
	})
	if err != nil {
		return err
	}

	return nil
}

func stopChannel(client *mediatailor.MediaTailor, channelName *string) error {
	_, err := client.StopChannel(&mediatailor.StopChannelInput{
		ChannelName: channelName,
	})
	if err != nil {
		return err
	}
	return nil
}

func checkStatusAndStartChannel(ctx context.Context, req resource.CreateRequest, client *mediatailor.MediaTailor) error {
	var plan resourceChannelModel

	_ = req.Plan.Get(ctx, &plan)

	if plan.ChannelState == types.StringValue("RUNNING") {
		if err := startChannel(client, plan.Name); err != nil {
			return err
		}
	}
	return nil
}

// OUTPUTS

func getOutputs(ctx context.Context, req resource.CreateRequest, resp resource.CreateResponse) []*mediatailor.RequestOutputItem {
	var plan resourceChannelModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return nil
	}

	var outputs []*mediatailor.RequestOutputItem

	if plan.Outputs != nil && len(plan.Outputs) > 0 {
		for _, o := range plan.Outputs {

			output := mediatailor.RequestOutputItem{}

			if *o.ManifestName != "" && o.ManifestName != nil {
				output.ManifestName = o.ManifestName
			}

			if *o.SourceGroup != "" {
				output.SourceGroup = o.SourceGroup
			}

			if o.HlsPlaylistSettings != nil && len(o.HlsPlaylistSettings) > 0 {
				outputsHls := mediatailor.HlsPlaylistSettings{}
				outputsHls.ManifestWindowSeconds = o.HlsPlaylistSettings[0].ManifestWindowsSeconds
				output.HlsPlaylistSettings = &outputsHls
			}
			if o.DashPlaylistSettings != nil && len(o.DashPlaylistSettings) > 0 {
				manifestWindowSecondsDash := o.DashPlaylistSettings[0].ManifestWindowsSeconds
				minBufferTimeSeconds := o.DashPlaylistSettings[0].MinBufferTimeSeconds
				minUpdatePeriodSeconds := o.DashPlaylistSettings[0].MinUpdatePeriodSeconds
				suggestedPresentationDelaySeconds := o.DashPlaylistSettings[0].SuggestedPresentationDelaySeconds
				outputsDash := mediatailor.DashPlaylistSettings{}
				outputsDash.ManifestWindowSeconds = manifestWindowSecondsDash
				outputsDash.MinBufferTimeSeconds = minBufferTimeSeconds
				outputsDash.MinUpdatePeriodSeconds = minUpdatePeriodSeconds
				outputsDash.SuggestedPresentationDelaySeconds = suggestedPresentationDelaySeconds
				output.DashPlaylistSettings = &outputsDash
			}

			outputs = append(outputs, &output)

		}
	}
	return outputs
}

func getOutputsUpdate(ctx context.Context, req resource.UpdateRequest, resp resource.UpdateResponse) []*mediatailor.RequestOutputItem {
	var plan resourceChannelModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return nil
	}

	var outputs []*mediatailor.RequestOutputItem

	if plan.Outputs != nil && len(plan.Outputs) > 0 {
		for _, o := range plan.Outputs {

			output := mediatailor.RequestOutputItem{}

			if *o.ManifestName != "" && o.ManifestName != nil {
				output.ManifestName = o.ManifestName
			}

			if *o.SourceGroup != "" {
				output.SourceGroup = o.SourceGroup
			}

			if o.HlsPlaylistSettings != nil && len(o.HlsPlaylistSettings) > 0 {
				outputsHls := mediatailor.HlsPlaylistSettings{}
				outputsHls.ManifestWindowSeconds = o.HlsPlaylistSettings[0].ManifestWindowsSeconds
				output.HlsPlaylistSettings = &outputsHls
			}
			if o.DashPlaylistSettings != nil && len(o.DashPlaylistSettings) > 0 {
				manifestWindowSecondsDash := o.DashPlaylistSettings[0].ManifestWindowsSeconds
				minBufferTimeSeconds := o.DashPlaylistSettings[0].MinBufferTimeSeconds
				minUpdatePeriodSeconds := o.DashPlaylistSettings[0].MinUpdatePeriodSeconds
				suggestedPresentationDelaySeconds := o.DashPlaylistSettings[0].SuggestedPresentationDelaySeconds
				outputsDash := mediatailor.DashPlaylistSettings{}
				outputsDash.ManifestWindowSeconds = manifestWindowSecondsDash
				outputsDash.MinBufferTimeSeconds = minBufferTimeSeconds
				outputsDash.MinUpdatePeriodSeconds = minUpdatePeriodSeconds
				outputsDash.SuggestedPresentationDelaySeconds = suggestedPresentationDelaySeconds
				output.DashPlaylistSettings = &outputsDash
			}

			outputs = append(outputs, &output)

		}
	}
	return outputs
}

// FILLER SLATE

func getFillerSlate(ctx context.Context, req resource.CreateRequest, resp resource.CreateResponse) *mediatailor.SlateSource {
	var plan resourceChannelModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return nil
	}

	fillerSlate := mediatailor.SlateSource{}

	if plan.FillerSlate != nil || len(plan.FillerSlate) > 0 {
		if plan.FillerSlate[0].SourceLocationName != nil {
			fillerSlate = mediatailor.SlateSource{
				SourceLocationName: plan.FillerSlate[0].SourceLocationName,
			}
		}
		if plan.FillerSlate[0].VodSourceName != nil || len(plan.FillerSlate) > 0 {
			fillerSlate = mediatailor.SlateSource{
				VodSourceName: plan.FillerSlate[0].VodSourceName,
			}
		}
		if plan.FillerSlate != nil && plan.FillerSlate[0].VodSourceName != nil || len(plan.FillerSlate) > 0 {
			fillerSlate = mediatailor.SlateSource{
				SourceLocationName: plan.FillerSlate[0].SourceLocationName,
				VodSourceName:      plan.FillerSlate[0].VodSourceName,
			}
		}
		return &fillerSlate
	}
	return nil
}

func getFillerSlateUpdate(ctx context.Context, req resource.UpdateRequest, resp resource.UpdateResponse) *mediatailor.SlateSource {
	var plan resourceChannelModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return nil
	}

	fillerSlate := mediatailor.SlateSource{}

	if plan.FillerSlate != nil || len(plan.FillerSlate) > 0 {
		if plan.FillerSlate[0].SourceLocationName != nil {
			fillerSlate = mediatailor.SlateSource{
				SourceLocationName: plan.FillerSlate[0].SourceLocationName,
			}
		}
		if plan.FillerSlate[0].VodSourceName != nil || len(plan.FillerSlate) > 0 {
			fillerSlate = mediatailor.SlateSource{
				VodSourceName: plan.FillerSlate[0].VodSourceName,
			}
		}
		if plan.FillerSlate != nil && plan.FillerSlate[0].VodSourceName != nil || len(plan.FillerSlate) > 0 {
			fillerSlate = mediatailor.SlateSource{
				SourceLocationName: plan.FillerSlate[0].SourceLocationName,
				VodSourceName:      plan.FillerSlate[0].VodSourceName,
			}
		}
		return &fillerSlate
	}
	return nil
}

// UPDATE CHANNEL

func updateTags(client *mediatailor.MediaTailor, arn *string, oldTagValue map[string]*string, newTagValue map[string]*string) error {
	var removedTags []string
	for k := range oldTagValue {
		if _, ok := (newTagValue)[k]; !ok {
			removedTags = append(removedTags, k)
		}
	}

	err := deleteTags(client, aws.StringValue(arn), removedTags)
	if err != nil {
		return err
	}

	if newTagValue != nil {
		var newTags = make(map[string]*string)
		for k, v := range newTagValue {
			val := v
			newTags[k] = val
		}
		tagInput := mediatailor.TagResourceInput{ResourceArn: arn, Tags: newTags}
		_, err := client.TagResource(&tagInput)
		if err != nil {
			return err
		}
	}
	return nil
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

func getUpdateChannelInput(ctx context.Context, req resource.UpdateRequest, resp resource.UpdateResponse) mediatailor.UpdateChannelInput {
	var plan resourceChannelModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return mediatailor.UpdateChannelInput{}
	}

	var params mediatailor.UpdateChannelInput

	params.ChannelName = plan.Name

	outputs := getOutputsUpdate(ctx, req, resp)
	if outputs != nil {
		params.Outputs = outputs
	}
	fillerSlate := getFillerSlateUpdate(ctx, req, resp)
	if fillerSlate != nil {
		params.FillerSlate = fillerSlate
	}

	return params
}

// SET STATE

func setStateUpdate(ctx context.Context, req resource.UpdateRequest, resp resource.UpdateResponse, channel mediatailor.UpdateChannelOutput) {
	var state resourceChannelModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.Arn = types.StringValue(*channel.Arn)
	state.Name = channel.ChannelName
	state.CreationTime = types.StringValue((channel.CreationTime).String())
	if state.FillerSlate != nil && len(state.FillerSlate) > 0 {
		state.FillerSlate = append(state.FillerSlate, resourceChannelFillerSlatesModel{
			SourceLocationName: channel.FillerSlate.SourceLocationName,
			VodSourceName:      channel.FillerSlate.VodSourceName,
		})
	}
	state.LastModifiedTime = types.StringValue((channel.LastModifiedTime).String())
	state.Outputs[0].PlaybackUrl = types.StringValue(*channel.Outputs[0].PlaybackUrl)
	if state.Outputs[0].DashPlaylistSettings != nil && len(state.Outputs[0].DashPlaylistSettings) > 0 {
		state.Outputs[0].DashPlaylistSettings = append(state.Outputs[0].DashPlaylistSettings, resourceChannelDashPlaylistSettingsModel{
			ManifestWindowsSeconds:            channel.Outputs[0].DashPlaylistSettings.ManifestWindowSeconds,
			MinBufferTimeSeconds:              channel.Outputs[0].DashPlaylistSettings.MinBufferTimeSeconds,
			MinUpdatePeriodSeconds:            channel.Outputs[0].DashPlaylistSettings.MinUpdatePeriodSeconds,
			SuggestedPresentationDelaySeconds: channel.Outputs[0].DashPlaylistSettings.SuggestedPresentationDelaySeconds,
		})
	}
	if state.Outputs[0].HlsPlaylistSettings != nil && len(state.Outputs[0].HlsPlaylistSettings) > 0 {
		state.Outputs[0].HlsPlaylistSettings = append(state.Outputs[0].HlsPlaylistSettings, resourceChannelHlsPlaylistSettingsModel{
			ManifestWindowsSeconds: channel.Outputs[0].HlsPlaylistSettings.ManifestWindowSeconds,
		})
	}
	if state.Outputs[0].ManifestName != nil {
		state.Outputs[0].ManifestName = channel.Outputs[0].ManifestName
	}
	if state.Outputs[0].SourceGroup != nil {
		state.Outputs[0].SourceGroup = channel.Outputs[0].SourceGroup
	}
	state.PlaybackMode = channel.PlaybackMode
	state.Tags = channel.Tags
	state.Tier = channel.Tier
}
