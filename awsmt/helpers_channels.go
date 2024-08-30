package awsmt

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/mediatailor/types"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"
	"strings"
	"time"
)

// functions to create MediaTailor inputs

func getCreateChannelInput(model channelModel) *mediatailor.CreateChannelInput {
	var input mediatailor.CreateChannelInput

	input.ChannelName, input.FillerSlate, input.Outputs = getSharedChannelInput(&model)

	if model.PlaybackMode != nil {
		var mode awsTypes.PlaybackMode
		switch *model.PlaybackMode {
		case "LINEAR":
			mode = awsTypes.PlaybackModeLinear
		default:
			mode = awsTypes.PlaybackModeLoop
		}
		input.PlaybackMode = mode
	}

	if model.Tags != nil && len(model.Tags) > 0 {
		input.Tags = model.Tags
	}

	if model.Tier != nil {
		var tier awsTypes.Tier
		switch *model.Tier {
		case "BASIC":
			tier = awsTypes.TierBasic
		default:
			tier = awsTypes.TierStandard
		}
		input.Tier = tier
	}

	return &input
}

func getUpdateChannelInput(model channelModel) *mediatailor.UpdateChannelInput {
	var input mediatailor.UpdateChannelInput

	input.ChannelName, input.FillerSlate, input.Outputs = getSharedChannelInput(&model)

	return &input
}

func getSharedChannelInput(model *channelModel) (name *string, source *awsTypes.SlateSource, outputItem []awsTypes.RequestOutputItem) {
	return model.Name, buildSlateSource(model), buildRequestOutputs(model)
}

func buildSlateSource(model *channelModel) *awsTypes.SlateSource {
	if model.FillerSlate == nil {
		return nil
	}
	temp := &awsTypes.SlateSource{}
	if model.FillerSlate.SourceLocationName != nil {
		temp.SourceLocationName = model.FillerSlate.SourceLocationName
	}
	if model.FillerSlate.VodSourceName != nil {
		temp.VodSourceName = model.FillerSlate.VodSourceName
	}
	return temp
}

func buildRequestOutputs(model *channelModel) []awsTypes.RequestOutputItem {
	var temp []awsTypes.RequestOutputItem

	for _, o := range model.Outputs {
		output := awsTypes.RequestOutputItem{}

		if o.DashPlaylistSettings != nil {
			output.DashPlaylistSettings = buildDashPlaylistSettings(o.DashPlaylistSettings)
		}

		if o.HlsPlaylistSettings != nil {
			output.HlsPlaylistSettings = buildHLSPlaylistSettings(o.HlsPlaylistSettings)
		}

		if o.ManifestName != nil {
			output.ManifestName = o.ManifestName
		}

		if o.SourceGroup != nil {
			output.SourceGroup = o.SourceGroup
		}

		temp = append(temp, output)
	}

	return temp
}

func buildDashPlaylistSettings(settings *dashPlaylistSettingsModel) *awsTypes.DashPlaylistSettings {
	dashSettings := &awsTypes.DashPlaylistSettings{}
	if settings.ManifestWindowSeconds != nil {
		manifestWindowSeconds := int32(*settings.ManifestWindowSeconds)
		dashSettings.ManifestWindowSeconds = &manifestWindowSeconds
	}
	if settings.MinBufferTimeSeconds != nil {
		minBufferTimeSeconds := int32(*settings.MinBufferTimeSeconds)
		dashSettings.MinBufferTimeSeconds = &minBufferTimeSeconds
	}
	if settings.MinUpdatePeriodSeconds != nil {
		minUpdatePeriodSeconds := int32(*settings.MinUpdatePeriodSeconds)
		dashSettings.MinUpdatePeriodSeconds = &minUpdatePeriodSeconds
	}
	if settings.SuggestedPresentationDelaySeconds != nil {
		suggestedPresentationDelaySeconds := int32(*settings.SuggestedPresentationDelaySeconds)
		dashSettings.SuggestedPresentationDelaySeconds = &suggestedPresentationDelaySeconds
	}

	return dashSettings
}

func buildHLSPlaylistSettings(settings *hlsPlaylistSettingsModel) *awsTypes.HlsPlaylistSettings {
	hlsSettings := &awsTypes.HlsPlaylistSettings{}

	if settings.AdMarkupType != nil && len(settings.AdMarkupType) > 0 {
		var adMarkupType []awsTypes.AdMarkupType
		for _, a := range settings.AdMarkupType {
			switch *a {
			case "SCTE35_ENHANCED":
				adMarkupType = append(adMarkupType, awsTypes.AdMarkupTypeScte35Enhanced)
			default:
				adMarkupType = append(adMarkupType, awsTypes.AdMarkupTypeDaterange)
			}
		}
		hlsSettings.AdMarkupType = append(hlsSettings.AdMarkupType, adMarkupType...)
	} else if settings.AdMarkupType == nil {
		hlsSettings.AdMarkupType = append(hlsSettings.AdMarkupType, awsTypes.AdMarkupTypeDaterange)
	}

	if settings.ManifestWindowSeconds != nil {
		manifestWindowSeconds := int32(*settings.ManifestWindowSeconds)
		hlsSettings.ManifestWindowSeconds = &manifestWindowSeconds
	}

	return hlsSettings
}

