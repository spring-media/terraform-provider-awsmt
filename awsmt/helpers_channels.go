package awsmt

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"
	"time"
)

// Functions used to edit the channel once created

func createChannelPolicy(channelName *string, policy *string, client *mediatailor.MediaTailor) error {
	putChannelPolicyParams := mediatailor.PutChannelPolicyInput{
		ChannelName: channelName,
		Policy:      policy,
	}
	_, err := client.PutChannelPolicy(&putChannelPolicyParams)
	if err != nil {
		return err
	}
	return err
}

func stopChannel(state *string, channelName *string, client *mediatailor.MediaTailor) error {
	if *state == "RUNNING" {
		_, err := client.StopChannel(&mediatailor.StopChannelInput{ChannelName: channelName})
		if err != nil {
			return err
		}
	}
	return nil
}

func updatePolicy(plan *channelModel, channelName *string, oldPolicy jsontypes.Normalized, newPolicy jsontypes.Normalized, client *mediatailor.MediaTailor) (channelModel, error) {
	if !reflect.DeepEqual(oldPolicy, newPolicy) {
		if !newPolicy.IsNull() {
			plan.Policy = newPolicy
			policy := newPolicy.ValueString()
			_, err := client.PutChannelPolicy(&mediatailor.PutChannelPolicyInput{ChannelName: channelName, Policy: &policy})
			if err != nil {
				return *plan, err
			}
		} else if newPolicy.IsNull() {
			plan.Policy = jsontypes.NewNormalizedNull()
			_, err := client.DeleteChannelPolicy(&mediatailor.DeleteChannelPolicyInput{ChannelName: channelName})
			if err != nil {
				return *plan, err
			}
		}
	} else {
		plan.Policy = oldPolicy
	}
	return *plan, nil
}

// Functions used to create the resource in MediaTailor

func newChannelInputBuilder(channelName *string, outputs []outputsModel, fillerSlate *fillerSlateModel) (*string, []*mediatailor.RequestOutputItem, *mediatailor.SlateSource) {
	theChannelName := channelName
	output := buildRequestOutput(outputs)
	fillerSlates := buildSlateSource(fillerSlate)
	return theChannelName, output, fillerSlates
}

func buildChannelInput(plan channelModel) mediatailor.CreateChannelInput {
	var input mediatailor.CreateChannelInput

	input.ChannelName, input.Outputs, input.FillerSlate = newChannelInputBuilder(plan.Name, plan.Outputs, plan.FillerSlate)

	if plan.PlaybackMode != nil {
		input.PlaybackMode = plan.PlaybackMode
	}
	if plan.Tags != nil {
		input.Tags = plan.Tags
	}

	if plan.Tier != nil {
		input.Tier = plan.Tier
	}

	return input
}

func buildUpdateChannelInput(plan channelModel) mediatailor.UpdateChannelInput {
	var input mediatailor.UpdateChannelInput
	input.ChannelName, input.Outputs, input.FillerSlate = newChannelInputBuilder(plan.Name, plan.Outputs, plan.FillerSlate)
	return input
}

func buildRequestOutput(outputsFromPlan []outputsModel) []*mediatailor.RequestOutputItem {
	var res []*mediatailor.RequestOutputItem

	for _, output := range outputsFromPlan {
		outputs := &mediatailor.RequestOutputItem{}

		if output.DashPlaylistSettings != nil {
			outputs.DashPlaylistSettings = buildDashPlaylistSettings(output.DashPlaylistSettings)
		}

		if output.HlsPlaylistSettings != nil {
			outputs.HlsPlaylistSettings = buildHLSPlaylistSettings(output.HlsPlaylistSettings)
		}

		if output.ManifestName != nil {
			outputs.ManifestName = output.ManifestName
		}

		if output.SourceGroup != nil {
			outputs.SourceGroup = output.SourceGroup
		}

		res = append(res, outputs)
	}

	return res
}

