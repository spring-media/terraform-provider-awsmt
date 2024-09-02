package awsmt

import (
	"context"
	"fmt"
	mediatailorV2 "github.com/aws/aws-sdk-go-v2/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"strings"
)

func untag(client *mediatailorV2.Client, oldTags map[string]string, resourceArn string) error {
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

func tag(client *mediatailorV2.Client, newTags map[string]string, resourceArn string) error {
	if len(newTags) == 0 {
		return nil
	}
	if _, err := client.TagResource(context.TODO(), &mediatailorV2.TagResourceInput{ResourceArn: &resourceArn, Tags: newTags}); err != nil {
		return err
	}
	return nil
}

func UpdatesTags(client *mediatailorV2.Client, oldTags map[string]string, newTags map[string]string, resourceArn string) error {
	if !tagsEqual(oldTags, newTags) {
		if err := untag(client, oldTags, resourceArn); err != nil {
			return err
		}
		if err := tag(client, newTags, resourceArn); err != nil {
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
