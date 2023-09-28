package awsmt

import "github.com/hashicorp/terraform-plugin-framework/types"

type sourceLocationModel struct {
	ID                                  types.String                              `tfsdk:"id"`
	AccessConfiguration                 *accessConfigurationModel                 `tfsdk:"access_configuration"`
	Arn                                 types.String                              `tfsdk:"arn"`
	CreationTime                        types.String                              `tfsdk:"creation_time"`
	DefaultSegmentDeliveryConfiguration *defaultSegmentDeliveryConfigurationModel `tfsdk:"default_segment_delivery_configuration"`
	HttpConfiguration                   *httpConfigurationModel                   `tfsdk:"http_configuration"`
	LastModifiedTime                    types.String                              `tfsdk:"last_modified_time"`
	SegmentDeliveryConfigurations       []segmentDeliveryConfigurationsModel      `tfsdk:"segment_delivery_configurations"`
	Name                                *string                                   `tfsdk:"name"`
	Tags                                map[string]*string                        `tfsdk:"tags"`
}

type accessConfigurationModel struct {
	AccessType                             *string                                      `tfsdk:"access_type"`
	SecretsManagerAccessTokenConfiguration *secretsManagerAccessTokenConfigurationModel `tfsdk:"smatc"`
}

type secretsManagerAccessTokenConfigurationModel struct {
	HeaderName      *string `tfsdk:"header_name"`
	SecretArn       *string `tfsdk:"secret_arn"`
	SecretStringKey *string `tfsdk:"secret_string_key"`
}

type defaultSegmentDeliveryConfigurationModel struct {
	BaseUrl *string `tfsdk:"base_url"`
}

type httpConfigurationModel struct {
	BaseUrl *string `tfsdk:"base_url"`
}

type segmentDeliveryConfigurationsModel struct {
	BaseUrl *string `tfsdk:"base_url"`
	SDCName *string `tfsdk:"name"`
}
