package awsmt

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"reflect"
)

func untagResource(client *mediatailor.MediaTailor, oldTags map[string]*string, newTags map[string]*string, resourceArn string) error {
	if !reflect.DeepEqual(oldTags, newTags) {
		if oldTags != nil && len(oldTags) > 0 {
			var removeTags []*string
			for k := range oldTags {
				removeTags = append(removeTags, aws.String(k))
			}
			_, err := client.UntagResource(&mediatailor.UntagResourceInput{ResourceArn: &resourceArn, TagKeys: removeTags})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func tagResource(client *mediatailor.MediaTailor, oldTags map[string]*string, newTags map[string]*string, resourceArn string) error {
	if !reflect.DeepEqual(oldTags, newTags) {
		if newTags != nil && len(newTags) > 0 {
			_, err := client.TagResource(&mediatailor.TagResourceInput{ResourceArn: &resourceArn, Tags: newTags})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func updatesTags(client *mediatailor.MediaTailor, oldTags map[string]*string, newTags map[string]*string, resourceArn string) error {
	if err := untagResource(client, oldTags, newTags, resourceArn); err != nil {
		return err
	}
	if err := tagResource(client, oldTags, newTags, resourceArn); err != nil {
		return err
	}
	return nil
}
