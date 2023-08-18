package awsmt

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

// CHANNEL

func channelInput(plan resourceChannelModel) mediatailor.CreateChannelInput {
	var input mediatailor.CreateChannelInput

	input.ChannelName = plan.ChannelName

	input.Outputs = getOutputsFromPlan(plan)

	input.FillerSlate = getFillerSlateFromPlan(plan)

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

func readChannelToPlan(plan resourceChannelModel, channel mediatailor.CreateChannelOutput) resourceChannelModel {

	plan = readChannelComputedValuesToPlan(plan, channel.Arn, channel.ChannelName, channel.CreationTime, channel.LastModifiedTime)

	plan = readFillerSlateToPlan(plan, channel.FillerSlate)

	plan = readOutputsToPlan(plan, channel.Outputs)

	plan = readOptionalValuesToPlan(plan, channel.PlaybackMode, channel.Tags, channel.Tier)

	return plan
}

func readChannelToData(data dataSourceChannelModel, channel mediatailor.DescribeChannelOutput) dataSourceChannelModel {
	data.ID = types.StringValue(*channel.ChannelName)
	if channel.Arn != nil {
		data.Arn = types.StringValue(*channel.Arn)
	}

	if channel.ChannelName != nil {
		data.ChannelName = channel.ChannelName
	}

	if channel.ChannelState != nil {
		data.ChannelState = types.StringValue(*channel.ChannelState)
	}

	if channel.CreationTime != nil {
		data.CreationTime = types.StringValue((aws.TimeValue(channel.CreationTime)).String())
	}

	if channel.FillerSlate != nil {
		data.FillerSlate = &fillerSlateDSModel{}
		if channel.FillerSlate.SourceLocationName != nil {
			data.FillerSlate.SourceLocationName = types.StringValue(*channel.FillerSlate.SourceLocationName)
		}
		if channel.FillerSlate.VodSourceName != nil {
			data.FillerSlate.VodSourceName = types.StringValue(*channel.FillerSlate.VodSourceName)
		}
	}

	if channel.LastModifiedTime != nil {
		data.LastModifiedTime = types.StringValue((aws.TimeValue(channel.LastModifiedTime)).String())
	}

	if channel.Outputs != nil {
		data.Outputs = []outputsDSModel{}
		for _, output := range channel.Outputs {
			outputs := outputsDSModel{}
			if output.DashPlaylistSettings != nil {
				outputs.DashPlaylistSettings = &dashPlaylistSettingsDSModel{}
				if output.DashPlaylistSettings.ManifestWindowSeconds != nil {
					outputs.DashPlaylistSettings.ManifestWindowSeconds = types.Int64Value(*output.DashPlaylistSettings.ManifestWindowSeconds)
				}
				if output.DashPlaylistSettings.MinBufferTimeSeconds != nil {
					outputs.DashPlaylistSettings.MinBufferTimeSeconds = types.Int64Value(*output.DashPlaylistSettings.MinBufferTimeSeconds)
				}
				if output.DashPlaylistSettings.MinUpdatePeriodSeconds != nil {
					outputs.DashPlaylistSettings.MinUpdatePeriodSeconds = types.Int64Value(*output.DashPlaylistSettings.MinUpdatePeriodSeconds)
				}
				if output.DashPlaylistSettings.SuggestedPresentationDelaySeconds != nil {
					outputs.DashPlaylistSettings.SuggestedPresentationDelaySeconds = types.Int64Value(*output.DashPlaylistSettings.SuggestedPresentationDelaySeconds)
				}
			}
			if output.HlsPlaylistSettings != nil {
				outputs.HlsPlaylistSettings = &hlsPlaylistSettingsDSModel{}
				if output.HlsPlaylistSettings.AdMarkupType != nil && len(output.HlsPlaylistSettings.AdMarkupType) > 0 {
					outputs.HlsPlaylistSettings.AdMarkupType = []types.String{}
					output.HlsPlaylistSettings.AdMarkupType = append(output.HlsPlaylistSettings.AdMarkupType, output.HlsPlaylistSettings.AdMarkupType...)
				}
				if output.HlsPlaylistSettings.ManifestWindowSeconds != nil {
					outputs.HlsPlaylistSettings.ManifestWindowSeconds = types.Int64Value(*output.HlsPlaylistSettings.ManifestWindowSeconds)
				}
			}
			if output.ManifestName != nil {
				outputs.ManifestName = types.StringValue(*output.ManifestName)
			}
			if output.PlaybackUrl != nil {
				outputs.PlaybackUrl = types.StringValue(*output.PlaybackUrl)
			}
			if output.SourceGroup != nil {
				outputs.SourceGroup = types.StringValue(*output.SourceGroup)
			}
			data.Outputs = append(data.Outputs, outputs)
		}
	}

	if channel.PlaybackMode != nil {
		data.PlaybackMode = types.StringValue(*channel.PlaybackMode)
	}

	if channel.Tags != nil && len(channel.Tags) > 0 {
		data.Tags = make(map[string]*string)
		for key, value := range channel.Tags {
			data.Tags[key] = value
		}
	}

	if channel.Tier != nil {
		data.Tier = types.StringValue(*channel.Tier)
	}

	return data
}

func readChannelToState(state resourceChannelModel, channel mediatailor.DescribeChannelOutput) resourceChannelModel {
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
func getUpdateChannelInput(plan resourceChannelModel) mediatailor.UpdateChannelInput {
	var input mediatailor.UpdateChannelInput

	input.ChannelName = plan.ChannelName

	input.Outputs = getOutputsFromPlan(plan)

	input.FillerSlate = getFillerSlateFromPlan(plan)

	return input
}

// GET OUTPUTS FROM PLAN
func getOutputsFromPlan(plan resourceChannelModel) []*mediatailor.RequestOutputItem {
	var outputFromPlan []*mediatailor.RequestOutputItem
	if plan.Outputs != nil && len(plan.Outputs) > 0 {
		for _, output := range plan.Outputs {
			outputs := &mediatailor.RequestOutputItem{}
			if output.DashPlaylistSettings != nil {
				outputs.DashPlaylistSettings = &mediatailor.DashPlaylistSettings{}
				if output.DashPlaylistSettings.ManifestWindowSeconds != nil {
					outputs.DashPlaylistSettings.ManifestWindowSeconds = output.DashPlaylistSettings.ManifestWindowSeconds
				}
				if output.DashPlaylistSettings.MinBufferTimeSeconds != nil {
					outputs.DashPlaylistSettings.MinBufferTimeSeconds = output.DashPlaylistSettings.MinBufferTimeSeconds
				}
				if output.DashPlaylistSettings.MinUpdatePeriodSeconds != nil {
					outputs.DashPlaylistSettings.MinUpdatePeriodSeconds = output.DashPlaylistSettings.MinUpdatePeriodSeconds
				}
				if output.DashPlaylistSettings.SuggestedPresentationDelaySeconds != nil {
					outputs.DashPlaylistSettings.SuggestedPresentationDelaySeconds = output.DashPlaylistSettings.SuggestedPresentationDelaySeconds
				}
			}
			if output.HlsPlaylistSettings != nil {
				outputs.HlsPlaylistSettings = &mediatailor.HlsPlaylistSettings{}
				if output.HlsPlaylistSettings.AdMarkupType != nil && len(output.HlsPlaylistSettings.AdMarkupType) > 0 {
					outputs.HlsPlaylistSettings.AdMarkupType = output.HlsPlaylistSettings.AdMarkupType
				}
				if output.HlsPlaylistSettings.ManifestWindowSeconds != nil {
					outputs.HlsPlaylistSettings.ManifestWindowSeconds = output.HlsPlaylistSettings.ManifestWindowSeconds
				}
			}
			if output.ManifestName != nil {
				outputs.ManifestName = output.ManifestName
			}
			if output.SourceGroup != nil {
				outputs.SourceGroup = output.SourceGroup
			}
			outputFromPlan = append(outputFromPlan, outputs)
		}
	}
	return outputFromPlan
}

func getFillerSlateFromPlan(plan resourceChannelModel) *mediatailor.SlateSource {
	var slateSource *mediatailor.SlateSource
	if plan.FillerSlate != nil {
		slateSource = &mediatailor.SlateSource{}
		if plan.FillerSlate.SourceLocationName != nil {
			slateSource.SourceLocationName = plan.FillerSlate.SourceLocationName
		}
		if plan.FillerSlate.VodSourceName != nil {
			slateSource.VodSourceName = plan.FillerSlate.VodSourceName
		}
	}
	return slateSource
}

// READ COMPUTED VALUES TO PLAN
func readChannelComputedValuesToPlan(plan resourceChannelModel, arn *string, channelName *string, creationTime *time.Time, lastModifiedTime *time.Time) resourceChannelModel {
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
func readFillerSlateToPlan(plan resourceChannelModel, channel *mediatailor.SlateSource) resourceChannelModel {
	if channel != nil {
		plan.FillerSlate = &fillerSlateRModel{}
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
func readOutputsToPlan(plan resourceChannelModel, channel []*mediatailor.ResponseOutputItem) resourceChannelModel {
	if channel != nil {
		plan.Outputs = []outputsRModel{}
		for _, output := range channel {
			outputs := outputsRModel{}
			if output.DashPlaylistSettings != nil {
				outputs.DashPlaylistSettings = &dashPlaylistSettingsRModel{}
				if output.DashPlaylistSettings.ManifestWindowSeconds != nil {
					outputs.DashPlaylistSettings.ManifestWindowSeconds = output.DashPlaylistSettings.ManifestWindowSeconds
				}
				if output.DashPlaylistSettings.MinBufferTimeSeconds != nil {
					outputs.DashPlaylistSettings.MinBufferTimeSeconds = output.DashPlaylistSettings.MinBufferTimeSeconds
				}
				if output.DashPlaylistSettings.MinUpdatePeriodSeconds != nil {
					outputs.DashPlaylistSettings.MinUpdatePeriodSeconds = output.DashPlaylistSettings.MinUpdatePeriodSeconds
				}
				if output.DashPlaylistSettings.SuggestedPresentationDelaySeconds != nil {
					outputs.DashPlaylistSettings.SuggestedPresentationDelaySeconds = output.DashPlaylistSettings.SuggestedPresentationDelaySeconds
				}
			}
			if output.HlsPlaylistSettings != nil {
				outputs.HlsPlaylistSettings = &hlsPlaylistSettingsRModel{}
				if output.HlsPlaylistSettings.AdMarkupType != nil && len(output.HlsPlaylistSettings.AdMarkupType) > 0 {
					output.HlsPlaylistSettings.AdMarkupType = append(output.HlsPlaylistSettings.AdMarkupType, output.HlsPlaylistSettings.AdMarkupType...)
				}
				if output.HlsPlaylistSettings.ManifestWindowSeconds != nil {
					outputs.HlsPlaylistSettings.ManifestWindowSeconds = output.HlsPlaylistSettings.ManifestWindowSeconds
				}
			}
			if output.ManifestName != nil {
				outputs.ManifestName = output.ManifestName
			}
			if output.PlaybackUrl != nil {
				outputs.PlaybackUrl = types.StringValue(*output.PlaybackUrl)
			}
			if output.SourceGroup != nil {
				outputs.SourceGroup = output.SourceGroup
			}
			plan.Outputs = append(plan.Outputs, outputs)
		}
	}
	return plan
}

// READ OPTIONAL VALUES TO PLAN
func readOptionalValuesToPlan(plan resourceChannelModel, playbackMode *string, tags map[string]*string, tier *string) resourceChannelModel {
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
