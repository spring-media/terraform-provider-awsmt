package awsmt

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func liveSourceInput(plan resourceLiveSourceModel) mediatailor.CreateLiveSourceInput {
	var input mediatailor.CreateLiveSourceInput

	if plan.HttpPackageConfigurations != nil && len(plan.HttpPackageConfigurations) > 0 {
		for _, httpPackageConfiguration := range plan.HttpPackageConfigurations {
			httpPackageConfigurations := &mediatailor.HttpPackageConfiguration{}
			httpPackageConfigurations.Path = aws.String(httpPackageConfiguration.Path.String())
			httpPackageConfigurations.SourceGroup = aws.String(httpPackageConfiguration.SourceGroup.String())
			httpPackageConfigurations.Type = aws.String(httpPackageConfiguration.Type.String())
			input.HttpPackageConfigurations = append(input.HttpPackageConfigurations, httpPackageConfigurations)
		}
	}

	if !plan.LiveSourceName.IsUnknown() && !plan.LiveSourceName.IsNull() {
		input.LiveSourceName = aws.String(plan.LiveSourceName.String())
	}

	if !plan.SourceLocationName.IsUnknown() && !plan.SourceLocationName.IsNull() {
		input.SourceLocationName = aws.String(plan.SourceLocationName.String())
	}

	if plan.Tags != nil && len(plan.Tags) > 0 {
		input.Tags = plan.Tags
	}

	return input
}

func readLiveSourceToPlan(plan resourceLiveSourceModel, liveSource mediatailor.CreateLiveSourceOutput) resourceLiveSourceModel {
	plan.ID = types.StringValue(*liveSource.Arn)

	if liveSource.Arn != nil {
		plan.Arn = types.StringValue(*liveSource.Arn)
	}

	if liveSource.CreationTime != nil {
		plan.CreationTime = types.StringValue((aws.TimeValue(liveSource.CreationTime)).String())
	}

	if liveSource.HttpPackageConfigurations != nil && len(liveSource.HttpPackageConfigurations) > 0 {
		for _, httpPackageConfiguration := range liveSource.HttpPackageConfigurations {
			httpPackageConfigurations := httpPackageConfigurationsLSRModel{}
			httpPackageConfigurations.Path = types.StringValue(*httpPackageConfiguration.Path)
			httpPackageConfigurations.SourceGroup = types.StringValue(*httpPackageConfiguration.SourceGroup)
			httpPackageConfigurations.Type = types.StringValue(*httpPackageConfiguration.Type)
			plan.HttpPackageConfigurations = append(plan.HttpPackageConfigurations, httpPackageConfigurations)
		}
	}

	if liveSource.LastModifiedTime != nil {
		plan.LastModifiedTime = types.StringValue((aws.TimeValue(liveSource.LastModifiedTime)).String())
	}

	if liveSource.LiveSourceName != nil {
		plan.LiveSourceName = types.StringValue(*liveSource.LiveSourceName)
	}

	if liveSource.SourceLocationName != nil {
		plan.SourceLocationName = types.StringValue(*liveSource.SourceLocationName)
	}

	if liveSource.Tags != nil && len(liveSource.Tags) > 0 {
		plan.Tags = liveSource.Tags
	}

	return plan
}

