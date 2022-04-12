package cloud66

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloud66EnvVariable() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloud66EnvVariableRead,
		Schema: map[string]*schema.Schema{
			"stack_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"readonly": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func dataSourceCloud66EnvVariableRead(d *schema.ResourceData, meta interface{}) error {
	key := d.Get("key").(string)
	d.SetId(key)

	return resourceCloud66EnvVariableRead(d, meta)
}
