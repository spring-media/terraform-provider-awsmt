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
	Arn              types.String               `tfsdk:"arn"`
	Name             types.String               `tfsdk:"name"`
	ChannelState     types.String               `tfsdk:"channel_state"`
	CreationTime     types.String               `tfsdk:"creation_time"`
	FillerSlate      []channelFillerSlatesModel `tfsdk:"filler_slate"`
	LastModifiedTime types.String               `tfsdk:"last_modified_time"`
	Outputs          []channelOutputsModel      `tfsdk:"outputs"`
	PlaybackMode     *string                    `tfsdk:"playback_mode"`
	Policy           *string                    `tfsdk:"policy"`
	Tags             map[string]*string         `tfsdk:"tags"`
	Tier             *string                    `tfsdk:"tier"`
}

type channelFillerSlatesModel struct {
	SourceLocationName *string `tfsdk:"source_location_name"`
	VodSourceName      *string `tfsdk:"vod_source_name"`
}

type channelOutputsModel struct {
	DashPlaylistSettings []channelDashPlaylistSettingsModel `tfsdk:"dash_playlist_settings"`
	HlsPlaylistSettings  []channelHlsPlaylistSettingsModel  `tfsdk:"hls_playlist_settings"`
	ManifestName         *string                            `tfsdk:"manifest_name"`
	PlaybackUrl          types.String                       `tfsdk:"playback_url"`
	SourceGroup          *string                            `tfsdk:"source_group"`
}

type channelDashPlaylistSettingsModel struct {
	ManifestWindowsSeconds            types.Int64 `tfsdk:"manifest_windows_seconds"`
	MinBufferTimeSeconds              types.Int64 `tfsdk:"min_buffer_time_seconds"`
	MinUpdatePeriodSeconds            types.Int64 `tfsdk:"min_update_period_seconds"`
	SuggestedPresentationDelaySeconds types.Int64 `tfsdk:"suggested_presentation_delay_seconds"`
}

type channelHlsPlaylistSettingsModel struct {
	ManifestWindowsSeconds types.Int64 `tfsdk:"manifest_windows_seconds"`
}

func (d *dataSourceChannel) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_channel"
}

func (d *dataSourceChannel) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{

			"arn": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"channel_state": schema.StringAttribute{
				Computed: true,
			},
			"creation_time": schema.StringAttribute{
				Computed: true,
			},
			"filler_slate": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"source_location_name": schema.StringAttribute{
							Computed: true,
						},
						"vod_source_name": schema.StringAttribute{
							Computed: true,
						},
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
						"dash_playlist_settings": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"manifest_windows_seconds": schema.Int64Attribute{
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
						},
						"hls_playlist_settings": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"manifest_windows_seconds": schema.Int64Attribute{
										Computed: true,
									},
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

	channel, err := d.client.DescribeChannel(&mediatailor.DescribeChannelInput{ChannelName: aws.String(data.Name.ValueString())})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error describing channel "+err.Error(),
			err.Error(),
		)
		return
	}

	policy, err := d.client.GetChannelPolicy(&mediatailor.GetChannelPolicyInput{ChannelName: aws.String(data.Name.ValueString())})
	if err != nil && !strings.Contains(err.Error(), "NotFound") {
		resp.Diagnostics.AddError(
			"Error getting channel policy "+err.Error(),
			err.Error(),
		)
		return
	} else {
		if policy.Policy == nil {
			policy.Policy = aws.String("policy")
		}
	}

	data.Arn = types.StringValue(aws.StringValue(channel.Arn))
	data.Name = types.StringValue(aws.StringValue(channel.ChannelName))
	data.ChannelState = types.StringValue(aws.StringValue(channel.ChannelState))
	data.CreationTime = types.StringValue((aws.TimeValue(channel.CreationTime)).String())
	data.LastModifiedTime = types.StringValue((aws.TimeValue(channel.LastModifiedTime)).String())
	data.PlaybackMode = channel.PlaybackMode
	data.Policy = policy.Policy
	data.Tags = channel.Tags
	data.Tier = channel.Tier

	if data.FillerSlate != nil && len(data.FillerSlate) > 0 {
		data.FillerSlate = append(data.FillerSlate, channelFillerSlatesModel{
			SourceLocationName: channel.FillerSlate.SourceLocationName,
			VodSourceName:      channel.FillerSlate.VodSourceName,
		})
	}

	for _, o := range data.Outputs {
		output := channelOutputsModel{}
		if *o.ManifestName != "" && o.ManifestName != nil {
			output.ManifestName = o.ManifestName
		}

		if *o.SourceGroup != "" {
			output.SourceGroup = o.SourceGroup
		}

		if o.PlaybackUrl != types.StringValue("") {
			output.SourceGroup = o.SourceGroup
		}

		if o.HlsPlaylistSettings[0].ManifestWindowsSeconds != types.Int64Null() {
			outputsHls := channelHlsPlaylistSettingsModel{
				ManifestWindowsSeconds: o.HlsPlaylistSettings[0].ManifestWindowsSeconds,
			}
			output.HlsPlaylistSettings[0] = outputsHls
		}

		if o.DashPlaylistSettings[0].ManifestWindowsSeconds != types.Int64Null() || o.DashPlaylistSettings[0].MinBufferTimeSeconds != types.Int64Null() || o.DashPlaylistSettings[0].MinUpdatePeriodSeconds != types.Int64Null() || o.DashPlaylistSettings[0].SuggestedPresentationDelaySeconds != types.Int64Null() {
			manifestWindowSecondsDash := o.DashPlaylistSettings[0].ManifestWindowsSeconds
			minBufferTimeSeconds := o.DashPlaylistSettings[0].MinBufferTimeSeconds
			minUpdatePeriodSeconds := o.DashPlaylistSettings[0].MinUpdatePeriodSeconds
			suggestedPresentationDelaySeconds := o.DashPlaylistSettings[0].SuggestedPresentationDelaySeconds
			outputsDash := channelDashPlaylistSettingsModel{
				ManifestWindowsSeconds:            manifestWindowSecondsDash,
				MinBufferTimeSeconds:              minBufferTimeSeconds,
				MinUpdatePeriodSeconds:            minUpdatePeriodSeconds,
				SuggestedPresentationDelaySeconds: suggestedPresentationDelaySeconds,
			}

			output.DashPlaylistSettings[0] = outputsDash
		}
		data.Outputs = append(data.Outputs, output)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
