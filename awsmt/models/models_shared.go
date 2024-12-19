package models

type HttpPackageConfigurationsModel struct {
	Path        *string `tfsdk:"path"`
	SourceGroup *string `tfsdk:"source_group"`
	Type        *string `tfsdk:"type"`
}
