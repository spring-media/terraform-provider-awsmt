package awsmt

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceChannelPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChannelPolicyPut,
		ReadContext:   resourceChannelPolicyRead,
		UpdateContext: resourceChannelPolicyPut,
		DeleteContext: resourceChannelPolicyDelete,
		Schema: map[string]*schema.Schema{
			"channel_name": &requiredString,
			"policy":       &requiredString,
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceChannelPolicyPut(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	var putChannelPolicyParams = mediatailor.PutChannelPolicyInput{
		ChannelName: aws.String((d.Get("channel_name")).(string)),
		Policy:      aws.String((d.Get("policy")).(string)),
	}

	_, err := client.PutChannelPolicy(&putChannelPolicyParams)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while creating the policy: %v", err))
	}

	d.SetId(aws.StringValue(putChannelPolicyParams.ChannelName))
	return nil
}

func resourceChannelPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	channelName := d.Get("channel_name").(string)
	if len(channelName) == 0 && len(d.Id()) > 0 {
		channelName = d.Id()
	}

	res, err := client.GetChannelPolicy(&mediatailor.GetChannelPolicyInput{ChannelName: aws.String(channelName)})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while retrieving the policy: %v", err))
	}

	if err := d.Set("policy", res.Policy); err != nil {
		return diag.FromErr(fmt.Errorf("error while setting the policy: %v", err))
	}

	return nil
}

func resourceChannelPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	_, err := client.DeleteChannelPolicy(&mediatailor.DeleteChannelPolicyInput{ChannelName: aws.String(d.Get("channel_name").(string))})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error while deleting the policy: %v", err))
	}
	return nil
}
