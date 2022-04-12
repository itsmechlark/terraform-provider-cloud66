package cloud66

import (
	"fmt"
	"net/http"
	"regexp"

	api "github.com/cloud66-oss/cloud66"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type ProviderConfig struct {
	client      *api.Client
	accessToken *string
}

func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CLOUD66_URL", "https://app.cloud66.com/api/3"),
				Description: "The Cloud66 URL.",
			},
			"access_token": {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("CLOUD66_ACCESS_TOKEN", nil),
				Description:  "The Access Token for operations.",
				ValidateFunc: validation.StringMatch(regexp.MustCompile("[A-Za-z0-9-_]{43}"), "Access tokens must only contain characters a-z, A-Z, 0-9, hyphens and underscores"),
			},
		},

		ConfigureFunc: providerConfigure,

		DataSourcesMap: map[string]*schema.Resource{
			"cloud66_env_variable": dataSourceCloud66EnvVariable(),
			"cloud66_stack":        dataSourceCloud66Stack(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"cloud66_env_variable":    resourceCloud66EnvVariable(),
			"cloud66_ssl_certificate": resourceCloud66SslCertificate(),
		},
	}

	return provider
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	url := d.Get("url").(string)
	accessToken := d.Get("access_token").(string)
	if accessToken == "" {
		return nil, fmt.Errorf("failed to configure Cloud66 API: %s", "access_token is required")
	}

	client := api.Client{
		URL:       url,
		UserAgent: "Terraform",
		AdditionalHeaders: http.Header(map[string][]string{
			"Authorization": {"Bearer " + accessToken},
		}),
	}

	res := ProviderConfig{
		client:      &client,
		accessToken: &accessToken,
	}

	return res, nil
}
