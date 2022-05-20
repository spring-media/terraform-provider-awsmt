package awsmt

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func setAccessConfiguration(values *mediatailor.DescribeSourceLocationOutput, d *schema.ResourceData) error {
	if values.AccessConfiguration != nil && values.AccessConfiguration != &(mediatailor.AccessConfiguration{}) {
		temp := map[string]interface{}{}
		if values.AccessConfiguration.AccessType != nil {
			temp["access_type"] = values.AccessConfiguration
		}
		if values.AccessConfiguration.SecretsManagerAccessTokenConfiguration != nil && values.AccessConfiguration.SecretsManagerAccessTokenConfiguration != &(mediatailor.SecretsManagerAccessTokenConfiguration{}) {
			if values.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName != nil {
				temp["smatc_header_name"] = values.AccessConfiguration.SecretsManagerAccessTokenConfiguration.HeaderName
			}
			if values.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn != nil {
				temp["smatc_secret_arn"] = values.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretArn
			}
			if values.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey != nil {
				temp["smatc_secret_string_key"] = values.AccessConfiguration.SecretsManagerAccessTokenConfiguration.SecretStringKey
			}
		}
		if err := d.Set("access_configuration", []interface{}{temp}); err != nil {
			return fmt.Errorf("error while setting the access configuration: %w", err)
		}
	}
	return nil
}

func setSegmentDeliveryConfigurations(values *mediatailor.DescribeSourceLocationOutput, d *schema.ResourceData) error {
	var configurations []map[string]interface{}
	for _, c := range values.SegmentDeliveryConfigurations {
		temp := map[string]interface{}{}
		temp["base_url"] = c.BaseUrl
		temp["name"] = c.Name
		configurations = append(configurations, temp)
	}
	if err := d.Set("segment_delivery_configurations", configurations); err != nil {
		return fmt.Errorf("error while setting the segment delivery configurations: %w", err)
	}
	return nil
}

func setResourceLocation(values *mediatailor.DescribeSourceLocationOutput, d *schema.ResourceData) error {
	var errors []error

	errors = append(errors, setAccessConfiguration(values, d))
	errors = append(errors, d.Set("arn", values.Arn))
	errors = append(errors, d.Set("creation_time", values.CreationTime.String()))
	errors = append(errors, d.Set("default_segment_delivery_configuration_url", values.DefaultSegmentDeliveryConfiguration))
	if values.HttpConfiguration != nil && values.HttpConfiguration != &(mediatailor.HttpConfiguration{}) {
		if values.HttpConfiguration.BaseUrl != nil {
			errors = append(errors, d.Set("http_configuration_url", values.HttpConfiguration.BaseUrl))
		}
	}
	errors = append(errors, d.Set("last_modified_time", values.LastModifiedTime.String()))
	errors = append(errors, setSegmentDeliveryConfigurations(values, d))

	for _, e := range errors {
		if e != nil {
			return fmt.Errorf("the following error occured while setting the values: %w", e)
		}
	}

	return nil
}
