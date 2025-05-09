package awsmt

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/mediatailor/types"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-mediatailor/awsmt/models"
)

func getCreateSourceLocationInput(model models.SourceLocationModel) mediatailor.CreateSourceLocationInput {
	var params mediatailor.CreateSourceLocationInput

	// Access Configuration
	if model.AccessConfiguration != nil {
		params.AccessConfiguration = getAccessConfigurationInput(model.AccessConfiguration)
	}

	// Default Segment Delivery Configuration
	if model.DefaultSegmentDeliveryConfiguration != nil {
		params.DefaultSegmentDeliveryConfiguration = getDefaultSegmentDeliveryConfigurationInput(model.DefaultSegmentDeliveryConfiguration)
	}

	// HTTP Configuration
	params.HttpConfiguration = getHttpConfigurationInput(model.HttpConfiguration)

	// Source Location Name
	params.SourceLocationName = model.Name

	// Segment Delivery Configurations
	if len(model.SegmentDeliveryConfigurations) > 0 {
		params.SegmentDeliveryConfigurations = getSegmentDeliveryConfigurationsInput(model.SegmentDeliveryConfigurations)
	}

	// Tags
	if len(model.Tags) > 0 && model.Tags != nil {
		params.Tags = model.Tags
	}

	return params
}

func getAccessConfigurationInput(accessConfiguration *models.AccessConfigurationModel) *awsTypes.AccessConfiguration {
	if accessConfiguration == nil {
		return nil
	}

	temp := &awsTypes.AccessConfiguration{}

	if accessConfiguration.AccessType != nil {
		var accessType awsTypes.AccessType
		switch *accessConfiguration.AccessType {
		case "SECRETS_MANAGER_ACCESS_TOKEN":
			accessType = awsTypes.AccessTypeSecretsManagerAccessToken
		case "AUTODETECT_SIGV4":
			accessType = awsTypes.AccessTypeAutodetectSigv4
		default:
			accessType = awsTypes.AccessTypeS3Sigv4
		}
		temp.AccessType = accessType
	}

	if accessConfiguration.SecretsManagerAccessTokenConfiguration != nil {
		temp.SecretsManagerAccessTokenConfiguration = getSMATC(*accessConfiguration.SecretsManagerAccessTokenConfiguration)
	}

	return temp
}

func getSMATC(smatc models.SecretsManagerAccessTokenConfigurationModel) *awsTypes.SecretsManagerAccessTokenConfiguration {
	temp := &awsTypes.SecretsManagerAccessTokenConfiguration{}
	if smatc.HeaderName != nil && *smatc.HeaderName != "" {
		temp.HeaderName = smatc.HeaderName
	}
	if smatc.SecretArn != nil && *smatc.SecretArn != "" {
		temp.SecretArn = smatc.SecretArn
	}
	if smatc.SecretStringKey != nil && *smatc.SecretStringKey != "" {
		temp.SecretStringKey = smatc.SecretStringKey
	}
	return temp
}

func getDefaultSegmentDeliveryConfigurationInput(defaultSegmentDeliveryConfiguration *models.DefaultSegmentDeliveryConfigurationModel) *awsTypes.DefaultSegmentDeliveryConfiguration {
	if defaultSegmentDeliveryConfiguration.BaseUrl == nil || *defaultSegmentDeliveryConfiguration.BaseUrl == "" {
		return nil
	}
	temp := &awsTypes.DefaultSegmentDeliveryConfiguration{
		BaseUrl: defaultSegmentDeliveryConfiguration.BaseUrl,
	}

	return temp
}

func getHttpConfigurationInput(httpConfiguration *models.HttpConfigurationModel) *awsTypes.HttpConfiguration {
	if httpConfiguration == nil {
		return nil
	}
	if httpConfiguration.BaseUrl == nil || *httpConfiguration.BaseUrl == "" {
		return nil
	}
	temp := &awsTypes.HttpConfiguration{
		BaseUrl: httpConfiguration.BaseUrl,
	}
	return temp
}

