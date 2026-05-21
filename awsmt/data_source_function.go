package awsmt

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-mediatailor/awsmt/models"
)

var (
	_ datasource.DataSource              = &dataSourceFunction{}
	_ datasource.DataSourceWithConfigure = &dataSourceFunction{}
)

func DataSourceFunction() datasource.DataSource {
	return &dataSourceFunction{}
}

type dataSourceFunction struct {
	client *mediatailor.Client
}

func (d *dataSourceFunction) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_function"
}

func (d *dataSourceFunction) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id":            computedString,
			"arn":           computedString,
			"function_id":   requiredString,
			"function_type": computedString,
			"description":   computedString,
			"custom_output_configuration": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"runtime": computedString,
					"output":  computedMap,
				},
			},
			"http_request_configuration": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"method_type":                  computedString,
					"request_timeout_milliseconds": computedInt64,
					"runtime":                      computedString,
					"url":                          computedString,
					"body":                         computedString,
					"headers":                      computedMap,
					"output":                       computedMap,
				},
			},
			"sequential_executor_configuration": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"function_list": schema.ListNestedAttribute{
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"function_id":   computedString,
								"run_condition": computedString,
							},
						},
					},
					"runtime":              computedString,
					"timeout_milliseconds": computedInt64,
					"output":               computedMap,
				},
			},
			"tags": computedMap,
		},
	}
}

func (d *dataSourceFunction) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.client = req.ProviderData.(*mediatailor.Client)
}

func (d *dataSourceFunction) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.FunctionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	output, err := d.client.GetFunction(ctx, &mediatailor.GetFunctionInput{FunctionId: data.FunctionId})
	if err != nil {
		resp.Diagnostics.AddError("Error reading function", err.Error())
		return
	}

	data.ID = types.StringValue(*output.FunctionId)
	if output.Arn != nil {
		data.Arn = types.StringValue(*output.Arn)
	}
	ft := string(output.FunctionType)
	data.FunctionType = &ft
	data.Description = output.Description
	data.Tags = output.Tags

	mapCustomOutputToModel(&data, output.CustomOutputConfiguration)
	mapHttpRequestConfigToModel(&data, output.HttpRequestConfiguration)
	mapSequentialExecutorToModel(&data, output.SequentialExecutorConfiguration)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}
