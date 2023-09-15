package cloud66

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCloud66Firewall() *schema.Resource {
	return &schema.Resource{
		Create: resourceCloud66FirewallCreate,
		Read:   resourceCloud66FirewallRead,
		Delete: resourceCloud66FirewallDelete,
		Importer: &schema.ResourceImporter{
			State: resourceCloud66FirewallImport,
		},

		SchemaVersion: 2,
		Schema:        resourceCloud66FirewallSchema(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
			Update: schema.DefaultTimeout(30 * time.Second),
		},
	}
}

func resourceCloud66FirewallCreate(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(ProviderConfig)
	client := providerConfig.client

	stackID := d.Get("stack_id").(string)

	newRecord := Firewall{
		Port: d.Get("port").(int),
		Ttl:  d.Get("ttl").(int),
	}
	newRecord.SetProtocol(d.Get("protocol").(string))
	log.Printf("[INFO] ProtocolCode -> %s", newRecord.ProtocolCode)

	fromIp := d.Get("from_ip").(string)
	fromGroupId := d.Get("from_group_id").(int)
	fromServerId := d.Get("from_server_id").(int)
	if fromIp != "" { //nolint:golint,gocritic
		newRecord.FromIp = fromIp
	} else if fromGroupId != 0 {
		newRecord.FromGroupId = fromGroupId
	} else if fromServerId != 0 {
		newRecord.FromServerId = fromServerId
	}

	toIp := d.Get("to_ip").(string)
	toGroupId := d.Get("to_group_id").(int)
	toServerId := d.Get("to_server_id").(int)
	if toIp != "" { //nolint:golint,gocritic
		newRecord.ToIp = toIp
	} else if toGroupId != 0 {
		newRecord.ToGroupId = toGroupId
	} else if toServerId != 0 {
		newRecord.ToServerId = toServerId
	}

	log.Printf("[INFO] Creating Firewall for stack %s", stackID)

	res, err := CreateFirewall(client, stackID, &newRecord)

	if res != nil {
		records, err := ListFirewalls(client, stackID)
		if records != nil {
			for _, record := range records {
				if newRecord.FromIp == record.FromIp &&
					newRecord.FromGroupId == record.FromGroupId &&
					newRecord.FromServerId == record.FromServerId &&
					newRecord.ToIp == record.ToIp &&
					newRecord.ToGroupId == record.ToGroupId &&
					newRecord.ToServerId == record.ToServerId &&
					newRecord.Protocol() == record.Protocol() &&
					newRecord.Port == record.Port {
					setCloud66FirewallData(d, &record)
					break
				}
			}
		} else {
			return fmt.Errorf("error reading Firewall %q: %s", stackID, err)
		}
	} else {
		return fmt.Errorf("error creating Firewall %q: %s", stackID, err)
	}

	return nil
}

func resourceCloud66FirewallRead(d *schema.ResourceData, meta interface{}) error {
	providerConfig := meta.(ProviderConfig)
	client := providerConfig.client

	stackID := d.Get("stack_id").(string)
	id := d.Id()

	record, err := GetFirewall(client, stackID, id)
	if record != nil {
		setCloud66FirewallData(d, record)
	} else {
		return fmt.Errorf("error reading Firewall %q: %s", stackID, err)
	}

	return nil
}

func resourceCloud66FirewallDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceCloud66FirewallImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// split the id so we can lookup
	idAttr := strings.SplitN(d.Id(), "/", 2)
	if len(idAttr) != 2 {
		return nil, fmt.Errorf("invalid id (\"%s\") specified, should be in format \"stackID/id\"", d.Id())
	}

	stackID, id := idAttr[0], idAttr[1]

	log.Printf("[DEBUG] Importing Firewall %s for stack %s", id, stackID)

	d.Set("stack_id", stackID)
	d.SetId(id)

	resourceCloud66FirewallRead(d, meta) //nolint:golint,errcheck

	return []*schema.ResourceData{d}, nil
}

func setCloud66FirewallData(d *schema.ResourceData, firewall *Firewall) {
	d.SetId(fmt.Sprint(firewall.Id))
	d.Set("from_ip", firewall.FromIp)
	d.Set("from_group_id", firewall.FromGroupId)
	d.Set("from_server_id", firewall.FromServerId)
	d.Set("to_ip", firewall.ToIp)
	d.Set("to_group_id", firewall.ToGroupId)
	d.Set("to_server_id", firewall.ToServerId)
	d.Set("protocol", firewall.Protocol())
	d.Set("port", firewall.Port)
	d.Set("ttl", firewall.Ttl)
}