func getSegmentDeliveryConfigurationsInput(segmentDeliveryConfigurations []models.SegmentDeliveryConfigurationsModel) []awsTypes.SegmentDeliveryConfiguration {
	var params []awsTypes.SegmentDeliveryConfiguration
	for _, segmentDeliveryConfiguration := range segmentDeliveryConfigurations {
		params = append(params, awsTypes.SegmentDeliveryConfiguration{
			BaseUrl: segmentDeliveryConfiguration.BaseUrl,
			Name:    segmentDeliveryConfiguration.SDCName,
		})
	}
	return params
}

func getUpdateSourceLocationInput(model models.SourceLocationModel) mediatailor.UpdateSourceLocationInput {
	var params mediatailor.UpdateSourceLocationInput

	params.AccessConfiguration = getAccessConfigurationInput(model.AccessConfiguration)
	// Default Segment Delivery Configuration
	params.DefaultSegmentDeliveryConfiguration = getDefaultSegmentDeliveryConfigurationInput(model.DefaultSegmentDeliveryConfiguration)
	// HTTP Configuration
	params.HttpConfiguration = getHttpConfigurationInput(model.HttpConfiguration)

	// Segment Delivery Configurations
	if len(model.SegmentDeliveryConfigurations) > 0 && model.SegmentDeliveryConfigurations != nil {
		params.SegmentDeliveryConfigurations = getSegmentDeliveryConfigurationsInput(model.SegmentDeliveryConfigurations)
	}

	// Source Location Name
	params.SourceLocationName = model.Name

	return params
}

func writeSourceLocationToPlan(model models.SourceLocationModel, sourceLocation mediatailor.CreateSourceLocationOutput) models.SourceLocationModel {
	// Set state

	model = readSourceLocationComputedValues(model, sourceLocation)

	model = readAccessConfiguration(model, sourceLocation)

	model = readDefaultSegmentDeliveryConfiguration(model, sourceLocation)

	model = readHttpConfiguration(model, sourceLocation)

	model = readSegmentDeliveryConfigurations(model, sourceLocation)

	if len(sourceLocation.Tags) > 0 {
		model.Tags = sourceLocation.Tags
	}

	return model
}

func readSourceLocationComputedValues(model models.SourceLocationModel, sourceLocation mediatailor.CreateSourceLocationOutput) models.SourceLocationModel {
	model.ID = types.StringValue(*sourceLocation.SourceLocationName)

	if sourceLocation.Arn != nil && *sourceLocation.Arn != "" {
		model.Arn = types.StringValue(*sourceLocation.Arn)
	}

	if sourceLocation.CreationTime != nil {
		model.CreationTime = types.StringValue(sourceLocation.CreationTime.String())
	}

	if sourceLocation.LastModifiedTime != nil {
		model.LastModifiedTime = types.StringValue(sourceLocation.LastModifiedTime.String())
	}

	if sourceLocation.SourceLocationName != nil && *sourceLocation.SourceLocationName != "" {
		model.Name = sourceLocation.SourceLocationName
	}

	return model
}

func readAccessConfiguration(model models.SourceLocationModel, sourceLocation mediatailor.CreateSourceLocationOutput) models.SourceLocationModel {
	if sourceLocation.AccessConfiguration == nil {
		return model
	}

	model.AccessConfiguration = &models.AccessConfigurationModel{}

	if string(sourceLocation.AccessConfiguration.AccessType) != "" {
		accessType := string(sourceLocation.AccessConfiguration.AccessType)
		model.AccessConfiguration.AccessType = &accessType
	}

	if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration != nil {
		model = readSMATConfiguration(model, sourceLocation)
	}

	return model
}

func readSMATConfiguration(model models.SourceLocationModel, sourceLocation mediatailor.CreateSourceLocationOutput) models.SourceLocationModel {
	if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration == nil {
		return model
	}

	model.AccessConfiguration.SecretsManagerAccessTokenConfiguration = &models.SecretsManagerAccessTokenConfigurationModel{}
	if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName != nil && *sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName != "" {
		model.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName = sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName
	}
	if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn != nil && *sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn != "" {
		model.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn = sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn
	}
	if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey != nil && *sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey != "" {
		model.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey = sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey
	}

	return model
}

