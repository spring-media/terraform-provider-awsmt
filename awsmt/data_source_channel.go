package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

var (
	_ datasource.DataSource              = &dataSourceChannel{}
	_ datasource.DataSourceWithConfigure = &dataSourceChannel{}
)

func DataSourceChannel() datasource.DataSource {
	return &dataSourceChannel{}
}

type dataSourceChannel struct {
	client *mediatailor.MediaTailor
}
type dataSourceChannelModel struct {
	ID               types.String        `tfsdk:"id"`
	Arn              types.String        `tfsdk:"arn"`
	ChannelName      *string             `tfsdk:"channel_name"`
	ChannelState     types.String        `tfsdk:"channel_state"`
	CreationTime     types.String        `tfsdk:"creation_time"`
	FillerSlate      *fillerSlateDSModel `tfsdk:"filler_slate"`
	LastModifiedTime types.String        `tfsdk:"last_modified_time"`
	Outputs          []outputsDSModel    `tfsdk:"outputs"`
	PlaybackMode     types.String        `tfsdk:"playback_mode"`
	Policy           types.String        `tfsdk:"policy"`
	Tags             map[string]*string  `tfsdk:"tags"`
	Tier             types.String        `tfsdk:"tier"`
}

type fillerSlateDSModel struct {
	SourceLocationName types.String `tfsdk:"source_location_name"`
	VodSourceName      types.String `tfsdk:"vod_source_name"`
}

type outputsDSModel struct {
	DashPlaylistSettings *dashPlaylistSettingsDSModel `tfsdk:"dash_playlist_settings"`
	HlsPlaylistSettings  *hlsPlaylistSettingsDSModel  `tfsdk:"hls_playlist_settings"`
	ManifestName         types.String                 `tfsdk:"manifest_name"`
	PlaybackUrl          types.String                 `tfsdk:"playback_url"`
	SourceGroup          types.String                 `tfsdk:"source_group"`
}

type dashPlaylistSettingsDSModel struct {
	ManifestWindowSeconds             types.Int64 `tfsdk:"manifest_window_seconds"`
	MinBufferTimeSeconds              types.Int64 `tfsdk:"min_buffer_time_seconds"`
	MinUpdatePeriodSeconds            types.Int64 `tfsdk:"min_update_period_seconds"`
	SuggestedPresentationDelaySeconds types.Int64 `tfsdk:"suggested_presentation_delay_seconds"`
}
type hlsPlaylistSettingsDSModel struct {
	AdMarkupType          []types.String `tfsdk:"ad_markup_type"`
	ManifestWindowSeconds types.Int64    `tfsdk:"manifest_window_seconds"`
}

func (d *dataSourceChannel) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_channel"
}

func (d *dataSourceChannel) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"arn": schema.StringAttribute{
				Computed: true,
			},
			"channel_name": schema.StringAttribute{
				Required: true,
			},
			"channel_state": schema.StringAttribute{
				Computed: true,
			},
			"creation_time": schema.StringAttribute{
				Computed: true,
			},
			"filler_slate": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"source_location_name": schema.StringAttribute{
						Computed: true,
					},
					"vod_source_name": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"last_modified_time": schema.StringAttribute{
				Computed: true,
			},
			"outputs": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"dash_playlist_settings": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"manifest_window_seconds": schema.Int64Attribute{
									Computed: true,
								},
								"min_buffer_time_seconds": schema.Int64Attribute{
									Computed: true,
								},
								"min_update_period_seconds": schema.Int64Attribute{
									Computed: true,
								},
								"suggested_presentation_delay_seconds": schema.Int64Attribute{
									Computed: true,
								},
							},
						},
						"hls_playlist_settings": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"ad_markup_type": schema.ListAttribute{
									Computed:    true,
									ElementType: types.StringType,
								},
								"manifest_window_seconds": schema.Int64Attribute{
									Computed: true,
								},
							},
						},
						"manifest_name": schema.StringAttribute{
							Computed: true,
						},
						"playback_url": schema.StringAttribute{
							Computed: true,
						},
						"source_group": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"playback_mode": schema.StringAttribute{
				Computed: true,
			},
			"policy": schema.StringAttribute{
				Computed: true,
			},
			"tags": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"tier": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *dataSourceChannel) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*mediatailor.MediaTailor)
}

