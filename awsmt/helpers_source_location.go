package awsmt

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func sourceLocationInput(plan resourceSourceLocationModel) mediatailor.CreateSourceLocationInput {
	var params mediatailor.CreateSourceLocationInput

	emptyString := ""

	// Access Configuration
	if plan.AccessConfiguration != nil {
		params.AccessConfiguration = &mediatailor.AccessConfiguration{}
		if plan.AccessConfiguration.AccessType != nil && plan.AccessConfiguration.AccessType != &emptyString {
			params.AccessConfiguration.AccessType = plan.AccessConfiguration.AccessType
		}
		if plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration != nil {
			params.AccessConfiguration.SecretsManagerAccessTokenConfiguration = &mediatailor.SecretsManagerAccessTokenConfiguration{}
			params.AccessConfiguration.SecretsManagerAccessTokenConfiguration = getSMATC(*plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration)
		}
	}

	// Default Segment Delivery Configuration
	if plan.DefaultSegmentDeliveryConfiguration != nil {
		params.DefaultSegmentDeliveryConfiguration = &mediatailor.DefaultSegmentDeliveryConfiguration{}
		if plan.DefaultSegmentDeliveryConfiguration.BaseUrl != nil && plan.DefaultSegmentDeliveryConfiguration.BaseUrl != &emptyString {
			params.DefaultSegmentDeliveryConfiguration = &mediatailor.DefaultSegmentDeliveryConfiguration{
				BaseUrl: plan.DefaultSegmentDeliveryConfiguration.BaseUrl,
			}
		}
	}

	// HTTP Configuration
	if plan.HttpConfiguration != nil {
		params.HttpConfiguration = &mediatailor.HttpConfiguration{}
		if plan.HttpConfiguration.BaseUrl != nil && plan.HttpConfiguration.BaseUrl != &emptyString {
			params.HttpConfiguration.BaseUrl = plan.HttpConfiguration.BaseUrl
		}

	}

	// Source Location Name
	params.SourceLocationName = plan.SourceLocationName

	// Segment Delivery Configurations
	if len(plan.SegmentDeliveryConfigurations) > 0 && plan.SegmentDeliveryConfigurations != nil {
		params.SegmentDeliveryConfigurations = []*mediatailor.SegmentDeliveryConfiguration{}
		for _, segmentDeliveryConfiguration := range plan.SegmentDeliveryConfigurations {
			segmentDeliveryConfigurations := &mediatailor.SegmentDeliveryConfiguration{}
			segmentDeliveryConfigurations.BaseUrl = segmentDeliveryConfiguration.BaseUrl
			segmentDeliveryConfigurations.Name = segmentDeliveryConfiguration.SDCName
			params.SegmentDeliveryConfigurations = append(params.SegmentDeliveryConfigurations, segmentDeliveryConfigurations)
		}
	}

	// Tags
	if len(plan.Tags) > 0 && plan.Tags != nil {
		params.Tags = plan.Tags
	}

	return params
}

func updateSourceLocationInput(plan resourceSourceLocationModel) mediatailor.UpdateSourceLocationInput {
	var params mediatailor.UpdateSourceLocationInput

	emptyString := ""

	// Access Configuration
	if plan.AccessConfiguration != nil {
		params.AccessConfiguration = &mediatailor.AccessConfiguration{}
		if plan.AccessConfiguration.AccessType != nil && plan.AccessConfiguration.AccessType != &emptyString {
			params.AccessConfiguration = &mediatailor.AccessConfiguration{
				AccessType: plan.AccessConfiguration.AccessType,
			}
		}
		if plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration != nil {
			params.AccessConfiguration.SecretsManagerAccessTokenConfiguration = &mediatailor.SecretsManagerAccessTokenConfiguration{}
			if plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration != nil {
				params.AccessConfiguration.SecretsManagerAccessTokenConfiguration = &mediatailor.SecretsManagerAccessTokenConfiguration{}
				params.AccessConfiguration.SecretsManagerAccessTokenConfiguration = getSMATC(*plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration)
			}
		}
	}

	// Default Segment Delivery Configuration
	if plan.DefaultSegmentDeliveryConfiguration != nil {
		params.DefaultSegmentDeliveryConfiguration = &mediatailor.DefaultSegmentDeliveryConfiguration{}
		if plan.DefaultSegmentDeliveryConfiguration.BaseUrl != nil && plan.DefaultSegmentDeliveryConfiguration.BaseUrl != &emptyString {
			params.DefaultSegmentDeliveryConfiguration = &mediatailor.DefaultSegmentDeliveryConfiguration{
				BaseUrl: plan.DefaultSegmentDeliveryConfiguration.BaseUrl,
			}
		}
	}

	// HTTP Configuration
	if plan.HttpConfiguration != nil {
		params.HttpConfiguration = &mediatailor.HttpConfiguration{}
		if plan.HttpConfiguration.BaseUrl != nil {
			params.HttpConfiguration.BaseUrl = plan.HttpConfiguration.BaseUrl
		}
	}

	// Segment Delivery Configurations
	if len(plan.SegmentDeliveryConfigurations) > 0 && plan.SegmentDeliveryConfigurations != nil {
		params.SegmentDeliveryConfigurations = []*mediatailor.SegmentDeliveryConfiguration{}
		for _, segmentDeliveryConfiguration := range plan.SegmentDeliveryConfigurations {
			segmentDeliveryConfigurations := &mediatailor.SegmentDeliveryConfiguration{}
			segmentDeliveryConfigurations.BaseUrl = segmentDeliveryConfiguration.BaseUrl
			segmentDeliveryConfigurations.Name = segmentDeliveryConfiguration.SDCName
			params.SegmentDeliveryConfigurations = append(params.SegmentDeliveryConfigurations, segmentDeliveryConfigurations)
		}
	}

	// Source Location Name
	params.SourceLocationName = plan.SourceLocationName

	return params
}

