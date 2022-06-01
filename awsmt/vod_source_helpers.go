package awsmt

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func setHttpPackageConfigurations(values []*mediatailor.HttpPackageConfiguration, d *schema.ResourceData) error {
	var configurations []map[string]interface{}
	for _, c := range values {
		temp := map[string]interface{}{}
		temp["path"] = c.Path
		temp["source_group"] = c.SourceGroup
		temp["type"] = c.Type
		configurations = append(configurations, temp)
	}
	if err := d.Set("http_package_configurations", configurations); err != nil {
		return fmt.Errorf("error while setting the http package configurations: %w", err)
	}
	return nil
}

func setVodSource(values *mediatailor.DescribeVodSourceOutput, d *schema.ResourceData) error {
	var errs []error

	if values.Arn != nil {
		errs = append(errs, d.Set("arn", values.Arn))
	}
	if values.CreationTime != nil {
		errs = append(errs, d.Set("creation_time", values.CreationTime.String()))
	}
	errs = append(errs, setHttpPackageConfigurations(values.HttpPackageConfigurations, d))
	if values.LastModifiedTime != nil {
		errs = append(errs, d.Set("last_modified_time", values.LastModifiedTime.String()))
	}
	if values.SourceLocationName != nil {
		errs = append(errs, d.Set("source_location_name", values.SourceLocationName))
	}
	errs = append(errs, d.Set("tags", values.Tags))
	if values.VodSourceName != nil {
		errs = append(errs, d.Set("vod_source_name", values.VodSourceName))
	}
	for _, e := range errs {
		if e != nil {
			return fmt.Errorf("the following error occured while setting the values: %w", e)
		}
	}
	return nil
}

func getHttpPackageConfigurations(d *schema.ResourceData) []*mediatailor.HttpPackageConfiguration {
	if v, ok := d.GetOk("http_package_configurations"); ok && v.([]interface{})[0] != nil {
		configurations := v.([]interface{})

		var res []*mediatailor.HttpPackageConfiguration

		for _, c := range configurations {
			current := c.(map[string]interface{})
			temp := mediatailor.HttpPackageConfiguration{}

			if str, ok := current["path"]; ok {
				temp.Path = aws.String(str.(string))
			}
			if str, ok := current["source_group"]; ok {
				temp.SourceGroup = aws.String(str.(string))
			}
			if str, ok := current["type"]; ok {
				temp.Type = aws.String(str.(string))
			}

			res = append(res, &temp)
		}
		return res
	}
	return nil
}

func getCreateVodSourceInput(d *schema.ResourceData) mediatailor.CreateVodSourceInput {
	var vodSourceInputParams mediatailor.CreateVodSourceInput

	if c := getHttpPackageConfigurations(d); c != nil {
		vodSourceInputParams.HttpPackageConfigurations = c
	}

	if s, ok := d.GetOk("source_location_name"); ok {
		vodSourceInputParams.SourceLocationName = aws.String(s.(string))
	}

	outputMap := make(map[string]*string)
	if v, ok := d.GetOk("tags"); ok {
		val := v.(map[string]interface{})
		for k, value := range val {
			temp := value.(string)
			outputMap[k] = &temp
		}
	}
	vodSourceInputParams.Tags = outputMap

	if s, ok := d.GetOk("vod_source_name"); ok {
		vodSourceInputParams.VodSourceName = aws.String(s.(string))
	}

	return vodSourceInputParams
}

func getUpdateVodSourceInput(d *schema.ResourceData) mediatailor.UpdateVodSourceInput {
	var updatedVodSourceParams mediatailor.UpdateVodSourceInput

	if c := getHttpPackageConfigurations(d); c != nil {
		updatedVodSourceParams.HttpPackageConfigurations = c
	}
	if s, ok := d.GetOk("source_location_name"); ok {
		updatedVodSourceParams.SourceLocationName = aws.String(s.(string))
	}
	if s, ok := d.GetOk("vod_source_name"); ok {
		updatedVodSourceParams.VodSourceName = aws.String(s.(string))
	}
	return updatedVodSourceParams
}
