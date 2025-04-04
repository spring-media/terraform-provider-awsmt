package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"strings"
	"terraform-provider-mediatailor/awsmt/models"
)

var (
	_ datasource.DataSource              = &dataSourceChannel{}
	_ datasource.DataSourceWithConfigure = &dataSourceChannel{}
)

func DataSourceChannel() datasource.DataSource {
	return &dataSourceChannel{}
}

type dataSourceChannel struct {
	client *mediatailor.Client
}

func (d *dataSourceChannel) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_channel"
}

func (d *dataSourceChannel) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                 computedString,
			"arn":                computedString,
			"name":               requiredString,
			"channel_state":      computedString,
			"creation_time":      computedString,
			"enable_as_run_logs": computedBool,
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
								"ad_markup_type":          computedStringList,
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
				Computed:   true,
				CustomType: jsontypes.NormalizedType{},
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

	d.client = req.ProviderData.(*mediatailor.Client)
}

func (d *dataSourceChannel) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.ChannelModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	channelName := data.Name

	channel, err := d.client.DescribeChannel(ctx, &mediatailor.DescribeChannelInput{ChannelName: channelName})
	if err != nil {
		resp.Diagnostics.AddError("Error while describing channel "+*channelName, err.Error())
		return
	}

	policy, err := d.client.GetChannelPolicy(ctx, &mediatailor.GetChannelPolicyInput{ChannelName: channelName})
	if err != nil && !strings.Contains(err.Error(), "NotFound") {
		resp.Diagnostics.AddError(
			"Error while getting the channel policy "+err.Error(),
			err.Error(),
		)
		return
	}

	if policy != nil && policy.Policy != nil {
		data.Policy = jsontypes.NewNormalizedPointerValue(policy.Policy)
	}

	if channel.ChannelState != "" {
		channelState := string(channel.ChannelState)
		data.ChannelState = &channelState
	}

	data = writeChannelToState(data, *channel)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
