package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"reflect"
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
			"name": requiredString,
			"tags": optionalMap,
		},
	}
}

func (r *resourceSourceLocation) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(clients).v2
}

func (r *resourceSourceLocation) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sourceLocationModel

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

	plan = writeSourceLocationToPlan(plan, *sourceLocation, true)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSourceLocation) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sourceLocationModel
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

	state = writeSourceLocationToPlan(state, mediatailor.CreateSourceLocationOutput(*sourceLocation), true)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSourceLocation) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan sourceLocationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name

	sourceLocation, err := r.client.DescribeSourceLocation(ctx, &mediatailor.DescribeSourceLocationInput{SourceLocationName: name})
	if err != nil {
		resp.Diagnostics.AddError("Error while describing source location", "Could not describe the source location: "+*name+": "+err.Error())
		return
	}

	oldTags := sourceLocation.Tags
	newTags := plan.Tags

	// Check if tags are different
	if !reflect.DeepEqual(oldTags, newTags) {
		err = V2UpdatesTags(r.client, oldTags, newTags, *sourceLocation.Arn)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error while updating playback configuration tags"+err.Error(),
				err.Error(),
			)
		}
	}

	if !reflect.DeepEqual(sourceLocation.AccessConfiguration, plan.AccessConfiguration) {
		// delete source location
		name := plan.Name
		err := deleteSourceLocation(r.client, name)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error while deleting source location "+err.Error(),
				err.Error(),
			)
			return
		}

		// create source location
		params := getCreateSourceLocationInput(plan)
		sourceLocation, err := r.client.CreateSourceLocation(ctx, &params)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error while creating new source location with new access configuration "+err.Error(),
				err.Error(),
			)
			return
		}

		plan = writeSourceLocationToPlan(plan, *sourceLocation, true)
	}

	params := getUpdateSourceLocationInput(plan)

	sourceLocationUpdated, err := r.client.UpdateSourceLocation(ctx, &params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating source location "+err.Error(),
			err.Error(),
		)
		return
	}

	plan = writeSourceLocationToPlan(plan, mediatailor.CreateSourceLocationOutput(*sourceLocationUpdated), true)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *resourceSourceLocation) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state sourceLocationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := state.Name

	vodSourcesList, err := r.client.ListVodSources(ctx, &mediatailor.ListVodSourcesInput{SourceLocationName: name})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving vod sources "+err.Error(),
			err.Error(),
		)
		return
	}
	for _, vodSource := range vodSourcesList.Items {
		_, err := r.client.DeleteVodSource(ctx, &mediatailor.DeleteVodSourceInput{SourceLocationName: name, VodSourceName: vodSource.VodSourceName})
		if err != nil {
			resp.Diagnostics.AddError(
				"Error deleting vod sources "+err.Error(),
				err.Error(),
			)
			return
		}
	}

	liveSourcesList, err := r.client.ListLiveSources(ctx, &mediatailor.ListLiveSourcesInput{SourceLocationName: name})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error retrieving live sources "+err.Error(),
			err.Error(),
		)
		return
	}
	for _, liveSource := range liveSourcesList.Items {
		if _, err := r.client.DeleteLiveSource(ctx, &mediatailor.DeleteLiveSourceInput{LiveSourceName: liveSource.LiveSourceName, SourceLocationName: name}); err != nil {
			resp.Diagnostics.AddError(
				"Error deleting live sources "+err.Error(),
				err.Error(),
			)
			return
		}
	}

	_, err = r.client.DeleteSourceLocation(ctx, &mediatailor.DeleteSourceLocationInput{SourceLocationName: name})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting resource "+err.Error(),
			err.Error(),
		)
		return
	}
}

func (r *resourceSourceLocation) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
