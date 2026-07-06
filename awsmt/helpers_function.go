package awsmt

import (
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/mediatailor/types"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-mediatailor/awsmt/models"
)

func readFunctionToModel(model models.FunctionModel, output mediatailor.PutFunctionOutput) models.FunctionModel {
	model.ID = types.StringValue(*output.FunctionId)
	if output.Arn != nil {
		model.Arn = types.StringValue(*output.Arn)
	}
	model.FunctionId = output.FunctionId
	ft := string(output.FunctionType)
	model.FunctionType = &ft
	model.Description = output.Description
	if len(output.Tags) > 0 {
		model.Tags = output.Tags
	}

	mapCustomOutputToModel(&model, output.CustomOutputConfiguration)
	mapHttpRequestConfigToModel(&model, output.HttpRequestConfiguration)
	mapSequentialExecutorToModel(&model, output.SequentialExecutorConfiguration)

	return model
}

func buildPutFunctionInput(model models.FunctionModel) *mediatailor.PutFunctionInput {
	input := &mediatailor.PutFunctionInput{
		FunctionId:   model.FunctionId,
		FunctionType: awsTypes.FunctionType(*model.FunctionType),
	}

	if model.Description != nil {
		input.Description = model.Description
	}
	if model.Tags != nil {
		input.Tags = model.Tags
	}

	if model.CustomOutputConfiguration != nil {
		cfg := &awsTypes.CustomOutputConfiguration{
			Runtime: awsTypes.RuntimeType(*model.CustomOutputConfiguration.Runtime),
		}
		if model.CustomOutputConfiguration.Output != nil {
			cfg.Output = model.CustomOutputConfiguration.Output
		}
		input.CustomOutputConfiguration = cfg
	}

	if model.HttpRequestConfiguration != nil {
		cfg := &awsTypes.HttpRequestConfiguration{
			MethodType:                 awsTypes.MethodType(*model.HttpRequestConfiguration.MethodType),
			RequestTimeoutMilliseconds: model.HttpRequestConfiguration.RequestTimeoutMilliseconds,
			Runtime:                    awsTypes.RuntimeType(*model.HttpRequestConfiguration.Runtime),
			Url:                        model.HttpRequestConfiguration.Url,
		}
		if model.HttpRequestConfiguration.Body != nil {
			cfg.Body = model.HttpRequestConfiguration.Body
		}
		if model.HttpRequestConfiguration.Headers != nil {
			cfg.Headers = model.HttpRequestConfiguration.Headers
		}
		if model.HttpRequestConfiguration.Output != nil {
			cfg.Output = model.HttpRequestConfiguration.Output
		}
		input.HttpRequestConfiguration = cfg
	}

	if model.SequentialExecutorConfiguration != nil {
		cfg := &awsTypes.SequentialExecutorConfiguration{
			Runtime:             awsTypes.RuntimeType(*model.SequentialExecutorConfiguration.Runtime),
			TimeoutMilliseconds: model.SequentialExecutorConfiguration.TimeoutMilliseconds,
		}
		for _, ref := range model.SequentialExecutorConfiguration.FunctionList {
			fr := awsTypes.FunctionRef{}
			if ref.FunctionId != nil {
				fr.FunctionId = ref.FunctionId
			}
			if ref.RunCondition != nil {
				fr.RunCondition = ref.RunCondition
			}
			cfg.FunctionList = append(cfg.FunctionList, fr)
		}
		if model.SequentialExecutorConfiguration.Output != nil {
			cfg.Output = model.SequentialExecutorConfiguration.Output
		}
		input.SequentialExecutorConfiguration = cfg
	}

	return input
}

func mapCustomOutputToModel(model *models.FunctionModel, cfg *awsTypes.CustomOutputConfiguration) {
	if cfg == nil {
		return
	}
	model.CustomOutputConfiguration = &models.CustomOutputConfigurationModel{
		Output: cfg.Output,
	}
	rt := string(cfg.Runtime)
	model.CustomOutputConfiguration.Runtime = &rt
}

func mapHttpRequestConfigToModel(model *models.FunctionModel, cfg *awsTypes.HttpRequestConfiguration) {
	if cfg == nil {
		return
	}
	mt := string(cfg.MethodType)
	rt := string(cfg.Runtime)
	model.HttpRequestConfiguration = &models.HttpRequestConfigurationModel{
		MethodType:                 &mt,
		RequestTimeoutMilliseconds: cfg.RequestTimeoutMilliseconds,
		Runtime:                    &rt,
		Url:                        cfg.Url,
		Body:                       cfg.Body,
		Headers:                    cfg.Headers,
		Output:                     cfg.Output,
	}
}

func mapSequentialExecutorToModel(model *models.FunctionModel, cfg *awsTypes.SequentialExecutorConfiguration) {
	if cfg == nil {
		return
	}
	rt := string(cfg.Runtime)
	sec := &models.SequentialExecutorConfigurationModel{
		Runtime:             &rt,
		TimeoutMilliseconds: cfg.TimeoutMilliseconds,
		Output:              cfg.Output,
	}
	for _, ref := range cfg.FunctionList {
		sec.FunctionList = append(sec.FunctionList, models.FunctionRefModel{
			FunctionId:   ref.FunctionId,
			RunCondition: ref.RunCondition,
		})
	}
	model.SequentialExecutorConfiguration = sec
}
