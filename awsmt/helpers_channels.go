package awsmt

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CHANNEL

func channelInput(plan resourceChannelModel) mediatailor.CreateChannelInput {
	var input mediatailor.CreateChannelInput

	input.ChannelName = plan.ChannelName

	if plan.FillerSlate != nil {
		input.FillerSlate = &mediatailor.SlateSource{}
		if plan.FillerSlate.SourceLocationName != nil {
			input.FillerSlate.SourceLocationName = plan.FillerSlate.SourceLocationName
		}
		if plan.FillerSlate.VodSourceName != nil {
			input.FillerSlate.VodSourceName = plan.FillerSlate.VodSourceName
		}
	}
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
			input.Outputs = append(input.Outputs, outputs)
		}
	}

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
	plan.ID = types.StringValue(*channel.ChannelName)

	if channel.Arn != nil {
		plan.Arn = types.StringValue(*channel.Arn)
	}

	if channel.ChannelName != nil {
		plan.ChannelName = channel.ChannelName
	}

	if channel.CreationTime != nil {
		plan.CreationTime = types.StringValue((aws.TimeValue(channel.CreationTime)).String())
	}

	if channel.FillerSlate != nil {
		plan.FillerSlate = &fillerSlateRModel{}
		if channel.FillerSlate.SourceLocationName != nil {
			plan.FillerSlate.SourceLocationName = channel.FillerSlate.SourceLocationName
		}
		if channel.FillerSlate.VodSourceName != nil {
			plan.FillerSlate.VodSourceName = channel.FillerSlate.VodSourceName
		}
	}

	if channel.LastModifiedTime != nil {
		plan.LastModifiedTime = types.StringValue((aws.TimeValue(channel.LastModifiedTime)).String())
	}

	if channel.Outputs != nil {
		plan.Outputs = []outputsRModel{}
		for _, output := range channel.Outputs {
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

	if channel.PlaybackMode != nil {
		plan.PlaybackMode = channel.PlaybackMode
	}

	if channel.Tags != nil && len(channel.Tags) > 0 {
		plan.Tags = channel.Tags
	}

	if channel.Tier != nil {
		plan.Tier = channel.Tier
	}
	plan.ID = types.StringValue(*channel.ChannelName)

	return plan
}

func readChannelToState(state resourceChannelModel, channel mediatailor.DescribeChannelOutput) resourceChannelModel {
	state.ID = types.StringValue(*channel.ChannelName)

	if channel.Arn != nil {
		state.Arn = types.StringValue(*channel.Arn)
	}

	if channel.ChannelName != nil {
		state.ChannelName = channel.ChannelName
	}

	if channel.CreationTime != nil {
		state.CreationTime = types.StringValue((aws.TimeValue(channel.CreationTime)).String())
	}

	if channel.FillerSlate != nil {
		state.FillerSlate = &fillerSlateRModel{}
		if channel.FillerSlate.SourceLocationName != nil {
			state.FillerSlate.SourceLocationName = channel.FillerSlate.SourceLocationName
		}
		if channel.FillerSlate.VodSourceName != nil {
			state.FillerSlate.VodSourceName = channel.FillerSlate.VodSourceName
		}
	}

	if channel.LastModifiedTime != nil {
		state.LastModifiedTime = types.StringValue((aws.TimeValue(channel.LastModifiedTime)).String())
	}

	if channel.Outputs != nil {
		state.Outputs = []outputsRModel{}
		for _, output := range channel.Outputs {
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
			state.Outputs = append(state.Outputs, outputs)
		}
	}

	if channel.PlaybackMode != nil {
		state.PlaybackMode = channel.PlaybackMode
	}

	if channel.Tags != nil && len(channel.Tags) > 0 {
		state.Tags = channel.Tags
	}

	if channel.Tier != nil {
		state.Tier = channel.Tier
	}

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
			input.Outputs = append(input.Outputs, outputs)
		}
	}
	if plan.FillerSlate != nil {
		input.FillerSlate = &mediatailor.SlateSource{}
		if plan.FillerSlate.SourceLocationName != nil {
			input.FillerSlate.SourceLocationName = plan.FillerSlate.SourceLocationName
		}
		if plan.FillerSlate.VodSourceName != nil {
			input.FillerSlate.VodSourceName = plan.FillerSlate.VodSourceName
		}
	}

	return input
}
