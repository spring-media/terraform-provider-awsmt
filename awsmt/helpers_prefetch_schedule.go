package awsmt

import (
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/mediatailor/types"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-mediatailor/awsmt/models"
	"time"
)

func buildCreatePrefetchScheduleInput(model models.PrefetchScheduleModel) *mediatailor.CreatePrefetchScheduleInput {
	input := &mediatailor.CreatePrefetchScheduleInput{
		Name:                      model.Name,
		PlaybackConfigurationName: model.PlaybackConfigurationName,
	}

	if model.ScheduleType != nil {
		input.ScheduleType = awsTypes.PrefetchScheduleType(*model.ScheduleType)
	}
	if model.StreamId != nil {
		input.StreamId = model.StreamId
	}
	if model.Tags != nil {
		input.Tags = model.Tags
	}

	if model.Consumption != nil {
		consumption := &awsTypes.PrefetchConsumption{}
		if model.Consumption.EndTime != nil {
			t, _ := time.Parse(time.RFC3339, *model.Consumption.EndTime)
			consumption.EndTime = &t
		}
		if model.Consumption.StartTime != nil {
			t, _ := time.Parse(time.RFC3339, *model.Consumption.StartTime)
			consumption.StartTime = &t
		}
		for _, c := range model.Consumption.AvailMatchingCriteria {
			consumption.AvailMatchingCriteria = append(consumption.AvailMatchingCriteria, awsTypes.AvailMatchingCriteria{
				DynamicVariable: c.DynamicVariable,
				Operator:        awsTypes.Operator(*c.Operator),
			})
		}
		input.Consumption = consumption
	}

	if model.Retrieval != nil {
		retrieval := &awsTypes.PrefetchRetrieval{}
		if model.Retrieval.EndTime != nil {
			t, _ := time.Parse(time.RFC3339, *model.Retrieval.EndTime)
			retrieval.EndTime = &t
		}
		if model.Retrieval.StartTime != nil {
			t, _ := time.Parse(time.RFC3339, *model.Retrieval.StartTime)
			retrieval.StartTime = &t
		}
		if model.Retrieval.DynamicVariables != nil {
			retrieval.DynamicVariables = model.Retrieval.DynamicVariables
		}
		input.Retrieval = retrieval
	}

	if model.RecurringPrefetchConfiguration != nil {
		rpc := &awsTypes.RecurringPrefetchConfiguration{}
		if model.RecurringPrefetchConfiguration.EndTime != nil {
			t, _ := time.Parse(time.RFC3339, *model.RecurringPrefetchConfiguration.EndTime)
			rpc.EndTime = &t
		}
		if model.RecurringPrefetchConfiguration.StartTime != nil {
			t, _ := time.Parse(time.RFC3339, *model.RecurringPrefetchConfiguration.StartTime)
			rpc.StartTime = &t
		}
		if model.RecurringPrefetchConfiguration.RecurringConsumption != nil {
			rc := &awsTypes.RecurringConsumption{}
			if model.RecurringPrefetchConfiguration.RecurringConsumption.RetrievedAdExpirationSeconds != nil {
				exp := int32(*model.RecurringPrefetchConfiguration.RecurringConsumption.RetrievedAdExpirationSeconds)
				rc.RetrievedAdExpirationSeconds = &exp
			}
			for _, c := range model.RecurringPrefetchConfiguration.RecurringConsumption.AvailMatchingCriteria {
				rc.AvailMatchingCriteria = append(rc.AvailMatchingCriteria, awsTypes.AvailMatchingCriteria{
					DynamicVariable: c.DynamicVariable,
					Operator:        awsTypes.Operator(*c.Operator),
				})
			}
			rpc.RecurringConsumption = rc
		}
		if model.RecurringPrefetchConfiguration.RecurringRetrieval != nil {
			rr := &awsTypes.RecurringRetrieval{}
			if model.RecurringPrefetchConfiguration.RecurringRetrieval.DelayAfterAvailEndSeconds != nil {
				delay := int32(*model.RecurringPrefetchConfiguration.RecurringRetrieval.DelayAfterAvailEndSeconds)
				rr.DelayAfterAvailEndSeconds = &delay
			}
			if model.RecurringPrefetchConfiguration.RecurringRetrieval.DynamicVariables != nil {
				rr.DynamicVariables = model.RecurringPrefetchConfiguration.RecurringRetrieval.DynamicVariables
			}
			rpc.RecurringRetrieval = rr
		}
		input.RecurringPrefetchConfiguration = rpc
	}

	return input
}

func readPrefetchScheduleOutput(model models.PrefetchScheduleModel, output *mediatailor.CreatePrefetchScheduleOutput) models.PrefetchScheduleModel {
	if output.Arn != nil {
		model.Arn = types.StringValue(*output.Arn)
	}
	model.ID = types.StringValue(*output.Name)
	model.Name = output.Name
	model.PlaybackConfigurationName = output.PlaybackConfigurationName
	if output.ScheduleType != "" {
		st := string(output.ScheduleType)
		model.ScheduleType = &st
	}
	model.StreamId = output.StreamId
	model.Tags = output.Tags
	return model
}

func readGetPrefetchScheduleOutput(model models.PrefetchScheduleModel, output *mediatailor.GetPrefetchScheduleOutput) models.PrefetchScheduleModel {
	if output.Arn != nil {
		model.Arn = types.StringValue(*output.Arn)
	}
	model.ID = types.StringValue(*output.Name)
	model.Name = output.Name
	model.PlaybackConfigurationName = output.PlaybackConfigurationName
	if output.ScheduleType != "" {
		st := string(output.ScheduleType)
		model.ScheduleType = &st
	}
	model.StreamId = output.StreamId
	model.Tags = output.Tags
	return model
}
