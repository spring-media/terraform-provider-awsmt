package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type VodSourceModel struct {
	ID                               types.String                     `tfsdk:"id"`
	Arn                              types.String                     `tfsdk:"arn"`
	CreationTime                     types.String                     `tfsdk:"creation_time"`
	HttpPackageConfigurations        []HttpPackageConfigurationsModel `tfsdk:"http_package_configurations"`
	LastModifiedTime                 types.String                     `tfsdk:"last_modified_time"`
	SourceLocationName               *string                          `tfsdk:"source_location_name"`
	Tags                             map[string]string                `tfsdk:"tags"`
	Name                             *string                          `tfsdk:"name"`
	AdBreakOpportunitiesOffsetMillis types.List                       `tfsdk:"ad_break_opportunities_offset_millis"`
}
