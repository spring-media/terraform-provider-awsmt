package awsmt

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func sourceLocationInput(plan resourceSourceLocationModel) mediatailor.CreateSourceLocationInput {
	var params mediatailor.CreateSourceLocationInput
	// Access Configuration
	if !plan.AccessConfiguration.AccessType.IsUnknown() && !plan.AccessConfiguration.AccessType.IsNull() {
		params.AccessConfiguration = &mediatailor.AccessConfiguration{
			AccessType: aws.String(plan.AccessConfiguration.AccessType.String()),
		}
	}
	if !plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName.IsUnknown() && !plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName.IsNull() {
		params.AccessConfiguration.SecretsManagerAccessTokenConfiguration = &mediatailor.SecretsManagerAccessTokenConfiguration{
			HeaderName: aws.String(plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName.String()),
		}
	}
	if !plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn.IsUnknown() && !plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn.IsNull() {
		params.AccessConfiguration.SecretsManagerAccessTokenConfiguration = &mediatailor.SecretsManagerAccessTokenConfiguration{
			SecretArn: aws.String(plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn.String()),
		}
	}
	if !plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey.IsUnknown() && !plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey.IsNull() {
		params.AccessConfiguration.SecretsManagerAccessTokenConfiguration = &mediatailor.SecretsManagerAccessTokenConfiguration{
			SecretStringKey: aws.String(plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey.String()),
		}
	}

	// Default Segment Delivery Configuration
	if !plan.DefaultSegmentDeliveryConfiguration.BaseUrl.IsUnknown() && !plan.DefaultSegmentDeliveryConfiguration.BaseUrl.IsNull() {
		params.DefaultSegmentDeliveryConfiguration = &mediatailor.DefaultSegmentDeliveryConfiguration{
			BaseUrl: aws.String(plan.DefaultSegmentDeliveryConfiguration.BaseUrl.String()),
		}
	}

	// HTTP Configuration
	params.HttpConfiguration.BaseUrl = aws.String(plan.HttpConfiguration.BaseUrl.String())

	// Segment Delivery Configurations
	if len(plan.SegmentDeliveryConfigurations) > 0 && plan.SegmentDeliveryConfigurations != nil {
		params.SegmentDeliveryConfigurations = make([]*mediatailor.SegmentDeliveryConfiguration, len(plan.SegmentDeliveryConfigurations))
		for i, v := range plan.SegmentDeliveryConfigurations {
			params.SegmentDeliveryConfigurations[i] = &mediatailor.SegmentDeliveryConfiguration{
				BaseUrl: aws.String(v.BaseUrl.String()),
				Name:    aws.String(v.Name.String()),
			}
		}
	}

	// Source Location Name
	params.SourceLocationName = aws.String(plan.SourceLocationName.String())

	// Tags
	if len(plan.Tags) > 0 && plan.Tags != nil {
		for k, v := range plan.Tags {
			params.Tags[k] = aws.String(*v)
		}
	}

	return params
}