func buildDashPlaylistSettings(settings *dashPlaylistSettingsModel) *mediatailor.DashPlaylistSettings {
	dashSettings := &mediatailor.DashPlaylistSettings{}
	if settings.ManifestWindowSeconds != nil {
		dashSettings.ManifestWindowSeconds = settings.ManifestWindowSeconds
	}
	if settings.MinBufferTimeSeconds != nil {
		dashSettings.MinBufferTimeSeconds = settings.MinBufferTimeSeconds
	}
	if settings.MinUpdatePeriodSeconds != nil {
		dashSettings.MinUpdatePeriodSeconds = settings.MinUpdatePeriodSeconds
	}
	if settings.SuggestedPresentationDelaySeconds != nil {
		dashSettings.SuggestedPresentationDelaySeconds = settings.SuggestedPresentationDelaySeconds
	}

	return dashSettings
}

func buildHLSPlaylistSettings(settings *hlsPlaylistSettingsModel) *mediatailor.HlsPlaylistSettings {
	hlsSettings := &mediatailor.HlsPlaylistSettings{}
	if settings.AdMarkupType != nil && len(settings.AdMarkupType) > 0 {
		hlsSettings.AdMarkupType = append(hlsSettings.AdMarkupType, settings.AdMarkupType...)
	} else if settings.AdMarkupType == nil {
		temp := "DATERANGE"
		hlsSettings.AdMarkupType = append(hlsSettings.AdMarkupType, &temp)
	}
	if settings.ManifestWindowSeconds != nil {
		hlsSettings.ManifestWindowSeconds = settings.ManifestWindowSeconds
	}
	return hlsSettings
}

func buildSlateSource(fillerSlate *fillerSlateModel) *mediatailor.SlateSource {
	var slateSource *mediatailor.SlateSource
	if fillerSlate != nil {
		slateSource = &mediatailor.SlateSource{}
		if fillerSlate.SourceLocationName != nil {
			slateSource.SourceLocationName = fillerSlate.SourceLocationName
		}
		if fillerSlate.VodSourceName != nil {
			slateSource.VodSourceName = fillerSlate.VodSourceName
		}
	}
	return slateSource
}

// Functions used to read MediaTailor resources to plan and state

func readChannelComputedValues(plan channelModel, arn *string, channelName *string, creationTime *time.Time, lastModifiedTime *time.Time) channelModel {
	plan.ID = types.StringValue(*channelName)

	if arn != nil {
		plan.Arn = types.StringValue(*arn)
	}

	plan.Name = channelName

	if creationTime != nil {
		plan.CreationTime = types.StringValue((aws.TimeValue(creationTime)).String())
	}

	if lastModifiedTime != nil {
		plan.LastModifiedTime = types.StringValue((aws.TimeValue(lastModifiedTime)).String())
	}

	return plan
}

func readFillerSlate(plan channelModel, channel *mediatailor.SlateSource) channelModel {
	if channel != nil {
		plan.FillerSlate = &fillerSlateModel{}
		if channel.SourceLocationName != nil {
			plan.FillerSlate.SourceLocationName = channel.SourceLocationName
		}
		if channel.VodSourceName != nil {
			plan.FillerSlate.VodSourceName = channel.VodSourceName
		}
	}
	return plan
}

func readOutputs(plan channelModel, channel []*mediatailor.ResponseOutputItem) channelModel {

	if channel == nil {
		return plan
	}

	var tempOutputs []outputsModel
	for i, output := range channel {
		outputs := outputsModel{}
		if output.DashPlaylistSettings != nil {
			outputs.DashPlaylistSettings = readDashPlaylistConfigurationsToPlan(output)
		}
		if output.HlsPlaylistSettings != nil {
			if len(plan.Outputs) > 0 && i <= len(plan.Outputs) {
				outputs.HlsPlaylistSettings = readHlsPlaylistConfigurationsToPlan(output, plan.Outputs[i])
			} else {
				outputs.HlsPlaylistSettings = readHlsPlaylistConfigurationsToPlanDS(output)
			}

		}
		outputs.ManifestName, outputs.PlaybackUrl, outputs.SourceGroup = readRMPS(output.ManifestName, output.PlaybackUrl, output.SourceGroup)
		tempOutputs = append(tempOutputs, outputs)
	}
	plan.Outputs = tempOutputs

	return plan
}

