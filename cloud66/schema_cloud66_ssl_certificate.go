package cloud66

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceCloud66SslCertificateSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"stack_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"ca_name": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"type": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice([]string{"lets_encrypt", "manual"}, false),
		},
		"ssl_termination": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"server_group_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"server_names": {
			Type:     schema.TypeSet,
			Required: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"certificate": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"key": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"intermediate_certificate": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"has_intermediate_cert": {
			Type:     schema.TypeBool,
			Computed: true,
			Optional: true,
		},
		"status": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
	}
}
