package awsmt

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func vodSourceInput(plan resourceVodSourceModel) mediatailor.CreateVodSourceInput {
	var input mediatailor.CreateVodSourceInput

	if plan.HttpPackageConfigurations != nil && len(plan.HttpPackageConfigurations) > 0 {
		for _, httpPackageConfiguration := range plan.HttpPackageConfigurations {
			httpPackageConfigurations := &mediatailor.HttpPackageConfiguration{}
			httpPackageConfigurations.Path = aws.String(httpPackageConfiguration.Path.String())
			httpPackageConfigurations.SourceGroup = aws.String(httpPackageConfiguration.SourceGroup.String())
			httpPackageConfigurations.Type = aws.String(httpPackageConfiguration.Type.String())
			input.HttpPackageConfigurations = append(input.HttpPackageConfigurations, httpPackageConfigurations)
		}
	}

	if !plan.VodSourceName.IsUnknown() && !plan.VodSourceName.IsNull() {
		input.VodSourceName = aws.String(plan.VodSourceName.String())
	}

	if !plan.SourceLocationName.IsUnknown() && !plan.SourceLocationName.IsNull() {
		input.SourceLocationName = aws.String(plan.SourceLocationName.String())
	}

	if plan.Tags != nil && len(plan.Tags) > 0 {
		input.Tags = plan.Tags
	}

	return input
}

func readVodSourceToPlan(plan resourceVodSourceModel, vodSource mediatailor.CreateVodSourceOutput) resourceVodSourceModel {
	plan.ID = types.StringValue(*vodSource.Arn)

	if vodSource.Arn != nil {
		plan.Arn = types.StringValue(*vodSource.Arn)
	}

	if vodSource.CreationTime != nil {
		plan.CreationTime = types.StringValue((aws.TimeValue(vodSource.CreationTime)).String())
	}

	if vodSource.HttpPackageConfigurations != nil && len(vodSource.HttpPackageConfigurations) > 0 {
		for _, httpPackageConfiguration := range vodSource.HttpPackageConfigurations {
			httpPackageConfigurations := httpPackageConfigurationsVSRModel{}
			httpPackageConfigurations.Path = types.StringValue(*httpPackageConfiguration.Path)
			httpPackageConfigurations.SourceGroup = types.StringValue(*httpPackageConfiguration.SourceGroup)
			httpPackageConfigurations.Type = types.StringValue(*httpPackageConfiguration.Type)
			plan.HttpPackageConfigurations = append(plan.HttpPackageConfigurations, httpPackageConfigurations)
		}
	}

	if vodSource.LastModifiedTime != nil {
		plan.LastModifiedTime = types.StringValue((aws.TimeValue(vodSource.LastModifiedTime)).String())
	}

	if vodSource.VodSourceName != nil {
		plan.VodSourceName = types.StringValue(*vodSource.VodSourceName)
	}

	if vodSource.SourceLocationName != nil {
		plan.SourceLocationName = types.StringValue(*vodSource.SourceLocationName)
	}

	if vodSource.Tags != nil && len(vodSource.Tags) > 0 {
		plan.Tags = vodSource.Tags
	}

	return plan
}

