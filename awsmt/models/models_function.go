package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type FunctionModel struct {
	ID                              types.String                         `tfsdk:"id"`
	Arn                             types.String                         `tfsdk:"arn"`
	FunctionId                      *string                              `tfsdk:"function_id"`
	FunctionType                    *string                              `tfsdk:"function_type"`
	Description                     *string                              `tfsdk:"description"`
	CustomOutputConfiguration       *CustomOutputConfigurationModel      `tfsdk:"custom_output_configuration"`
	HttpRequestConfiguration        *HttpRequestConfigurationModel       `tfsdk:"http_request_configuration"`
	SequentialExecutorConfiguration *SequentialExecutorConfigurationModel `tfsdk:"sequential_executor_configuration"`
	Tags                            map[string]string                    `tfsdk:"tags"`
}

type CustomOutputConfigurationModel struct {
	Runtime *string           `tfsdk:"runtime"`
	Output  map[string]string `tfsdk:"output"`
}

type HttpRequestConfigurationModel struct {
	MethodType                 *string           `tfsdk:"method_type"`
	RequestTimeoutMilliseconds *int32            `tfsdk:"request_timeout_milliseconds"`
	Runtime                    *string           `tfsdk:"runtime"`
	Url                        *string           `tfsdk:"url"`
	Body                       *string           `tfsdk:"body"`
	Headers                    map[string]string `tfsdk:"headers"`
	Output                     map[string]string `tfsdk:"output"`
}

type SequentialExecutorConfigurationModel struct {
	FunctionList        []FunctionRefModel `tfsdk:"function_list"`
	Runtime             *string            `tfsdk:"runtime"`
	TimeoutMilliseconds *int32             `tfsdk:"timeout_milliseconds"`
	Output              map[string]string  `tfsdk:"output"`
}

type FunctionRefModel struct {
	FunctionId   *string `tfsdk:"function_id"`
	RunCondition *string `tfsdk:"run_condition"`
}
