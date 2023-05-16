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
	ID               types.String               `tfsdk:"id"`
	Arn              types.String               `tfsdk:"arn"`
	Name             *string                    `tfsdk:"name"`
	ChannelState     types.String               `tfsdk:"channel_state"`
	CreationTime     types.String               `tfsdk:"creation_time"`
	FillerSlate      []channelFillerSlatesModel `tfsdk:"filler_slate"`
	LastModifiedTime types.String               `tfsdk:"last_modified_time"`
	Outputs          []channelOutputsModel      `tfsdk:"outputs"`
	PlaybackMode     *string                    `tfsdk:"playback_mode"`
	Policy           types.String               `tfsdk:"policy"`
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
	ManifestWindowsSeconds            *int64 `tfsdk:"manifest_windows_seconds"`
	MinBufferTimeSeconds              *int64 `tfsdk:"min_buffer_time_seconds"`
	MinUpdatePeriodSeconds            *int64 `tfsdk:"min_update_period_seconds"`
	SuggestedPresentationDelaySeconds *int64 `tfsdk:"suggested_presentation_delay_seconds"`
}

type channelHlsPlaylistSettingsModel struct {
	ManifestWindowsSeconds *int64 `tfsdk:"manifest_windows_seconds"`
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

	name := data.Name

	channel, err := d.client.DescribeChannel(&mediatailor.DescribeChannelInput{ChannelName: name})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while retrieving the channel "+err.Error(),
			err.Error(),
		)
		return
	}

	policy, err := d.client.GetChannelPolicy(&mediatailor.GetChannelPolicyInput{ChannelName: name})
	if err != nil && !strings.Contains(err.Error(), "NotFound") {
		resp.Diagnostics.AddError(
			"Error while getting the channel policy "+err.Error(),
			err.Error(),
		)
		return
	}

	data.Arn = types.StringValue(aws.StringValue(channel.Arn))
	data.Name = channel.ChannelName
	data.ChannelState = types.StringValue(aws.StringValue(channel.ChannelState))
	data.CreationTime = types.StringValue((aws.TimeValue(channel.CreationTime)).String())
	data.LastModifiedTime = types.StringValue((aws.TimeValue(channel.LastModifiedTime)).String())
	data.PlaybackMode = channel.PlaybackMode
	data.Policy = types.StringValue(*policy.Policy)
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
			output.PlaybackUrl = o.PlaybackUrl
		}

		if o.HlsPlaylistSettings[0].ManifestWindowsSeconds != nil {
			var outputsHls []channelHlsPlaylistSettingsModel
			outputsHls[0].ManifestWindowsSeconds = o.HlsPlaylistSettings[0].ManifestWindowsSeconds
			output.HlsPlaylistSettings = outputsHls
		}

		var outputsDash []channelDashPlaylistSettingsModel
		if o.DashPlaylistSettings[0].ManifestWindowsSeconds != nil {
			outputsDash[0].ManifestWindowsSeconds = o.DashPlaylistSettings[0].ManifestWindowsSeconds
		}
		if o.DashPlaylistSettings[0].MinBufferTimeSeconds != nil {
			outputsDash[0].MinBufferTimeSeconds = o.DashPlaylistSettings[0].MinBufferTimeSeconds
		}
		if o.DashPlaylistSettings[0].MinUpdatePeriodSeconds != nil {
			outputsDash[0].MinUpdatePeriodSeconds = o.DashPlaylistSettings[0].MinUpdatePeriodSeconds
		}
		if o.DashPlaylistSettings[0].SuggestedPresentationDelaySeconds != nil {
			outputsDash[0].SuggestedPresentationDelaySeconds = o.DashPlaylistSettings[0].SuggestedPresentationDelaySeconds
		}

		output.DashPlaylistSettings = outputsDash
		data.Outputs = append(data.Outputs, output)
	}

	data.ID = types.StringValue(aws.StringValue(channel.ChannelName))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
