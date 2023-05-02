package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"reflect"
)

var (
	_ resource.Resource              = &resourceChannel{}
	_ resource.ResourceWithConfigure = &resourceChannel{}
)

func ResourceChannel() resource.Resource {
	return &resourceChannel{}
}

type resourceChannel struct {
	client *mediatailor.MediaTailor
}

type resourceChannelModel struct {
	Arn              types.String                       `tfsdk:"arn"`
	Name             *string                            `tfsdk:"name"`
	ChannelState     types.String                       `tfsdk:"channel_state"`
	CreationTime     types.String                       `tfsdk:"creation_time"`
	FillerSlate      []resourceChannelFillerSlatesModel `tfsdk:"filler_slate"`
	LastModifiedTime types.String                       `tfsdk:"last_modified_time"`
	Outputs          []resourceChannelOutputsModel      `tfsdk:"outputs"`
	PlaybackMode     *string                            `tfsdk:"playback_mode"`
	Policy           types.String                       `tfsdk:"policy"`
	Tags             map[string]*string                 `tfsdk:"tags"`
	Tier             *string                            `tfsdk:"tier"`
}

type resourceChannelFillerSlatesModel struct {
	SourceLocationName *string `tfsdk:"source_location_name"`
	VodSourceName      *string `tfsdk:"vod_source_name"`
}

type resourceChannelOutputsModel struct {
	DashPlaylistSettings []resourceChannelDashPlaylistSettingsModel `tfsdk:"dash_playlist_settings"`
	HlsPlaylistSettings  []resourceChannelHlsPlaylistSettingsModel  `tfsdk:"hls_playlist_settings"`
	ManifestName         *string                                    `tfsdk:"manifest_name"`
	PlaybackUrl          types.String                               `tfsdk:"playback_url"`
	SourceGroup          *string                                    `tfsdk:"source_group"`
}

type resourceChannelDashPlaylistSettingsModel struct {
	ManifestWindowsSeconds            *int64 `tfsdk:"manifest_windows_seconds"`
	MinBufferTimeSeconds              *int64 `tfsdk:"min_buffer_time_seconds"`
	MinUpdatePeriodSeconds            *int64 `tfsdk:"min_update_period_seconds"`
	SuggestedPresentationDelaySeconds *int64 `tfsdk:"suggested_presentation_delay_seconds"`
}

type resourceChannelHlsPlaylistSettingsModel struct {
	ManifestWindowsSeconds *int64 `tfsdk:"manifest_windows_seconds"`
}

func (r *resourceChannel) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_channel"
}

func (r *resourceChannel) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"arn": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			// @ADR
			// Context: We cannot test the deletion of a running channel if we cannot set the channel_state property
			// through the provider
			// Decision: We decided to turn the channel_state property into an optional string and call the SDK to
			//start/stop the channel accordingly.
			// Consequences: The schema of the object differs from that of the SDK and we need to make additional
			// SDK calls.
			"channel_state": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("RUNNING", "STOPPED"),
				},
			},
			"creation_time": schema.StringAttribute{
				Computed: true,
			},
			"filler_slate": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"source_location_name": schema.StringAttribute{
							Optional: true,
						},
						"vod_source_name": schema.StringAttribute{
							Optional: true,
						},
					},
				},
			},
			"last_modified_time": schema.StringAttribute{
				Computed: true,
			},
			"outputs": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"dash_playlist_settings": schema.ListNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"manifest_windows_seconds": schema.Int64Attribute{
										Optional: true,
									},
									"min_buffer_time_seconds": schema.Int64Attribute{
										Optional: true,
									},
									"min_update_period_seconds": schema.Int64Attribute{
										Optional: true,
									},
									"suggested_presentation_delay_seconds": schema.Int64Attribute{
										Optional: true,
									},
								},
							},
						},
						"hls_playlist_settings": schema.ListNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"manifest_windows_seconds": schema.Int64Attribute{
										Optional: true,
									},
								},
							},
						},
						"manifest_name": schema.StringAttribute{
							Required: true,
						},
						"playback_url": schema.StringAttribute{
							Computed: true,
						},
						"source_group": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
			"playback_mode": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("LINEAR", "LOOP"),
				},
			},
			// @ADR
			// Context: The provider needs to support channel policies, but such resources do not have an ARN
			// Decision: We decided to incorporate the channel policy resource in the channel resource and not to develop
			// a standalone resource.
			// Consequences: The CRUD functions for the channel resource now have to perform more than 1 API calls,
			// increasing the chances of error. Also, and the policy requires the developer to specify the ARN for the channel
			// it refers to, even if it is not known while declaring the resource, forcing the developer to create the
			// ARN themselves using the account ID and resource name.
			"policy": schema.StringAttribute{
				Optional: true,
			},
			"tags": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"tier": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("BASIC", "STANDARD"),
				},
			},
		},
	}
}

func (r *resourceChannel) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*mediatailor.MediaTailor)
}

