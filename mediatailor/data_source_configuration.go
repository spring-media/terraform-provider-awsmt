package mediatailor

import (
	"context"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConfigurationRead,
		// schema based on https://docs.aws.amazon.com/sdk-for-go/api/service/mediatailor/#GetPlaybackConfigurationOutput
		// with types found on https://sourcegraph.com/github.com/aws/aws-sdk-go/-/docs/service/mediatailor
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"configuration": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ad_decision_server_url": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"cdn_configuration": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ad_segment_url_prefix": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"content_segment_url_prefix": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"dash_configuration": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"manifest_endpoint_prefix": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"mpd_location": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"origin_manifest_type": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"hls_configuration": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"manifest_endpoint_prefix": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"playback_configuration_arn": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"playback_endpoint_prefix": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"session_initialization_endpoint_prefix": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"slate_ad_url": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": &schema.Schema{
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"transcode_profile_name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"video_content_source_url": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceConfigurationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*mediatailor.MediaTailor)

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	if name != "" {
		output, err := getSinglePlaybackConfiguration(client, name)
		if err != nil {
			return diag.FromErr(err)
		}
		flatOutput := flatten(output)
		if err := d.Set("configuration", []interface{}{flatOutput}); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(uuid.New().String())
	} else {
		output, err := listPlaybackConfigurations(client)
		if err != nil {
			return diag.FromErr(err)
		}
		var flatOutputList []interface{}
		for _, c := range output {
			flatOutputList = append(flatOutputList, flatten(c))
		}
		if err := d.Set("configuration", flatOutputList); err != nil {
			return diag.FromErr(err)
		}
		d.SetId(uuid.New().String())
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

func listPlaybackConfigurations(c *mediatailor.MediaTailor) ([]*mediatailor.PlaybackConfiguration, error) {
	output, err := c.ListPlaybackConfigurations(&mediatailor.ListPlaybackConfigurationsInput{})
	if err != nil {
		return nil, err
	}
	return output.Items, nil
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
