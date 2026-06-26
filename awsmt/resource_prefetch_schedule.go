package awsmt

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strings"
	"terraform-provider-mediatailor/awsmt/models"
)

var (
	_ resource.Resource                = &resourcePrefetchSchedule{}
	_ resource.ResourceWithConfigure   = &resourcePrefetchSchedule{}
	_ resource.ResourceWithImportState = &resourcePrefetchSchedule{}
)

func ResourcePrefetchSchedule() resource.Resource {
	return &resourcePrefetchSchedule{}
}

type resourcePrefetchSchedule struct {
	client *mediatailor.Client
}

func (r *resourcePrefetchSchedule) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prefetch_schedule"
}

func (r *resourcePrefetchSchedule) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":  computedStringWithStateForUnknown,
			"arn": computedStringWithStateForUnknown,
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"playback_configuration_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"schedule_type": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf("SINGLE", "RECURRING"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"stream_id": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"consumption": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"end_time":   requiredString,
					"start_time": optionalString,
					"avail_matching_criteria": schema.ListNestedAttribute{
						Optional: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"dynamic_variable": requiredString,
								"operator": schema.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.OneOf("EQUALS"),
									},
								},
							},
						},
					},
				},
			},
			"retrieval": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"end_time":          requiredString,
					"start_time":        optionalString,
					"dynamic_variables": optionalMap,
				},
			},
			"recurring_prefetch_configuration": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"end_time":   requiredString,
					"start_time": optionalString,
					"recurring_consumption": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"avail_matching_criteria": schema.ListNestedAttribute{
								Optional: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"dynamic_variable": requiredString,
										"operator": schema.StringAttribute{
											Required: true,
											Validators: []validator.String{
												stringvalidator.OneOf("EQUALS"),
											},
										},
									},
								},
							},
							"retrieved_ad_expiration_seconds": optionalInt64,
						},
					},
					"recurring_retrieval": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"delay_after_avail_end_seconds": optionalInt64,
							"dynamic_variables":             optionalMap,
						},
					},
				},
			},
			"tags": optionalMap,
		},
	}
}

func (r *resourcePrefetchSchedule) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*mediatailor.Client)
}

func (r *resourcePrefetchSchedule) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.PrefetchScheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := buildCreatePrefetchScheduleInput(plan)
	output, err := r.client.CreatePrefetchSchedule(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error creating prefetch schedule", err.Error())
		return
	}

	plan = readPrefetchScheduleOutput(plan, output)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *resourcePrefetchSchedule) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.PrefetchScheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := r.client.GetPrefetchSchedule(ctx, &mediatailor.GetPrefetchScheduleInput{
		Name:                      state.Name,
		PlaybackConfigurationName: state.PlaybackConfigurationName,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error reading prefetch schedule", err.Error())
		return
	}

	state = readGetPrefetchScheduleOutput(state, output)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *resourcePrefetchSchedule) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	// PrefetchSchedule does not support updates - all mutable fields force replacement
	resp.Diagnostics.AddError("Update not supported", "Prefetch schedules cannot be updated. All changes require replacement.")
}

func (r *resourcePrefetchSchedule) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.PrefetchScheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeletePrefetchSchedule(ctx, &mediatailor.DeletePrefetchScheduleInput{
		Name:                      state.Name,
		PlaybackConfigurationName: state.PlaybackConfigurationName,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error deleting prefetch schedule", err.Error())
		return
	}
}

// Import uses composite ID: playback_configuration_name/prefetch_schedule_name
func (r *resourcePrefetchSchedule) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) != 2 {
		resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf("Expected format: playback_configuration_name/prefetch_schedule_name, got: %s", req.ID))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("playback_configuration_name"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), parts[1])...)
}
