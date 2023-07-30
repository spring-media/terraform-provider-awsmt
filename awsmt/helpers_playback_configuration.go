package awsmt

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func playbackConfigurationInput(plan resourcePlaybackConfigurationModel) mediatailor.PutPlaybackConfigurationInput {

	input := &mediatailor.PutPlaybackConfigurationInput{}

	input.AdDecisionServerUrl = plan.AdDecisionServerUrl

	if plan.AvailSupression != nil {
		if plan.AvailSupression.Mode != nil {
			input.AvailSuppression.Mode = plan.AvailSupression.Mode
		}
		if plan.AvailSupression.Value != nil {
			input.AvailSuppression.Value = plan.AvailSupression.Value
		}
	}

	if plan.Bumper != nil {
		if plan.Bumper.EndUrl != nil {
			input.Bumper.EndUrl = plan.Bumper.EndUrl
		}
		if plan.Bumper.StartUrl != nil {
			input.Bumper.StartUrl = plan.Bumper.StartUrl
		}
	}

	if plan.CdnConfiguration != nil {
		if plan.CdnConfiguration.AdSegmentUrlPrefix != nil {
			input.CdnConfiguration.AdSegmentUrlPrefix = plan.CdnConfiguration.AdSegmentUrlPrefix
		}
		if plan.CdnConfiguration.ContentSegmentUrlPrefix != nil {
			input.CdnConfiguration.ContentSegmentUrlPrefix = plan.CdnConfiguration.ContentSegmentUrlPrefix
		}
	}

	if plan.ConfigurationAliases != nil {
		input.ConfigurationAliases = plan.ConfigurationAliases
	}

	if plan.DashConfiguration != nil {
		if plan.DashConfiguration.MpdLocation != nil {
			input.DashConfiguration.MpdLocation = plan.DashConfiguration.MpdLocation
		}
		if plan.DashConfiguration.OriginManifestType != nil {
			input.DashConfiguration.OriginManifestType = plan.DashConfiguration.OriginManifestType
		}
	}

	if plan.LivePreRollConfiguration != nil {
		if plan.LivePreRollConfiguration.AdDecisionServerUrl != nil {
			input.LivePreRollConfiguration.AdDecisionServerUrl = plan.LivePreRollConfiguration.AdDecisionServerUrl
		}
		if plan.LivePreRollConfiguration.MaxDurationSeconds != nil {
			input.LivePreRollConfiguration.MaxDurationSeconds = plan.LivePreRollConfiguration.MaxDurationSeconds
		}
	}

	if plan.ManifestProcessingRules != nil {
		input.ManifestProcessingRules = &mediatailor.ManifestProcessingRules{
			AdMarkerPassthrough: &mediatailor.AdMarkerPassthrough{
				Enabled: plan.ManifestProcessingRules.AdMarkerPassthrough.Enabled,
			},
		}
	}

	input.Name = aws.String(plan.Name.String())

	if plan.PersonalizationThresholdSeconds != nil {
		input.PersonalizationThresholdSeconds = plan.PersonalizationThresholdSeconds
	}

	if plan.SlateAdUrl != nil && *plan.SlateAdUrl != "" {
		input.SlateAdUrl = plan.SlateAdUrl
	}

	if plan.Tags != nil {
		input.Tags = plan.Tags
	}

	if plan.TranscodeProfileName != nil && *plan.TranscodeProfileName != "" {
		input.TranscodeProfileName = plan.TranscodeProfileName
	}

	if plan.VideoContentSourceUrl != nil && *plan.VideoContentSourceUrl != "" {
		input.VideoContentSourceUrl = plan.VideoContentSourceUrl
	}

	return *input
}

