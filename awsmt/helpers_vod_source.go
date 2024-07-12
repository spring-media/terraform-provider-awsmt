package awsmt

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/mediatailor/types"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getCreateVodSourceInput(plan vodSourceModel) mediatailor.CreateVodSourceInput {
	var input mediatailor.CreateVodSourceInput

	input.HttpPackageConfigurations, input.VodSourceName, input.SourceLocationName = getBasicVodSourceInput(&plan)

	if plan.Tags != nil && len(plan.Tags) > 0 {
		input.Tags = plan.Tags
	}

	return input
}

func getBasicVodSourceInput(plan *vodSourceModel) ([]awsTypes.HttpPackageConfiguration, *string, *string) {
	var httpPackageConfigurations []awsTypes.HttpPackageConfiguration
	var vodSourceName *string
	var sourceLocationName *string

	if plan.HttpPackageConfigurations != nil && len(plan.HttpPackageConfigurations) > 0 {
		httpPackageConfigurations = buildHttpPackageConfigurations(plan.HttpPackageConfigurations)
	}

	if plan.Name != nil {
		vodSourceName = plan.Name
	}

	if plan.SourceLocationName != nil {
		sourceLocationName = plan.SourceLocationName
	}
	return httpPackageConfigurations, vodSourceName, sourceLocationName
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
		plan.CreationTime = types.StringValue(vodSource.CreationTime.String())
	}

	if vodSource.HttpPackageConfigurations != nil && len(vodSource.HttpPackageConfigurations) > 0 {
		plan.HttpPackageConfigurations = []httpPackageConfigurationsModel{}
		for _, c := range vodSource.HttpPackageConfigurations {
			plan.HttpPackageConfigurations = append(plan.HttpPackageConfigurations, httpPackageConfigurationsModel{
				Type:        aws.String(string(c.Type)),
				Path:        c.Path,
				SourceGroup: c.SourceGroup,
			})
		}
	}

	if vodSource.LastModifiedTime != nil {
		plan.LastModifiedTime = types.StringValue(vodSource.LastModifiedTime.String())
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

func readVodSourceToState(plan vodSourceModel, vodSource mediatailor.DescribeVodSourceOutput) vodSourceModel {
	vodSourceName := *vodSource.VodSourceName
	sourceLocationName := *vodSource.SourceLocationName
	idNames := sourceLocationName + "," + vodSourceName

	plan.ID = types.StringValue(idNames)

	if vodSource.AdBreakOpportunities != nil {
		for _, value := range vodSource.AdBreakOpportunities {
			plan.AdBreakOpportunitiesOffsetMillis = append(plan.AdBreakOpportunitiesOffsetMillis, aws.Int64(value.OffsetMillis))
		}
	}

	if vodSource.HttpPackageConfigurations != nil && len(vodSource.HttpPackageConfigurations) > 0 {
		plan.HttpPackageConfigurations = []httpPackageConfigurationsModel{}
		for _, c := range vodSource.HttpPackageConfigurations {
			plan.HttpPackageConfigurations = append(plan.HttpPackageConfigurations, httpPackageConfigurationsModel{
				Type:        aws.String(string(c.Type)),
				Path:        c.Path,
				SourceGroup: c.SourceGroup,
			})
		}
	}
	if vodSource.Arn != nil {
		plan.Arn = types.StringValue(*vodSource.Arn)
	}

	if vodSource.LastModifiedTime != nil {
		plan.LastModifiedTime = types.StringValue(vodSource.LastModifiedTime.String())
	}

	if vodSource.CreationTime != nil {
		plan.CreationTime = types.StringValue(vodSource.CreationTime.String())
	}

	if vodSource.SourceLocationName != nil {
		plan.SourceLocationName = vodSource.SourceLocationName
	}

	if vodSource.VodSourceName != nil {
		plan.Name = vodSource.VodSourceName
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
