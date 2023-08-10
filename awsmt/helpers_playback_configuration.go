package awsmt

import (
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func playbackConfigurationInput(plan resourcePlaybackConfigurationModel) mediatailor.PutPlaybackConfigurationInput {

	input := &mediatailor.PutPlaybackConfigurationInput{}

	input.AdDecisionServerUrl = plan.AdDecisionServerUrl

	if plan.AvailSupression != nil {
		input.AvailSuppression = &mediatailor.AvailSuppression{}
		if plan.AvailSupression.Mode != nil {
			input.AvailSuppression.Mode = plan.AvailSupression.Mode
		}
		if plan.AvailSupression.Value != nil {
			input.AvailSuppression.Value = plan.AvailSupression.Value
		}
		if plan.AvailSupression.FillPolicy != nil {
			input.AvailSuppression.FillPolicy = plan.AvailSupression.FillPolicy
		}
	}

	if plan.Bumper != nil {
		input.Bumper = &mediatailor.Bumper{}
		if plan.Bumper.EndUrl != nil {
			input.Bumper.EndUrl = plan.Bumper.EndUrl
		}
		if plan.Bumper.StartUrl != nil {
			input.Bumper.StartUrl = plan.Bumper.StartUrl
		}
	}

	if plan.CdnConfiguration != nil {
		input.CdnConfiguration = &mediatailor.CdnConfiguration{}
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
		input.DashConfiguration = &mediatailor.DashConfigurationForPut{}
		if plan.DashConfiguration.MpdLocation != nil {
			input.DashConfiguration.MpdLocation = plan.DashConfiguration.MpdLocation
		}
		if plan.DashConfiguration.OriginManifestType != nil {
			input.DashConfiguration.OriginManifestType = plan.DashConfiguration.OriginManifestType
		}
	}

	if plan.LivePreRollConfiguration != nil {
		input.LivePreRollConfiguration = &mediatailor.LivePreRollConfiguration{}
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

	input.Name = plan.Name

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
		if playbackConfiguration.AvailSuppression.FillPolicy != nil {
			plan.AvailSupression.FillPolicy = playbackConfiguration.AvailSuppression.FillPolicy
		}
	}
	// BUMPER
	if playbackConfiguration.Bumper != nil && (playbackConfiguration.Bumper.EndUrl != nil || playbackConfiguration.Bumper.StartUrl != nil) {
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
		plan.DashConfiguration = &resourceDashConfigurationModel{}
		plan.DashConfiguration.ManifestEndpointPrefix = types.StringValue(*playbackConfiguration.DashConfiguration.MpdLocation)
		if playbackConfiguration.DashConfiguration.MpdLocation != nil {
			plan.DashConfiguration.MpdLocation = playbackConfiguration.DashConfiguration.MpdLocation
		}
		if playbackConfiguration.DashConfiguration.OriginManifestType != nil {
			plan.DashConfiguration.OriginManifestType = playbackConfiguration.DashConfiguration.OriginManifestType
		}
	}

	// HLS CONFIGURATION
	if playbackConfiguration.HlsConfiguration != nil && playbackConfiguration.HlsConfiguration.ManifestEndpointPrefix != nil {
		plan.HlsConfigurationManifestEndpointPrefix = types.StringValue(*playbackConfiguration.HlsConfiguration.ManifestEndpointPrefix)
	}

	// LOG CONFIGURATION
	if playbackConfiguration.LogConfiguration != nil {
		plan.LogConfigurationPercentEnabled = types.Int64Value(*playbackConfiguration.LogConfiguration.PercentEnabled)
	} else {
		plan.LogConfigurationPercentEnabled = types.Int64Value(0)
	}

	// LIVE PRE ROLL CONFIGURATION
	if playbackConfiguration.LivePreRollConfiguration != nil && (playbackConfiguration.LivePreRollConfiguration.AdDecisionServerUrl != nil || playbackConfiguration.LivePreRollConfiguration.MaxDurationSeconds != nil) {
		plan.LivePreRollConfiguration = &resourceLivePreRollConfigurationModel{}
		if playbackConfiguration.LivePreRollConfiguration.AdDecisionServerUrl != nil {
			plan.LivePreRollConfiguration.AdDecisionServerUrl = playbackConfiguration.LivePreRollConfiguration.AdDecisionServerUrl
		}
		if playbackConfiguration.LivePreRollConfiguration.MaxDurationSeconds != nil {
			plan.LivePreRollConfiguration.MaxDurationSeconds = playbackConfiguration.LivePreRollConfiguration.MaxDurationSeconds
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
	plan.Name = playbackConfiguration.Name
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
