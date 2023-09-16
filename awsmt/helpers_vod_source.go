package awsmt

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func vodSourceInput(plan resourceVodSourceModel) mediatailor.CreateVodSourceInput {
	var input mediatailor.CreateVodSourceInput

	input.HttpPackageConfigurations, input.VodSourceName, input.SourceLocationName = getBasicVodSourceInput(&plan)

	if plan.Tags != nil && len(plan.Tags) > 0 {
		input.Tags = plan.Tags
	}

	return input
}

func readVodSourceToPlan(plan resourceVodSourceModel, vodSource mediatailor.CreateVodSourceOutput) resourceVodSourceModel {
	vodSourceName := *vodSource.VodSourceName
	sourceLocationName := *vodSource.SourceLocationName
	idNames := sourceLocationName + "," + vodSourceName

	plan.ID = types.StringValue(idNames)

	if vodSource.Arn != nil {
		plan.Arn = types.StringValue(*vodSource.Arn)
	}

	if vodSource.CreationTime != nil {
		plan.CreationTime = types.StringValue((aws.TimeValue(vodSource.CreationTime)).String())
	}

	if vodSource.HttpPackageConfigurations != nil && len(vodSource.HttpPackageConfigurations) > 0 {
		plan.HttpPackageConfigurations = []httpPackageConfigurationsVSRModel{}
		for _, httpPackageConfiguration := range vodSource.HttpPackageConfigurations {
			httpPackageConfigurations := httpPackageConfigurationsVSRModel{}
			httpPackageConfigurations.Path = httpPackageConfiguration.Path
			httpPackageConfigurations.SourceGroup = httpPackageConfiguration.SourceGroup
			httpPackageConfigurations.Type = httpPackageConfiguration.Type
			plan.HttpPackageConfigurations = append(plan.HttpPackageConfigurations, httpPackageConfigurations)
		}
	}

	if vodSource.LastModifiedTime != nil {
		plan.LastModifiedTime = types.StringValue((aws.TimeValue(vodSource.LastModifiedTime)).String())
	}

	if vodSource.VodSourceName != nil {
		plan.VodSourceName = vodSource.VodSourceName
	}

	if vodSource.SourceLocationName != nil {
		plan.SourceLocationName = vodSource.SourceLocationName
	}

	if vodSource.Tags != nil && len(vodSource.Tags) > 0 {
		plan.Tags = vodSource.Tags
	}

	return plan
}

func vodSourceUpdateInput(plan resourceVodSourceModel) mediatailor.UpdateVodSourceInput {
	var input mediatailor.UpdateVodSourceInput

	input.HttpPackageConfigurations, input.VodSourceName, input.SourceLocationName = getBasicVodSourceInput(&plan)

	return input
}

func getBasicVodSourceInput(plan *resourceVodSourceModel) ([]*mediatailor.HttpPackageConfiguration, *string, *string) {
	var httpPackageConfigurations []*mediatailor.HttpPackageConfiguration
	var vodSourceName *string
	var sourceLocationName *string

	if plan.HttpPackageConfigurations != nil && len(plan.HttpPackageConfigurations) > 0 {
		httpPackageConfigurations = getHttpInput(plan.HttpPackageConfigurations)
	}

	if plan.VodSourceName != nil {
		vodSourceName = plan.VodSourceName
	}

	if plan.SourceLocationName != nil {
		sourceLocationName = plan.SourceLocationName
	}
	return httpPackageConfigurations, vodSourceName, sourceLocationName
}

func getHttpInput(plan []httpPackageConfigurationsVSRModel) []*mediatailor.HttpPackageConfiguration {
	var input mediatailor.CreateVodSourceInput
	if len(plan) > 0 {
		input.HttpPackageConfigurations = []*mediatailor.HttpPackageConfiguration{}
		for _, httpPackageConfiguration := range plan {
			httpPackageConfigurations := &mediatailor.HttpPackageConfiguration{}
			httpPackageConfigurations.Path = httpPackageConfiguration.Path
			httpPackageConfigurations.SourceGroup = httpPackageConfiguration.SourceGroup
			httpPackageConfigurations.Type = httpPackageConfiguration.Type
			input.HttpPackageConfigurations = append(input.HttpPackageConfigurations, httpPackageConfigurations)
		}
	}
	return input.HttpPackageConfigurations
}