func readSourceLocationToPlan(plan resourceSourceLocationModel, sourceLocation mediatailor.CreateSourceLocationOutput) resourceSourceLocationModel {
	// Set state
	plan.ID = types.StringValue(*sourceLocation.SourceLocationName)
	if sourceLocation.AccessConfiguration != nil {
		plan.AccessConfiguration = &accessConfigurationRModel{}
		if sourceLocation.AccessConfiguration.AccessType != nil && *sourceLocation.AccessConfiguration.AccessType != "" {
			plan.AccessConfiguration.AccessType = sourceLocation.AccessConfiguration.AccessType
		}
		if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration != nil {
			plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration = &secretsManagerAccessTokenConfigurationRModel{}
			if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName != nil && *sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName != "" {
				plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName = sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName
			}
			if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn != nil && *sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn != "" {
				plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn = sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn
			}
			if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey != nil && *sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey != "" {
				plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey = sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey
			}
		}
	}
	if sourceLocation.Arn != nil && *sourceLocation.Arn != "" {
		plan.Arn = types.StringValue(*sourceLocation.Arn)
	}
	if sourceLocation.CreationTime != nil {
		plan.CreationTime = types.StringValue((aws.TimeValue(sourceLocation.CreationTime)).String())
	}
	if sourceLocation.DefaultSegmentDeliveryConfiguration != nil {
		plan.DefaultSegmentDeliveryConfiguration = &defaultSegmentDeliveryConfigurationRModel{}
		if sourceLocation.DefaultSegmentDeliveryConfiguration.BaseUrl != nil && *sourceLocation.DefaultSegmentDeliveryConfiguration.BaseUrl != "" {
			plan.DefaultSegmentDeliveryConfiguration.BaseUrl = sourceLocation.DefaultSegmentDeliveryConfiguration.BaseUrl
		}
	}
	if sourceLocation.HttpConfiguration != nil {
		plan.HttpConfiguration = &httpConfigurationRModel{}
		if sourceLocation.HttpConfiguration.BaseUrl != nil && *sourceLocation.HttpConfiguration.BaseUrl != "" {
			plan.HttpConfiguration.BaseUrl = sourceLocation.HttpConfiguration.BaseUrl
		}
	}
	if sourceLocation.LastModifiedTime != nil {
		plan.LastModifiedTime = types.StringValue((aws.TimeValue(sourceLocation.LastModifiedTime)).String())
	}
	if sourceLocation.SegmentDeliveryConfigurations != nil && len(sourceLocation.SegmentDeliveryConfigurations) > 0 {
		plan.SegmentDeliveryConfigurations = []segmentDeliveryConfigurationsRModel{}
		for _, segmentDeliveryConfiguration := range sourceLocation.SegmentDeliveryConfigurations {
			segmentDeliveryConfigurations := segmentDeliveryConfigurationsRModel{}
			if segmentDeliveryConfiguration.BaseUrl != nil && *segmentDeliveryConfiguration.BaseUrl != "" {
				segmentDeliveryConfigurations.BaseUrl = segmentDeliveryConfiguration.BaseUrl
			}
			if segmentDeliveryConfiguration.Name != nil && *segmentDeliveryConfiguration.Name != "" {
				segmentDeliveryConfigurations.SDCName = segmentDeliveryConfiguration.Name
			}
			plan.SegmentDeliveryConfigurations = append(plan.SegmentDeliveryConfigurations, segmentDeliveryConfigurations)
		}
	}

	if sourceLocation.SourceLocationName != nil && *sourceLocation.SourceLocationName != "" {
		plan.SourceLocationName = sourceLocation.SourceLocationName
	}
	if sourceLocation.Tags != nil && len(sourceLocation.Tags) > 0 {
		plan.Tags = sourceLocation.Tags
	}

	return plan
}

func deleteSourceLocation(client *mediatailor.MediaTailor, name *string) error {
	vodSourcesList, err := client.ListVodSources(&mediatailor.ListVodSourcesInput{SourceLocationName: name})
	if err != nil {
		return err
	}
	for _, vodSource := range vodSourcesList.Items {
		if _, err := client.DeleteVodSource(&mediatailor.DeleteVodSourceInput{VodSourceName: vodSource.VodSourceName, SourceLocationName: name}); err != nil {
			return err
		}
	}
	liveSourcesList, err := client.ListLiveSources(&mediatailor.ListLiveSourcesInput{SourceLocationName: name})
	if err != nil {

		return err
	}
	for _, liveSource := range liveSourcesList.Items {
		if _, err := client.DeleteLiveSource(&mediatailor.DeleteLiveSourceInput{LiveSourceName: liveSource.LiveSourceName, SourceLocationName: name}); err != nil {

			return err
		}
	}
	_, err = client.DeleteSourceLocation(&mediatailor.DeleteSourceLocationInput{SourceLocationName: name})
	if err != nil {

		return err
	}

	return nil
}

func getSMATC(plan secretsManagerAccessTokenConfigurationRModel) *mediatailor.SecretsManagerAccessTokenConfiguration {
	params := &mediatailor.SecretsManagerAccessTokenConfiguration{}
	if plan.HeaderName != nil && *plan.HeaderName != "" {
		params.HeaderName = plan.HeaderName
	}
	if plan.SecretArn != nil && *plan.SecretArn != "" {
		params.SecretArn = plan.SecretArn
	}
	if plan.SecretStringKey != nil && *plan.SecretStringKey != "" {
		params.SecretStringKey = plan.SecretStringKey
	}
	return params
}
