package awsmt

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func createBaseList(fields map[string]*schema.Schema) *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeList,
		Elem: &schema.Resource{
			Schema: fields,
		},
	}
}

func createOptionalList(fields map[string]*schema.Schema) *schema.Schema {
	s := createBaseList(fields)
	s.Optional = true
	s.MaxItems = 1
	return s
}

func createRequiredList(fields map[string]*schema.Schema) *schema.Schema {
	s := createBaseList(fields)
	s.Required = true
	s.MaxItems = 1
	return s
}

func createComputedList(fields map[string]*schema.Schema) *schema.Schema {
	s := createBaseList(fields)
	s.Computed = true
	return s
}

func updateTags(client *mediatailor.MediaTailor, arn *string, oldTagValue, newTagValue interface{}) error {

	var removedTags []string
	for k := range oldTagValue.(map[string]interface{}) {
		if _, ok := (newTagValue.(map[string]interface{}))[k]; !ok {
			removedTags = append(removedTags, k)
		}
	}

	err := deleteTags(client, aws.StringValue(arn), removedTags)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if newTagValue != nil {
		var newTags = make(map[string]*string)
		for k, v := range newTagValue.(map[string]interface{}) {
			val := v.(string)
			newTags[k] = &val
		}
		tagInput := mediatailor.TagResourceInput{ResourceArn: arn, Tags: newTags}
		_, err := client.TagResource(&tagInput)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}
	return nil
}

func deleteTags(client *mediatailor.MediaTailor, resourceArn string, removedTags []string) error {
	if len(removedTags) != 0 {

		var removedValuesPointer []*string
		for i := range removedTags {
			removedValuesPointer = append(removedValuesPointer, &removedTags[i])
		}

		untagInput := mediatailor.UntagResourceInput{ResourceArn: aws.String(resourceArn), TagKeys: removedValuesPointer}
		_, err := client.UntagResource(&untagInput)
		if err != nil {
			return err
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
