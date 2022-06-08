package awsmt

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceChannelPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceChannelPolicyRead,
		Schema: map[string]*schema.Schema{
			"channel_name": &requiredString,
			"policy":       &computedString,
		},
	}
}
func dataSourceChannelPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*mediatailor.MediaTailor)

	name := d.Get("channel_name").(string)
	res, err := client.GetChannelPolicy(&mediatailor.GetChannelPolicyInput{ChannelName: aws.String(name)})

	if err != nil {
		return diag.FromErr(fmt.Errorf("error while reading the channel policy: %v", err))
	}

	if err := d.Set("policy", res.Policy); err != nil {
		return diag.FromErr(fmt.Errorf("error while setting the policy: %v", err))
	}
	d.SetId(name)

	return nil
}
