package awsmt

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/mediatailor/types"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type putPlaybackConfigurationInputBuilder struct {
	model playbackConfigurationModel
	input *mediatailor.PutPlaybackConfigurationInput
}

type putPlaybackConfigurationModelbuilder struct {
	model      *playbackConfigurationModel
	output     mediatailor.PutPlaybackConfigurationOutput
	isResource bool
}

func (i *putPlaybackConfigurationInputBuilder) getInput() *mediatailor.PutPlaybackConfigurationInput {

	i.addAvailSuppressionToInput()
	i.addBumperToInput()
	i.addCdnConfigurationToInput()
	i.addDashConfigurationToInput()
	i.addLivePreRollConfigurationToInput()
	i.addManifestProcessingRulesToInput()
	i.addOptionalFieldsToInput()
	i.addRequiredFieldsToInput()

	return i.input
}

func (i *putPlaybackConfigurationInputBuilder) addAvailSuppressionToInput() {
	if i.model.AvailSuppression == nil {
		return
	}
	temp := &awsTypes.AvailSuppression{}
	if i.model.AvailSuppression.Mode != nil {
		var mode awsTypes.Mode
		switch *i.model.AvailSuppression.Mode {
		case "BEHIND_LIVE_EDGE":
			mode = awsTypes.ModeBehindLiveEdge
		case "AFTER_LIVE_EDGE":
			mode = awsTypes.ModeAfterLiveEdge
		default:
			mode = awsTypes.ModeOff
		}
		temp.Mode = mode
	}
	if i.model.AvailSuppression.Value != nil {
		temp.Value = i.model.AvailSuppression.Value
	}
	if i.model.AvailSuppression.FillPolicy != nil {
		var policy awsTypes.FillPolicy
		if *i.model.AvailSuppression.FillPolicy == "FULL_AVAIL_ONLY" {
			policy = awsTypes.FillPolicyFullAvailOnly
		} else {
			policy = awsTypes.FillPolicyPartialAvail
		}
		temp.FillPolicy = policy
	}
	i.input.AvailSuppression = temp
}

func (i *putPlaybackConfigurationInputBuilder) addBumperToInput() {
	if i.model.Bumper == nil {
		return
	}
	temp := &awsTypes.Bumper{}
	if i.model.Bumper.EndUrl != nil {
		temp.EndUrl = i.model.Bumper.EndUrl
	}
	if i.model.Bumper.StartUrl != nil {
		temp.StartUrl = i.model.Bumper.StartUrl
	}
	i.input.Bumper = temp

}

func (i *putPlaybackConfigurationInputBuilder) addCdnConfigurationToInput() {
	if i.model.CdnConfiguration == nil {
		return
	}
	temp := &awsTypes.CdnConfiguration{}
	if i.model.CdnConfiguration != nil {
		if i.model.CdnConfiguration.AdSegmentUrlPrefix != nil {
			temp.AdSegmentUrlPrefix = i.model.CdnConfiguration.AdSegmentUrlPrefix
		}
		if i.model.CdnConfiguration.ContentSegmentUrlPrefix != nil {
			temp.ContentSegmentUrlPrefix = i.model.CdnConfiguration.ContentSegmentUrlPrefix
		}
	}
	i.input.CdnConfiguration = temp
}

func (i *putPlaybackConfigurationInputBuilder) addDashConfigurationToInput() {
	if i.model.DashConfiguration == nil {
		return
	}
	temp := &awsTypes.DashConfigurationForPut{}
	if i.model.DashConfiguration.MpdLocation != nil {
		temp.MpdLocation = i.model.DashConfiguration.MpdLocation
	}
	if i.model.DashConfiguration.OriginManifestType != nil {
		var manifestType awsTypes.OriginManifestType
		if *i.model.DashConfiguration.OriginManifestType == "SINGLE_PERIOD" {
			manifestType = awsTypes.OriginManifestTypeSinglePeriod
		} else {
			manifestType = awsTypes.OriginManifestTypeMultiPeriod
		}
		temp.OriginManifestType = manifestType
	}
	i.input.DashConfiguration = temp
}

