package awsmt

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func sourceLocationInput(plan resourceSourceLocationModel) mediatailor.CreateSourceLocationInput {
	var params mediatailor.CreateSourceLocationInput

	// Access Configuration
	if plan.AccessConfiguration != nil {
		params.AccessConfiguration = getAccessConfigurationInput(plan.AccessConfiguration)
	}
	// Default Segment Delivery Configuration
	if plan.DefaultSegmentDeliveryConfiguration != nil {
		params.DefaultSegmentDeliveryConfiguration = getDefaultSegmentDeliveryConfigurationInput(plan.DefaultSegmentDeliveryConfiguration)
	}

	// HTTP Configuration
	if plan.HttpConfiguration != nil {
		params.HttpConfiguration = getHttpConfigurationInput(plan.HttpConfiguration)
	}

	// Source Location Name
	params.SourceLocationName = plan.SourceLocationName

	// Segment Delivery Configurations
	if len(plan.SegmentDeliveryConfigurations) > 0 && plan.SegmentDeliveryConfigurations != nil {
		params.SegmentDeliveryConfigurations = getSegmentDeliveryConfigurationsInput(plan.SegmentDeliveryConfigurations)
	}

	// Tags
	if len(plan.Tags) > 0 && plan.Tags != nil {
		params.Tags = plan.Tags
	}

	return params
}

