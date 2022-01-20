package mediatailor

import (
	"context"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceConfigurationRead,
		// schema based on https://docs.aws.amazon.com/sdk-for-go/api/service/mediatailor/#GetPlaybackConfigurationOutput
		// with types found on https://sourcegraph.com/github.com/aws/aws-sdk-go/-/docs/service/mediatailor
		Schema: map[string]*schema.Schema{
			"Configuration": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"AdDecisionServerUrl": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"AvailSuppression": &schema.Schema{
							Type:     schema.TypeList,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"Mode": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"Value": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"Bumper": &schema.Schema{
							Type:     schema.TypeList,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"EndUrl": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"StartUrl": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"CdnConfiguration": &schema.Schema{
							Type:     schema.TypeList,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"AdSegmentUrlPrefix": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"ContentSegmentUrlPrefix": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"ConfigurationAliases": &schema.Schema{
							Type: schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"DashConfiguration": &schema.Schema{
							Type:     schema.TypeList,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ManifestEndpointPrefix": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"MpdLocation": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"OriginManifestType": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"HlsConfiguration": &schema.Schema{
							Type:     schema.TypeList,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ManifestEndpointPrefix": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"LivePreRollConfiguration": &schema.Schema{
							Type:     schema.TypeList,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"AdDecisionServerUrl": &schema.Schema{
										Type:     schema.TypeString,
										Computed: true,
									},
									"MaxDurationSeconds": &schema.Schema{
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"LogConfiguration": &schema.Schema{
							Type:     schema.TypeList,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"PercentEnabled": &schema.Schema{
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"ManifestProcessingRules": &schema.Schema{
							Type:     schema.TypeList,
							MaxItems: 1,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"AdMarkerPassthrough": &schema.Schema{
										Type:     schema.TypeList,
										MaxItems: 1,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"Enabled": &schema.Schema{
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"Name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"PersonalizationThresholdSeconds": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
							Default:  1,
						},
						"PlaybackConfigurationString": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"PlaybackEndpointPrefix": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"SessionInitializationEndpointPrefix": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"SlateAdUrl": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"Tags": &schema.Schema{
							Type: schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"TranscodeProfileName": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"VideoContentSourceUrl": &schema.Schema{
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
	output, err := client.GetPlaybackConfiguration(&mediatailor.GetPlaybackConfigurationInput{Name: &name})
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("configuration", output); err != nil {
		return diag.FromErr(err)
	}
	return diags
}
