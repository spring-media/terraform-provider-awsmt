package awsmt

import (
	"github.com/aws/aws-sdk-go/service/mediatailor"
)

func readHttpPackageConfigurations(configurations []*mediatailor.HttpPackageConfiguration) []httpPackageConfigurationsModel {
	var httpPackageConfigurationsRead []httpPackageConfigurationsModel
	if len(configurations) > 0 {
		for _, httpPackageConfiguration := range configurations {
			httpPackageConfigurations := httpPackageConfigurationsModel{}
			httpPackageConfigurations.Path = httpPackageConfiguration.Path
			httpPackageConfigurations.SourceGroup = httpPackageConfiguration.SourceGroup
			httpPackageConfigurations.Type = httpPackageConfiguration.Type
			httpPackageConfigurationsRead = append(httpPackageConfigurationsRead, httpPackageConfigurations)
		}
	}
	return httpPackageConfigurationsRead
}

func getHttpInput(plan []httpPackageConfigurationsModel) []*mediatailor.HttpPackageConfiguration {
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
