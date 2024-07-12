package awsmt

import (
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getLiveSourceInput(plan liveSourceModel) mediatailor.CreateLiveSourceInput {
	var input mediatailor.CreateLiveSourceInput

	if plan.HttpPackageConfigurations != nil && len(plan.HttpPackageConfigurations) > 0 {
		input.HttpPackageConfigurations = buildHttpPackageConfigurations(plan.HttpPackageConfigurations)
	}

	if plan.Name != nil {
		input.LiveSourceName = plan.Name
	}

	if plan.SourceLocationName != nil {
		input.SourceLocationName = plan.SourceLocationName
	}

	if plan.Tags != nil && len(plan.Tags) > 0 {
		input.Tags = plan.Tags
	}

	return input
}

func readLiveSourceToPlan(plan liveSourceModel, liveSource mediatailor.CreateLiveSourceOutput) liveSourceModel {
	liveSourceName := *liveSource.LiveSourceName
	sourceLocationName := *liveSource.SourceLocationName
	idNames := sourceLocationName + "," + liveSourceName

	plan.ID = types.StringValue(idNames)

	if liveSource.Arn != nil {
		plan.Arn = types.StringValue(*liveSource.Arn)
	}

	if liveSource.CreationTime != nil {
		plan.CreationTime = types.StringValue(liveSource.CreationTime.String())
	}

	plan.HttpPackageConfigurations = readHttpPackageConfigurations(liveSource.HttpPackageConfigurations)

	if liveSource.LastModifiedTime != nil {
		plan.LastModifiedTime = types.StringValue(liveSource.LastModifiedTime.String())
	}

	if liveSource.LiveSourceName != nil {
		plan.Name = liveSource.LiveSourceName
	}

	if liveSource.SourceLocationName != nil {
		plan.SourceLocationName = liveSource.SourceLocationName
	}

	if len(liveSource.Tags) > 0 {
		plan.Tags = liveSource.Tags
	}

	return plan
}

func liveSourceUpdateInput(plan liveSourceModel) mediatailor.UpdateLiveSourceInput {
	var input mediatailor.UpdateLiveSourceInput

	if plan.HttpPackageConfigurations != nil && len(plan.HttpPackageConfigurations) > 0 {
		input.HttpPackageConfigurations = buildHttpPackageConfigurations(plan.HttpPackageConfigurations)
	}

	if plan.Name != nil {
		input.LiveSourceName = plan.Name
	}

	if plan.SourceLocationName != nil {
		input.SourceLocationName = plan.SourceLocationName
	}

	return input
}
