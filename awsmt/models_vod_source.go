package awsmt

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type vodSourceModel struct {
	ID                               types.String                     `tfsdk:"id"`
	Arn                              types.String                     `tfsdk:"arn"`
	CreationTime                     types.String                     `tfsdk:"creation_time"`
	HttpPackageConfigurations        []httpPackageConfigurationsModel `tfsdk:"http_package_configurations"`
	LastModifiedTime                 types.String                     `tfsdk:"last_modified_time"`
	SourceLocationName               *string                          `tfsdk:"source_location_name"`
	Tags                             map[string]string                `tfsdk:"tags"`
	Name                             *string                          `tfsdk:"name"`
	AdBreakOpportunitiesOffsetMillis []*int64                         `tfsdk:"ad_break_opportunities_offset_millis"`
}
