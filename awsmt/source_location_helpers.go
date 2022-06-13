package awsmt

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

func setSourceLocation(values *mediatailor.DescribeSourceLocationOutput, d *schema.ResourceData) error {
	var errors []error

	errors = append(errors, d.Set("arn", values.Arn))
	errors = append(errors, d.Set("creation_time", values.CreationTime.String()))
	if values.DefaultSegmentDeliveryConfiguration != nil && values.DefaultSegmentDeliveryConfiguration != &(mediatailor.DefaultSegmentDeliveryConfiguration{}) {
		errors = append(errors, d.Set("default_segment_delivery_configuration_url", values.DefaultSegmentDeliveryConfiguration.BaseUrl))
	}
	if values.HttpConfiguration != nil && values.HttpConfiguration != &(mediatailor.HttpConfiguration{}) {
		if values.HttpConfiguration.BaseUrl != nil {
			errors = append(errors, d.Set("http_configuration_url", values.HttpConfiguration.BaseUrl))
		}
	}
	errors = append(errors, d.Set("last_modified_time", values.LastModifiedTime.String()))
	errors = append(errors, setSegmentDeliveryConfigurations(values, d))
	errors = append(errors, d.Set("name", values.SourceLocationName))
	errors = append(errors, d.Set("tags", values.Tags))

	for _, e := range errors {
		if e != nil {
			return fmt.Errorf("the following error occured while setting the values: %w", e)
		}
	}

	return nil
}

func getSegmentDeliveryConfigurations(d *schema.ResourceData) []*mediatailor.SegmentDeliveryConfiguration {
	if v, ok := d.GetOk("segment_delivery_configurations"); ok && v.([]interface{})[0] != nil {
		configurations := v.([]interface{})

		var res []*mediatailor.SegmentDeliveryConfiguration

		for _, c := range configurations {
			current := c.(map[string]interface{})
			temp := mediatailor.SegmentDeliveryConfiguration{}

			if str, ok := current["base_url"]; ok {
				temp.BaseUrl = aws.String(str.(string))
			}
			if str, ok := current["name"]; ok {
				temp.Name = aws.String(str.(string))
			}
			res = append(res, &temp)
		}
		return res
	}
	return nil
}

func getCreateSourceLocationInput(d *schema.ResourceData) mediatailor.CreateSourceLocationInput {
	var inputParams mediatailor.CreateSourceLocationInput

	if v, ok := d.GetOk("default_segment_delivery_configuration_url"); ok {
		inputParams.DefaultSegmentDeliveryConfiguration = &mediatailor.DefaultSegmentDeliveryConfiguration{BaseUrl: aws.String(v.(string))}
	}

	if v, ok := d.GetOk("http_configuration_url"); ok {
		inputParams.HttpConfiguration = &mediatailor.HttpConfiguration{BaseUrl: aws.String(v.(string))}
	}

	if s := getSegmentDeliveryConfigurations(d); s != nil {
		inputParams.SegmentDeliveryConfigurations = s
	}

	if v, ok := d.GetOk("name"); ok {
		inputParams.SourceLocationName = aws.String(v.(string))
	}

	outputMap := make(map[string]*string)
	if v, ok := d.GetOk("tags"); ok {
		val := v.(map[string]interface{})
		for k, value := range val {
			temp := value.(string)
			outputMap[k] = &temp
		}
	}
	inputParams.Tags = outputMap

	return inputParams

}

func getUpdateSourceLocationInput(d *schema.ResourceData) mediatailor.UpdateSourceLocationInput {
	var updateParams mediatailor.UpdateSourceLocationInput

	if v, ok := d.GetOk("default_segment_delivery_configuration_url"); ok {
		updateParams.DefaultSegmentDeliveryConfiguration = &mediatailor.DefaultSegmentDeliveryConfiguration{BaseUrl: aws.String(v.(string))}
	}

	if v, ok := d.GetOk("http_configuration_url"); ok {
		updateParams.HttpConfiguration = &mediatailor.HttpConfiguration{BaseUrl: aws.String(v.(string))}
	}

	if s := getSegmentDeliveryConfigurations(d); s != nil {
		updateParams.SegmentDeliveryConfigurations = s
	}

	if v, ok := d.GetOk("name"); ok {
		updateParams.SourceLocationName = aws.String(v.(string))
	}

	return updateParams
}