// functions to manipulate a channel once it is created

func createChannelPolicy(channelName *string, policy *string, client *mediatailor.Client) error {
	putChannelPolicyParams := mediatailor.PutChannelPolicyInput{
		ChannelName: channelName,
		Policy:      policy,
	}
	_, err := client.PutChannelPolicy(context.TODO(), &putChannelPolicyParams)
	if err != nil {
		return err
	}
	return err
}

func stopChannel(state awsTypes.ChannelState, channelName *string, client *mediatailor.Client) error {
	if state == awsTypes.ChannelStateRunning {
		_, err := client.StopChannel(context.TODO(), &mediatailor.StopChannelInput{ChannelName: channelName})
		if err != nil {
			return err
		}
	}
	return nil
}

func handlePolicyUpdate(context context.Context, client *mediatailor.Client, plan channelModel) error {
	var normalizedOldPolicy jsontypes.Normalized

	oldPolicy, err := client.GetChannelPolicy(context, &mediatailor.GetChannelPolicyInput{ChannelName: plan.Name})
	if err != nil && !strings.Contains(err.Error(), "NotFound") {
		return fmt.Errorf("error getting policy %v", err)
	}

	if oldPolicy != nil && oldPolicy.Policy != nil {
		normalizedOldPolicy = jsontypes.NewNormalizedPointerValue(oldPolicy.Policy)
	} else {
		normalizedOldPolicy = jsontypes.NewNormalizedNull()
	}

	plan, err = updatePolicy(&plan, plan.Name, normalizedOldPolicy, plan.Policy, client)
	if err != nil {
		return fmt.Errorf("error updating policy %v", err)
	}
	return nil
}

func updatePolicy(model *channelModel, channelName *string, oldPolicy jsontypes.Normalized, newPolicy jsontypes.Normalized, client *mediatailor.Client) (channelModel, error) {
	if !reflect.DeepEqual(oldPolicy, newPolicy) {
		if !newPolicy.IsNull() {
			model.Policy = newPolicy
			policy := newPolicy.ValueString()
			_, err := client.PutChannelPolicy(context.TODO(), &mediatailor.PutChannelPolicyInput{ChannelName: channelName, Policy: &policy})
			if err != nil {
				return *model, err
			}
		} else if newPolicy.IsNull() {
			model.Policy = jsontypes.NewNormalizedNull()
			_, err := client.DeleteChannelPolicy(context.TODO(), &mediatailor.DeleteChannelPolicyInput{ChannelName: channelName})
			if err != nil {
				return *model, err
			}
		}
	} else {
		model.Policy = oldPolicy
	}
	return *model, nil
}

// Functions used to read MediaTailor resources to plan and state

func readChannelComputedValues(model channelModel, arn *string, channelName *string, creationTime *time.Time, lastModifiedTime *time.Time) channelModel {
	model.ID = types.StringValue(*channelName)

	if arn != nil {
		model.Arn = types.StringValue(*arn)
	}

	model.Name = channelName

	if creationTime != nil {
		model.CreationTime = types.StringValue(creationTime.String())
	}

	if lastModifiedTime != nil {
		model.LastModifiedTime = types.StringValue(lastModifiedTime.String())
	}

	return model
}

func readFillerSlate(plan channelModel, fillerSlate *awsTypes.SlateSource) channelModel {
	if fillerSlate != nil {
		plan.FillerSlate = &fillerSlateModel{}
		if fillerSlate.SourceLocationName != nil {
			plan.FillerSlate.SourceLocationName = fillerSlate.SourceLocationName
		}
		if fillerSlate.VodSourceName != nil {
			plan.FillerSlate.VodSourceName = fillerSlate.VodSourceName
		}
	}
	return plan
}

