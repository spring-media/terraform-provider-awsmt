package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type SourceLocationModel struct {
	ID                                  types.String                              `tfsdk:"id"`
	AccessConfiguration                 *AccessConfigurationModel                 `tfsdk:"access_configuration"`
	Arn                                 types.String                              `tfsdk:"arn"`
	CreationTime                        types.String                              `tfsdk:"creation_time"`
	DefaultSegmentDeliveryConfiguration *DefaultSegmentDeliveryConfigurationModel `tfsdk:"default_segment_delivery_configuration"`
	HttpConfiguration                   *HttpConfigurationModel                   `tfsdk:"http_configuration"`
	LastModifiedTime                    types.String                              `tfsdk:"last_modified_time"`
	SegmentDeliveryConfigurations       []SegmentDeliveryConfigurationsModel      `tfsdk:"segment_delivery_configurations"`
	Name                                *string                                   `tfsdk:"name"`
	Tags                                map[string]string                         `tfsdk:"tags"`
}

type AccessConfigurationModel struct {
	AccessType                             *string                                      `tfsdk:"access_type"`
	SecretsManagerAccessTokenConfiguration *SecretsManagerAccessTokenConfigurationModel `tfsdk:"smatc"`
}

func (a *AccessConfigurationModel) Equal(other *AccessConfigurationModel) bool {
	if a == nil && other == nil {
		return true
	}
	if a == nil || other == nil {
		return false
	}

	if !equalStringPointers(a.AccessType, other.AccessType) {
		return false
	}

	if a.SecretsManagerAccessTokenConfiguration == nil && other.SecretsManagerAccessTokenConfiguration == nil {
		return true
	}
	if a.SecretsManagerAccessTokenConfiguration == nil || other.SecretsManagerAccessTokenConfiguration == nil {
		return false
	}
	return a.SecretsManagerAccessTokenConfiguration.Equal(other.SecretsManagerAccessTokenConfiguration)
}

type SecretsManagerAccessTokenConfigurationModel struct {
	HeaderName      *string `tfsdk:"header_name"`
	SecretArn       *string `tfsdk:"secret_arn"`
	SecretStringKey *string `tfsdk:"secret_string_key"`
}

func (s *SecretsManagerAccessTokenConfigurationModel) Equal(other *SecretsManagerAccessTokenConfigurationModel) bool {
	if s == nil && other == nil {
		return true
	}
	if s == nil || other == nil {
		return false
	}

	return equalStringPointers(s.HeaderName, other.HeaderName) &&
		equalStringPointers(s.SecretArn, other.SecretArn) &&
		equalStringPointers(s.SecretStringKey, other.SecretStringKey)
}

type DefaultSegmentDeliveryConfigurationModel struct {
	BaseUrl *string `tfsdk:"base_url"`
}

type HttpConfigurationModel struct {
	BaseUrl *string `tfsdk:"base_url"`
}

type SegmentDeliveryConfigurationsModel struct {
	BaseUrl *string `tfsdk:"base_url"`
	SDCName *string `tfsdk:"name"`
}

func equalStringPointers(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
