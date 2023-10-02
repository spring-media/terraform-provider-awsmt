package awsmt

import (
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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

func buildDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Attributes: map[string]datasourceSchema.Attribute{
			"id":            computedString,
			"arn":           computedString,
			"creation_time": computedString,
			"http_package_configurations": datasourceSchema.ListNestedAttribute{
				Computed: true,
				NestedObject: datasourceSchema.NestedAttributeObject{
					Attributes: map[string]datasourceSchema.Attribute{
						"path":         computedString,
						"source_group": computedString,
						"type": datasourceSchema.StringAttribute{
							Computed: true,
							Validators: []validator.String{
								stringvalidator.OneOf("HLS", "DASH"),
							},
						},
					},
				},
			},
			"last_modified_time":   computedString,
			"name":                 requiredString,
			"source_location_name": requiredString,
			"tags":                 computedMap,
		},
	}
}

func buildResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Attributes: map[string]resourceSchema.Attribute{
			"id":            computedString,
			"arn":           computedString,
			"creation_time": computedString,
			"http_package_configurations": resourceSchema.ListNestedAttribute{
				Required: true,
				NestedObject: resourceSchema.NestedAttributeObject{
					Attributes: map[string]resourceSchema.Attribute{
						"path":         requiredString,
						"source_group": requiredString,
						"type": resourceSchema.StringAttribute{
							Required: true,
							Validators: []validator.String{
								stringvalidator.OneOf("HLS", "DASH"),
							},
						},
					},
				},
			},
			"last_modified_time":   computedString,
			"source_location_name": requiredString,
			"tags":                 optionalMap,
			"name":                 requiredString,
		},
	}
}