func (i *putPlaybackConfigurationInputBuilder) addLivePreRollConfigurationToInput() {
	if i.model.LivePreRollConfiguration == nil {
		return
	}
	temp := &awsTypes.LivePreRollConfiguration{}
	if i.model.LivePreRollConfiguration.AdDecisionServerUrl != nil {
		temp.AdDecisionServerUrl = i.model.LivePreRollConfiguration.AdDecisionServerUrl
	}
	if i.model.LivePreRollConfiguration.MaxDurationSeconds != nil {
		temp.MaxDurationSeconds = i.model.LivePreRollConfiguration.MaxDurationSeconds
	}
	i.input.LivePreRollConfiguration = temp
}

func (i *putPlaybackConfigurationInputBuilder) addManifestProcessingRulesToInput() {
	if i.model.ManifestProcessingRules == nil {
		return
	}
	temp := &awsTypes.ManifestProcessingRules{}
	if i.model.ManifestProcessingRules.AdMarkerPassthrough != nil {
		temp.AdMarkerPassthrough = &awsTypes.AdMarkerPassthrough{
			Enabled: i.model.ManifestProcessingRules.AdMarkerPassthrough.Enabled,
		}
	}
	i.input.ManifestProcessingRules = temp
}

func (i *putPlaybackConfigurationInputBuilder) addOptionalFieldsToInput() {
	if i.model.ConfigurationAliases != nil {
		i.input.ConfigurationAliases = i.model.ConfigurationAliases
	}

	if i.model.PersonalizationThresholdSeconds != nil {
		i.input.PersonalizationThresholdSeconds = i.model.PersonalizationThresholdSeconds
	}

	if i.model.SlateAdUrl != nil {
		i.input.SlateAdUrl = i.model.SlateAdUrl
	}

	if i.model.Tags != nil {
		i.input.Tags = i.model.Tags
	}

	if i.model.TranscodeProfileName != nil {
		i.input.TranscodeProfileName = i.model.TranscodeProfileName
	}

	if i.model.VideoContentSourceUrl != nil {
		i.input.VideoContentSourceUrl = i.model.VideoContentSourceUrl
	}
}

func (i *putPlaybackConfigurationInputBuilder) addRequiredFieldsToInput() {
	i.input.AdDecisionServerUrl = i.model.AdDecisionServerUrl
	i.input.Name = i.model.Name
}

func (m *putPlaybackConfigurationModelbuilder) getModel() playbackConfigurationModel {

	m.addAvailSuppressionToModel()
	m.addBumperToModel()
	m.addCdnConfigurationToModel()
	m.addDashConfigurationToModel()
	m.addOptionalFieldsToModel()
	m.addLivePreRollConfigurationToModel()
	m.addManifestProcessingRulesToModel()
	m.addRequiredFieldsToModel()

	return *m.model
}

func (m *putPlaybackConfigurationModelbuilder) addAvailSuppressionToModel() {
	if m.output.AvailSuppression == nil {
		return
	}
	if m.model.AvailSuppression == nil && m.isResource {
		return
	}
	m.model.AvailSuppression = &availSuppressionModel{}

	m.model.AvailSuppression.Mode = aws.String(string(m.output.AvailSuppression.Mode))
	m.model.AvailSuppression.FillPolicy = aws.String(string(m.output.AvailSuppression.FillPolicy))
	if m.output.AvailSuppression.Value != nil {
		m.model.AvailSuppression.Value = m.output.AvailSuppression.Value
	}
}

func (m *putPlaybackConfigurationModelbuilder) addBumperToModel() {
	if m.output.Bumper == nil || (m.output.Bumper.EndUrl == nil || m.output.Bumper.StartUrl == nil) {
		return
	}
	m.model.Bumper = &bumperModel{}
	if m.output.Bumper.EndUrl != nil {
		m.model.Bumper.EndUrl = m.output.Bumper.EndUrl
	}
	if m.output.Bumper.StartUrl != nil {
		m.model.Bumper.StartUrl = m.output.Bumper.StartUrl
	}
}

func (m *putPlaybackConfigurationModelbuilder) addCdnConfigurationToModel() {
	if m.output.CdnConfiguration == nil {
		return
	}
	if m.model.CdnConfiguration == nil && m.isResource {
		return
	}
	m.model.CdnConfiguration = &cdnConfigurationModel{}
	if m.output.CdnConfiguration.AdSegmentUrlPrefix != nil {
		m.model.CdnConfiguration.AdSegmentUrlPrefix = m.output.CdnConfiguration.AdSegmentUrlPrefix
	}
	if m.output.CdnConfiguration.ContentSegmentUrlPrefix != nil {
		m.model.CdnConfiguration.ContentSegmentUrlPrefix = m.output.CdnConfiguration.ContentSegmentUrlPrefix
	}
}

