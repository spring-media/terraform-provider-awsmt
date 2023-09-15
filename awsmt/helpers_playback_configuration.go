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

func getAvailSuppressionInput(availSuppression *resourceAvailSupressionModel) *mediatailor.AvailSuppression {
	params := &mediatailor.AvailSuppression{}
	if availSuppression != nil {
		if availSuppression.Mode != nil && *availSuppression.Mode != "" {
			params.Mode = availSuppression.Mode
		}
		if availSuppression.Value != nil && *availSuppression.Value != "" {
			params.Value = availSuppression.Value
		}
		if availSuppression.FillPolicy != nil && *availSuppression.FillPolicy != "" {
			params.FillPolicy = availSuppression.FillPolicy
		}
	}
	return params
}

func getBumperInput(bumper *resourceBumperModel) *mediatailor.Bumper {
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

func getCdnConfigurationInput(cdnConfiguration *resourceCdnConfigurationModel) *mediatailor.CdnConfiguration {
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

func getDashConfigurationInput(dashConfiguration *resourceDashConfigurationModel) *mediatailor.DashConfigurationForPut {
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

func getLivePreRollConfigurationInput(livePreRollConfiguration *resourceLivePreRollConfigurationModel) *mediatailor.LivePreRollConfiguration {
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

func readPlaybackConfigToPlan(plan resourcePlaybackConfigurationModel, playbackConfiguration mediatailor.PutPlaybackConfigurationOutput) resourcePlaybackConfigurationModel {
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
		plan = readManifestProcessingRules(plan, playbackConfiguration)
	}

	plan.Name, plan.PersonalizationThresholdSeconds, plan.PlaybackEndpointPrefix, plan.SessionInitializationEndpointPrefix, plan.SlateAdUrl, plan.TranscodeProfileName, plan.VideoContentSourceUrl, plan.Tags = readPlaybackConfigurationTemps(plan, playbackConfiguration)

	plan.ID = types.StringValue(*playbackConfiguration.Name)

	return plan
}

func readPlaybackConfigurationTemps(plan resourcePlaybackConfigurationModel, playbackConfiguration mediatailor.PutPlaybackConfigurationOutput) (*string, *int64, types.String, types.String, *string, *string, *string, map[string]*string) {
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
	if playbackConfiguration.Tags != nil {
		plan.Tags = playbackConfiguration.Tags
	}
	return plan.Name, plan.PersonalizationThresholdSeconds, plan.PlaybackEndpointPrefix, plan.SessionInitializationEndpointPrefix, plan.SlateAdUrl, plan.TranscodeProfileName, plan.VideoContentSourceUrl, plan.Tags
}

func readAvailSuppression(plan resourcePlaybackConfigurationModel, playbackConfiguration mediatailor.PutPlaybackConfigurationOutput) resourcePlaybackConfigurationModel {
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
	return plan
}

func readBumper(plan resourcePlaybackConfigurationModel, playbackConfiguration mediatailor.PutPlaybackConfigurationOutput) resourcePlaybackConfigurationModel {
	if playbackConfiguration.Bumper != nil && (playbackConfiguration.Bumper.EndUrl != nil || playbackConfiguration.Bumper.StartUrl != nil) {
		plan.Bumper = &resourceBumperModel{}
		if playbackConfiguration.Bumper.EndUrl != nil {
			plan.Bumper.EndUrl = playbackConfiguration.Bumper.EndUrl
		}
		if playbackConfiguration.Bumper.StartUrl != nil {
			plan.Bumper.StartUrl = playbackConfiguration.Bumper.StartUrl
		}
	}
	return plan
}

func readCdnConfiguration(plan resourcePlaybackConfigurationModel, playbackConfiguration mediatailor.PutPlaybackConfigurationOutput) resourcePlaybackConfigurationModel {
	if playbackConfiguration.CdnConfiguration != nil {
		plan.CdnConfiguration = &resourceCdnConfigurationModel{}
		if playbackConfiguration.CdnConfiguration.AdSegmentUrlPrefix != nil {
			plan.CdnConfiguration.AdSegmentUrlPrefix = playbackConfiguration.CdnConfiguration.AdSegmentUrlPrefix
		}
		if playbackConfiguration.CdnConfiguration.ContentSegmentUrlPrefix != nil {
			plan.CdnConfiguration.ContentSegmentUrlPrefix = playbackConfiguration.CdnConfiguration.ContentSegmentUrlPrefix
		}
	}
	return plan
}

func readDashConfiguration(plan resourcePlaybackConfigurationModel, playbackConfiguration mediatailor.PutPlaybackConfigurationOutput) resourcePlaybackConfigurationModel {
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
	return plan
}

func readLivePreRollConfiguration(plan resourcePlaybackConfigurationModel, playbackConfiguration mediatailor.PutPlaybackConfigurationOutput) resourcePlaybackConfigurationModel {
	if playbackConfiguration.LivePreRollConfiguration != nil && (playbackConfiguration.LivePreRollConfiguration.AdDecisionServerUrl != nil || playbackConfiguration.LivePreRollConfiguration.MaxDurationSeconds != nil) {
		plan.LivePreRollConfiguration = &resourceLivePreRollConfigurationModel{}
		if playbackConfiguration.LivePreRollConfiguration.AdDecisionServerUrl != nil {
			plan.LivePreRollConfiguration.AdDecisionServerUrl = playbackConfiguration.LivePreRollConfiguration.AdDecisionServerUrl
		}
		if playbackConfiguration.LivePreRollConfiguration.MaxDurationSeconds != nil {
			plan.LivePreRollConfiguration.MaxDurationSeconds = playbackConfiguration.LivePreRollConfiguration.MaxDurationSeconds
		}
	}
	return plan
}

func readManifestProcessingRules(plan resourcePlaybackConfigurationModel, playbackConfiguration mediatailor.PutPlaybackConfigurationOutput) resourcePlaybackConfigurationModel {
	if playbackConfiguration.ManifestProcessingRules != nil {
		plan.ManifestProcessingRules = &resourceManifestProcessingRulesModel{}
		if playbackConfiguration.ManifestProcessingRules.AdMarkerPassthrough != nil && playbackConfiguration.ManifestProcessingRules.AdMarkerPassthrough.Enabled != nil {
			plan.ManifestProcessingRules.AdMarkerPassthrough = &resourceAdMarkerPassthroughModel{}
			plan.ManifestProcessingRules.AdMarkerPassthrough.Enabled = playbackConfiguration.ManifestProcessingRules.AdMarkerPassthrough.Enabled
		}
	}
	return plan
}

func readPlaybackConfigToData(data dataSourcePlaybackConfigurationModel, playbackConfiguration mediatailor.GetPlaybackConfigurationOutput) dataSourcePlaybackConfigurationModel {
	data.AdDecisionServerUrl = playbackConfiguration.AdDecisionServerUrl

	// AVAIL SUPRESSION
	if playbackConfiguration.AvailSuppression != nil {
		data = readAvailSuppressionToData(data, playbackConfiguration)
	}
	// BUMPER
	if playbackConfiguration.Bumper != nil {
		data = readBumperToData(data, playbackConfiguration)
	}
	// CDN CONFIGURATION
	if playbackConfiguration.CdnConfiguration != nil {
		data = readCdnConfigurationToData(data, playbackConfiguration)
	}

	if playbackConfiguration.ConfigurationAliases != nil {
		data.ConfigurationAliases = playbackConfiguration.ConfigurationAliases
	}
	// DASH CONFIGURATION
	if playbackConfiguration.DashConfiguration != nil {
		data = readDashConfigurationToData(data, playbackConfiguration)
	}
	// HLS CONFIGURATION
	if playbackConfiguration.HlsConfiguration != nil {
		data.HlsConfiguration = &hlsConfigurationModel{}
		if playbackConfiguration.HlsConfiguration.ManifestEndpointPrefix != nil {
			data.HlsConfiguration.ManifestEndpointPrefix = playbackConfiguration.HlsConfiguration.ManifestEndpointPrefix
		}
	}
	// LIVE PRE ROLL CONFIGURATION
	if playbackConfiguration.LivePreRollConfiguration != nil {
		data = readLivePreRollConfigurationToData(data, playbackConfiguration)
	}
	// LOG CONFIGURATION
	if playbackConfiguration.LogConfiguration != nil {
		data.LogConfiguration = &logConfigurationModel{}
		if playbackConfiguration.LogConfiguration.PercentEnabled != nil {
			data.LogConfiguration.PercentEnabled = playbackConfiguration.LogConfiguration.PercentEnabled
		}
	} else {
		data.LogConfiguration = &logConfigurationModel{}
		data.LogConfiguration.PercentEnabled = aws.Int64(0)
	}
	// MANIFEST PROCESSING RULES
	if playbackConfiguration.ManifestProcessingRules != nil {
		data = readManifestProcessingRulesToData(data, playbackConfiguration)
	}

	data.Name = playbackConfiguration.Name
	data.PersonalizationThresholdSeconds = playbackConfiguration.PersonalizationThresholdSeconds
	data.PlaybackConfigurationArn = playbackConfiguration.PlaybackConfigurationArn
	data.PlaybackEndpointPrefix = playbackConfiguration.PlaybackEndpointPrefix
	data.SessionInitializationEndpointPrefix = playbackConfiguration.SessionInitializationEndpointPrefix
	data.SlateAdUrl = playbackConfiguration.SlateAdUrl
	data.Tags = playbackConfiguration.Tags
	data.TranscodeProfileName = playbackConfiguration.TranscodeProfileName
	data.VideoContentSourceUrl = playbackConfiguration.VideoContentSourceUrl
	data.ID = types.StringValue(*playbackConfiguration.Name)

	return data
}

func readAvailSuppressionToData(data dataSourcePlaybackConfigurationModel, playbackConfiguration mediatailor.GetPlaybackConfigurationOutput) dataSourcePlaybackConfigurationModel {
	if playbackConfiguration.AvailSuppression != nil {
		data.AvailSupression = &availSupressionModel{}
		if playbackConfiguration.AvailSuppression.Mode != nil {
			data.AvailSupression.Mode = playbackConfiguration.AvailSuppression.Mode
		}
		if playbackConfiguration.AvailSuppression.Value != nil {
			data.AvailSupression.Value = playbackConfiguration.AvailSuppression.Value
		}
		if playbackConfiguration.AvailSuppression.FillPolicy != nil {
			data.AvailSupression.FillPolicy = playbackConfiguration.AvailSuppression.FillPolicy
		}
	}
	return data
}

func readBumperToData(data dataSourcePlaybackConfigurationModel, playbackConfiguration mediatailor.GetPlaybackConfigurationOutput) dataSourcePlaybackConfigurationModel {
	if playbackConfiguration.Bumper != nil && (playbackConfiguration.Bumper.EndUrl != nil || playbackConfiguration.Bumper.StartUrl != nil) {
		data.Bumper = &bumperModel{}
		if playbackConfiguration.Bumper.EndUrl != nil {
			data.Bumper.EndUrl = playbackConfiguration.Bumper.EndUrl
		}
		if playbackConfiguration.Bumper.StartUrl != nil {
			data.Bumper.StartUrl = playbackConfiguration.Bumper.StartUrl
		}
	}
	return data
}

func readCdnConfigurationToData(data dataSourcePlaybackConfigurationModel, playbackConfiguration mediatailor.GetPlaybackConfigurationOutput) dataSourcePlaybackConfigurationModel {
	if playbackConfiguration.CdnConfiguration != nil {
		data.CdnConfiguration = &cdnConfigurationModel{}
		if playbackConfiguration.CdnConfiguration.AdSegmentUrlPrefix != nil {
			data.CdnConfiguration.AdSegmentUrlPrefix = playbackConfiguration.CdnConfiguration.AdSegmentUrlPrefix
		}
		if playbackConfiguration.CdnConfiguration.ContentSegmentUrlPrefix != nil {
			data.CdnConfiguration.ContentSegmentUrlPrefix = playbackConfiguration.CdnConfiguration.ContentSegmentUrlPrefix
		}
	}
	return data
}

func readDashConfigurationToData(data dataSourcePlaybackConfigurationModel, playbackConfiguration mediatailor.GetPlaybackConfigurationOutput) dataSourcePlaybackConfigurationModel {
	if playbackConfiguration.DashConfiguration != nil {
		data.DashConfiguration = &dashConfigurationModel{}
		data.DashConfiguration.ManifestEndpointPrefix = playbackConfiguration.DashConfiguration.MpdLocation
		if playbackConfiguration.DashConfiguration.MpdLocation != nil {
			data.DashConfiguration.MpdLocation = playbackConfiguration.DashConfiguration.MpdLocation
		}
		if playbackConfiguration.DashConfiguration.OriginManifestType != nil {
			data.DashConfiguration.OriginManifestType = playbackConfiguration.DashConfiguration.OriginManifestType
		}
	}
	return data
}

func readLivePreRollConfigurationToData(data dataSourcePlaybackConfigurationModel, playbackConfiguration mediatailor.GetPlaybackConfigurationOutput) dataSourcePlaybackConfigurationModel {
	if playbackConfiguration.LivePreRollConfiguration != nil && (playbackConfiguration.LivePreRollConfiguration.AdDecisionServerUrl != nil || playbackConfiguration.LivePreRollConfiguration.MaxDurationSeconds != nil) {
		data.LivePreRollConfiguration = &livePreRollConfigurationModel{}
		if playbackConfiguration.LivePreRollConfiguration.AdDecisionServerUrl != nil {
			data.LivePreRollConfiguration.AdDecisionServerUrl = playbackConfiguration.LivePreRollConfiguration.AdDecisionServerUrl
		}
		if playbackConfiguration.LivePreRollConfiguration.MaxDurationSeconds != nil {
			data.LivePreRollConfiguration.MaxDurationSeconds = playbackConfiguration.LivePreRollConfiguration.MaxDurationSeconds
		}
	}
	return data
}

func readManifestProcessingRulesToData(data dataSourcePlaybackConfigurationModel, playbackConfiguration mediatailor.GetPlaybackConfigurationOutput) dataSourcePlaybackConfigurationModel {
	if playbackConfiguration.ManifestProcessingRules != nil {
		data.ManifestProcessingRules = &manifestProcessingRulesModel{}
		if playbackConfiguration.ManifestProcessingRules.AdMarkerPassthrough != nil && playbackConfiguration.ManifestProcessingRules.AdMarkerPassthrough.Enabled != nil {
			data.ManifestProcessingRules.AdMarkerPassthrough = &adMarkerPassthroughModel{}
			data.ManifestProcessingRules.AdMarkerPassthrough.Enabled = playbackConfiguration.ManifestProcessingRules.AdMarkerPassthrough.Enabled
		}
	}
	return data
}
