package awsmt

import (
	"context"
	"fmt"
	mediatailorV2 "github.com/aws/aws-sdk-go-v2/service/mediatailor"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"reflect"
	"strings"
)

func untagResource(client *mediatailor.MediaTailor, oldTags map[string]*string, resourceArn string) error {
	var removeTags []*string
	for k := range oldTags {
		removeTags = append(removeTags, aws.String(k))
	}
	_, err := client.UntagResource(&mediatailor.UntagResourceInput{ResourceArn: &resourceArn, TagKeys: removeTags})
	if err != nil {
		return err
	}
	return nil
}

func tagResource(client *mediatailor.MediaTailor, newTags map[string]*string, resourceArn string) error {
	_, err := client.TagResource(&mediatailor.TagResourceInput{ResourceArn: &resourceArn, Tags: newTags})
	if err != nil {
		return err
	}
	return nil
}

func updatesTags(client *mediatailor.MediaTailor, oldTags map[string]*string, newTags map[string]*string, resourceArn string) error {
	if !reflect.DeepEqual(oldTags, newTags) {
		if err := untagResource(client, oldTags, resourceArn); err != nil {
			return err
		}
		if err := tagResource(client, newTags, resourceArn); err != nil {
			return err
		}
	}
	return nil
}

func v2Untag(client *mediatailorV2.Client, oldTags map[string]string, resourceArn string) error {
	var removeTags []string
	for k := range oldTags {
		removeTags = append(removeTags, k)
	}
	if len(removeTags) == 0 {
		return nil
	}
	if _, err := client.UntagResource(context.TODO(), &mediatailorV2.UntagResourceInput{ResourceArn: &resourceArn, TagKeys: removeTags}); err != nil {
		return err
	}
	return nil
}

func v2Tag(client *mediatailorV2.Client, newTags map[string]string, resourceArn string) error {
	if len(newTags) == 0 {
		return nil
	}
	if _, err := client.TagResource(context.TODO(), &mediatailorV2.TagResourceInput{ResourceArn: &resourceArn, Tags: newTags}); err != nil {
		return err
	}
	return nil
}

func V2UpdatesTags(client *mediatailorV2.Client, oldTags map[string]string, newTags map[string]string, resourceArn string) error {
	if !tagsEqual(oldTags, newTags) {
		if err := v2Untag(client, oldTags, resourceArn); err != nil {
			return err
		}
		if err := v2Tag(client, newTags, resourceArn); err != nil {
			return err
		}
	}
	return nil
}

func tagsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			return false
		}
	}
	return true
}

func importStateForContentSources(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: source_location_name,name. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("source_location_name"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
}
