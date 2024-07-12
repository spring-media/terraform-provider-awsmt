package awsmt

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/mediatailor/types"
)

func readHttpPackageConfigurations(configurations []awsTypes.HttpPackageConfiguration) []httpPackageConfigurationsModel {
	var httpPackageConfiguration []httpPackageConfigurationsModel
	if len(configurations) > 0 {
		for _, c := range configurations {
			httpPackageConfiguration = append(httpPackageConfiguration, httpPackageConfigurationsModel{
				Path:        c.Path,
				SourceGroup: c.SourceGroup,
				Type:        aws.String(string(c.Type)),
			})
		}
	}
	return httpPackageConfiguration
}

func buildHttpPackageConfigurations(plan []httpPackageConfigurationsModel) []awsTypes.HttpPackageConfiguration {
	var tmp []awsTypes.HttpPackageConfiguration
	if len(plan) > 0 {
		for _, c := range plan {
			var cType awsTypes.Type
			if *c.Type == "DASH" {
				cType = awsTypes.TypeDash
			} else {
				cType = awsTypes.TypeHls
			}
			tmp = append(tmp, awsTypes.HttpPackageConfiguration{
				Path:        c.Path,
				SourceGroup: c.SourceGroup,
				Type:        cType,
			})
		}
	}
	return tmp
}