func readOutputs(plan channelModel, responseOutputItems []awsTypes.ResponseOutputItem) channelModel {

	if responseOutputItems == nil {
		return plan
	}

	var tempOutputs []outputsModel
	for i, output := range responseOutputItems {
		outputs := outputsModel{}
		if output.DashPlaylistSettings != nil {
			outputs.DashPlaylistSettings = readDashPlaylistConfigurationsToPlan(&output)
		}
		if output.HlsPlaylistSettings != nil {
			if len(plan.Outputs) > 0 && i <= len(plan.Outputs) {
				outputs.HlsPlaylistSettings = readHlsPlaylistConfigurationsToPlan(&output, plan.Outputs[i])
			} else {
				outputs.HlsPlaylistSettings = readHlsPlaylistConfigurationsToPlanDS(&output)
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

func readDashPlaylistConfigurationsToPlan(output *awsTypes.ResponseOutputItem) *dashPlaylistSettingsModel {
	outputs := &dashPlaylistSettingsModel{}
	if output.DashPlaylistSettings.ManifestWindowSeconds != nil {
		manifestWindowSeconds := int64(*output.DashPlaylistSettings.ManifestWindowSeconds)
		outputs.ManifestWindowSeconds = &manifestWindowSeconds
	}
	if output.DashPlaylistSettings.MinBufferTimeSeconds != nil {
		minBufferTimeSeconds := int64(*output.DashPlaylistSettings.MinBufferTimeSeconds)
		outputs.MinBufferTimeSeconds = &minBufferTimeSeconds
	}
	if output.DashPlaylistSettings.MinUpdatePeriodSeconds != nil {
		minUpdatePeriodSeconds := int64(*output.DashPlaylistSettings.MinUpdatePeriodSeconds)
		outputs.MinUpdatePeriodSeconds = &minUpdatePeriodSeconds
	}
	if output.DashPlaylistSettings.SuggestedPresentationDelaySeconds != nil {
		suggestedPresentationDelaySeconds := int64(*output.DashPlaylistSettings.SuggestedPresentationDelaySeconds)
		outputs.SuggestedPresentationDelaySeconds = &suggestedPresentationDelaySeconds
	}
	return outputs
}

func readHlsPlaylistConfigurationsToPlan(output *awsTypes.ResponseOutputItem, stateOutput outputsModel) *hlsPlaylistSettingsModel {
	outputs := &hlsPlaylistSettingsModel{}
	if stateOutput.HlsPlaylistSettings.AdMarkupType != nil && output.HlsPlaylistSettings.AdMarkupType != nil && len(output.HlsPlaylistSettings.AdMarkupType) > 0 {
		var adMarkupTypes []*string
		for _, a := range output.HlsPlaylistSettings.AdMarkupType {
			adMarkupType := string(a)
			adMarkupTypes = append(adMarkupTypes, &adMarkupType)
		}
		outputs.AdMarkupType = append(outputs.AdMarkupType, adMarkupTypes...)
	}
	if stateOutput.HlsPlaylistSettings.ManifestWindowSeconds != nil && output.HlsPlaylistSettings.ManifestWindowSeconds != nil {
		manifestWindowSeconds := int64(*output.HlsPlaylistSettings.ManifestWindowSeconds)
		outputs.ManifestWindowSeconds = &manifestWindowSeconds
	}
	return outputs
}

func readHlsPlaylistConfigurationsToPlanDS(output *awsTypes.ResponseOutputItem) *hlsPlaylistSettingsModel {
	outputs := &hlsPlaylistSettingsModel{}
	if output.HlsPlaylistSettings.AdMarkupType != nil && len(output.HlsPlaylistSettings.AdMarkupType) > 0 {
		var adMarkupTypes []*string
		for _, a := range output.HlsPlaylistSettings.AdMarkupType {
			adMarkupType := string(a)
			adMarkupTypes = append(adMarkupTypes, &adMarkupType)
		}
		outputs.AdMarkupType = append(outputs.AdMarkupType, adMarkupTypes...)
	}
	if output.HlsPlaylistSettings.ManifestWindowSeconds != nil {
		manifestWindowSeconds := int64(*output.HlsPlaylistSettings.ManifestWindowSeconds)
		outputs.ManifestWindowSeconds = &manifestWindowSeconds
	}
	return outputs
}

func readOptionalValues(plan channelModel, playbackMode *string, tags map[string]string, tier *string) channelModel {
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

func writeChannelToPlan(model channelModel, channel mediatailor.CreateChannelOutput) channelModel {

	model = readChannelComputedValues(model, channel.Arn, channel.ChannelName, channel.CreationTime, channel.LastModifiedTime)

	model = readFillerSlate(model, channel.FillerSlate)

	model = readOutputs(model, channel.Outputs)

	model = readOptionalValues(model, channel.PlaybackMode, channel.Tags, channel.Tier)

	return model
}

func writeChannelToState(model channelModel, channel mediatailor.DescribeChannelOutput) channelModel {

	model = readChannelComputedValues(model, channel.Arn, channel.ChannelName, channel.CreationTime, channel.LastModifiedTime)

	model = readFillerSlate(model, channel.FillerSlate)

	model = readOutputs(model, channel.Outputs)

	model = readOptionalValues(model, channel.PlaybackMode, channel.Tags, channel.Tier)

	return model
}

// helper functions to simplify update function logic
func shouldStartChannel(previousState awsTypes.ChannelState, newState *string) bool {
	wasRunning := previousState == awsTypes.ChannelStateRunning
	shouldRun := newState != nil && *newState == "RUNNING"
	return (newState == nil && wasRunning) || shouldRun
}