func readRMPS(manifestName *string, playbackUrl *string, sourceGroup *string) (*string, types.String, *string) {
	outputs := outputsModel{}
	if manifestName != nil {
		outputs.ManifestName = manifestName
	}
	if playbackUrl != nil {
		outputs.PlaybackUrl = types.StringValue(*playbackUrl)
	}
	if sourceGroup != nil {
		outputs.SourceGroup = sourceGroup
	}
	return outputs.ManifestName, outputs.PlaybackUrl, outputs.SourceGroup
}

func readDashPlaylistConfigurationsToPlan(output *mediatailor.ResponseOutputItem) *dashPlaylistSettingsModel {
	outputs := &dashPlaylistSettingsModel{}
	if output.DashPlaylistSettings.ManifestWindowSeconds != nil {
		outputs.ManifestWindowSeconds = output.DashPlaylistSettings.ManifestWindowSeconds
	}
	if output.DashPlaylistSettings.MinBufferTimeSeconds != nil {
		outputs.MinBufferTimeSeconds = output.DashPlaylistSettings.MinBufferTimeSeconds
	}
	if output.DashPlaylistSettings.MinUpdatePeriodSeconds != nil {
		outputs.MinUpdatePeriodSeconds = output.DashPlaylistSettings.MinUpdatePeriodSeconds
	}
	if output.DashPlaylistSettings.SuggestedPresentationDelaySeconds != nil {
		outputs.SuggestedPresentationDelaySeconds = output.DashPlaylistSettings.SuggestedPresentationDelaySeconds
	}
	return outputs
}

func readHlsPlaylistConfigurationsToPlan(output *mediatailor.ResponseOutputItem, stateOutput outputsModel) *hlsPlaylistSettingsModel {
	outputs := &hlsPlaylistSettingsModel{}
	if stateOutput.HlsPlaylistSettings.AdMarkupType != nil && output.HlsPlaylistSettings.AdMarkupType != nil && len(output.HlsPlaylistSettings.AdMarkupType) > 0 {
		outputs.AdMarkupType = append(outputs.AdMarkupType, output.HlsPlaylistSettings.AdMarkupType...)
	}
	if stateOutput.HlsPlaylistSettings.ManifestWindowSeconds != nil && output.HlsPlaylistSettings.ManifestWindowSeconds != nil {
		outputs.ManifestWindowSeconds = output.HlsPlaylistSettings.ManifestWindowSeconds
	}
	return outputs
}

func readHlsPlaylistConfigurationsToPlanDS(output *mediatailor.ResponseOutputItem) *hlsPlaylistSettingsModel {
	outputs := &hlsPlaylistSettingsModel{}
	if output.HlsPlaylistSettings.AdMarkupType != nil && len(output.HlsPlaylistSettings.AdMarkupType) > 0 {
		outputs.AdMarkupType = append(outputs.AdMarkupType, output.HlsPlaylistSettings.AdMarkupType...)
	}
	if output.HlsPlaylistSettings.ManifestWindowSeconds != nil {
		outputs.ManifestWindowSeconds = output.HlsPlaylistSettings.ManifestWindowSeconds
	}
	return outputs
}

func readOptionalValues(plan channelModel, playbackMode *string, tags map[string]*string, tier *string) channelModel {
	if playbackMode != nil {
		plan.PlaybackMode = playbackMode
	}

	if len(tags) > 0 {
		plan.Tags = tags
	}

	if tier != nil {
		plan.Tier = tier
	}
	return plan
}

func writeChannelToPlan(plan channelModel, channel mediatailor.CreateChannelOutput) channelModel {

	plan = readChannelComputedValues(plan, channel.Arn, channel.ChannelName, channel.CreationTime, channel.LastModifiedTime)

	plan = readFillerSlate(plan, channel.FillerSlate)

	plan = readOutputs(plan, channel.Outputs)

	plan = readOptionalValues(plan, channel.PlaybackMode, channel.Tags, channel.Tier)

	return plan
}

func writeChannelToState(state channelModel, channel mediatailor.DescribeChannelOutput) channelModel {

	state = readChannelComputedValues(state, channel.Arn, channel.ChannelName, channel.CreationTime, channel.LastModifiedTime)

	state = readFillerSlate(state, channel.FillerSlate)

	state = readOutputs(state, channel.Outputs)

	state = readOptionalValues(state, channel.PlaybackMode, channel.Tags, channel.Tier)

	return state
}
