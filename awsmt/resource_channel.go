package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strings"
)

var (
	_ resource.Resource                = &resourceChannel{}
	_ resource.ResourceWithConfigure   = &resourceChannel{}
	_ resource.ResourceWithImportState = &resourceChannel{}
)

func ResourceChannel() resource.Resource {
	return &resourceChannel{}
}

type resourceChannel struct {
	client *mediatailor.MediaTailor
}

func (r *resourceChannel) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_channel"
}

func (r *resourceChannel) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":   computedString,
			"arn":  computedString,
			"name": requiredString,
			// @ADR
			// Context: We cannot test the deletion of a running channel if we cannot set the channel_state property
			// through the provider
			// Decision: We decided to turn the channel_state property into an optional string and call the SDK to
			// start/stop the channel accordingly.
			// Consequences: The schema of the object differs from that of the SDK and we need to make additional
			// SDK calls.
			"channel_state": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("RUNNING", "STOPPED"),
				},
			},
			"creation_time": computedString,
			"filler_slate": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"source_location_name": optionalString,
					"vod_source_name":      optionalString,
				},
			},
			"last_modified_time": computedString,
			"outputs": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"dash_playlist_settings": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"manifest_window_seconds":              optionalInt64,
								"min_buffer_time_seconds":              optionalInt64,
								"min_update_period_seconds":            optionalInt64,
								"suggested_presentation_delay_seconds": optionalInt64,
							},
						},
						"hls_playlist_settings": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"ad_markup_type":          optionalUnknownList,
								"manifest_window_seconds": optionalUnknownInt64,
							},
						},
						"manifest_name": requiredString,
						"playback_url":  computedString,
						"source_group":  requiredString,
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
				Optional:   true,
				CustomType: jsontypes.NormalizedType{},
			},
			"tags": optionalMap,
			"tier": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"BASIC", "STANDARD"}...),
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
	var plan channelModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := buildChannelInput(plan)

	channel, err := r.client.CreateChannel(&input)
	if err != nil {
		resp.Diagnostics.AddError("Error while creating channel "+*input.ChannelName, err.Error())
		return
	}

	if plan.ChannelState != nil && *plan.ChannelState == "RUNNING" {
		_, err := r.client.StartChannel(&mediatailor.StartChannelInput{ChannelName: plan.Name})
		if err != nil {
			resp.Diagnostics.AddError("Error while starting the channel "+*channel.ChannelName, err.Error())
			return
		}
	}

	if !plan.Policy.IsNull() {
		policy := plan.Policy.ValueString()
		if err := createChannelPolicy(plan.Name, &policy, r.client); err != nil {
			resp.Diagnostics.AddError("Error while creating the channel policy for channel "+*channel.ChannelName, err.Error())
			return
		}
	}

	newPlan := writeChannelToPlan(plan, *channel)

	diags = resp.State.Set(ctx, newPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *resourceChannel) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state channelModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	channel, err := r.client.DescribeChannel(&mediatailor.DescribeChannelInput{ChannelName: state.Name})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while describing channel "+err.Error(),
			err.Error(),
		)
	}

	policy, err := r.client.GetChannelPolicy(&mediatailor.GetChannelPolicyInput{ChannelName: state.Name})
	if err != nil && !strings.Contains(err.Error(), "NotFound") {
		resp.Diagnostics.AddError(
			"Error while getting channel policy "+err.Error(),
			err.Error(),
		)
	}

	if policy.Policy != nil {
		state.Policy = jsontypes.NewNormalizedPointerValue(policy.Policy)

	} else {
		state.Policy = jsontypes.NewNormalizedNull()
	}

	state = writeChannelToState(state, *channel)

	if state.ChannelState != nil {
		state.ChannelState = channel.ChannelState
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *resourceChannel) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan channelModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	channelName := plan.Name

	channel, err := r.client.DescribeChannel(&mediatailor.DescribeChannelInput{ChannelName: channelName})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while describing channel "+err.Error(),
			err.Error(),
		)
	}

	err = updatesTags(r.client, channel.Tags, plan.Tags, *channel.Arn)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while updating channel tags"+err.Error(),
			err.Error(),
		)
	}

	previousState := channel.ChannelState
	newState := plan.ChannelState

	err = stopChannel(previousState, channelName, r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while stopping running channel "+*channelName+err.Error(),
			err.Error(),
		)
	}

	oldPolicy, err := r.client.GetChannelPolicy(&mediatailor.GetChannelPolicyInput{ChannelName: channelName})
	if err != nil && !strings.Contains(err.Error(), "NotFound") {
		resp.Diagnostics.AddError(
			"Error while getting channel policy "+err.Error(),
			err.Error(),
		)
	}

	policy := jsontypes.NewNormalizedPointerValue(oldPolicy.Policy)

	newPolicy := plan.Policy

	plan, err = updatePolicy(&plan, channelName, policy, newPolicy, r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while updating channel policy "+err.Error(),
			err.Error(),
		)
	}

	var params = buildUpdateChannelInput(plan)
	updatedChannel, err := r.client.UpdateChannel(&params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while updating channel "+*channel.ChannelName+err.Error(),
			err.Error(),
		)
	}

	wasRunning := previousState != nil && *previousState == "RUNNING"
	shouldRun := newState != nil && *newState == "RUNNING"
	if (newState == nil && wasRunning) || shouldRun {
		_, err := r.client.StartChannel(&mediatailor.StartChannelInput{ChannelName: channelName})
		if err != nil {
			resp.Diagnostics.AddError("Error while starting the channel "+*channelName, err.Error())
			return
		}
	}

	plan.ChannelState = newState

	plan = writeChannelToPlan(plan, mediatailor.CreateChannelOutput(*updatedChannel))

	// @ADR
	// Context: The official AWS Mediatailor Go SDK states that the PlaybackMode is part of the UpdateChannelOutput,
	// but it is not. As tested, the PlaybackMode is only returned when describing a channel.
	// Decision: We decided to use the previous API call to describe the channel and get the PlaybackMode from there.
	// Consequences: The PlaybackMode is not updated when updating the channel.

	plan.PlaybackMode = channel.PlaybackMode

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceChannel) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state channelModel
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

	if _, err := r.client.DeleteChannelPolicy(&mediatailor.DeleteChannelPolicyInput{ChannelName: state.Name}); err != nil {
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

func (r *resourceChannel) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
