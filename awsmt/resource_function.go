package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"terraform-provider-mediatailor/awsmt/models"
)

var (
	_ resource.Resource                = &resourceFunction{}
	_ resource.ResourceWithConfigure   = &resourceFunction{}
	_ resource.ResourceWithImportState = &resourceFunction{}
)

func ResourceFunction() resource.Resource {
	return &resourceFunction{}
}

type resourceFunction struct {
	client *mediatailor.Client
}

func (r *resourceFunction) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_function"
}

func (r *resourceFunction) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":  computedStringWithStateForUnknown,
			"arn": computedStringWithStateForUnknown,
			"function_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"function_type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("CUSTOM_OUTPUT", "HTTP_REQUEST", "SEQUENTIAL_EXECUTOR"),
				},
			},
			"description": optionalString,
			"custom_output_configuration": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"runtime": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("JSONATA"),
						},
					},
					"output": optionalMap,
				},
			},
			"http_request_configuration": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"method_type": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("GET", "POST"),
						},
					},
					"request_timeout_milliseconds": schema.Int64Attribute{
						Required: true,
					},
					"runtime": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("JSONATA"),
						},
					},
					"url":     requiredString,
					"body":    optionalString,
					"headers": optionalMap,
					"output":  optionalMap,
				},
			},
			"sequential_executor_configuration": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"function_list": schema.ListNestedAttribute{
						Required: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"function_id":   optionalString,
								"run_condition": optionalString,
							},
						},
					},
					"runtime": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							stringvalidator.OneOf("JSONATA"),
						},
					},
					"timeout_milliseconds": schema.Int64Attribute{
						Required: true,
					},
					"output": optionalMap,
				},
			},
			"tags": optionalMap,
		},
	}
}

func (r *resourceFunction) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*mediatailor.Client)
}

func (r *resourceFunction) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.FunctionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := buildPutFunctionInput(plan)
	output, err := r.client.PutFunction(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error creating function", err.Error())
		return
	}

	plan = readFunctionOutput(plan, output)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *resourceFunction) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.FunctionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := r.client.GetFunction(ctx, &mediatailor.GetFunctionInput{FunctionId: state.FunctionId})
	if err != nil {
		resp.Diagnostics.AddError("Error reading function", err.Error())
		return
	}

	state = readGetFunctionOutput(state, output)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *resourceFunction) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.FunctionModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle tag removal
	var state models.FunctionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	getOutput, err := r.client.GetFunction(ctx, &mediatailor.GetFunctionInput{FunctionId: plan.FunctionId})
	if err != nil {
		resp.Diagnostics.AddError("Error reading function for update", err.Error())
		return
	}

	if err := UpdatesTags(r.client, getOutput.Tags, plan.Tags, *getOutput.Arn); err != nil {
		resp.Diagnostics.AddError("Error updating function tags", err.Error())
		return
	}

	input := buildPutFunctionInput(plan)
	output, err := r.client.PutFunction(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error updating function", err.Error())
		return
	}

	plan = readFunctionOutput(plan, output)
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *resourceFunction) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.FunctionModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.DeleteFunction(ctx, &mediatailor.DeleteFunctionInput{FunctionId: state.FunctionId})
	if err != nil {
		resp.Diagnostics.AddError("Error deleting function", err.Error())
		return
	}
}

func (r *resourceFunction) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("function_id"), req, resp)
}

func readFunctionOutput(model models.FunctionModel, output *mediatailor.PutFunctionOutput) models.FunctionModel {
	return readFunctionToModel(model, *output)
}

func readGetFunctionOutput(model models.FunctionModel, output *mediatailor.GetFunctionOutput) models.FunctionModel {
	return readFunctionToModel(model, mediatailor.PutFunctionOutput{
		Arn:                             output.Arn,
		FunctionId:                      output.FunctionId,
		FunctionType:                    output.FunctionType,
		Description:                     output.Description,
		CustomOutputConfiguration:       output.CustomOutputConfiguration,
		HttpRequestConfiguration:        output.HttpRequestConfiguration,
		SequentialExecutorConfiguration: output.SequentialExecutorConfiguration,
		Tags:                            output.Tags,
	})
}