func (d *dataSourceChannel) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dataSourceChannelModel
	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	channelName := data.ChannelName

	channel, err := d.client.DescribeChannel(&mediatailor.DescribeChannelInput{ChannelName: channelName})
	if err != nil {
		resp.Diagnostics.AddError("Error while describing channel "+*channelName, err.Error())
		return
	}

	policy, err := d.client.GetChannelPolicy(&mediatailor.GetChannelPolicyInput{ChannelName: channelName})
	if err != nil && !strings.Contains(err.Error(), "NotFound") {
		resp.Diagnostics.AddError(
			"Error while getting the channel policy "+err.Error(),
			err.Error(),
		)
		return
	}

	data.ID = types.StringValue(*channel.ChannelName)
	if channel.Arn != nil {
		data.Arn = types.StringValue(*channel.Arn)
	}

	if channel.ChannelName != nil {
		data.ChannelName = channel.ChannelName
	}

	if channel.ChannelState != nil {
		data.ChannelState = types.StringValue(*channel.ChannelState)
	}

	if channel.CreationTime != nil {
		data.CreationTime = types.StringValue((aws.TimeValue(channel.CreationTime)).String())
	}

	if channel.FillerSlate != nil {
		data.FillerSlate = &fillerSlateDSModel{}
		if channel.FillerSlate.SourceLocationName != nil {
			data.FillerSlate.SourceLocationName = types.StringValue(*channel.FillerSlate.SourceLocationName)
		}
		if channel.FillerSlate.VodSourceName != nil {
			data.FillerSlate.VodSourceName = types.StringValue(*channel.FillerSlate.VodSourceName)
		}
	}

	if channel.LastModifiedTime != nil {
		data.LastModifiedTime = types.StringValue((aws.TimeValue(channel.LastModifiedTime)).String())
	}

	if channel.Outputs != nil {
		data.Outputs = []outputsDSModel{}
		for _, output := range channel.Outputs {
			outputs := outputsDSModel{}
			if output.DashPlaylistSettings != nil {
				outputs.DashPlaylistSettings = &dashPlaylistSettingsDSModel{}
				if output.DashPlaylistSettings.ManifestWindowSeconds != nil {
					outputs.DashPlaylistSettings.ManifestWindowSeconds = types.Int64Value(*output.DashPlaylistSettings.ManifestWindowSeconds)
				}
				if output.DashPlaylistSettings.MinBufferTimeSeconds != nil {
					outputs.DashPlaylistSettings.MinBufferTimeSeconds = types.Int64Value(*output.DashPlaylistSettings.MinBufferTimeSeconds)
				}
				if output.DashPlaylistSettings.MinUpdatePeriodSeconds != nil {
					outputs.DashPlaylistSettings.MinUpdatePeriodSeconds = types.Int64Value(*output.DashPlaylistSettings.MinUpdatePeriodSeconds)
				}
				if output.DashPlaylistSettings.SuggestedPresentationDelaySeconds != nil {
					outputs.DashPlaylistSettings.SuggestedPresentationDelaySeconds = types.Int64Value(*output.DashPlaylistSettings.SuggestedPresentationDelaySeconds)
				}
			}
			if output.HlsPlaylistSettings != nil {
				outputs.HlsPlaylistSettings = &hlsPlaylistSettingsDSModel{}
				if output.HlsPlaylistSettings.AdMarkupType != nil && len(output.HlsPlaylistSettings.AdMarkupType) > 0 {
					outputs.HlsPlaylistSettings.AdMarkupType = []types.String{}
					output.HlsPlaylistSettings.AdMarkupType = append(output.HlsPlaylistSettings.AdMarkupType, output.HlsPlaylistSettings.AdMarkupType...)
				}
				if output.HlsPlaylistSettings.ManifestWindowSeconds != nil {
					outputs.HlsPlaylistSettings.ManifestWindowSeconds = types.Int64Value(*output.HlsPlaylistSettings.ManifestWindowSeconds)
				}
			}
			if output.ManifestName != nil {
				outputs.ManifestName = types.StringValue(*output.ManifestName)
			}
			if output.PlaybackUrl != nil {
				outputs.PlaybackUrl = types.StringValue(*output.PlaybackUrl)
			}
			if output.SourceGroup != nil {
				outputs.SourceGroup = types.StringValue(*output.SourceGroup)
			}
			data.Outputs = append(data.Outputs, outputs)
		}
	}

	if channel.PlaybackMode != nil {
		data.PlaybackMode = types.StringValue(*channel.PlaybackMode)
	}

	if policy.Policy != nil {
		data.Policy = types.StringValue(*policy.Policy)
	}

	if channel.Tags != nil && len(channel.Tags) > 0 {
		data.Tags = make(map[string]*string)
		for key, value := range channel.Tags {
			data.Tags[key] = value
		}
	}

	if channel.Tier != nil {
		data.Tier = types.StringValue(*channel.Tier)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
