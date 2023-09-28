package awsmt

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"
	"time"
)

// CHANNEL

func newChannelInputBuilder(channelName *string, outputs []outputsModel, fillerSlate *fillerSlateModel) (*string, []*mediatailor.RequestOutputItem, *mediatailor.SlateSource) {
	theChannelName := channelName
	output := getOutputsFromPlan(outputs)
	fillerSlates := getFillerSlateFromPlan(fillerSlate)
	return theChannelName, output, fillerSlates
}

func channelInput(plan channelModel) mediatailor.CreateChannelInput {
	var input mediatailor.CreateChannelInput

	input.ChannelName, input.Outputs, input.FillerSlate = newChannelInputBuilder(plan.ChannelName, plan.Outputs, plan.FillerSlate)

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

func readChannelToPlan(plan channelModel, channel mediatailor.CreateChannelOutput) channelModel {

	plan = readChannelComputedValuesToPlan(plan, channel.Arn, channel.ChannelName, channel.CreationTime, channel.LastModifiedTime)

	plan = readFillerSlateToPlan(plan, channel.FillerSlate)

	plan = readOutputsToPlan(plan, channel.Outputs)

	plan = readOptionalValuesToPlan(plan, channel.PlaybackMode, channel.Tags, channel.Tier)

	return plan
}

func readChannelToState(state channelModel, channel mediatailor.DescribeChannelOutput) channelModel {

	state = readChannelComputedValuesToPlan(state, channel.Arn, channel.ChannelName, channel.CreationTime, channel.LastModifiedTime)

	state = readFillerSlateToPlan(state, channel.FillerSlate)

	state = readOutputsToPlan(state, channel.Outputs)

	state = readOptionalValuesToPlan(state, channel.PlaybackMode, channel.Tags, channel.Tier)

	return state
}

// POLICY
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

// UPDATE CHANNEL
func getUpdateChannelInput(plan channelModel) mediatailor.UpdateChannelInput {
	var input mediatailor.UpdateChannelInput
	input.ChannelName, input.Outputs, input.FillerSlate = newChannelInputBuilder(plan.ChannelName, plan.Outputs, plan.FillerSlate)
	return input
}

// GET OUTPUTS FROM PLAN
func getOutputsFromPlan(outputsFromPlan []outputsModel) []*mediatailor.RequestOutputItem {
	var outputFromPlan []*mediatailor.RequestOutputItem

	for _, output := range outputsFromPlan {
		outputs := &mediatailor.RequestOutputItem{}

		if output.DashPlaylistSettings != nil {
			outputs.DashPlaylistSettings = getDashPlaylistSettings(output.DashPlaylistSettings)
		}

		if output.HlsPlaylistSettings != nil {
			outputs.HlsPlaylistSettings = getHLSPlaylistSettings(output.HlsPlaylistSettings)
		}

		if output.ManifestName != nil {
			outputs.ManifestName = output.ManifestName
		}
		if output.SourceGroup != nil {
			outputs.SourceGroup = output.SourceGroup
		}

		outputFromPlan = append(outputFromPlan, outputs)
	}

	return outputFromPlan
}

func getDashPlaylistSettings(settings *dashPlaylistSettingsModel) *mediatailor.DashPlaylistSettings {
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

func getHLSPlaylistSettings(settings *hlsPlaylistSettingsModel) *mediatailor.HlsPlaylistSettings {
	hlsSettings := &mediatailor.HlsPlaylistSettings{}
	if settings.AdMarkupType != nil && len(settings.AdMarkupType) > 0 {
		for _, value := range settings.AdMarkupType {
			temp := value
			hlsSettings.AdMarkupType = append(hlsSettings.AdMarkupType, temp)
		}
	} else if settings.AdMarkupType == nil {
		temp := "DATERANGE"
		hlsSettings.AdMarkupType = append(hlsSettings.AdMarkupType, &temp)
	}
	if settings.ManifestWindowSeconds != nil {
		hlsSettings.ManifestWindowSeconds = settings.ManifestWindowSeconds
	}
	return hlsSettings
}

func getFillerSlateFromPlan(fillerSlate *fillerSlateModel) *mediatailor.SlateSource {
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

// READ COMPUTED VALUES TO PLAN
func readChannelComputedValuesToPlan(plan channelModel, arn *string, channelName *string, creationTime *time.Time, lastModifiedTime *time.Time) channelModel {
	plan.ID = types.StringValue(*channelName)

	if arn != nil {
		plan.Arn = types.StringValue(*arn)
	}

	plan.ChannelName = channelName

	if creationTime != nil {
		plan.CreationTime = types.StringValue((aws.TimeValue(creationTime)).String())
	}

	if lastModifiedTime != nil {
		plan.LastModifiedTime = types.StringValue((aws.TimeValue(lastModifiedTime)).String())
	}

	return plan
}

// READ FILLER SLATE TO PLAN
func readFillerSlateToPlan(plan channelModel, channel *mediatailor.SlateSource) channelModel {
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

// READ OUTPUTS TO PLAN
func readOutputsToPlan(plan channelModel, channel []*mediatailor.ResponseOutputItem) channelModel {
	if channel != nil {
		plan.Outputs = []outputsModel{}
		for _, output := range channel {
			outputs := outputsModel{}
			if output.DashPlaylistSettings != nil {
				outputs.DashPlaylistSettings = &dashPlaylistSettingsModel{}
				dashPlaylistSettings := readDashPlaylistConfigurationsToPlan(output)
				outputs.DashPlaylistSettings = dashPlaylistSettings

			}
			if output.HlsPlaylistSettings != nil {
				outputs.HlsPlaylistSettings = &hlsPlaylistSettingsModel{}
				hlsPlaylistSettings := readHlsPlaylistConfigurationsToPlan(output)
				outputs.HlsPlaylistSettings = hlsPlaylistSettings

			}
			outputs.ManifestName, outputs.PlaybackUrl, outputs.SourceGroup = readRMPS(output.ManifestName, output.PlaybackUrl, output.SourceGroup)

			plan.Outputs = append(plan.Outputs, outputs)
		}
	}
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

func readHlsPlaylistConfigurationsToPlan(output *mediatailor.ResponseOutputItem) *hlsPlaylistSettingsModel {
	outputs := &hlsPlaylistSettingsModel{}
	if output.HlsPlaylistSettings.AdMarkupType != nil && len(output.HlsPlaylistSettings.AdMarkupType) > 0 {
		for _, value := range output.HlsPlaylistSettings.AdMarkupType {
			temp := value
			outputs.AdMarkupType = append(outputs.AdMarkupType, temp)
		}
	}
	if output.HlsPlaylistSettings.ManifestWindowSeconds != nil {
		outputs.ManifestWindowSeconds = output.HlsPlaylistSettings.ManifestWindowSeconds
	}
	return outputs
}

// READ OPTIONAL VALUES TO PLAN
func readOptionalValuesToPlan(plan channelModel, playbackMode *string, tags map[string]*string, tier *string) channelModel {
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

func stopChannel(state *string, channelName *string, client *mediatailor.MediaTailor) error {
	if *state == "RUNNING" {
		_, err := client.StopChannel(&mediatailor.StopChannelInput{ChannelName: channelName})
		if err != nil {
			return err
		}
	}
	return nil
}

func updatePolicy(plan *channelModel, channelName *string, oldPolicy *mediatailor.GetChannelPolicyOutput, newPolicy jsontypes.Normalized, client *mediatailor.MediaTailor) (channelModel, error) {
	if !reflect.DeepEqual(oldPolicy, newPolicy) {
		if !newPolicy.IsNull() {
			plan.Policy = jsontypes.NewNormalizedPointerValue(aws.String(newPolicy.String()))
			_, err := client.PutChannelPolicy(&mediatailor.PutChannelPolicyInput{ChannelName: channelName, Policy: aws.String(newPolicy.String())})
			if err != nil {
				return *plan, err
			}
		} else {
			plan.Policy = jsontypes.NewNormalizedNull()
			_, err := client.DeleteChannelPolicy(&mediatailor.DeleteChannelPolicyInput{ChannelName: channelName})
			if err != nil {
				return *plan, err
			}
		}
	} else {
		plan.Policy = jsontypes.NewNormalizedPointerValue(oldPolicy.Policy)
	}
	return *plan, nil
}
