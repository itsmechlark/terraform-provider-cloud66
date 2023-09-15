package cloud66

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	api "github.com/itsmechlark/cloud66"
)

type Firewall struct {
	Id           int    `json:"id"`
	FromIp       string `json:"from_ip"`
	FromGroupId  int    `json:"from_group_id"`
	FromServerId int    `json:"from_server_id"`
	ToIp         string `json:"to_ip"`
	ToGroupId    int    `json:"to_group_id"`
	ToServerId   int    `json:"to_server_id"`
	ProtocolCode string `json:"protocol"`
	Port         int    `json:"port"`
	Ttl          int    `json:"ttl"`
}

func (firewall *Firewall) Protocol() string {
	switch firewall.ProtocolCode {
	case "1":
		return "tcp"
	case "2":
		return "udp"
	case "3":
		return "icmp"
	default:
		return firewall.ProtocolCode
	}
}

func (firewall *Firewall) SetProtocol(protocol string) {
	switch protocol {
	case "tcp":
		firewall.ProtocolCode = "1"
	case "udp":
		firewall.ProtocolCode = "2"
	case "icmp":
		firewall.ProtocolCode = "3"
	default:
	}
}

func ListFirewalls(client *api.Client, stackUID string) ([]Firewall, error) {
	queryStrings := make(map[string]string)
	queryStrings["page"] = "1"

	var p api.Pagination
	var result []Firewall
	var firewallRes []Firewall

	for {
		req, err := client.NewRequest("GET", "/stacks/"+stackUID+"/firewalls.json", nil, queryStrings)
		if err != nil {
			return nil, err
		}

		firewallRes = nil
		err = client.DoReq(req, &firewallRes, &p)
		if err != nil {
			return nil, err
		}

		result = append(result, firewallRes...)
		if p.Current < p.Next {
			queryStrings["page"] = strconv.Itoa(p.Next)
		} else {
			break
		}
	}

	return result, nil
}

func CreateFirewall(client *api.Client, stackUID string, firewall *Firewall) (*api.AsyncResult, error) {
	req, err := client.NewRequest("POST", fmt.Sprintf("/stacks/%s/firewalls.json", stackUID), firewall, nil)
	if err != nil {
		return nil, err
	}

	var asyncResult *api.AsyncResult
	return asyncResult, client.DoReq(req, &asyncResult, nil)
}

func GetFirewall(client *api.Client, stackUID string, id string) (*Firewall, error) {
	req, err := client.NewRequest("GET", fmt.Sprintf("/stacks/%s/firewalls/%s.json", stackUID, id), nil, nil)
	if err != nil {
		return nil, err
	}

	var firewallResult *Firewall
	return firewallResult, client.DoReq(req, &firewallResult, nil)
}

func resourceCloud66FirewallSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"stack_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"from_ip": {
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"from_group_id", "from_server_id"},
		},
		"from_group_id": {
			Type:          schema.TypeInt,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"from_ip", "from_server_id"},
		},
		"from_server_id": {
			Type:          schema.TypeInt,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"from_ip", "from_group_id"},
		},
		"to_ip": {
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"to_group_id", "to_server_id"},
		},
		"to_group_id": {
			Type:          schema.TypeInt,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"to_ip", "to_server_id"},
		},
		"to_server_id": {
			Type:          schema.TypeInt,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"to_ip", "to_group_id"},
		},
		"protocol": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"tcp", "udp", "icmp"}, false),
			ForceNew:     true,
		},
		"port": {
			Type:     schema.TypeInt,
			Required: true,
			ForceNew: true,
		},
		"ttl": {
			Type:     schema.TypeInt,
			Optional: true,
			ForceNew: true,
		},
	}
}