func (r *resourceChannel) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceChannelModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var params mediatailor.CreateChannelInput

	// get channel name
	var channelName *string

	if *plan.Name != "" {
		channelName = plan.Name
	}

	// get filler slate

	fillerSlate := getFillerSlate(ctx, req, *resp)
	if fillerSlate != nil {
		params.FillerSlate = fillerSlate
	}

	// get outputs

	outputs := getOutputs(ctx, req, *resp)

	// get playback mode
	var playbackMode *string

	if *plan.PlaybackMode != "" {
		playbackMode = plan.PlaybackMode
	}

	// get tags

	tags := map[string]*string{}
	indTag := plan.Tags
	for k, value := range indTag {
		temp := *value
		tags[k] = &temp
	}

	// get Tier
	var tier *string
	if *plan.Tier != "" {
		tier = plan.Tier
	}

	params.ChannelName = channelName
	params.Outputs = outputs
	params.PlaybackMode = playbackMode
	params.Tags = tags
	params.Tier = tier

	// create channel
	channel, err := r.client.CreateChannel(&params)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error while creating channel "+err.Error(),
			err.Error(),
		)
	}
	if err = createChannelPolicy(ctx, req, r.client); err != nil {
		resp.Diagnostics.AddError(
			"Error while creating channel policy "+err.Error(),
			err.Error(),
		)
	}

	if err := checkStatusAndStartChannel(ctx, req, r.client); err != nil {
		resp.Diagnostics.AddError(
			"Error while starting channel "+err.Error(),
			err.Error(),
		)
	}

	plan.Arn = types.StringValue(aws.StringValue(channel.Arn))
	plan.Name = channel.ChannelName
	plan.CreationTime = types.StringValue((aws.TimeValue(channel.CreationTime)).String())
	plan.LastModifiedTime = types.StringValue((aws.TimeValue(channel.LastModifiedTime)).String())
	plan.PlaybackMode = channel.PlaybackMode
	plan.Tier = channel.Tier
	plan.Outputs[0].PlaybackUrl = types.StringValue(*channel.Outputs[0].PlaybackUrl)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *resourceChannel) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceChannelModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	channel, err := r.client.DescribeChannel(&mediatailor.DescribeChannelInput{
		ChannelName: state.Name,
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while reading channel "+err.Error(),
			err.Error(),
		)
	}

	policy, err := r.client.GetChannelPolicy(&mediatailor.GetChannelPolicyInput{ChannelName: state.Name})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while reading channel policy "+err.Error(),
			err.Error(),
		)
	}
	if err := setChannelPolicy(policy); err != nil {
		diag.FromErr(err)
	}

	state.Arn = types.StringValue(aws.StringValue(channel.Arn))
	state.Name = channel.ChannelName
	state.CreationTime = types.StringValue((channel.CreationTime).String())
	state.LastModifiedTime = types.StringValue((channel.LastModifiedTime).String())
	state.PlaybackMode = channel.PlaybackMode
	state.Tier = channel.Tier
	state.Outputs[0].PlaybackUrl = types.StringValue(*channel.Outputs[0].PlaybackUrl)

	if state.FillerSlate != nil && len(state.FillerSlate) > 0 {
		state.FillerSlate = append(state.FillerSlate, resourceChannelFillerSlatesModel{
			SourceLocationName: channel.FillerSlate.SourceLocationName,
			VodSourceName:      channel.FillerSlate.VodSourceName,
		})
	}

	if state.ChannelState != types.StringNull() {
		state.ChannelState = types.StringValue(*channel.ChannelState)
	}

	if state.Tags != nil {
		state.Tags = channel.Tags
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *resourceChannel) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state resourceChannelModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resourceName := plan.Name

	tagsChanged := reflect.DeepEqual(plan.Tags, state.Tags)
	if tagsChanged {
		// tags attribute was changed
		res, err := r.client.DescribeChannel(&mediatailor.DescribeChannelInput{ChannelName: resourceName})
		if err != nil {
			return
		}
		if err := updateTags(r.client, res.Arn, state.Tags, plan.Tags); err != nil {
			return
		}
	}

	res, err := r.client.DescribeChannel(&mediatailor.DescribeChannelInput{ChannelName: resourceName})
	if err != nil {
		return
	}

	previousState := res.ChannelState
	newState := plan.ChannelState

	if *previousState == "RUNNING" {
		if err := stopChannel(r.client, resourceName); err != nil {
			return
		}
	}

	var params = getUpdateChannelInput(ctx, req, *resp)
	channel, err := r.client.UpdateChannel(&params)
	if err != nil {
		return
	}

	if (*previousState == "RUNNING" || newState == types.StringValue("RUNNING")) && newState != types.StringValue("STOPPED") {
		if err := startChannel(r.client, resourceName); err != nil {
			return
		}
	}

	setStateUpdate(ctx, req, *resp, *channel)
	state.ChannelState = plan.ChannelState

	if state.Policy != plan.Policy {
		if err := updatePolicy(r.client, resourceName, (state.Policy).String(), (plan.Policy).String()); err != nil {
			return
		}
	}

	// right now policy update is not supported
	state.Policy = plan.Policy

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceChannel) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceChannelModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if _, err := r.client.StopChannel(&mediatailor.StopChannelInput{ChannelName: state.Name}); err != nil {
		resp.Diagnostics.AddError(
			"error while stopping the channel "+err.Error(),
			err.Error(),
		)
		return
	}

	if _, err := r.client.DeleteChannelPolicy(&mediatailor.DeleteChannelPolicyInput{
		ChannelName: state.Name,
	}); err != nil {
		resp.Diagnostics.AddError(
			"error while deleting the channel policy "+err.Error(),
			err.Error(),
		)
		return
	}

	if _, err := r.client.DeleteChannel(&mediatailor.DeleteChannelInput{ChannelName: state.Name}); err != nil {
		resp.Diagnostics.AddError(
			"error while deleting the channel "+err.Error(),
			err.Error(),
		)
		return
	}
}
