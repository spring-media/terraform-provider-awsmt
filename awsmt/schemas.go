package awsmt

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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

var computedStringList = schema.ListAttribute{
	Computed:    true,
	ElementType: types.StringType,
}

var computedInt64List = schema.ListAttribute{
	Computed:    true,
	ElementType: types.Int64Type,
}

var computedBool = schema.BoolAttribute{
	Computed: true,
}

var optionalComputedBool = schema.BoolAttribute{
	Optional: true,
	Computed: true,
	Default:  booldefault.StaticBool(false),
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

var computedStringWithStateForUnknown = schema.StringAttribute{
	Computed: true,
	PlanModifiers: []planmodifier.String{
		stringplanmodifier.UseStateForUnknown(),
	},
}

var requiredStringWithRequiresReplace = schema.StringAttribute{
	Required: true,
	PlanModifiers: []planmodifier.String{
		stringplanmodifier.RequiresReplace(),
	},
}

var httpPackageConfigurationsResourceSchema = schema.ListNestedAttribute{
	Required: true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"path":         requiredString,
			"source_group": requiredString,
			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.OneOf("HLS", "DASH"),
				},
			},
		},
	},
}

var httpPackageConfigurationsDataSourceSchema = schema.ListNestedAttribute{
	Computed: true,
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"path":         computedString,
			"source_group": computedString,
			"type": schema.StringAttribute{
				Computed: true,
				Validators: []validator.String{
					stringvalidator.OneOf("HLS", "DASH"),
				},
			},
		},
	},
}
