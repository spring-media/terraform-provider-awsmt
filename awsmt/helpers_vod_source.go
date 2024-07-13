package awsmt

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/mediatailor/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func getCreateVodSourceInput(plan vodSourceModel) *mediatailor.CreateVodSourceInput {
	var input mediatailor.CreateVodSourceInput

	input.HttpPackageConfigurations, input.VodSourceName, input.SourceLocationName = getSharedVodSourceInput(&plan)

	if plan.Tags != nil && len(plan.Tags) > 0 {
		input.Tags = plan.Tags
	}

	return &input
}

func getUpdateVodSourceInput(model vodSourceModel) mediatailor.UpdateVodSourceInput {
	var input mediatailor.UpdateVodSourceInput

	input.HttpPackageConfigurations, input.VodSourceName, input.SourceLocationName = getSharedVodSourceInput(&model)

	return input
}

func getSharedVodSourceInput(model *vodSourceModel) ([]awsTypes.HttpPackageConfiguration, *string, *string) {
	var httpPackageConfigurations []awsTypes.HttpPackageConfiguration
	var vodSourceName *string
	var sourceLocationName *string

	if model.HttpPackageConfigurations != nil && len(model.HttpPackageConfigurations) > 0 {
		httpPackageConfigurations = getHttpPackageConfigurations(model.HttpPackageConfigurations)
	}

	if model.Name != nil {
		vodSourceName = model.Name
	}

	if model.SourceLocationName != nil {
		sourceLocationName = model.SourceLocationName
	}
	return httpPackageConfigurations, vodSourceName, sourceLocationName
}

// the readVodSourceToPlan is used to convert the output from the create and update operations to the plan
func readVodSourceToPlan(model vodSourceModel, vodSource mediatailor.CreateVodSourceOutput) vodSourceModel {
	vodSourceName := *vodSource.VodSourceName
	sourceLocationName := *vodSource.SourceLocationName
	idNames := sourceLocationName + "," + vodSourceName

	model.ID = types.StringValue(idNames)

	if vodSource.Arn != nil {
		model.Arn = types.StringValue(*vodSource.Arn)
	}

	if vodSource.CreationTime != nil {
		model.CreationTime = types.StringValue(vodSource.CreationTime.String())
	}

	if vodSource.HttpPackageConfigurations != nil && len(vodSource.HttpPackageConfigurations) > 0 {
		model.HttpPackageConfigurations = []httpPackageConfigurationsModel{}
		for _, c := range vodSource.HttpPackageConfigurations {
			model.HttpPackageConfigurations = append(model.HttpPackageConfigurations, httpPackageConfigurationsModel{
				Type:        aws.String(string(c.Type)),
				Path:        c.Path,
				SourceGroup: c.SourceGroup,
			})
		}
	}

	if vodSource.LastModifiedTime != nil {
		model.LastModifiedTime = types.StringValue(vodSource.LastModifiedTime.String())
	}

	if vodSource.SourceLocationName != nil {
		model.SourceLocationName = vodSource.SourceLocationName
	}

	if vodSource.VodSourceName != nil {
		model.Name = vodSource.VodSourceName
	}

	if len(vodSource.Tags) > 0 {
		model.Tags = vodSource.Tags
	}

	model.AdBreakOpportunitiesOffsetMillis, _ = types.ListValue(types.Int64Type, []attr.Value{})

	return model
}

// the readVodSourceToState is used to convert the output from the describe operation to the state
func readVodSourceToState(model vodSourceModel, vodSource mediatailor.DescribeVodSourceOutput) vodSourceModel {

	model = readVodSourceToPlan(model, mediatailor.CreateVodSourceOutput{
		Arn:                       vodSource.Arn,
		CreationTime:              vodSource.CreationTime,
		HttpPackageConfigurations: vodSource.HttpPackageConfigurations,
		LastModifiedTime:          vodSource.LastModifiedTime,
		SourceLocationName:        vodSource.SourceLocationName,
		VodSourceName:             vodSource.VodSourceName,
		Tags:                      vodSource.Tags,
	})

	if vodSource.AdBreakOpportunities != nil && len(vodSource.AdBreakOpportunities) > 0 {
		var offsetMillisElements []attr.Value
		for _, value := range vodSource.AdBreakOpportunities {
			offsetMillisElements = append(offsetMillisElements, types.Int64Value(value.OffsetMillis))
		}
		offsetMillisList, _ := types.ListValue(types.Int64Type, offsetMillisElements)

		model.AdBreakOpportunitiesOffsetMillis = offsetMillisList
	}

	return model
}
