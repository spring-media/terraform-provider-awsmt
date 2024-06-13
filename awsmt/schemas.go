package awsmt

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var requiredString = schema.StringAttribute{
	Required: true,
}

var computedString = schema.StringAttribute{
	Computed: true,
}

var computedInt64 = schema.Int64Attribute{
	Computed: true,
}

var computedMap = schema.MapAttribute{
	Computed:    true,
	ElementType: types.StringType,
}

var computedList = schema.ListAttribute{
	Computed:    true,
	ElementType: types.StringType,
}

var computedBool = schema.BoolAttribute{
	Computed: true,
}

var optionalString = schema.StringAttribute{
	Optional: true,
}

var optionalInt64 = schema.Int64Attribute{
	Optional: true,
}

var optionalUnknownInt64 = schema.Int64Attribute{
	Optional: true,
	PlanModifiers: []planmodifier.Int64{
		int64planmodifier.UseStateForUnknown(),
	},
}

var optionalMap = schema.MapAttribute{
	Optional:    true,
	ElementType: types.StringType,
}

var optionalUnknownList = schema.ListAttribute{
	Optional:    true,
	ElementType: types.StringType,
	PlanModifiers: []planmodifier.List{
		listplanmodifier.UseStateForUnknown(),
	},
}

var optionalBool = schema.BoolAttribute{
	Optional: true,
}