func readDefaultSegmentDeliveryConfiguration(plan models.SourceLocationModel, sourceLocation mediatailor.CreateSourceLocationOutput) models.SourceLocationModel {
	if sourceLocation.DefaultSegmentDeliveryConfiguration == nil {
		return plan
	}
	plan.DefaultSegmentDeliveryConfiguration = &models.DefaultSegmentDeliveryConfigurationModel{}
	if sourceLocation.DefaultSegmentDeliveryConfiguration.BaseUrl != nil && *sourceLocation.DefaultSegmentDeliveryConfiguration.BaseUrl != "" {
		plan.DefaultSegmentDeliveryConfiguration.BaseUrl = sourceLocation.DefaultSegmentDeliveryConfiguration.BaseUrl
	}

	return plan
}

func readSegmentDeliveryConfigurations(model models.SourceLocationModel, sourceLocation mediatailor.CreateSourceLocationOutput) models.SourceLocationModel {
	if len(sourceLocation.SegmentDeliveryConfigurations) == 0 {
		return model
	}
	model.SegmentDeliveryConfigurations = []models.SegmentDeliveryConfigurationsModel{}
	for _, segmentDeliveryConfiguration := range sourceLocation.SegmentDeliveryConfigurations {
		segmentDeliveryConfigurations := models.SegmentDeliveryConfigurationsModel{}
		if segmentDeliveryConfiguration.BaseUrl != nil && *segmentDeliveryConfiguration.BaseUrl != "" {
			segmentDeliveryConfigurations.BaseUrl = segmentDeliveryConfiguration.BaseUrl
		}
		if segmentDeliveryConfiguration.Name != nil && *segmentDeliveryConfiguration.Name != "" {
			segmentDeliveryConfigurations.SDCName = segmentDeliveryConfiguration.Name
		}
		model.SegmentDeliveryConfigurations = append(model.SegmentDeliveryConfigurations, segmentDeliveryConfigurations)
	}

	return model
}

func readHttpConfiguration(model models.SourceLocationModel, sourceLocation mediatailor.CreateSourceLocationOutput) models.SourceLocationModel {
	if sourceLocation.HttpConfiguration == nil || sourceLocation.HttpConfiguration.BaseUrl == nil || *sourceLocation.HttpConfiguration.BaseUrl == "" {
		return model
	}

	model.HttpConfiguration = &models.HttpConfigurationModel{
		BaseUrl: sourceLocation.HttpConfiguration.BaseUrl,
	}

	return model
}

func deleteSourceLocation(client *mediatailor.Client, name *string) error {
	vodSourcesList, err := client.ListVodSources(context.TODO(), &mediatailor.ListVodSourcesInput{SourceLocationName: name})
	if err != nil {
		return err
	}
	for _, vodSource := range vodSourcesList.Items {
		if _, err := client.DeleteVodSource(context.TODO(), &mediatailor.DeleteVodSourceInput{VodSourceName: vodSource.VodSourceName, SourceLocationName: name}); err != nil {
			return err
		}
	}
	liveSourcesList, err := client.ListLiveSources(context.TODO(), &mediatailor.ListLiveSourcesInput{SourceLocationName: name})
	if err != nil {

		return err
	}
	for _, liveSource := range liveSourcesList.Items {
		if _, err := client.DeleteLiveSource(context.TODO(), &mediatailor.DeleteLiveSourceInput{LiveSourceName: liveSource.LiveSourceName, SourceLocationName: name}); err != nil {

			return err
		}
	}
	_, err = client.DeleteSourceLocation(context.TODO(), &mediatailor.DeleteSourceLocationInput{SourceLocationName: name})
	if err != nil {

		return err
	}

	return nil
}

func recreateSourceLocation(client *mediatailor.Client, plan models.SourceLocationModel) (*models.SourceLocationModel, error) {
	err := deleteSourceLocation(client, plan.Name)
	if err != nil {
		return nil, err
	}

	params := getCreateSourceLocationInput(plan)
	sourceLocation, err := client.CreateSourceLocation(context.TODO(), &params)
	if err != nil {
		return nil, fmt.Errorf("error while creating new source location with new access configuration %v", err.Error())
	}
	model := writeSourceLocationToPlan(plan, *sourceLocation)
	return &model, nil
}
