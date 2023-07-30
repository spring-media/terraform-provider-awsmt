package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
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

type resourcePlaybackConfigurationModel struct {
	ID                                  types.String                           `tfsdk:"id"`
	AdDecisionServerUrl                 *string                                `tfsdk:"ad_decision_server_url"`
	AvailSupression                     *resourceAvailSupressionModel          `tfsdk:"avail_supression"`
	Bumper                              *resourceBumperModel                   `tfsdk:"bumper"`
	CdnConfiguration                    *resourceCdnConfigurationModel         `tfsdk:"cdn_configuration"`
	ConfigurationAliases                map[string]map[string]*string          `tfsdk:"configuration_aliases"`
	DashConfiguration                   *resourceDashConfigurationModel        `tfsdk:"dash_configuration"`
	HlsConfiguration                    *resourceHlsConfigurationModel         `tfsdk:"hls_configuration"`
	LivePreRollConfiguration            *resourceLivePreRollConfigurationModel `tfsdk:"live_pre_roll_configuration"`
	LogConfiguration                    *resourceLogConfigurationModel         `tfsdk:"log_configuration"`
	ManifestProcessingRules             *resourceManifestProcessingRulesModel  `tfsdk:"manifest_processing_rules"`
	Name                                types.String                           `tfsdk:"name"`
	PersonalizationThresholdSeconds     *int64                                 `tfsdk:"personalization_threshold_seconds"`
	PlaybackConfigurationArn            types.String                           `tfsdk:"playback_configuration_arn"`
	PlaybackEndpointPrefix              types.String                           `tfsdk:"playback_endpoint_prefix"`
	SessionInitializationEndpointPrefix types.String                           `tfsdk:"session_initialization_endpoint_prefix"`
	SlateAdUrl                          *string                                `tfsdk:"slate_ad_url"`
	Tags                                map[string]*string                     `tfsdk:"tags"`
	TranscodeProfileName                *string                                `tfsdk:"transcode_profile_name"`
	VideoContentSourceUrl               *string                                `tfsdk:"video_content_source_url"`
}

type resourceAvailSupressionModel struct {
	FillPolicy *string `tfsdk:"fill_policy"`
	Mode       *string `tfsdk:"mode"`
	Value      *string `tfsdk:"value"`
}

type resourceBumperModel struct {
	EndUrl   *string `tfsdk:"end_url"`
	StartUrl *string `tfsdk:"start_url"`
}

type resourceCdnConfigurationModel struct {
	AdSegmentUrlPrefix      *string `tfsdk:"ad_segment_url_prefix"`
	ContentSegmentUrlPrefix *string `tfsdk:"content_segment_url_prefix"`
}

type resourceDashConfigurationModel struct {
	ManifestEndpointPrefix types.String `tfsdk:"manifest_endpoint_prefix"`
	MpdLocation            *string      `tfsdk:"mpd_location"`
	OriginManifestType     *string      `tfsdk:"origin_manifest_type"`
}

type resourceHlsConfigurationModel struct {
	ManifestEndpointPrefix types.String `tfsdk:"manifest_endpoint_prefix"`
}

type resourceLivePreRollConfigurationModel struct {
	AdDecisionServerUrl *string `tfsdk:"ad_decision_server_url"`
	MaxDurationSeconds  *int64  `tfsdk:"max_duration_seconds"`
}

type resourceLogConfigurationModel struct {
	PercentEnabled types.Int64 `tfsdk:"percent_enabled"`
}

type resourceManifestProcessingRulesModel struct {
	AdMarkerPassthrough *resourceAdMarkerPassthroughModel `tfsdk:"ad_marker_passthrough"`
}
type resourceAdMarkerPassthroughModel struct {
	Enabled *bool `tfsdk:"enabled"`
}

func (r *resourcePlaybackConfiguration) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_playback_configuration"
}

func (r *resourcePlaybackConfiguration) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"ad_decision_server_url": schema.StringAttribute{
				Required: true,
			},
			"avail_supression": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"fill_policy": schema.StringAttribute{
						Optional:   true,
						CustomType: types.StringType,
					},
					"mode": schema.StringAttribute{
						Optional:   true,
						CustomType: types.StringType,
					},
					"value": schema.StringAttribute{
						Optional:   true,
						CustomType: types.StringType,
					},
				},
			},
			"bumper": schema.MapNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"end_url": schema.StringAttribute{
							Optional: true,
						},
						"start_url": schema.StringAttribute{
							Optional: true,
						},
					},
				},
			},
			"cdn_configuration": schema.SingleNestedAttribute{
				Optional: true,

				Attributes: map[string]schema.Attribute{
					"ad_segment_url_prefix": schema.StringAttribute{
						Optional: true,
					},
					"content_segment_url_prefix": schema.StringAttribute{
						Optional: true,
					},
				},
			},
			"configuration_aliases": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"dash_configuration": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"manifest_endpoint_prefix": schema.StringAttribute{
						Computed: true,
					},
					"mpd_location": schema.StringAttribute{
						Optional: true,
					},
					"origin_manifest_type": schema.StringAttribute{
						Optional: true,
					},
				},
			},
			"hls_configuration": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"manifest_endpoint_prefix": schema.StringAttribute{
						Computed:   true,
						CustomType: types.StringType,
					},
				},
			},
			"live_pre_roll_configuration": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"ad_decision_server_url": schema.StringAttribute{
						Optional: true,
					},
					"max_duration_seconds": schema.Int64Attribute{
						Optional: true,
					},
				},
			},
			"log_configuration": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"percent_enabled": schema.Int64Attribute{
						Computed:   true,
						CustomType: types.Int64Type,
					},
				},
			},
			"manifest_processing_rules": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"ad_marker_passthrough": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Optional: true,
							},
						},
					},
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"personalization_threshold_seconds": schema.Int64Attribute{
				Optional: true,
			},
			"playback_configuration_arn": schema.StringAttribute{
				Computed: true,
			},
			"playback_endpoint_prefix": schema.StringAttribute{
				Computed: true,
			},
			"session_initialization_endpoint_prefix": schema.StringAttribute{
				Computed: true,
			},
			"slate_ad_url": schema.StringAttribute{
				Optional: true,
			},
			"tags": schema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"transcode_profile_name": schema.StringAttribute{
				Optional: true,
			},
			"video_content_source_url": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *resourcePlaybackConfiguration) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*mediatailor.MediaTailor)
}

func (r *resourcePlaybackConfiguration) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourcePlaybackConfigurationModel

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
	var state resourcePlaybackConfigurationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := aws.String(state.Name.String())

	// Get the playback configuration
	playbackConfiguration, err := r.client.GetPlaybackConfiguration(&mediatailor.GetPlaybackConfigurationInput{Name: name})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while retrieving playback configuration "+err.Error(),
			err.Error(),
		)
		return
	}

	state = readPlaybackConfigToState(state, *playbackConfiguration)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourcePlaybackConfiguration) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resourcePlaybackConfigurationModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// retrieve the resource playbackConfiguration
	name := aws.String(plan.Name.String())

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
		err = untagResource(r.client, oldTags, newTags, *playbackConfiguration.PlaybackConfigurationArn)
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
	var state resourcePlaybackConfigurationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	name := aws.String(state.Name.String())
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
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
