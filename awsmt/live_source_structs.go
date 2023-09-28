package awsmt

import "github.com/hashicorp/terraform-plugin-framework/types"

type liveSourceModel struct {
	ID                        types.String                     `tfsdk:"id"`
	Arn                       types.String                     `tfsdk:"arn"`
	CreationTime              types.String                     `tfsdk:"creation_time"`
	HttpPackageConfigurations []httpPackageConfigurationsModel `tfsdk:"http_package_configurations"`
	LastModifiedTime          types.String                     `tfsdk:"last_modified_time"`
	LiveSourceName            *string                          `tfsdk:"live_source_name"`
	SourceLocationName        *string                          `tfsdk:"source_location_name"`
	Tags                      map[string]*string               `tfsdk:"tags"`
}

type httpPackageConfigurationsModel struct {
	Path        *string `tfsdk:"path"`
	SourceGroup *string `tfsdk:"source_group"`
	Type        *string `tfsdk:"type"`
}
