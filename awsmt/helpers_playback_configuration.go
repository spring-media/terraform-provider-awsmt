package awsmt

import (
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func playbackConfigurationInput(plan playbackConfigurationModel) mediatailor.PutPlaybackConfigurationInput {

	input := &mediatailor.PutPlaybackConfigurationInput{}

	input.AdDecisionServerUrl = plan.AdDecisionServerUrl

	if plan.AvailSupression != nil {
		input.AvailSuppression = getAvailSuppressionInput(plan.AvailSupression)
	}

	if plan.Bumper != nil {
		input.Bumper = getBumperInput(plan.Bumper)
	}

	if plan.CdnConfiguration != nil {
		input.CdnConfiguration = getCdnConfigurationInput(plan.CdnConfiguration)
	}

	if plan.ConfigurationAliases != nil {
		input.ConfigurationAliases = plan.ConfigurationAliases
	}

	if plan.DashConfiguration != nil {
		input.DashConfiguration = getDashConfigurationInput(plan.DashConfiguration)
	}

	if plan.LivePreRollConfiguration != nil {
		input.LivePreRollConfiguration = getLivePreRollConfigurationInput(plan.LivePreRollConfiguration)
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

func getAvailSuppressionInput(availSuppression *availSupressionModel) *mediatailor.AvailSuppression {
	params := &mediatailor.AvailSuppression{}
	if availSuppression != nil {
		if availSuppression.Mode != nil {
			params.Mode = availSuppression.Mode
		}
		if availSuppression.Value != nil {
			params.Value = availSuppression.Value
		}
		if availSuppression.FillPolicy != nil {
			params.FillPolicy = availSuppression.FillPolicy
		}
	}
	return params
}

func getBumperInput(bumper *bumperModel) *mediatailor.Bumper {
	params := &mediatailor.Bumper{}
	if bumper != nil {
		if bumper.EndUrl != nil && *bumper.EndUrl != "" {
			params.EndUrl = bumper.EndUrl
		}
		if bumper.StartUrl != nil && *bumper.StartUrl != "" {
			params.StartUrl = bumper.StartUrl
		}
	}
	return params
}

func getCdnConfigurationInput(cdnConfiguration *cdnConfigurationModel) *mediatailor.CdnConfiguration {
	params := &mediatailor.CdnConfiguration{}
	if cdnConfiguration != nil {
		if cdnConfiguration.AdSegmentUrlPrefix != nil && *cdnConfiguration.AdSegmentUrlPrefix != "" {
			params.AdSegmentUrlPrefix = cdnConfiguration.AdSegmentUrlPrefix
		}
		if cdnConfiguration.ContentSegmentUrlPrefix != nil && *cdnConfiguration.ContentSegmentUrlPrefix != "" {
			params.ContentSegmentUrlPrefix = cdnConfiguration.ContentSegmentUrlPrefix
		}
	}
	return params
}

func getDashConfigurationInput(dashConfiguration *dashConfigurationModel) *mediatailor.DashConfigurationForPut {
	input := &mediatailor.DashConfigurationForPut{}
	if dashConfiguration != nil {
		if dashConfiguration.MpdLocation != nil {
			input.MpdLocation = dashConfiguration.MpdLocation
		}
		if dashConfiguration.OriginManifestType != nil {
			input.OriginManifestType = dashConfiguration.OriginManifestType
		}
	}
	return input
}

func getLivePreRollConfigurationInput(livePreRollConfiguration *livePreRollConfigurationModel) *mediatailor.LivePreRollConfiguration {
	input := &mediatailor.LivePreRollConfiguration{}
	if livePreRollConfiguration != nil {
		if livePreRollConfiguration.AdDecisionServerUrl != nil {
			input.AdDecisionServerUrl = livePreRollConfiguration.AdDecisionServerUrl
		}
		if livePreRollConfiguration.MaxDurationSeconds != nil {
			input.MaxDurationSeconds = livePreRollConfiguration.MaxDurationSeconds
		}
	}
	return input
}

func readPlaybackConfigToPlan(plan playbackConfigurationModel, playbackConfiguration mediatailor.PutPlaybackConfigurationOutput) playbackConfigurationModel {
	plan.PlaybackConfigurationArn = types.StringValue(*playbackConfiguration.PlaybackConfigurationArn)
	plan.AdDecisionServerUrl = playbackConfiguration.AdDecisionServerUrl
	// AVAIL SUPRESSION
	if playbackConfiguration.AvailSuppression != nil {
		plan = readAvailSuppression(plan, playbackConfiguration)
	}
	// BUMPER
	if playbackConfiguration.Bumper != nil && (playbackConfiguration.Bumper.EndUrl != nil || playbackConfiguration.Bumper.StartUrl != nil) {
		plan = readBumper(plan, playbackConfiguration)
	}
	// CDN CONFIGURATION
	if playbackConfiguration.CdnConfiguration != nil {
		plan = readCdnConfiguration(plan, playbackConfiguration)
	}
	// CONFIGURATION ALIASES
	if playbackConfiguration.ConfigurationAliases != nil {
		plan.ConfigurationAliases = playbackConfiguration.ConfigurationAliases
	}
	// DASH CONFIGURATION
	if playbackConfiguration.DashConfiguration != nil {
		plan = readDashConfiguration(plan, playbackConfiguration)
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
		plan = readLivePreRollConfiguration(plan, playbackConfiguration)
	}

	// MANIFEST PROCESSING RULES
	if playbackConfiguration.ManifestProcessingRules != nil {
		var mpr = playbackConfiguration.ManifestProcessingRules
		if mpr.AdMarkerPassthrough != nil && mpr.AdMarkerPassthrough.Enabled != nil && *mpr.AdMarkerPassthrough.Enabled {
			plan.ManifestProcessingRules = &manifestProcessingRulesModel{AdMarkerPassthrough: &adMarkerPassthroughModel{Enabled: mpr.AdMarkerPassthrough.Enabled}}
		}
	}

	plan.Name, plan.PersonalizationThresholdSeconds, plan.PlaybackEndpointPrefix, plan.SessionInitializationEndpointPrefix, plan.SlateAdUrl, plan.TranscodeProfileName, plan.VideoContentSourceUrl, plan.Tags = readPlaybackConfigurationTemps(plan, playbackConfiguration)

	plan.ID = types.StringValue(*playbackConfiguration.Name)

	return plan
}

func readPlaybackConfigurationTemps(plan playbackConfigurationModel, playbackConfiguration mediatailor.PutPlaybackConfigurationOutput) (*string, *int64, types.String, types.String, *string, *string, *string, map[string]*string) {
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
	// TRANSCODE PROFILE NAME
	if playbackConfiguration.TranscodeProfileName != nil {
		plan.TranscodeProfileName = playbackConfiguration.TranscodeProfileName
	}
	// VIDEO CONTENT SOURCE URL
	if playbackConfiguration.VideoContentSourceUrl != nil {
		plan.VideoContentSourceUrl = playbackConfiguration.VideoContentSourceUrl
	}

	// TAGS
	if len(playbackConfiguration.Tags) > 0 {
		plan.Tags = playbackConfiguration.Tags
	}
	return plan.Name, plan.PersonalizationThresholdSeconds, plan.PlaybackEndpointPrefix, plan.SessionInitializationEndpointPrefix, plan.SlateAdUrl, plan.TranscodeProfileName, plan.VideoContentSourceUrl, plan.Tags
}

func readAvailSuppression(plan playbackConfigurationModel, playbackConfiguration mediatailor.PutPlaybackConfigurationOutput) playbackConfigurationModel {
	if playbackConfiguration.AvailSuppression != nil && *playbackConfiguration.AvailSuppression.Mode != "OFF" {
		plan.AvailSupression = &availSupressionModel{}
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
	return plan
}

func readBumper(plan playbackConfigurationModel, playbackConfiguration mediatailor.PutPlaybackConfigurationOutput) playbackConfigurationModel {
	if playbackConfiguration.Bumper != nil && (playbackConfiguration.Bumper.EndUrl != nil || playbackConfiguration.Bumper.StartUrl != nil) {
		plan.Bumper = &bumperModel{}
		if playbackConfiguration.Bumper.EndUrl != nil {
			plan.Bumper.EndUrl = playbackConfiguration.Bumper.EndUrl
		}
		if playbackConfiguration.Bumper.StartUrl != nil {
			plan.Bumper.StartUrl = playbackConfiguration.Bumper.StartUrl
		}
	}
	return plan
}

func readCdnConfiguration(plan playbackConfigurationModel, playbackConfiguration mediatailor.PutPlaybackConfigurationOutput) playbackConfigurationModel {
	if playbackConfiguration.CdnConfiguration != nil {
		plan.CdnConfiguration = &cdnConfigurationModel{}
		if playbackConfiguration.CdnConfiguration.AdSegmentUrlPrefix != nil {
			plan.CdnConfiguration.AdSegmentUrlPrefix = playbackConfiguration.CdnConfiguration.AdSegmentUrlPrefix
		}
		if playbackConfiguration.CdnConfiguration.ContentSegmentUrlPrefix != nil {
			plan.CdnConfiguration.ContentSegmentUrlPrefix = playbackConfiguration.CdnConfiguration.ContentSegmentUrlPrefix
		}
	}
	return plan
}

func readDashConfiguration(plan playbackConfigurationModel, playbackConfiguration mediatailor.PutPlaybackConfigurationOutput) playbackConfigurationModel {
	if playbackConfiguration.DashConfiguration != nil {
		plan.DashConfiguration = &dashConfigurationModel{}
		plan.DashConfiguration.ManifestEndpointPrefix = types.StringValue(*playbackConfiguration.DashConfiguration.MpdLocation)
		if playbackConfiguration.DashConfiguration.MpdLocation != nil {
			plan.DashConfiguration.MpdLocation = playbackConfiguration.DashConfiguration.MpdLocation
		}
		if playbackConfiguration.DashConfiguration.OriginManifestType != nil {
			plan.DashConfiguration.OriginManifestType = playbackConfiguration.DashConfiguration.OriginManifestType
		}
	}
	return plan
}

func readLivePreRollConfiguration(plan playbackConfigurationModel, playbackConfiguration mediatailor.PutPlaybackConfigurationOutput) playbackConfigurationModel {
	if playbackConfiguration.LivePreRollConfiguration != nil && (playbackConfiguration.LivePreRollConfiguration.AdDecisionServerUrl != nil || playbackConfiguration.LivePreRollConfiguration.MaxDurationSeconds != nil) {
		plan.LivePreRollConfiguration = &livePreRollConfigurationModel{}
		if playbackConfiguration.LivePreRollConfiguration.AdDecisionServerUrl != nil {
			plan.LivePreRollConfiguration.AdDecisionServerUrl = playbackConfiguration.LivePreRollConfiguration.AdDecisionServerUrl
		}
		if playbackConfiguration.LivePreRollConfiguration.MaxDurationSeconds != nil {
			plan.LivePreRollConfiguration.MaxDurationSeconds = playbackConfiguration.LivePreRollConfiguration.MaxDurationSeconds
		}
	}
	return plan
}
