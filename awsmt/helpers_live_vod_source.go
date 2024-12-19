package awsmt

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/mediatailor/types"
	"terraform-provider-mediatailor/awsmt/models"
)

func readHttpPackageConfigurations(configurations []awsTypes.HttpPackageConfiguration) []models.HttpPackageConfigurationsModel {
	var httpPackageConfiguration []models.HttpPackageConfigurationsModel
	if len(configurations) > 0 {
		for _, c := range configurations {
			httpPackageConfiguration = append(httpPackageConfiguration, models.HttpPackageConfigurationsModel{
				Path:        c.Path,
				SourceGroup: c.SourceGroup,
				Type:        aws.String(string(c.Type)),
			})
		}
	}
	return httpPackageConfiguration
}

func getHttpPackageConfigurations(plan []models.HttpPackageConfigurationsModel) []awsTypes.HttpPackageConfiguration {
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
