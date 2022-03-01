package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"
)

func resourcePlaybackConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePlaybackConfigurationCreate,
		UpdateContext: resourcePlaybackConfigurationUpdate,
		ReadContext:   resourcePlaybackConfigurationRead,
		DeleteContext: resourcePlaybackConfigurationDelete,
		// schema based on: https://docs.aws.amazon.com/mediatailor/latest/apireference/playbackconfiguration.html#playbackconfiguration-prop-putplaybackconfigurationrequest-personalizationthresholdseconds
		// and https://sourcegraph.com/github.com/aws/aws-sdk-go/-/docs/service/mediatailor#PutPlaybackConfigurationInput
		Schema: map[string]*schema.Schema{
			"ad_decision_server_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cdn_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ad_segment_url_prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"content_segment_url_prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"dash_configuration": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mpd_location": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"origin_manifest_type": {
							//SINGLE_PERIOD | MULTI_PERIOD
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"slate_ad_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"transcode_profile_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"video_content_source_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_updated": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourcePlaybackConfigurationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mediatailor.MediaTailor)
	var diags diag.Diagnostics

	input := getPlaybackConfigurationInput(d)

	_, err := client.PutPlaybackConfiguration(&input)
	if err != nil {
		return diag.FromErr(err)
	}
	resourcePlaybackConfigurationRead(ctx, d, m)
	d.SetId(*input.Name)
	return diags
}

func resourcePlaybackConfigurationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	if d.HasChanges("ad_decision_server_url", "cdn_configuration", "dash_configuration", "slate_ad_url", "tags", "transcode_profile_name", "video_content_source_url") {
		client := m.(*mediatailor.MediaTailor)
		input := getPlaybackConfigurationInput(d)
		_, err := client.PutPlaybackConfiguration(&input)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("last_updated", time.Now().Format(time.RFC850)); err != nil {
			return diag.FromErr(err)
		}
	}
	resourcePlaybackConfigurationRead(ctx, d, m)

	return diags
}

func resourcePlaybackConfigurationRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mediatailor.MediaTailor)
	var diags diag.Diagnostics

	name := d.Get("name").(string)
	if len(name) == 0 && len(d.Id()) > 0 {
		name = d.Id()
	}
	res, err := client.GetPlaybackConfiguration(&mediatailor.GetPlaybackConfigurationInput{Name: &name})
	if err != nil {
		return diag.FromErr(err)
	}

	output := flatterResourcePlaybackConfiguration(res)
	returnPlaybackConfigurationResource(d, output, diags)
	return diags
}

func resourcePlaybackConfigurationDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mediatailor.MediaTailor)
	var diags diag.Diagnostics
	name := d.Get("name").(string)
	_, err := client.DeletePlaybackConfiguration(&mediatailor.DeletePlaybackConfigurationInput{Name: &name})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}

func getPlaybackConfigurationInput(d *schema.ResourceData) mediatailor.PutPlaybackConfigurationInput {
	input := mediatailor.PutPlaybackConfigurationInput{}
	if v, ok := d.GetOk("ad_decision_server_url"); ok {
		val := v.(string)
		input.AdDecisionServerUrl = &val
	}
	if v, ok := d.GetOk("cdn_configuration"); ok {
		val := v.([]interface{})[0].(map[string]interface{})
		output := mediatailor.CdnConfiguration{}
		if str, ok := val["ad_segment_url_prefix"]; ok {
			converted := str.(string)
			output.AdSegmentUrlPrefix = &converted
		}
		if str, ok := val["content_segment_url_prefix"]; ok {
			converted := str.(string)
			output.ContentSegmentUrlPrefix = &converted
		}
		input.CdnConfiguration = &output
	}
	if v, ok := d.GetOk("dash_configuration"); ok {
		val := v.([]interface{})[0].(map[string]interface{})
		output := mediatailor.DashConfigurationForPut{}
		if str, ok := val["mpd_location"]; ok {
			converted := str.(string)
			output.MpdLocation = &converted
		}
		if str, ok := val["origin_manifest_type"]; ok {
			converted := str.(string)
			output.OriginManifestType = &converted
		}
		input.DashConfiguration = &output
	}
	if v, ok := d.GetOk("name"); ok {
		val := v.(string)
		input.Name = &val
	}
	if v, ok := d.GetOk("slate_ad_url"); ok {
		val := v.(string)
		input.SlateAdUrl = &val
	}
	if v, ok := d.GetOk("tags"); ok {
		outputMap := make(map[string]*string)
		val := v.(map[string]interface{})
		for k, value := range val {
			temp := value.(string)
			outputMap[k] = &temp
		}
		input.Tags = outputMap
	}
	if v, ok := d.GetOk("transcode_profile_name"); ok {
		val := v.(string)
		input.TranscodeProfileName = &val
	}
	if v, ok := d.GetOk("video_content_source_url"); ok {
		val := v.(string)
		input.VideoContentSourceUrl = &val
	}
	return input
}

func flatterResourcePlaybackConfiguration(configuration *mediatailor.GetPlaybackConfigurationOutput) map[string]interface{} {
	if configuration != nil {
		output := make(map[string]interface{})

		output["ad_decision_server_url"] = configuration.AdDecisionServerUrl
		output["cdn_configuration"] = []interface{}{map[string]interface{}{
			"ad_segment_url_prefix":      configuration.CdnConfiguration.AdSegmentUrlPrefix,
			"content_segment_url_prefix": configuration.CdnConfiguration.ContentSegmentUrlPrefix,
		}}
		output["dash_configuration"] = []interface{}{map[string]interface{}{
			"mpd_location":         configuration.DashConfiguration.MpdLocation,
			"origin_manifest_type": configuration.DashConfiguration.OriginManifestType,
		}}
		output["name"] = configuration.Name
		output["slate_ad_url"] = configuration.SlateAdUrl
		output["tags"] = configuration.Tags
		output["transcode_profile_name"] = configuration.TranscodeProfileName
		output["video_content_source_url"] = configuration.VideoContentSourceUrl
		return output
	}
	return map[string]interface{}{}
}

func returnPlaybackConfigurationResource(d *schema.ResourceData, values map[string]interface{}, diags diag.Diagnostics) diag.Diagnostics {
	for k := range values {
		setSingleValue(d, values, diags, k)
	}
	return diags
}

func setSingleValue(d *schema.ResourceData, values map[string]interface{}, diags diag.Diagnostics, name string) diag.Diagnostics {
	err := d.Set(name, values[name])
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