func readPlaybackConfigToPlan(plan resourcePlaybackConfigurationModel, playbackConfiguration mediatailor.PutPlaybackConfigurationOutput) resourcePlaybackConfigurationModel {
	plan.PlaybackConfigurationArn = types.StringValue(*playbackConfiguration.PlaybackConfigurationArn)
	plan.AdDecisionServerUrl = playbackConfiguration.AdDecisionServerUrl
	// AVAIL SUPRESSION
	if playbackConfiguration.AvailSuppression != nil {
		plan.AvailSupression = &resourceAvailSupressionModel{}
		if playbackConfiguration.AvailSuppression.Mode != nil {
			plan.AvailSupression.Mode = playbackConfiguration.AvailSuppression.Mode
		}
		if playbackConfiguration.AvailSuppression.Value != nil {
			plan.AvailSupression.Value = playbackConfiguration.AvailSuppression.Value
		}
	}
	// BUMPER
	if playbackConfiguration.Bumper != nil {
		plan.Bumper = &resourceBumperModel{}
		if playbackConfiguration.Bumper.EndUrl != nil {
			plan.Bumper.EndUrl = playbackConfiguration.Bumper.EndUrl
		}
		if playbackConfiguration.Bumper.StartUrl != nil {
			plan.Bumper.StartUrl = playbackConfiguration.Bumper.StartUrl
		}
	}
	// CDN CONFIGURATION
	if playbackConfiguration.CdnConfiguration != nil {
		plan.CdnConfiguration = &resourceCdnConfigurationModel{}
		if playbackConfiguration.CdnConfiguration.AdSegmentUrlPrefix != nil {
			plan.CdnConfiguration.AdSegmentUrlPrefix = playbackConfiguration.CdnConfiguration.AdSegmentUrlPrefix
		}
		if playbackConfiguration.CdnConfiguration.ContentSegmentUrlPrefix != nil {
			plan.CdnConfiguration.ContentSegmentUrlPrefix = playbackConfiguration.CdnConfiguration.ContentSegmentUrlPrefix
		}
	}
	// CONFIGURATION ALIASES
	if playbackConfiguration.ConfigurationAliases != nil {
		plan.ConfigurationAliases = playbackConfiguration.ConfigurationAliases
	}
	// DASH CONFIGURATION
	if playbackConfiguration.DashConfiguration != nil {
		plan.DashConfiguration.ManifestEndpointPrefix = types.StringValue(*playbackConfiguration.DashConfiguration.MpdLocation)
		if playbackConfiguration.DashConfiguration.MpdLocation != nil {
			plan.DashConfiguration.MpdLocation = playbackConfiguration.DashConfiguration.MpdLocation
		}
		if playbackConfiguration.DashConfiguration.OriginManifestType != nil {
			plan.DashConfiguration.OriginManifestType = playbackConfiguration.DashConfiguration.OriginManifestType
		}
	}
	// HLS CONFIGURATION
	if playbackConfiguration.HlsConfiguration != nil {
		plan.HlsConfiguration = &resourceHlsConfigurationModel{}
		if playbackConfiguration.HlsConfiguration.ManifestEndpointPrefix != nil {
			plan.HlsConfiguration.ManifestEndpointPrefix = types.StringValue(*playbackConfiguration.HlsConfiguration.ManifestEndpointPrefix)
		}
	}
	// LIVE PRE ROLL CONFIGURATION
	if playbackConfiguration.LivePreRollConfiguration != nil {
		plan.LivePreRollConfiguration = &resourceLivePreRollConfigurationModel{}
		if playbackConfiguration.LivePreRollConfiguration.AdDecisionServerUrl != nil {
			plan.LivePreRollConfiguration.AdDecisionServerUrl = playbackConfiguration.LivePreRollConfiguration.AdDecisionServerUrl
		}
		if playbackConfiguration.LivePreRollConfiguration.MaxDurationSeconds != nil {
			plan.LivePreRollConfiguration.MaxDurationSeconds = playbackConfiguration.LivePreRollConfiguration.MaxDurationSeconds
		}
	}
	// LOG CONFIGURATION
	if playbackConfiguration.LogConfiguration != nil {
		plan.LogConfiguration = &resourceLogConfigurationModel{}
		if playbackConfiguration.LogConfiguration.PercentEnabled != nil {
			plan.LogConfiguration.PercentEnabled = types.Int64Value(*playbackConfiguration.LogConfiguration.PercentEnabled)
		}
	}
	// MANIFEST PROCESSING RULES
	if playbackConfiguration.ManifestProcessingRules != nil {
		plan.ManifestProcessingRules = &resourceManifestProcessingRulesModel{}
		if playbackConfiguration.ManifestProcessingRules.AdMarkerPassthrough != nil {
			plan.ManifestProcessingRules.AdMarkerPassthrough = &resourceAdMarkerPassthroughModel{}
			if playbackConfiguration.ManifestProcessingRules.AdMarkerPassthrough.Enabled != nil {
				plan.ManifestProcessingRules.AdMarkerPassthrough.Enabled = playbackConfiguration.ManifestProcessingRules.AdMarkerPassthrough.Enabled
			}
		}
	}
	plan.Name = types.StringValue(*playbackConfiguration.Name)
	// PERSONALIZATION THRESHOLD SECONDS
	if playbackConfiguration.PersonalizationThresholdSeconds != nil {
		plan.PersonalizationThresholdSeconds = playbackConfiguration.PersonalizationThresholdSeconds
	}
	// PLAYBACK ENDPOINT PREFIX
	plan.PlaybackEndpointPrefix = types.StringValue(*playbackConfiguration.PlaybackEndpointPrefix)
	// SESSION INITIALIZATION ENDPOINT PREFIX
	plan.SessionInitializationEndpointPrefix = types.StringValue(*playbackConfiguration.SessionInitializationEndpointPrefix)
	// SLATE AD URL
	if playbackConfiguration.SlateAdUrl != nil {
		plan.SlateAdUrl = playbackConfiguration.SlateAdUrl
	}
	// TAGS
	if playbackConfiguration.Tags != nil {
		plan.Tags = playbackConfiguration.Tags
	}
	// TRANSCODE PROFILE NAME
	if playbackConfiguration.TranscodeProfileName != nil {
		plan.TranscodeProfileName = playbackConfiguration.TranscodeProfileName
	}
	// VIDEO CONTENT SOURCE URL
	if playbackConfiguration.VideoContentSourceUrl != nil {
		plan.VideoContentSourceUrl = playbackConfiguration.VideoContentSourceUrl
	}

	plan.ID = types.StringValue(*playbackConfiguration.Name)

	return plan
}

