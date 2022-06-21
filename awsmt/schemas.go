package awsmt

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var computedString = schema.Schema{
	Type:     schema.TypeString,
	Computed: true,
}

var optionalString = schema.Schema{
	Type:     schema.TypeString,
	Optional: true,
}

var requiredString = schema.Schema{
	Type:     schema.TypeString,
	Required: true,
}

var computedInt = schema.Schema{
	Type:     schema.TypeInt,
	Computed: true,
}

var optionalInt = schema.Schema{
	Type:     schema.TypeInt,
	Optional: true,
}

var computedBool = schema.Schema{
	Type:     schema.TypeBool,
	Computed: true,
}

var optionalBool = schema.Schema{
	Type:     schema.TypeBool,
	Optional: true,
}

var computedTags = schema.Schema{
	Type:     schema.TypeMap,
	Computed: true,
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
}

var optionalTags = schema.Schema{
	Type:     schema.TypeMap,
	Optional: true,
	Elem: &schema.Schema{
		Type: schema.TypeString,
	},
}
