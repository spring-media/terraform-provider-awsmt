package awsmt

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func setLiveSource(values *mediatailor.DescribeLiveSourceOutput, d *schema.ResourceData) error {
	var errors []error

	if values.Arn != nil {
		errors = append(errors, d.Set("arn", values.Arn))
	}
	if values.CreationTime != nil {
		errors = append(errors, d.Set("creation_time", values.CreationTime.String()))
	}
	errors = append(errors, setHttpPackageConfigurations(values.HttpPackageConfigurations, d))
	if values.LastModifiedTime != nil {
		errors = append(errors, d.Set("last_modified_time", values.LastModifiedTime.String()))
	}
	if values.LiveSourceName != nil {
		errors = append(errors, d.Set("live_source_name", values.LiveSourceName))
	}
	if values.SourceLocationName != nil {
		errors = append(errors, d.Set("source_location_name", values.SourceLocationName))
	}
	errors = append(errors, d.Set("tags", values.Tags))
	for _, e := range errors {
		if e != nil {
			return fmt.Errorf("the following error occured while setting the values: %w", e)
		}
	}
	return nil
}

func getCreateLiveSourceInput(d *schema.ResourceData) mediatailor.CreateLiveSourceInput {
	var inputParams mediatailor.CreateLiveSourceInput

	if c := getHttpPackageConfigurations(d); c != nil {
		inputParams.HttpPackageConfigurations = c
	}

	if s, ok := d.GetOk("source_location_name"); ok {
		inputParams.SourceLocationName = aws.String(s.(string))
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

	if s, ok := d.GetOk("live_source_name"); ok {
		inputParams.LiveSourceName = aws.String(s.(string))
	}

	return inputParams
}

func getUpdateLiveSourceInput(d *schema.ResourceData) mediatailor.UpdateLiveSourceInput {
	var updateParams mediatailor.UpdateLiveSourceInput

	if c := getHttpPackageConfigurations(d); c != nil {
		updateParams.HttpPackageConfigurations = c
	}
	if s, ok := d.GetOk("source_location_name"); ok {
		updateParams.SourceLocationName = aws.String(s.(string))
	}
	if s, ok := d.GetOk("vod_source_name"); ok {
		updateParams.LiveSourceName = aws.String(s.(string))
	}
	return updateParams
}
