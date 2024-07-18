package awsmt

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/mediatailor/types"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getPutPlaybackConfigurationInput(model playbackConfigurationModel) mediatailor.PutPlaybackConfigurationInput {

	input := &mediatailor.PutPlaybackConfigurationInput{}

	input.AdDecisionServerUrl = model.AdDecisionServerUrl

	input.Name = model.Name

	addAvailSuppressionToInput(input, model)
	addBumperToInput(input, model)
	addCdnConfigurationToInput(input, model)
	addDashConfigurationToInput(input, model)
	addLivePreRollConfigurationToInput(input, model)
	addManifestProcessingRulesToInput(input, model)

	if model.ConfigurationAliases != nil {
		input.ConfigurationAliases = model.ConfigurationAliases
	}

	if model.PersonalizationThresholdSeconds != nil {
		input.PersonalizationThresholdSeconds = model.PersonalizationThresholdSeconds
	}

	if model.SlateAdUrl != nil {
		input.SlateAdUrl = model.SlateAdUrl
	}

	if model.Tags != nil {
		input.Tags = model.Tags
	}

	if model.TranscodeProfileName != nil {
		input.TranscodeProfileName = model.TranscodeProfileName
	}

	if model.VideoContentSourceUrl != nil {
		input.VideoContentSourceUrl = model.VideoContentSourceUrl
	}

	return *input
}

func addAvailSuppressionToInput(input *mediatailor.PutPlaybackConfigurationInput, model playbackConfigurationModel) {
	if model.AvailSuppression == nil {
		return
	}
	temp := &awsTypes.AvailSuppression{}
	if model.AvailSuppression.Mode != nil {
		var mode awsTypes.Mode
		switch *model.AvailSuppression.Mode {
		case "BEHIND_LIVE_EDGE":
			mode = awsTypes.ModeBehindLiveEdge
		case "AFTER_LIVE_EDGE":
			mode = awsTypes.ModeAfterLiveEdge
		default:
			mode = awsTypes.ModeOff
		}
		temp.Mode = mode
	}
	if model.AvailSuppression.Value != nil {
		temp.Value = model.AvailSuppression.Value
	}
	if model.AvailSuppression.FillPolicy != nil {
		var policy awsTypes.FillPolicy
		if *model.AvailSuppression.FillPolicy == "FULL_AVAIL_ONLY" {
			policy = awsTypes.FillPolicyFullAvailOnly
		} else {
			policy = awsTypes.FillPolicyPartialAvail
		}
		temp.FillPolicy = policy
	}
	input.AvailSuppression = temp
}

func addBumperToInput(input *mediatailor.PutPlaybackConfigurationInput, model playbackConfigurationModel) {
	if model.Bumper == nil {
		return
	}
	temp := &awsTypes.Bumper{}
	if model.Bumper.EndUrl != nil {
		temp.EndUrl = model.Bumper.EndUrl
	}
	if model.Bumper.StartUrl != nil {
		temp.StartUrl = model.Bumper.StartUrl
	}
	input.Bumper = temp

}

func addCdnConfigurationToInput(input *mediatailor.PutPlaybackConfigurationInput, model playbackConfigurationModel) {
	if model.CdnConfiguration == nil {
		return
	}
	temp := &awsTypes.CdnConfiguration{}
	if model.CdnConfiguration != nil {
		if model.CdnConfiguration.AdSegmentUrlPrefix != nil {
			temp.AdSegmentUrlPrefix = model.CdnConfiguration.AdSegmentUrlPrefix
		}
		if model.CdnConfiguration.ContentSegmentUrlPrefix != nil {
			temp.ContentSegmentUrlPrefix = model.CdnConfiguration.ContentSegmentUrlPrefix
		}
	}
	input.CdnConfiguration = temp
}

func addDashConfigurationToInput(input *mediatailor.PutPlaybackConfigurationInput, model playbackConfigurationModel) {
	if model.DashConfiguration == nil {
		return
	}
	temp := &awsTypes.DashConfigurationForPut{}
	if model.DashConfiguration.MpdLocation != nil {
		temp.MpdLocation = model.DashConfiguration.MpdLocation
	}
	if model.DashConfiguration.OriginManifestType != nil {
		var manifestType awsTypes.OriginManifestType
		if *model.DashConfiguration.OriginManifestType == "SINGLE_PERIOD" {
			manifestType = awsTypes.OriginManifestTypeSinglePeriod
		} else {
			manifestType = awsTypes.OriginManifestTypeMultiPeriod
		}
		temp.OriginManifestType = manifestType
	}
	input.DashConfiguration = temp
}

