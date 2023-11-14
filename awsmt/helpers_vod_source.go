package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func vodSourceInput(plan vodSourceModel) mediatailor.CreateVodSourceInput {
	var input mediatailor.CreateVodSourceInput

	input.HttpPackageConfigurations, input.VodSourceName, input.SourceLocationName = getBasicVodSourceInput(&plan)

	if plan.Tags != nil && len(plan.Tags) > 0 {
		input.Tags = plan.Tags
	}

	return input
}

func readVodSourceToPlan(plan vodSourceModel, vodSource mediatailor.CreateVodSourceOutput) vodSourceModel {
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
		plan.HttpPackageConfigurations = []httpPackageConfigurationsModel{}
		for _, httpPackageConfiguration := range vodSource.HttpPackageConfigurations {
			httpPackageConfigurations := httpPackageConfigurationsModel{}
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
		plan.Name = vodSource.VodSourceName
	}

	if vodSource.SourceLocationName != nil {
		plan.SourceLocationName = vodSource.SourceLocationName
	}

	if len(vodSource.Tags) > 0 {
		plan.Tags = vodSource.Tags
	}

	return plan
}

func readVodSourceToState(ctx context.Context, plan vodSourceModel, vodSource mediatailor.DescribeVodSourceOutput) vodSourceModel {
	vodSourceName := *vodSource.VodSourceName
	sourceLocationName := *vodSource.SourceLocationName
	idNames := sourceLocationName + "," + vodSourceName

	plan.ID = types.StringValue(idNames)

	if vodSource.AdBreakOpportunities != nil {
		for _, value := range vodSource.AdBreakOpportunities {
			plan.AdBreakOpportunitiesOffsetMillis = append(plan.AdBreakOpportunitiesOffsetMillis, value.OffsetMillis)
		}
	}

	if vodSource.Arn != nil {
		plan.Arn = types.StringValue(*vodSource.Arn)
	}

	if vodSource.CreationTime != nil {
		plan.CreationTime = types.StringValue((aws.TimeValue(vodSource.CreationTime)).String())
	}

	if vodSource.HttpPackageConfigurations != nil && len(vodSource.HttpPackageConfigurations) > 0 {
		plan.HttpPackageConfigurations = []httpPackageConfigurationsModel{}
		for _, httpPackageConfiguration := range vodSource.HttpPackageConfigurations {
			httpPackageConfigurations := httpPackageConfigurationsModel{}
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
		plan.Name = vodSource.VodSourceName
	}

	if vodSource.SourceLocationName != nil {
		plan.SourceLocationName = vodSource.SourceLocationName
	}

	if len(vodSource.Tags) > 0 {
		plan.Tags = vodSource.Tags
	}

	return plan
}

func vodSourceUpdateInput(plan vodSourceModel) mediatailor.UpdateVodSourceInput {
	var input mediatailor.UpdateVodSourceInput

	input.HttpPackageConfigurations, input.VodSourceName, input.SourceLocationName = getBasicVodSourceInput(&plan)

	return input
}

func getBasicVodSourceInput(plan *vodSourceModel) ([]*mediatailor.HttpPackageConfiguration, *string, *string) {
	var httpPackageConfigurations []*mediatailor.HttpPackageConfiguration
	var vodSourceName *string
	var sourceLocationName *string

	if plan.HttpPackageConfigurations != nil && len(plan.HttpPackageConfigurations) > 0 {
		httpPackageConfigurations = getHttpInput(plan.HttpPackageConfigurations)
	}

	if plan.Name != nil {
		vodSourceName = plan.Name
	}

	if plan.SourceLocationName != nil {
		sourceLocationName = plan.SourceLocationName
	}
	return httpPackageConfigurations, vodSourceName, sourceLocationName
}