func getAccessConfigurationInput(accessConfiguration *accessConfigurationRModel) *mediatailor.AccessConfiguration {
	params := &mediatailor.AccessConfiguration{}
	if accessConfiguration != nil {
		if accessConfiguration.AccessType != nil && *accessConfiguration.AccessType != "" {
			params.AccessType = accessConfiguration.AccessType
		}
		if accessConfiguration.SecretsManagerAccessTokenConfiguration != nil {
			params.SecretsManagerAccessTokenConfiguration = &mediatailor.SecretsManagerAccessTokenConfiguration{}
			params.SecretsManagerAccessTokenConfiguration = getSMATC(*accessConfiguration.SecretsManagerAccessTokenConfiguration)
		}
	}
	return params
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

func getDefaultSegmentDeliveryConfigurationInput(defaultSegmentDeliveryConfiguration *defaultSegmentDeliveryConfigurationRModel) *mediatailor.DefaultSegmentDeliveryConfiguration {
	params := &mediatailor.DefaultSegmentDeliveryConfiguration{}
	if defaultSegmentDeliveryConfiguration.BaseUrl != nil && *defaultSegmentDeliveryConfiguration.BaseUrl != "" {
		params.BaseUrl = defaultSegmentDeliveryConfiguration.BaseUrl
	}
	return params
}

func getHttpConfigurationInput(httpConfiguration *httpConfigurationRModel) *mediatailor.HttpConfiguration {
	params := &mediatailor.HttpConfiguration{}
	if httpConfiguration != nil {
		if httpConfiguration.BaseUrl != nil && *httpConfiguration.BaseUrl != "" {
			params.BaseUrl = httpConfiguration.BaseUrl
		}
	}
	return params
}

func getSegmentDeliveryConfigurationsInput(segmentDeliveryConfigurations []segmentDeliveryConfigurationsRModel) []*mediatailor.SegmentDeliveryConfiguration {
	var params []*mediatailor.SegmentDeliveryConfiguration
	for _, segmentDeliveryConfiguration := range segmentDeliveryConfigurations {
		segmentDeliveryConfigurations := &mediatailor.SegmentDeliveryConfiguration{}
		segmentDeliveryConfigurations.BaseUrl = segmentDeliveryConfiguration.BaseUrl
		segmentDeliveryConfigurations.Name = segmentDeliveryConfiguration.SDCName
		params = append(params, segmentDeliveryConfigurations)
	}
	return params
}

func updateSourceLocationInput(plan resourceSourceLocationModel) mediatailor.UpdateSourceLocationInput {
	var params mediatailor.UpdateSourceLocationInput

	// Access Configuration
	if plan.AccessConfiguration != nil {
		params.AccessConfiguration = getAccessConfigurationInput(plan.AccessConfiguration)
	}
	// Default Segment Delivery Configuration
	if plan.DefaultSegmentDeliveryConfiguration != nil {
		params.DefaultSegmentDeliveryConfiguration = getDefaultSegmentDeliveryConfigurationInput(plan.DefaultSegmentDeliveryConfiguration)
	}

	// HTTP Configuration
	if plan.HttpConfiguration != nil {
		params.HttpConfiguration = getHttpConfigurationInput(plan.HttpConfiguration)
	}

	// Segment Delivery Configurations
	if len(plan.SegmentDeliveryConfigurations) > 0 && plan.SegmentDeliveryConfigurations != nil {
		params.SegmentDeliveryConfigurations = getSegmentDeliveryConfigurationsInput(plan.SegmentDeliveryConfigurations)
	}

	// Source Location Name
	params.SourceLocationName = plan.SourceLocationName

	return params
}

func readSourceLocationToPlan(plan resourceSourceLocationModel, sourceLocation mediatailor.CreateSourceLocationOutput) resourceSourceLocationModel {
	// Set state
	plan.ID = types.StringValue(*sourceLocation.SourceLocationName)
	if sourceLocation.AccessConfiguration != nil {
		plan = readAccessConfiguration(plan, sourceLocation)
	}
	if sourceLocation.Arn != nil && *sourceLocation.Arn != "" {
		plan.Arn = types.StringValue(*sourceLocation.Arn)
	}
	if sourceLocation.CreationTime != nil {
		plan.CreationTime = types.StringValue((aws.TimeValue(sourceLocation.CreationTime)).String())
	}
	if sourceLocation.DefaultSegmentDeliveryConfiguration != nil {
		plan = readDefaultSegmentDeliveryConfiguration(plan, sourceLocation)
	}
	if sourceLocation.HttpConfiguration != nil {
		plan = readHttpConfiguration(plan, sourceLocation)
	}
	if sourceLocation.LastModifiedTime != nil {
		plan.LastModifiedTime = types.StringValue((aws.TimeValue(sourceLocation.LastModifiedTime)).String())
	}
	if sourceLocation.SegmentDeliveryConfigurations != nil && len(sourceLocation.SegmentDeliveryConfigurations) > 0 {
		plan = readSegmentDeliveryConfigurations(plan, sourceLocation)
	}

	if sourceLocation.SourceLocationName != nil && *sourceLocation.SourceLocationName != "" {
		plan.SourceLocationName = sourceLocation.SourceLocationName
	}
	if sourceLocation.Tags != nil && len(sourceLocation.Tags) > 0 {
		plan.Tags = sourceLocation.Tags
	}

	return plan
}

func readAccessConfiguration(plan resourceSourceLocationModel, sourceLocation mediatailor.CreateSourceLocationOutput) resourceSourceLocationModel {
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
	return plan
}

func readDefaultSegmentDeliveryConfiguration(plan resourceSourceLocationModel, sourceLocation mediatailor.CreateSourceLocationOutput) resourceSourceLocationModel {
	if sourceLocation.DefaultSegmentDeliveryConfiguration != nil {
		plan.DefaultSegmentDeliveryConfiguration = &defaultSegmentDeliveryConfigurationRModel{}
		if sourceLocation.DefaultSegmentDeliveryConfiguration.BaseUrl != nil && *sourceLocation.DefaultSegmentDeliveryConfiguration.BaseUrl != "" {
			plan.DefaultSegmentDeliveryConfiguration.BaseUrl = sourceLocation.DefaultSegmentDeliveryConfiguration.BaseUrl
		}
	}
	return plan
}
func readSegmentDeliveryConfigurations(plan resourceSourceLocationModel, sourceLocation mediatailor.CreateSourceLocationOutput) resourceSourceLocationModel {
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
	return plan
}

func readHttpConfiguration(plan resourceSourceLocationModel, sourceLocation mediatailor.CreateSourceLocationOutput) resourceSourceLocationModel {
	if sourceLocation.HttpConfiguration != nil {
		plan.HttpConfiguration = &httpConfigurationRModel{}
		if sourceLocation.HttpConfiguration.BaseUrl != nil && *sourceLocation.HttpConfiguration.BaseUrl != "" {
			plan.HttpConfiguration.BaseUrl = sourceLocation.HttpConfiguration.BaseUrl
		}
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
