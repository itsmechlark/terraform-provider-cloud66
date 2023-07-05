package cloud66

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCloud66Servers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCloud66StackServers,

		Schema: map[string]*schema.Schema{
			"uid": {
				Type:     schema.TypeString,
				Required: true,
			},
			"servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ext_ipv4": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ext_ipv6": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"int_ipv4": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"int_ipv6": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dns_record": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"roles": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceCloud66StackServers(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(ProviderConfig)
	client := providerConfig.client

	log.Printf("[DEBUG] Getting servers...")
	stackID := d.Get("uid").(string)

	servers := make([]map[string]interface{}, 0)
	serversResp, err := client.Servers(stackID)
	if serversResp != nil {
		for _, serverResp := range serversResp {
			server := make(map[string]interface{})
			server["uid"] = serverResp.Uid
			server["name"] = serverResp.Name
			server["address"] = serverResp.Address
			server["ext_ipv4"] = serverResp.ExtIpV4
			server["ext_ipv6"] = serverResp.ExtIpV6
			server["int_ipv4"] = serverResp.IntIpV4
			server["int_ipv6"] = serverResp.IntIpV6
			server["dns_record"] = serverResp.DnsRecord
			server["server_type"] = serverResp.ServerType
			server["roles"] = serverResp.Roles
			servers = append(servers, server)
		}
	} else {
		return fmt.Errorf("failed to retrieve servers: %s", err)
	}

	d.SetId(stackID)
	d.Set("servers", servers)

	return nil
}
