package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"reflect"
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
	client *mediatailor.MediaTailor
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
					"fill_policy": optionalString,
					"mode":        optionalString,
					"value":       optionalString,
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
				Required: true,
				Attributes: map[string]schema.Attribute{
					"manifest_endpoint_prefix": computedString,
					"mpd_location":             optionalString,
					"origin_manifest_type":     optionalString,
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
				Attributes: map[string]schema.Attribute{
					"ad_marker_passthrough": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"enabled": optionalBool,
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

	r.client = req.ProviderData.(clients).v1
}

func (r *resourcePlaybackConfiguration) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan playbackConfigurationModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := playbackConfigurationInput(plan)

	playbackConfiguration, err := r.client.PutPlaybackConfiguration(&input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while creating playback configuration "+err.Error(),
			err.Error(),
		)
		return
	}

	plan = readPlaybackConfigToPlan(plan, *playbackConfiguration)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourcePlaybackConfiguration) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state playbackConfigurationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := state.Name

	// Get the playback configuration
	playbackConfiguration, err := r.client.GetPlaybackConfiguration(&mediatailor.GetPlaybackConfigurationInput{Name: name})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while retrieving playback configuration "+err.Error(),
			err.Error(),
		)
		return
	}

	state = readPlaybackConfigToPlan(state, mediatailor.PutPlaybackConfigurationOutput(*playbackConfiguration))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourcePlaybackConfiguration) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan playbackConfigurationModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// retrieve the resource playbackConfiguration
	name := plan.Name

	// Get the playback configuration
	playbackConfiguration, err := r.client.GetPlaybackConfiguration(&mediatailor.GetPlaybackConfigurationInput{Name: name})
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

	oldTags := playbackConfiguration.Tags
	newTags := plan.Tags

	// Check if tags are different
	if !reflect.DeepEqual(oldTags, newTags) {
		err = untagResource(r.client, oldTags, *playbackConfiguration.PlaybackConfigurationArn)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error while untaging playback configuration tags"+err.Error(),
				err.Error(),
			)
		}
	}

	input := playbackConfigurationInput(plan)

	// Update the playback configuration
	playbackConfigurationUpdate, err := r.client.PutPlaybackConfiguration(&input)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while updating playback configuration "+err.Error(),
			err.Error(),
		)
		return
	}

	plan = readPlaybackConfigToPlan(plan, *playbackConfigurationUpdate)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourcePlaybackConfiguration) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state playbackConfigurationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	name := state.Name
	_, err := r.client.DeletePlaybackConfiguration(&mediatailor.DeletePlaybackConfigurationInput{Name: name})
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