func updateSourceLocationInput(plan resourceSourceLocationModel) mediatailor.UpdateSourceLocationInput {
	var params mediatailor.UpdateSourceLocationInput
	// Access Configuration
	if !plan.AccessConfiguration.AccessType.IsUnknown() && !plan.AccessConfiguration.AccessType.IsNull() {
		params.AccessConfiguration = &mediatailor.AccessConfiguration{
			AccessType: aws.String(plan.AccessConfiguration.AccessType.String()),
		}
	}
	if !plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName.IsUnknown() && !plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName.IsNull() {
		params.AccessConfiguration.SecretsManagerAccessTokenConfiguration = &mediatailor.SecretsManagerAccessTokenConfiguration{
			HeaderName: aws.String(plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName.String()),
		}
	}
	if !plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn.IsUnknown() && !plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn.IsNull() {
		params.AccessConfiguration.SecretsManagerAccessTokenConfiguration = &mediatailor.SecretsManagerAccessTokenConfiguration{
			SecretArn: aws.String(plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn.String()),
		}
	}
	if !plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey.IsUnknown() && !plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey.IsNull() {
		params.AccessConfiguration.SecretsManagerAccessTokenConfiguration = &mediatailor.SecretsManagerAccessTokenConfiguration{
			SecretStringKey: aws.String(plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey.String()),
		}
	}

	// Default Segment Delivery Configuration
	if !plan.DefaultSegmentDeliveryConfiguration.BaseUrl.IsUnknown() && !plan.DefaultSegmentDeliveryConfiguration.BaseUrl.IsNull() {
		params.DefaultSegmentDeliveryConfiguration = &mediatailor.DefaultSegmentDeliveryConfiguration{
			BaseUrl: aws.String(plan.DefaultSegmentDeliveryConfiguration.BaseUrl.String()),
		}
	}

	// HTTP Configuration
	params.HttpConfiguration.BaseUrl = aws.String(plan.HttpConfiguration.BaseUrl.String())

	// Segment Delivery Configurations
	if len(plan.SegmentDeliveryConfigurations) > 0 && plan.SegmentDeliveryConfigurations != nil {
		params.SegmentDeliveryConfigurations = make([]*mediatailor.SegmentDeliveryConfiguration, len(plan.SegmentDeliveryConfigurations))
		for i, v := range plan.SegmentDeliveryConfigurations {
			params.SegmentDeliveryConfigurations[i] = &mediatailor.SegmentDeliveryConfiguration{
				BaseUrl: aws.String(v.BaseUrl.String()),
				Name:    aws.String(v.Name.String()),
			}
		}
	}

	// Source Location Name
	params.SourceLocationName = aws.String(plan.SourceLocationName.String())

	return params
}

