package awsmt

import (
	"context"
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
			"id":            computedString,
			"arn":           computedString,
			"channel_name":  requiredString,
			"channel_state": computedString,
			"creation_time": computedString,
			"filler_slate": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"source_location_name": computedString,
					"vod_source_name":      computedString,
				},
			},
			"last_modified_time": computedString,
			"outputs": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"dash_playlist_settings": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"manifest_window_seconds":              computedInt64,
								"min_buffer_time_seconds":              computedInt64,
								"min_update_period_seconds":            computedInt64,
								"suggested_presentation_delay_seconds": computedInt64,
							},
						},
						"hls_playlist_settings": schema.SingleNestedAttribute{
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"ad_markup_type":          computedList,
								"manifest_window_seconds": computedInt64,
							},
						},
						"manifest_name": computedString,
						"playback_url":  computedString,
						"source_group":  computedString,
					},
				},
			},
			"playback_mode": schema.StringAttribute{
				Computed: true,
			},
			"policy": schema.StringAttribute{
				Computed: true,
			},
			"tags": computedMap,
			"tier": computedString,
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

	if policy.Policy != nil {
		data.Policy = types.StringValue(*policy.Policy)
	}

	data = readChannelToData(data, *channel)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
