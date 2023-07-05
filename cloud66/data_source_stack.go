package cloud66

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	api "github.com/itsmechlark/cloud66"
)

func dataSourceCloud66Stack() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloud66StackRead,

		Schema: map[string]*schema.Schema{
			"uid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"environment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fqdn": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"language": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"framework": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceCloud66StackRead(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(ProviderConfig)
	client := providerConfig.client

	log.Printf("[DEBUG] Getting stack...")
	stackID := d.Get("uid").(string)
	stackName := d.Get("name").(string)
	stackEnv := d.Get("environment").(string)

	var stack api.Stack
	if stackID != "" {
		stackResp, err := client.FindStackByUid(stackID)
		if stackResp != nil {
			stack = *stackResp
		} else {
			return fmt.Errorf("failed to find stack: %s", err)
		}
	} else if stackName != "" {
		stackResp, err := client.FindStackByName(stackName, stackEnv)
		if stackResp != nil {
			stack = *stackResp
		} else {
			return fmt.Errorf("failed to find stack: %s", err)
		}
	} else {
		return fmt.Errorf("uid or name must be set")
	}

	d.SetId(stack.Uid)
	d.Set("name", stack.Name)
	d.Set("environment", stack.Environment)
	d.Set("fqdn", stack.Fqdn)
	d.Set("language", stack.Language)
	d.Set("framework", stack.Framework)

	return nil
}