func (m *putPlaybackConfigurationModelbuilder) addDashConfigurationToModel() {
	if m.output.DashConfiguration == nil {
		return
	}
	if m.model.DashConfiguration == nil && m.isResource {
		return
	}
	m.model.DashConfiguration = &dashConfigurationModel{}
	m.model.DashConfiguration.ManifestEndpointPrefix = types.StringValue(*m.output.DashConfiguration.MpdLocation)
	if m.output.DashConfiguration.MpdLocation != nil {
		m.model.DashConfiguration.MpdLocation = m.output.DashConfiguration.MpdLocation
	}
	m.model.DashConfiguration.OriginManifestType = aws.String(string(m.output.DashConfiguration.OriginManifestType))
}

func (m *putPlaybackConfigurationModelbuilder) addLivePreRollConfigurationToModel() {
	if m.output.LivePreRollConfiguration == nil || (m.output.LivePreRollConfiguration.AdDecisionServerUrl == nil || m.output.LivePreRollConfiguration.MaxDurationSeconds == nil) {
		return
	}
	m.model.LivePreRollConfiguration = &livePreRollConfigurationModel{}
	if m.output.LivePreRollConfiguration.AdDecisionServerUrl != nil {
		m.model.LivePreRollConfiguration.AdDecisionServerUrl = m.output.LivePreRollConfiguration.AdDecisionServerUrl
	}
	if m.output.LivePreRollConfiguration.MaxDurationSeconds != nil {
		m.model.LivePreRollConfiguration.MaxDurationSeconds = m.output.LivePreRollConfiguration.MaxDurationSeconds
	}
}

func (m *putPlaybackConfigurationModelbuilder) addManifestProcessingRulesToModel() {
	if m.output.ManifestProcessingRules == nil || m.output.ManifestProcessingRules.AdMarkerPassthrough == nil {
		return
	}
	if m.model.ManifestProcessingRules == nil && m.isResource {
		return
	}
	m.model.ManifestProcessingRules = &manifestProcessingRulesModel{
		AdMarkerPassthrough: &adMarkerPassthroughModel{
			Enabled: m.output.ManifestProcessingRules.AdMarkerPassthrough.Enabled,
		},
	}
}

func (m *putPlaybackConfigurationModelbuilder) addRequiredFieldsToModel() {
	m.model.AdDecisionServerUrl = m.output.AdDecisionServerUrl
	m.model.ID = types.StringValue(*m.output.Name)
	m.model.Name = m.output.Name
	m.model.PlaybackConfigurationArn = types.StringValue(*m.output.PlaybackConfigurationArn)
	m.model.PlaybackEndpointPrefix = types.StringValue(*m.output.PlaybackEndpointPrefix)
	m.model.SessionInitializationEndpointPrefix = types.StringValue(*m.output.SessionInitializationEndpointPrefix)
}

func (m *putPlaybackConfigurationModelbuilder) addOptionalFieldsToModel() {
	if m.output.ConfigurationAliases != nil {
		m.model.ConfigurationAliases = m.output.ConfigurationAliases
	}

	if m.output.HlsConfiguration != nil && m.output.HlsConfiguration.ManifestEndpointPrefix != nil {
		m.model.HlsConfigurationManifestEndpointPrefix = types.StringValue(*m.output.HlsConfiguration.ManifestEndpointPrefix)
	}

	if m.output.LogConfiguration != nil {
		m.model.LogConfigurationPercentEnabled = types.Int64Value(int64(m.output.LogConfiguration.PercentEnabled))
	} else {
		m.model.LogConfigurationPercentEnabled = types.Int64Value(0)
	}

	if m.output.PersonalizationThresholdSeconds != nil {
		m.model.PersonalizationThresholdSeconds = m.output.PersonalizationThresholdSeconds
	}

	if m.output.SlateAdUrl != nil {
		m.model.SlateAdUrl = m.output.SlateAdUrl
	}

	if m.output.VideoContentSourceUrl != nil {
		m.model.VideoContentSourceUrl = m.output.VideoContentSourceUrl
	}

	if m.output.TranscodeProfileName != nil {
		m.model.TranscodeProfileName = m.output.TranscodeProfileName
	}

	if len(m.output.Tags) > 0 {
		m.model.Tags = m.output.Tags
	}
}