func readSourceLocationToPlan(plan resourceSourceLocationModel, sourceLocation mediatailor.CreateSourceLocationOutput) resourceSourceLocationModel {
	// Set state
	plan.ID = types.StringValue(*sourceLocation.SourceLocationName)
	if sourceLocation.AccessConfiguration != nil {
		if sourceLocation.AccessConfiguration.AccessType != nil && *sourceLocation.AccessConfiguration.AccessType != "" {
			plan.AccessConfiguration.AccessType = types.StringValue(*sourceLocation.AccessConfiguration.AccessType)
		}
		if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration != nil {
			if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName != nil && *sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName != "" {
				plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName = types.StringValue(*sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName)
			}
			if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn != nil && *sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn != "" {
				plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn = types.StringValue(*sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn)
			}
			if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey != nil && *sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey != "" {
				plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey = types.StringValue(*sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey)
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
		if sourceLocation.DefaultSegmentDeliveryConfiguration.BaseUrl != nil && *sourceLocation.DefaultSegmentDeliveryConfiguration.BaseUrl != "" {
			plan.DefaultSegmentDeliveryConfiguration.BaseUrl = types.StringValue(*sourceLocation.DefaultSegmentDeliveryConfiguration.BaseUrl)
		}
	}
	if sourceLocation.HttpConfiguration != nil {
		if sourceLocation.HttpConfiguration.BaseUrl != nil && *sourceLocation.HttpConfiguration.BaseUrl != "" {
			plan.HttpConfiguration.BaseUrl = types.StringValue(*sourceLocation.HttpConfiguration.BaseUrl)
		}
	}
	if sourceLocation.LastModifiedTime != nil {
		plan.LastModifiedTime = types.StringValue((aws.TimeValue(sourceLocation.LastModifiedTime)).String())
	}
	if sourceLocation.SegmentDeliveryConfigurations != nil && len(sourceLocation.SegmentDeliveryConfigurations) > 0 {
		for _, segmentDeliveryConfiguration := range sourceLocation.SegmentDeliveryConfigurations {
			if segmentDeliveryConfiguration.BaseUrl != nil && *segmentDeliveryConfiguration.BaseUrl != "" {
				plan.SegmentDeliveryConfigurations = append(plan.SegmentDeliveryConfigurations, segmentDeliveryConfigurationsRModel{
					BaseUrl: types.StringValue(*segmentDeliveryConfiguration.BaseUrl),
				})
			}
			if segmentDeliveryConfiguration.Name != nil && *segmentDeliveryConfiguration.Name != "" {
				plan.SegmentDeliveryConfigurations = append(plan.SegmentDeliveryConfigurations, segmentDeliveryConfigurationsRModel{
					Name: types.StringValue(*segmentDeliveryConfiguration.Name),
				})
			}
		}
	}
	if sourceLocation.SourceLocationName != nil && *sourceLocation.SourceLocationName != "" {
		plan.SourceLocationName = types.StringValue(*sourceLocation.SourceLocationName)
	}
	if sourceLocation.Tags != nil && len(sourceLocation.Tags) > 0 {
		for key, value := range sourceLocation.Tags {
			plan.Tags[key] = value
		}
	}

	return plan
}

func readSourceLocationToPlanUpdate(plan resourceSourceLocationModel, sourceLocation mediatailor.UpdateSourceLocationOutput) resourceSourceLocationModel {
	// Set state
	plan.ID = types.StringValue(*sourceLocation.SourceLocationName)
	if sourceLocation.AccessConfiguration != nil {
		if sourceLocation.AccessConfiguration.AccessType != nil && *sourceLocation.AccessConfiguration.AccessType != "" {
			plan.AccessConfiguration.AccessType = types.StringValue(*sourceLocation.AccessConfiguration.AccessType)
		}
		if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration != nil {
			if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName != nil && *sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName != "" {
				plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName = types.StringValue(*sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName)
			}
			if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn != nil && *sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn != "" {
				plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn = types.StringValue(*sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn)
			}
			if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey != nil && *sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey != "" {
				plan.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey = types.StringValue(*sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey)
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
		if sourceLocation.DefaultSegmentDeliveryConfiguration.BaseUrl != nil && *sourceLocation.DefaultSegmentDeliveryConfiguration.BaseUrl != "" {
			plan.DefaultSegmentDeliveryConfiguration.BaseUrl = types.StringValue(*sourceLocation.DefaultSegmentDeliveryConfiguration.BaseUrl)
		}
	}
	if sourceLocation.HttpConfiguration != nil {
		if sourceLocation.HttpConfiguration.BaseUrl != nil && *sourceLocation.HttpConfiguration.BaseUrl != "" {
			plan.HttpConfiguration.BaseUrl = types.StringValue(*sourceLocation.HttpConfiguration.BaseUrl)
		}
	}
	if sourceLocation.LastModifiedTime != nil {
		plan.LastModifiedTime = types.StringValue((aws.TimeValue(sourceLocation.LastModifiedTime)).String())
	}
	if sourceLocation.SegmentDeliveryConfigurations != nil && len(sourceLocation.SegmentDeliveryConfigurations) > 0 {
		for _, segmentDeliveryConfiguration := range sourceLocation.SegmentDeliveryConfigurations {
			if segmentDeliveryConfiguration.BaseUrl != nil && *segmentDeliveryConfiguration.BaseUrl != "" {
				plan.SegmentDeliveryConfigurations = append(plan.SegmentDeliveryConfigurations, segmentDeliveryConfigurationsRModel{
					BaseUrl: types.StringValue(*segmentDeliveryConfiguration.BaseUrl),
				})
			}
			if segmentDeliveryConfiguration.Name != nil && *segmentDeliveryConfiguration.Name != "" {
				plan.SegmentDeliveryConfigurations = append(plan.SegmentDeliveryConfigurations, segmentDeliveryConfigurationsRModel{
					Name: types.StringValue(*segmentDeliveryConfiguration.Name),
				})
			}
		}
	}
	if sourceLocation.SourceLocationName != nil && *sourceLocation.SourceLocationName != "" {
		plan.SourceLocationName = types.StringValue(*sourceLocation.SourceLocationName)
	}
	if sourceLocation.Tags != nil && len(sourceLocation.Tags) > 0 {
		for key, value := range sourceLocation.Tags {
			plan.Tags[key] = value
		}
	}

	return plan
}

func readSourceLocationToState(state resourceSourceLocationModel, sourceLocation mediatailor.DescribeSourceLocationOutput) resourceSourceLocationModel {
	state.ID = types.StringValue(*sourceLocation.SourceLocationName)
	if sourceLocation.AccessConfiguration != nil {
		if sourceLocation.AccessConfiguration.AccessType != nil && *sourceLocation.AccessConfiguration.AccessType != "" {
			state.AccessConfiguration.AccessType = types.StringValue(*sourceLocation.AccessConfiguration.AccessType)
		}
		if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration != nil {
			if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName != nil && *sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName != "" {
				state.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName = types.StringValue(*sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName)
			}
			if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn != nil && *sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn != "" {
				state.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn = types.StringValue(*sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn)
			}
			if sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey != nil && *sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey != "" {
				state.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey = types.StringValue(*sourceLocation.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey)
			}
		}
	}
	if sourceLocation.Arn != nil && *sourceLocation.Arn != "" {
		state.Arn = types.StringValue(*sourceLocation.Arn)
	}
	if sourceLocation.CreationTime != nil {
		state.CreationTime = types.StringValue((aws.TimeValue(sourceLocation.CreationTime)).String())
	}
	if sourceLocation.DefaultSegmentDeliveryConfiguration != nil {
		if sourceLocation.DefaultSegmentDeliveryConfiguration.BaseUrl != nil && *sourceLocation.DefaultSegmentDeliveryConfiguration.BaseUrl != "" {
			state.DefaultSegmentDeliveryConfiguration.BaseUrl = types.StringValue(*sourceLocation.DefaultSegmentDeliveryConfiguration.BaseUrl)
		}
	}
	if sourceLocation.HttpConfiguration != nil {
		if sourceLocation.HttpConfiguration.BaseUrl != nil && *sourceLocation.HttpConfiguration.BaseUrl != "" {
			state.HttpConfiguration.BaseUrl = types.StringValue(*sourceLocation.HttpConfiguration.BaseUrl)
		}
	}
	if sourceLocation.LastModifiedTime != nil {
		state.LastModifiedTime = types.StringValue((aws.TimeValue(sourceLocation.LastModifiedTime)).String())
	}
	if sourceLocation.SegmentDeliveryConfigurations != nil && len(sourceLocation.SegmentDeliveryConfigurations) > 0 {
		for _, segmentDeliveryConfiguration := range sourceLocation.SegmentDeliveryConfigurations {
			if segmentDeliveryConfiguration.BaseUrl != nil && *segmentDeliveryConfiguration.BaseUrl != "" {
				state.SegmentDeliveryConfigurations = append(state.SegmentDeliveryConfigurations, segmentDeliveryConfigurationsRModel{
					BaseUrl: types.StringValue(*segmentDeliveryConfiguration.BaseUrl),
				})
			}
			if segmentDeliveryConfiguration.Name != nil && *segmentDeliveryConfiguration.Name != "" {
				state.SegmentDeliveryConfigurations = append(state.SegmentDeliveryConfigurations, segmentDeliveryConfigurationsRModel{
					Name: types.StringValue(*segmentDeliveryConfiguration.Name),
				})
			}
		}
	}
	if sourceLocation.SourceLocationName != nil && *sourceLocation.SourceLocationName != "" {
		state.SourceLocationName = types.StringValue(*sourceLocation.SourceLocationName)
	}
	if sourceLocation.Tags != nil && len(sourceLocation.Tags) > 0 {
		for key, value := range sourceLocation.Tags {
			state.Tags[key] = value
		}
	}

	return state
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
