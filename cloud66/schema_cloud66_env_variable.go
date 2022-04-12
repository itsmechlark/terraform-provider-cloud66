package cloud66

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceCloud66EnvVariableSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"stack_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"key": {
			Type:     schema.TypeString,
			Required: true,
		},
		"value": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"readonly": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"apply_strategy": {
			Type:         schema.TypeString,
			ValidateFunc: validation.StringInSlice([]string{"deployment", "immediately"}, false),
			Optional:     true,
			Default:      "deployment",
		},
	}
}
