package cloud66

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_token": {
				Type:         schema.TypeString,
				Optional:     true,
				DefaultFunc:  schema.EnvDefaultFunc("CLOUD66_ACCESS_TOKEN", nil),
				Description:  "The Access Token for operations.",
				ValidateFunc: validation.StringMatch(regexp.MustCompile("[A-Za-z0-9-_]{40}"), "Access tokens must only contain characters a-z, A-Z, 0-9, hyphens and underscores"),
			},
		},
	}

	return provider
}
