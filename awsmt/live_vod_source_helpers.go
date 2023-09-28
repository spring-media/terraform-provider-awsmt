package awsmt

import (
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func readHttpPackageConfigurations(configurations []*mediatailor.HttpPackageConfiguration) []httpPackageConfigurationsModel {
	var httpPackageConfigurationsRead []httpPackageConfigurationsModel
	if len(configurations) > 0 {
		for _, httpPackageConfiguration := range configurations {
			httpPackageConfigurations := httpPackageConfigurationsModel{}
			httpPackageConfigurations.Path = types.StringValue(*httpPackageConfiguration.Path)
			httpPackageConfigurations.SourceGroup = types.StringValue(*httpPackageConfiguration.SourceGroup)
			httpPackageConfigurations.Type = types.StringValue(*httpPackageConfiguration.Type)
			httpPackageConfigurationsRead = append(httpPackageConfigurationsRead, httpPackageConfigurations)
		}
	}
	return httpPackageConfigurationsRead
}
