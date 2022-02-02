package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePlaybackConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePlaybackConfigurationRead,
		// schema based on https://docs.aws.amazon.com/sdk-for-go/api/service/mediatailor/#GetPlaybackConfigurationOutput
		// with types found on https://sourcegraph.com/github.com/aws/aws-sdk-go/-/docs/service/mediatailor
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ad_decision_server_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cdn_configuration": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ad_segment_url_prefix": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"content_segment_url_prefix": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"dash_configuration": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"manifest_endpoint_prefix": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"mpd_location": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"origin_manifest_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"hls_configuration": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"manifest_endpoint_prefix": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"playback_configuration_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"playback_endpoint_prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"session_initialization_endpoint_prefix": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slate_ad_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"transcode_profile_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"video_content_source_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePlaybackConfigurationRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mediatailor.MediaTailor)
	var diags diag.Diagnostics
	var output []interface{}

	name := d.Get("name").(string)

	if name != "" {
		res, err := getSinglePlaybackConfiguration(client, name)
		if err != nil {
			return diag.FromErr(err)
		}
		output = []interface{}{flatten(res)}
		returnValues(d, output, diags)
	} else {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "`name` parameter required",
			Detail:   "You need to specify a `name` parameter in your configuration",
		})
	}
	return diags
}

func getSinglePlaybackConfiguration(c *mediatailor.MediaTailor, name string) (*mediatailor.PlaybackConfiguration, error) {
	output, err := c.GetPlaybackConfiguration(&mediatailor.GetPlaybackConfigurationInput{Name: &name})
	if err != nil {
		return nil, err
	}
	return (*mediatailor.PlaybackConfiguration)(output), nil
}

func returnValues(d *schema.ResourceData, values []interface{}, diags diag.Diagnostics) diag.Diagnostics {
	if err := d.Set("configuration", values); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(uuid.New().String())
	return diags
}

func flatten(configuration *mediatailor.PlaybackConfiguration) map[string]interface{} {
	if configuration != nil {
		output := make(map[string]interface{})

		output["ad_decision_server_url"] = configuration.AdDecisionServerUrl
		output["cdn_configuration"] = []interface{}{map[string]interface{}{
			"ad_segment_url_prefix":      configuration.CdnConfiguration.AdSegmentUrlPrefix,
			"content_segment_url_prefix": configuration.CdnConfiguration.ContentSegmentUrlPrefix,
		}}
		output["dash_configuration"] = []interface{}{map[string]interface{}{
			"manifest_endpoint_prefix": configuration.DashConfiguration.ManifestEndpointPrefix,
			"mpd_location":             configuration.DashConfiguration.MpdLocation,
			"origin_manifest_type":     configuration.DashConfiguration.OriginManifestType,
		}}
		output["hls_configuration"] = []interface{}{map[string]interface{}{
			"manifest_endpoint_prefix": configuration.HlsConfiguration.ManifestEndpointPrefix,
		}}
		output["name"] = configuration.Name
		output["playback_configuration_arn"] = configuration.PlaybackConfigurationArn
		output["playback_endpoint_prefix"] = configuration.PlaybackEndpointPrefix
		output["session_initialization_endpoint_prefix"] = configuration.SessionInitializationEndpointPrefix
		output["slate_ad_url"] = configuration.SlateAdUrl
		output["tags"] = configuration.Tags
		output["transcode_profile_name"] = configuration.TranscodeProfileName
		output["video_content_source_url"] = configuration.VideoContentSourceUrl
		return output
	}
	return map[string]interface{}{}
}