func readVodSourceToState(state resourceVodSourceModel, vodSource mediatailor.DescribeVodSourceOutput) resourceVodSourceModel {
	state.ID = types.StringValue(*vodSource.Arn)

	if vodSource.Arn != nil {
		state.Arn = types.StringValue(*vodSource.Arn)
	}

	if vodSource.CreationTime != nil {
		state.CreationTime = types.StringValue((aws.TimeValue(vodSource.CreationTime)).String())
	}

	if vodSource.HttpPackageConfigurations != nil && len(vodSource.HttpPackageConfigurations) > 0 {
		for _, httpPackageConfiguration := range vodSource.HttpPackageConfigurations {
			httpPackageConfigurations := httpPackageConfigurationsVSRModel{}
			httpPackageConfigurations.Path = types.StringValue(*httpPackageConfiguration.Path)
			httpPackageConfigurations.SourceGroup = types.StringValue(*httpPackageConfiguration.SourceGroup)
			httpPackageConfigurations.Type = types.StringValue(*httpPackageConfiguration.Type)
			state.HttpPackageConfigurations = append(state.HttpPackageConfigurations, httpPackageConfigurations)
		}
	}

	if vodSource.LastModifiedTime != nil {
		state.LastModifiedTime = types.StringValue((aws.TimeValue(vodSource.LastModifiedTime)).String())
	}

	if vodSource.VodSourceName != nil {
		state.VodSourceName = types.StringValue(*vodSource.VodSourceName)
	}

	if vodSource.SourceLocationName != nil {
		state.SourceLocationName = types.StringValue(*vodSource.SourceLocationName)
	}

	if vodSource.Tags != nil && len(vodSource.Tags) > 0 {
		state.Tags = vodSource.Tags
	}

	return state
}

func vodSourceUpdateInput(plan resourceVodSourceModel) mediatailor.UpdateVodSourceInput {
	var input mediatailor.UpdateVodSourceInput

	if plan.HttpPackageConfigurations != nil && len(plan.HttpPackageConfigurations) > 0 {
		for _, httpPackageConfiguration := range plan.HttpPackageConfigurations {
			httpPackageConfigurations := &mediatailor.HttpPackageConfiguration{}
			httpPackageConfigurations.Path = aws.String(httpPackageConfiguration.Path.String())
			httpPackageConfigurations.SourceGroup = aws.String(httpPackageConfiguration.SourceGroup.String())
			httpPackageConfigurations.Type = aws.String(httpPackageConfiguration.Type.String())
			input.HttpPackageConfigurations = append(input.HttpPackageConfigurations, httpPackageConfigurations)
		}
	}

	if !plan.VodSourceName.IsUnknown() && !plan.VodSourceName.IsNull() {
		input.VodSourceName = aws.String(plan.VodSourceName.String())
	}

	if !plan.SourceLocationName.IsUnknown() && !plan.SourceLocationName.IsNull() {
		input.SourceLocationName = aws.String(plan.SourceLocationName.String())
	}

	return input
}

func readUpdatedVodSourceToPlan(plan resourceVodSourceModel, vodSource mediatailor.UpdateVodSourceOutput) resourceVodSourceModel {
	plan.ID = types.StringValue(*vodSource.Arn)

	if vodSource.Arn != nil {
		plan.Arn = types.StringValue(*vodSource.Arn)
	}

	if vodSource.CreationTime != nil {
		plan.CreationTime = types.StringValue((aws.TimeValue(vodSource.CreationTime)).String())
	}

	if vodSource.HttpPackageConfigurations != nil && len(vodSource.HttpPackageConfigurations) > 0 {
		for _, httpPackageConfiguration := range vodSource.HttpPackageConfigurations {
			httpPackageConfigurations := httpPackageConfigurationsVSRModel{}
			httpPackageConfigurations.Path = types.StringValue(*httpPackageConfiguration.Path)
			httpPackageConfigurations.SourceGroup = types.StringValue(*httpPackageConfiguration.SourceGroup)
			httpPackageConfigurations.Type = types.StringValue(*httpPackageConfiguration.Type)
			plan.HttpPackageConfigurations = append(plan.HttpPackageConfigurations, httpPackageConfigurations)
		}
	}

	if vodSource.LastModifiedTime != nil {
		plan.LastModifiedTime = types.StringValue((aws.TimeValue(vodSource.LastModifiedTime)).String())
	}

	if vodSource.VodSourceName != nil {
		plan.VodSourceName = types.StringValue(*vodSource.VodSourceName)
	}

	if vodSource.SourceLocationName != nil {
		plan.SourceLocationName = types.StringValue(*vodSource.SourceLocationName)
	}

	if vodSource.Tags != nil && len(vodSource.Tags) > 0 {
		plan.Tags = vodSource.Tags
	}

	return plan
}