func readPlaybackConfigToState(state resourcePlaybackConfigurationModel, playbackConfiguration mediatailor.GetPlaybackConfigurationOutput) resourcePlaybackConfigurationModel {
	state.Name = types.StringValue(*playbackConfiguration.Name)
	state.PlaybackConfigurationArn = types.StringValue(*playbackConfiguration.PlaybackConfigurationArn)
	state.AdDecisionServerUrl = playbackConfiguration.AdDecisionServerUrl
	// AVAIL SUPRESSION
	if playbackConfiguration.AvailSuppression != nil {
		state.AvailSupression = &resourceAvailSupressionModel{}
		if playbackConfiguration.AvailSuppression.Mode != nil {
			state.AvailSupression.Mode = playbackConfiguration.AvailSuppression.Mode
		}
		if playbackConfiguration.AvailSuppression.Value != nil {
			state.AvailSupression.Value = playbackConfiguration.AvailSuppression.Value
		}
	}
	// BUMPER
	if playbackConfiguration.Bumper != nil {
		state.Bumper = &resourceBumperModel{}
		if playbackConfiguration.Bumper.EndUrl != nil {
			state.Bumper.EndUrl = playbackConfiguration.Bumper.EndUrl
		}
		if playbackConfiguration.Bumper.StartUrl != nil {
			state.Bumper.StartUrl = playbackConfiguration.Bumper.StartUrl
		}
	}
	// CDN CONFIGURATION
	if playbackConfiguration.CdnConfiguration != nil {
		state.CdnConfiguration = &resourceCdnConfigurationModel{}
		if playbackConfiguration.CdnConfiguration.AdSegmentUrlPrefix != nil {
			state.CdnConfiguration.AdSegmentUrlPrefix = playbackConfiguration.CdnConfiguration.AdSegmentUrlPrefix
		}
		if playbackConfiguration.CdnConfiguration.ContentSegmentUrlPrefix != nil {
			state.CdnConfiguration.ContentSegmentUrlPrefix = playbackConfiguration.CdnConfiguration.ContentSegmentUrlPrefix
		}
	}
	// CONFIGURATION ALIASES
	if playbackConfiguration.ConfigurationAliases != nil {
		state.ConfigurationAliases = playbackConfiguration.ConfigurationAliases
	}
	// DASH CONFIGURATION
	if playbackConfiguration.DashConfiguration != nil {
		state.DashConfiguration.ManifestEndpointPrefix = types.StringValue(*playbackConfiguration.DashConfiguration.MpdLocation)
		if playbackConfiguration.DashConfiguration.MpdLocation != nil {
			state.DashConfiguration.MpdLocation = playbackConfiguration.DashConfiguration.MpdLocation
		}
		if playbackConfiguration.DashConfiguration.OriginManifestType != nil {
			state.DashConfiguration.OriginManifestType = playbackConfiguration.DashConfiguration.OriginManifestType
		}
	}
	// HLS CONFIGURATION
	if playbackConfiguration.HlsConfiguration != nil {
		state.HlsConfiguration = &resourceHlsConfigurationModel{}
		if playbackConfiguration.HlsConfiguration.ManifestEndpointPrefix != nil {
			state.HlsConfiguration.ManifestEndpointPrefix = types.StringValue(*playbackConfiguration.HlsConfiguration.ManifestEndpointPrefix)
		}
	}
	// LIVE PRE ROLL CONFIGURATION
	if playbackConfiguration.LivePreRollConfiguration != nil {
		state.LivePreRollConfiguration = &resourceLivePreRollConfigurationModel{}
		if playbackConfiguration.LivePreRollConfiguration.AdDecisionServerUrl != nil {
			state.LivePreRollConfiguration.AdDecisionServerUrl = playbackConfiguration.LivePreRollConfiguration.AdDecisionServerUrl
		}
		if playbackConfiguration.LivePreRollConfiguration.MaxDurationSeconds != nil {
			state.LivePreRollConfiguration.MaxDurationSeconds = playbackConfiguration.LivePreRollConfiguration.MaxDurationSeconds
		}
	}
	// LOG CONFIGURATION
	if playbackConfiguration.LogConfiguration != nil {
		state.LogConfiguration.PercentEnabled = types.Int64Value(*playbackConfiguration.LogConfiguration.PercentEnabled)
	}
	// MANIFEST PROCESSING RULES
	if playbackConfiguration.ManifestProcessingRules != nil {
		state.ManifestProcessingRules = &resourceManifestProcessingRulesModel{}
		if playbackConfiguration.ManifestProcessingRules.AdMarkerPassthrough != nil {
			state.ManifestProcessingRules.AdMarkerPassthrough = &resourceAdMarkerPassthroughModel{}
			if playbackConfiguration.ManifestProcessingRules.AdMarkerPassthrough.Enabled != nil {
				state.ManifestProcessingRules.AdMarkerPassthrough.Enabled = playbackConfiguration.ManifestProcessingRules.AdMarkerPassthrough.Enabled
			}
		}
	}
	// PERSONALIZATION THRESHOLD SECONDS
	if playbackConfiguration.PersonalizationThresholdSeconds != nil {
		state.PersonalizationThresholdSeconds = playbackConfiguration.PersonalizationThresholdSeconds
	}
	// PLAYBACK ENDPOINT PREFIX
	state.PlaybackEndpointPrefix = types.StringValue(*playbackConfiguration.PlaybackEndpointPrefix)
	// SESSION INITIALIZATION ENDPOINT PREFIX
	state.SessionInitializationEndpointPrefix = types.StringValue(*playbackConfiguration.SessionInitializationEndpointPrefix)
	// SLATE AD URL
	if playbackConfiguration.SlateAdUrl != nil {
		state.SlateAdUrl = playbackConfiguration.SlateAdUrl
	}
	// TAGS
	if playbackConfiguration.Tags != nil {
		state.Tags = playbackConfiguration.Tags
	}
	// TRANSCODE PROFILE NAME
	if playbackConfiguration.TranscodeProfileName != nil {
		state.TranscodeProfileName = playbackConfiguration.TranscodeProfileName
	}
	// VIDEO CONTENT SOURCE URL
	if playbackConfiguration.VideoContentSourceUrl != nil {
		state.VideoContentSourceUrl = playbackConfiguration.VideoContentSourceUrl
	}

	state.ID = types.StringValue(*playbackConfiguration.Name)

	return state
}
