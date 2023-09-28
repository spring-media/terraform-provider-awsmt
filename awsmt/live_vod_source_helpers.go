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
