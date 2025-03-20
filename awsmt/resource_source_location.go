package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strings"
	"terraform-provider-mediatailor/awsmt/models"
)

var (
	_ resource.Resource                = &resourceSourceLocation{}
	_ resource.ResourceWithConfigure   = &resourceSourceLocation{}
	_ resource.ResourceWithImportState = &resourceSourceLocation{}
)

func ResourceSourceLocation() resource.Resource {
	return &resourceSourceLocation{}
}

type resourceSourceLocation struct {
	client *mediatailor.Client
}

func (r *resourceSourceLocation) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source_location"
}

func (r *resourceSourceLocation) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": computedString,
			"access_configuration": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"access_type": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("S3_SIGV4", "SECRETS_MANAGER_ACCESS_TOKEN", "AUTODETECT_SIGV4"),
						},
					},
					"smatc": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"header_name":       optionalString,
							"secret_arn":        optionalString,
							"secret_string_key": optionalString,
						},
					},
				},
			},
			"arn":           computedString,
			"creation_time": computedString,
			"default_segment_delivery_configuration": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"base_url": optionalString,
				},
			},
			"http_configuration": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"base_url": requiredString,
				},
			},
			"last_modified_time": computedString,
			"segment_delivery_configurations": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"base_url": optionalString,
						"name":     optionalString,
					},
				},
			},
			"name": requiredStringWithRequiresReplace,
			"tags": optionalMap,
		},
	}
}

func (r *resourceSourceLocation) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*mediatailor.Client)
}

func (r *resourceSourceLocation) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.SourceLocationModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := getCreateSourceLocationInput(plan)

	// Create Source Location
	sourceLocation, err := r.client.CreateSourceLocation(ctx, &params)
	if err != nil {
		resp.Diagnostics.AddError("Error while creating source location", err.Error())
		return
	}

	plan = writeSourceLocationToPlan(plan, *sourceLocation)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSourceLocation) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.SourceLocationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := state.Name

	sourceLocation, err := r.client.DescribeSourceLocation(ctx, &mediatailor.DescribeSourceLocationInput{SourceLocationName: name})
	if err != nil {
		resp.Diagnostics.AddError("Error while describing source location", "Could not describe the source location: "+*name+": "+err.Error())
		return
	}

	state = writeSourceLocationToPlan(state, mediatailor.CreateSourceLocationOutput(*sourceLocation))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSourceLocation) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var currentState, plan models.SourceLocationModel

	resp.Diagnostics.Append(req.State.Get(ctx, &currentState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name

	sourceLocation, err := r.client.DescribeSourceLocation(ctx, &mediatailor.DescribeSourceLocationInput{SourceLocationName: name})
	if err != nil {
		resp.Diagnostics.AddError("Error while describing source location", "Could not describe the source location: "+*name+": "+err.Error())
		return
	}

	err = UpdatesTags(r.client, sourceLocation.Tags, plan.Tags, *sourceLocation.Arn)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while updating playback configuration tags"+err.Error(),
			err.Error(),
		)
		return
	}
	if !currentState.AccessConfiguration.Equal(plan.AccessConfiguration) {
		updatedSourceLocation, err := recreateSourceLocation(r.client, plan)
		if err != nil {
			resp.Diagnostics.AddError("Error while recreating source location "+err.Error(), err.Error())
			return
		}
		plan = *updatedSourceLocation

	} else {
		params := getUpdateSourceLocationInput(plan)
		sourceLocationUpdated, err := r.client.UpdateSourceLocation(ctx, &params)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating source location. "+err.Error(),
				err.Error(),
			)
			return
		}
		plan = writeSourceLocationToPlan(plan, mediatailor.CreateSourceLocationOutput(*sourceLocationUpdated))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSourceLocation) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.SourceLocationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := state.Name

	err := deleteSourceLocation(r.client, name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error while deleting source location",
			err.Error(),
		)
		return
	}
}

func (r *resourceSourceLocation) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Split the import ID to support various import formats
	idParts := strings.Split(req.ID, "/")

	if len(idParts) == 1 {
		resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
		return
	}

	// Support ARN import format
	if strings.HasPrefix(req.ID, "arn:aws:mediatailor:") {
		arnParts := strings.Split(req.ID, ":")
		if len(arnParts) >= 6 {
			resourcePath := arnParts[5]
			resourceParts := strings.Split(resourcePath, "/")
			if len(resourceParts) >= 2 && resourceParts[0] == "sourceLocation" {
				resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), resourceParts[1])...)
				return
			}
		}
	}

	resp.Diagnostics.AddError(
		"Invalid import ID",
		"Expected import ID to be either the source location name or the full ARN",
	)
}
