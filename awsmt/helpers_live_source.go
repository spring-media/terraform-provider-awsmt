package awsmt

import (
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/mediatailor/types"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getCreateLiveSourceInput(model liveSourceModel) *mediatailor.CreateLiveSourceInput {
	var input mediatailor.CreateLiveSourceInput

	input.HttpPackageConfigurations, input.LiveSourceName, input.SourceLocationName = getSharedLiveSourceInput(&model)

	if model.Tags != nil && len(model.Tags) > 0 {
		input.Tags = model.Tags
	}

	return &input
}

func getUpdateLiveSourceInput(model liveSourceModel) mediatailor.UpdateLiveSourceInput {
	var input mediatailor.UpdateLiveSourceInput

	input.HttpPackageConfigurations, input.LiveSourceName, input.SourceLocationName = getSharedLiveSourceInput(&model)

	return input
}

func getSharedLiveSourceInput(model *liveSourceModel) ([]awsTypes.HttpPackageConfiguration, *string, *string) {
	var httpPackageConfigurations []awsTypes.HttpPackageConfiguration
	var liveSourceName *string
	var sourceLocationName *string

	if model.HttpPackageConfigurations != nil && len(model.HttpPackageConfigurations) > 0 {
		httpPackageConfigurations = getHttpPackageConfigurations(model.HttpPackageConfigurations)
	}

	if model.Name != nil {
		liveSourceName = model.Name
	}

	if model.SourceLocationName != nil {
		sourceLocationName = model.SourceLocationName
	}
	return httpPackageConfigurations, liveSourceName, sourceLocationName
}

// readLiveSource is used for both plan and state since the output from create/update and describe is compatible
func readLiveSource(model liveSourceModel, liveSource mediatailor.CreateLiveSourceOutput) liveSourceModel {
	liveSourceName := *liveSource.LiveSourceName
	sourceLocationName := *liveSource.SourceLocationName
	idNames := sourceLocationName + "," + liveSourceName

	model.ID = types.StringValue(idNames)

	if liveSource.Arn != nil {
		model.Arn = types.StringValue(*liveSource.Arn)
	}

	if liveSource.CreationTime != nil {
		model.CreationTime = types.StringValue(liveSource.CreationTime.String())
	}

	model.HttpPackageConfigurations = readHttpPackageConfigurations(liveSource.HttpPackageConfigurations)

	if liveSource.LastModifiedTime != nil {
		model.LastModifiedTime = types.StringValue(liveSource.LastModifiedTime.String())
	}

	if liveSource.LiveSourceName != nil {
		model.Name = liveSource.LiveSourceName
	}

	if liveSource.SourceLocationName != nil {
		model.SourceLocationName = liveSource.SourceLocationName
	}

	if len(liveSource.Tags) > 0 {
		model.Tags = liveSource.Tags
	}

	return model
}