func readLiveSourceToState(state resourceLiveSourceModel, liveSource mediatailor.DescribeLiveSourceOutput) resourceLiveSourceModel {
	state.ID = types.StringValue(*liveSource.Arn)

	if liveSource.Arn != nil {
		state.Arn = types.StringValue(*liveSource.Arn)
	}

	if liveSource.CreationTime != nil {
		state.CreationTime = types.StringValue((aws.TimeValue(liveSource.CreationTime)).String())
	}

	if liveSource.HttpPackageConfigurations != nil && len(liveSource.HttpPackageConfigurations) > 0 {
		for _, httpPackageConfiguration := range liveSource.HttpPackageConfigurations {
			httpPackageConfigurations := httpPackageConfigurationsLSRModel{}
			httpPackageConfigurations.Path = types.StringValue(*httpPackageConfiguration.Path)
			httpPackageConfigurations.SourceGroup = types.StringValue(*httpPackageConfiguration.SourceGroup)
			httpPackageConfigurations.Type = types.StringValue(*httpPackageConfiguration.Type)
			state.HttpPackageConfigurations = append(state.HttpPackageConfigurations, httpPackageConfigurations)
		}
	}

	if liveSource.LastModifiedTime != nil {
		state.LastModifiedTime = types.StringValue((aws.TimeValue(liveSource.LastModifiedTime)).String())
	}

	if liveSource.LiveSourceName != nil {
		state.LiveSourceName = types.StringValue(*liveSource.LiveSourceName)
	}

	if liveSource.SourceLocationName != nil {
		state.SourceLocationName = types.StringValue(*liveSource.SourceLocationName)
	}

	if liveSource.Tags != nil && len(liveSource.Tags) > 0 {
		state.Tags = liveSource.Tags
	}

	return state
}

func liveSourceUpdateInput(plan resourceLiveSourceModel) mediatailor.UpdateLiveSourceInput {
	var input mediatailor.UpdateLiveSourceInput

	if plan.HttpPackageConfigurations != nil && len(plan.HttpPackageConfigurations) > 0 {
		for _, httpPackageConfiguration := range plan.HttpPackageConfigurations {
			httpPackageConfigurations := &mediatailor.HttpPackageConfiguration{}
			httpPackageConfigurations.Path = aws.String(httpPackageConfiguration.Path.String())
			httpPackageConfigurations.SourceGroup = aws.String(httpPackageConfiguration.SourceGroup.String())
			httpPackageConfigurations.Type = aws.String(httpPackageConfiguration.Type.String())
			input.HttpPackageConfigurations = append(input.HttpPackageConfigurations, httpPackageConfigurations)
		}
	}

	if !plan.LiveSourceName.IsUnknown() && !plan.LiveSourceName.IsNull() {
		input.LiveSourceName = aws.String(plan.LiveSourceName.String())
	}

	if !plan.SourceLocationName.IsUnknown() && !plan.SourceLocationName.IsNull() {
		input.SourceLocationName = aws.String(plan.SourceLocationName.String())
	}

	return input
}

func readUpdatedLiveSourceToPlan(plan resourceLiveSourceModel, liveSource mediatailor.UpdateLiveSourceOutput) resourceLiveSourceModel {
	plan.ID = types.StringValue(*liveSource.Arn)

	if liveSource.Arn != nil {
		plan.Arn = types.StringValue(*liveSource.Arn)
	}

	if liveSource.CreationTime != nil {
		plan.CreationTime = types.StringValue((aws.TimeValue(liveSource.CreationTime)).String())
	}

	if liveSource.HttpPackageConfigurations != nil && len(liveSource.HttpPackageConfigurations) > 0 {
		for _, httpPackageConfiguration := range liveSource.HttpPackageConfigurations {
			httpPackageConfigurations := httpPackageConfigurationsLSRModel{}
			httpPackageConfigurations.Path = types.StringValue(*httpPackageConfiguration.Path)
			httpPackageConfigurations.SourceGroup = types.StringValue(*httpPackageConfiguration.SourceGroup)
			httpPackageConfigurations.Type = types.StringValue(*httpPackageConfiguration.Type)
			plan.HttpPackageConfigurations = append(plan.HttpPackageConfigurations, httpPackageConfigurations)
		}
	}

	if liveSource.LastModifiedTime != nil {
		plan.LastModifiedTime = types.StringValue((aws.TimeValue(liveSource.LastModifiedTime)).String())
	}

	if liveSource.LiveSourceName != nil {
		plan.LiveSourceName = types.StringValue(*liveSource.LiveSourceName)
	}

	if liveSource.SourceLocationName != nil {
		plan.SourceLocationName = types.StringValue(*liveSource.SourceLocationName)
	}

	if liveSource.Tags != nil && len(liveSource.Tags) > 0 {
		plan.Tags = liveSource.Tags
	}

	return plan
}