func addLivePreRollConfigurationToInput(input *mediatailor.PutPlaybackConfigurationInput, model playbackConfigurationModel) {
	if model.LivePreRollConfiguration == nil {
		return
	}
	temp := &awsTypes.LivePreRollConfiguration{}
	if model.LivePreRollConfiguration.AdDecisionServerUrl != nil {
		temp.AdDecisionServerUrl = model.LivePreRollConfiguration.AdDecisionServerUrl
	}
	if model.LivePreRollConfiguration.MaxDurationSeconds != nil {
		temp.MaxDurationSeconds = model.LivePreRollConfiguration.MaxDurationSeconds
	}
	input.LivePreRollConfiguration = temp
}

func addManifestProcessingRulesToInput(input *mediatailor.PutPlaybackConfigurationInput, model playbackConfigurationModel) {
	if model.ManifestProcessingRules == nil {
		return
	}
	temp := &awsTypes.ManifestProcessingRules{}
	if model.ManifestProcessingRules.AdMarkerPassthrough != nil {
		temp.AdMarkerPassthrough = &awsTypes.AdMarkerPassthrough{
			Enabled: model.ManifestProcessingRules.AdMarkerPassthrough.Enabled,
		}
	}
	input.ManifestProcessingRules = temp
}

func readPlaybackConfig(plan playbackConfigurationModel, playbackConfiguration mediatailor.PutPlaybackConfigurationOutput) playbackConfigurationModel {

	plan.AdDecisionServerUrl = playbackConfiguration.AdDecisionServerUrl
	plan.ID = types.StringValue(*playbackConfiguration.Name)
	plan.Name = playbackConfiguration.Name
	plan.PlaybackConfigurationArn = types.StringValue(*playbackConfiguration.PlaybackConfigurationArn)
	plan.PlaybackEndpointPrefix = types.StringValue(*playbackConfiguration.PlaybackEndpointPrefix)
	plan.SessionInitializationEndpointPrefix = types.StringValue(*playbackConfiguration.SessionInitializationEndpointPrefix)

	addAvailSuppressionToModel(&plan, playbackConfiguration)
	addBumperToModel(&plan, playbackConfiguration)
	addCdnConfigurationToModel(&plan, playbackConfiguration)
	addDashConfigurationToModel(&plan, playbackConfiguration)
	addLivePreRollConfigurationToModel(&plan, playbackConfiguration)
	addManifestProcessingRulesToModel(&plan, playbackConfiguration)

	if playbackConfiguration.ConfigurationAliases != nil {
		plan.ConfigurationAliases = playbackConfiguration.ConfigurationAliases
	}

	if playbackConfiguration.HlsConfiguration != nil && playbackConfiguration.HlsConfiguration.ManifestEndpointPrefix != nil {
		plan.HlsConfigurationManifestEndpointPrefix = types.StringValue(*playbackConfiguration.HlsConfiguration.ManifestEndpointPrefix)
	}

	if playbackConfiguration.LogConfiguration != nil {
		plan.LogConfigurationPercentEnabled = types.Int64Value(int64(playbackConfiguration.LogConfiguration.PercentEnabled))
	} else {
		plan.LogConfigurationPercentEnabled = types.Int64Value(0)
	}

	if playbackConfiguration.PersonalizationThresholdSeconds != nil {
		plan.PersonalizationThresholdSeconds = playbackConfiguration.PersonalizationThresholdSeconds
	}

	if playbackConfiguration.SlateAdUrl != nil {
		plan.SlateAdUrl = playbackConfiguration.SlateAdUrl
	}

	if playbackConfiguration.VideoContentSourceUrl != nil {
		plan.VideoContentSourceUrl = playbackConfiguration.VideoContentSourceUrl
	}

	if playbackConfiguration.TranscodeProfileName != nil {
		plan.TranscodeProfileName = playbackConfiguration.TranscodeProfileName
	}

	if len(playbackConfiguration.Tags) > 0 {
		plan.Tags = playbackConfiguration.Tags
	}

	return plan
}

