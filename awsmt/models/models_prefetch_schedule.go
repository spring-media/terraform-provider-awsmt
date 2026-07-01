package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type PrefetchScheduleModel struct {
	ID                             types.String                        `tfsdk:"id"`
	Arn                            types.String                        `tfsdk:"arn"`
	Name                           *string                             `tfsdk:"name"`
	PlaybackConfigurationName      *string                             `tfsdk:"playback_configuration_name"`
	ScheduleType                   *string                             `tfsdk:"schedule_type"`
	StreamId                       *string                             `tfsdk:"stream_id"`
	Consumption                    *PrefetchConsumptionModel           `tfsdk:"consumption"`
	Retrieval                      *PrefetchRetrievalModel             `tfsdk:"retrieval"`
	RecurringPrefetchConfiguration *RecurringPrefetchConfigurationModel `tfsdk:"recurring_prefetch_configuration"`
	Tags                           map[string]string                   `tfsdk:"tags"`
}

type PrefetchConsumptionModel struct {
	EndTime               *string                     `tfsdk:"end_time"`
	StartTime             *string                     `tfsdk:"start_time"`
	AvailMatchingCriteria []AvailMatchingCriteriaModel `tfsdk:"avail_matching_criteria"`
}

type PrefetchRetrievalModel struct {
	EndTime          *string           `tfsdk:"end_time"`
	StartTime        *string           `tfsdk:"start_time"`
	DynamicVariables map[string]string `tfsdk:"dynamic_variables"`
}

type AvailMatchingCriteriaModel struct {
	DynamicVariable *string `tfsdk:"dynamic_variable"`
	Operator        *string `tfsdk:"operator"`
}

type RecurringPrefetchConfigurationModel struct {
	EndTime              *string                    `tfsdk:"end_time"`
	StartTime            *string                    `tfsdk:"start_time"`
	RecurringConsumption *RecurringConsumptionModel `tfsdk:"recurring_consumption"`
	RecurringRetrieval   *RecurringRetrievalModel   `tfsdk:"recurring_retrieval"`
}

type RecurringConsumptionModel struct {
	AvailMatchingCriteria        []AvailMatchingCriteriaModel `tfsdk:"avail_matching_criteria"`
	RetrievedAdExpirationSeconds *int64                       `tfsdk:"retrieved_ad_expiration_seconds"`
}

type RecurringRetrievalModel struct {
	DelayAfterAvailEndSeconds *int64            `tfsdk:"delay_after_avail_end_seconds"`
	DynamicVariables          map[string]string `tfsdk:"dynamic_variables"`
}
