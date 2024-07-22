package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &resourcePlaybackConfiguration{}
	_ resource.ResourceWithConfigure   = &resourcePlaybackConfiguration{}
	_ resource.ResourceWithImportState = &resourcePlaybackConfiguration{}
)

func ResourcePlaybackConfiguration() resource.Resource {
	return &resourcePlaybackConfiguration{}
}

type resourcePlaybackConfiguration struct {
	client *mediatailor.Client
}

func (r *resourcePlaybackConfiguration) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_playback_configuration"
}

func (r *resourcePlaybackConfiguration) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":                     computedStringWithStateForUnknown,
			"ad_decision_server_url": requiredString,
			"avail_suppression": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"fill_policy": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("FULL_AVAIL_ONLY", "PARTIAL_AVAIL"),
						},
					},
					"mode": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("OFF", "BEHIND_LIVE_EDGE", "AFTER_LIVE_EDGE"),
						},
					},
					"value": optionalString,
				},
			},
			"bumper": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"end_url":   optionalString,
					"start_url": optionalString,
				},
			},
			"cdn_configuration": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"ad_segment_url_prefix":      optionalString,
					"content_segment_url_prefix": optionalString,
				},
			},
			"configuration_aliases": schema.MapAttribute{
				Optional: true,
				ElementType: types.MapType{
					ElemType: types.StringType,
				},
			},
			"dash_configuration": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"manifest_endpoint_prefix": computedStringWithStateForUnknown,
					"mpd_location": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("DISABLED", "EMT_DEFAULT"),
						},
					},
					"origin_manifest_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf("SINGLE_PERIOD", "MULTI_PERIOD"),
						},
					},
				},
			},
			"hls_configuration_manifest_endpoint_prefix": computedStringWithStateForUnknown,
			"log_configuration_percent_enabled":          computedInt64WithStateForUnknown,
			"live_pre_roll_configuration": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"ad_decision_server_url": optionalString,
					"max_duration_seconds":   optionalInt64,
				},
			},
			"manifest_processing_rules": schema.SingleNestedAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"ad_marker_passthrough": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Optional: true,
								Computed: true,
							},
						},
					},
				},
			},
			"name":                                   requiredString,
			"personalization_threshold_seconds":      optionalInt64,
			"playback_configuration_arn":             computedStringWithStateForUnknown,
			"playback_endpoint_prefix":               computedStringWithStateForUnknown,
			"session_initialization_endpoint_prefix": computedStringWithStateForUnknown,
			"slate_ad_url":                           optionalString,
			"tags":                                   optionalMap,
			"transcode_profile_name":                 optionalString,
			"video_content_source_url":               requiredString,
		},
	}
}

func (r *resourcePlaybackConfiguration) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(clients).v2
}

func (r *resourcePlaybackConfiguration) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan playbackConfigurationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	p := putPlaybackConfigurationInputBuilder{input: &mediatailor.PutPlaybackConfigurationInput{}, model: plan}

	playbackConfiguration, err := r.client.PutPlaybackConfiguration(context.TODO(), p.getInput())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while creating playback configuration "+err.Error(),
			err.Error(),
		)
		return
	}

	m := putPlaybackConfigurationModelbuilder{model: &plan, output: *playbackConfiguration, isResource: true}

	resp.Diagnostics.Append(resp.State.Set(ctx, m.getModel())...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourcePlaybackConfiguration) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state playbackConfigurationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := state.Name

	// Get the playback configuration
	playbackConfiguration, err := r.client.GetPlaybackConfiguration(context.TODO(), &mediatailor.GetPlaybackConfigurationInput{Name: name})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while retrieving playback configuration "+err.Error(),
			err.Error(),
		)
		return
	}

	m := putPlaybackConfigurationModelbuilder{model: &state, output: mediatailor.PutPlaybackConfigurationOutput(*playbackConfiguration), isResource: true}

	// Set refreshed state
	resp.Diagnostics.Append(resp.State.Set(ctx, m.getModel())...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourcePlaybackConfiguration) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan playbackConfigurationModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// retrieve the resource playbackConfiguration
	name := plan.Name

	// Get the playback configuration
	playbackConfiguration, err := r.client.GetPlaybackConfiguration(context.TODO(), &mediatailor.GetPlaybackConfigurationInput{Name: name})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while retrieving playback configuration "+err.Error(),
			err.Error(),
		)
		return
	}

	// @ADR
	// Context: Updating tags using the PutPlaybackConfiguration method does not allow to remove them.
	// Decision: We decided to check for removed tags and remove them using the UntagResource method, while we still use
	// the PutPlaybackConfiguration method to add and update tags. We use this approach for every resource in the provider.
	// Consequences: The Update function logic is now more complicated, but tag removal is supported.

	err = V2UpdatesTags(r.client, playbackConfiguration.Tags, plan.Tags, *playbackConfiguration.PlaybackConfigurationArn)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while updating playback configuration tags"+err.Error(),
			err.Error(),
		)
	}

	p := putPlaybackConfigurationInputBuilder{input: &mediatailor.PutPlaybackConfigurationInput{}, model: plan}

	// Update the playback configuration
	playbackConfigurationUpdate, err := r.client.PutPlaybackConfiguration(context.TODO(), p.getInput())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while updating playback configuration "+err.Error(),
			err.Error(),
		)
		return
	}

	m := putPlaybackConfigurationModelbuilder{model: &plan, output: *playbackConfigurationUpdate, isResource: true}

	resp.Diagnostics.Append(resp.State.Set(ctx, m.getModel())...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourcePlaybackConfiguration) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state playbackConfigurationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	name := state.Name
	_, err := r.client.DeletePlaybackConfiguration(context.TODO(), &mediatailor.DeletePlaybackConfigurationInput{Name: name})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while deleting playback configuration "+err.Error(),
			err.Error(),
		)
		return
	}

}

func (r *resourcePlaybackConfiguration) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
