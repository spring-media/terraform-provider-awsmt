package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type LiveSourceModel struct {
	ID                        types.String                     `tfsdk:"id"`
	Arn                       types.String                     `tfsdk:"arn"`
	CreationTime              types.String                     `tfsdk:"creation_time"`
	HttpPackageConfigurations []HttpPackageConfigurationsModel `tfsdk:"http_package_configurations"`
	LastModifiedTime          types.String                     `tfsdk:"last_modified_time"`
	Name                      *string                          `tfsdk:"name"`
	SourceLocationName        *string                          `tfsdk:"source_location_name"`
	Tags                      map[string]string                `tfsdk:"tags"`
}
