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