func addAvailSuppressionToModel(plan *playbackConfigurationModel, playbackConfiguration mediatailor.PutPlaybackConfigurationOutput) {
	if playbackConfiguration.AvailSuppression == nil {
		return
	}
	plan.AvailSuppression = &availSuppressionModel{}
	plan.AvailSuppression.Mode = aws.String(string(playbackConfiguration.AvailSuppression.Mode))
	if playbackConfiguration.AvailSuppression.Value != nil {
		plan.AvailSuppression.Value = playbackConfiguration.AvailSuppression.Value
	}
	if plan.AvailSuppression.FillPolicy != nil {
		plan.AvailSuppression.FillPolicy = aws.String(string(playbackConfiguration.AvailSuppression.FillPolicy))
	}
}

func addBumperToModel(plan *playbackConfigurationModel, playbackConfiguration mediatailor.PutPlaybackConfigurationOutput) {
	if playbackConfiguration.Bumper == nil || (playbackConfiguration.Bumper.EndUrl == nil || playbackConfiguration.Bumper.StartUrl == nil) {
		return
	}
	plan.Bumper = &bumperModel{}
	if playbackConfiguration.Bumper.EndUrl != nil {
		plan.Bumper.EndUrl = playbackConfiguration.Bumper.EndUrl
	}
	if playbackConfiguration.Bumper.StartUrl != nil {
		plan.Bumper.StartUrl = playbackConfiguration.Bumper.StartUrl
	}
}

func addCdnConfigurationToModel(plan *playbackConfigurationModel, playbackConfiguration mediatailor.PutPlaybackConfigurationOutput) {
	if playbackConfiguration.CdnConfiguration == nil {
		return
	}
	plan.CdnConfiguration = &cdnConfigurationModel{}
	if playbackConfiguration.CdnConfiguration.AdSegmentUrlPrefix != nil {
		plan.CdnConfiguration.AdSegmentUrlPrefix = playbackConfiguration.CdnConfiguration.AdSegmentUrlPrefix
	}
	if playbackConfiguration.CdnConfiguration.ContentSegmentUrlPrefix != nil {
		plan.CdnConfiguration.ContentSegmentUrlPrefix = playbackConfiguration.CdnConfiguration.ContentSegmentUrlPrefix
	}
}

func addDashConfigurationToModel(plan *playbackConfigurationModel, playbackConfiguration mediatailor.PutPlaybackConfigurationOutput) {
	if playbackConfiguration.DashConfiguration == nil {
		return
	}
	plan.DashConfiguration = &dashConfigurationModel{}
	plan.DashConfiguration.ManifestEndpointPrefix = types.StringValue(*playbackConfiguration.DashConfiguration.MpdLocation)
	if playbackConfiguration.DashConfiguration.MpdLocation != nil {
		plan.DashConfiguration.MpdLocation = playbackConfiguration.DashConfiguration.MpdLocation
	}
	plan.DashConfiguration.OriginManifestType = aws.String(string(playbackConfiguration.DashConfiguration.OriginManifestType))
}

func addLivePreRollConfigurationToModel(plan *playbackConfigurationModel, playbackConfiguration mediatailor.PutPlaybackConfigurationOutput) {
	if playbackConfiguration.LivePreRollConfiguration == nil || (playbackConfiguration.LivePreRollConfiguration.AdDecisionServerUrl == nil || playbackConfiguration.LivePreRollConfiguration.MaxDurationSeconds == nil) {
		return
	}
	plan.LivePreRollConfiguration = &livePreRollConfigurationModel{}
	if playbackConfiguration.LivePreRollConfiguration.AdDecisionServerUrl != nil {
		plan.LivePreRollConfiguration.AdDecisionServerUrl = playbackConfiguration.LivePreRollConfiguration.AdDecisionServerUrl
	}
	if playbackConfiguration.LivePreRollConfiguration.MaxDurationSeconds != nil {
		plan.LivePreRollConfiguration.MaxDurationSeconds = playbackConfiguration.LivePreRollConfiguration.MaxDurationSeconds
	}
}

func addManifestProcessingRulesToModel(plan *playbackConfigurationModel, playbackConfiguration mediatailor.PutPlaybackConfigurationOutput) {
	if playbackConfiguration.ManifestProcessingRules == nil || playbackConfiguration.ManifestProcessingRules.AdMarkerPassthrough == nil {
		return
	}
	plan.ManifestProcessingRules = &manifestProcessingRulesModel{
		AdMarkerPassthrough: &adMarkerPassthroughModel{
			Enabled: playbackConfiguration.ManifestProcessingRules.AdMarkerPassthrough.Enabled,
		},
	}
}
