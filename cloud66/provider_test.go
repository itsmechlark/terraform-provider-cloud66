package cloud66

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jarcoal/httpmock"
)

const (
	// Provider name for single configuration testing
	ProviderNameCloud66 = "cloud66"
)

var (
	testAccProviders map[string]*schema.Provider
	testAccProvider  *schema.Provider

	// Integration test access token
	testAccCloud66AccessToken string = "bEVxqlsN800QT0UqVnDBRPXKbaPbkpOAQOMCbZTVV9u"
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		ProviderNameCloud66: testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

type preCheckFunc = func(*testing.T)

func generateRandomResourceName() string {
	return acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
}

func generateRandomUid() string {
	return acctest.RandStringFromCharSet(32, acctest.CharSetAlphaNum)
}

func generateRandomEnvKey() string {
	return acctest.RandString(10)
}

func generateRandomEnvValue() string {
	return acctest.RandString(10)
}

func init() {
	httpmock.Activate()
}

func testAccCloud66Stacks(uid string, name string) {
	data := fmt.Sprintf(`
	{
		"response": [
			{
				"uid": "%[1]s",
				"name": "%[2]s",
				"git": "http://github.com/cloud66-samples/awesome-app.git",
				"git_branch": "fig",
				"environment": "production",
				"cloud": "DigitalOcean",
				"fqdn": "awesome-app.dev.c66.me",
				"language": "ruby",
				"framework": "rails",
				"status": 1,
				"health": 3,
				"last_activity": "2014-08-14T01:46:53+00:00",
				"last_activity_iso": "2014-08-14T01:46:53+00:00",
				"maintenance_mode": false,
				"has_loadbalancer": false,
				"created_at": "2014-08-14 00:38:14 UTC",
				"updated_at": "2014-08-14 01:46:52 UTC",
				"deploy_directory": "/var/deploy/awesome_app",
				"cloud_status": "partial",
				"created_at_iso": "2014-08-14T00:38:14Z",
				"updated_at_iso": "2014-08-14T01:46:52Z",
				"redeploy_hook": "http://hooks.cloud66.com/stacks/redeploy/b806f1c3344eb3aa2a024b23254b75b3/6d677352a6b2eefec6e345ee2b491521"
			}
		],
		"count": 1,
		"pagination": {
			"previous": null,
			"next": null,
			"current": 1,
			"per_page": 30,
			"count": 1,
			"pages": 1
		}
	}
`, uid, name)

	httpmock.RegisterResponder("GET", "https://app.cloud66.com/api/3/stacks.json", httpmock.NewStringResponder(200, data))
}

func testAccCloud66Stack(uid string, name string) {
	data := fmt.Sprintf(`
	{
		"response": {
			"uid": "%[1]s",
			"name": "%[2]s",
			"git": "http://github.com/cloud66-samples/awesome-app.git",
			"git_branch": "fig",
			"environment": "production",
			"cloud": "DigitalOcean",
			"fqdn": "awesome-app.dev.c66.me",
			"language": "ruby",
			"framework": "rails",
			"status": 1,
			"health": 3,
			"last_activity": "2014-08-14T01:46:53+00:00",
			"last_activity_iso": "2014-08-14T01:46:53+00:00",
			"maintenance_mode": false,
			"has_loadbalancer": false,
			"created_at": "2014-08-14 00:38:14 UTC",
			"updated_at": "2014-08-14 01:46:52 UTC",
			"deploy_directory": "/var/deploy/awesome_app",
			"cloud_status": "partial",
			"created_at_iso": "2014-08-14T00:38:14Z",
			"updated_at_iso": "2014-08-14T01:46:52Z",
			"redeploy_hook": "http://hooks.cloud66.com/stacks/redeploy/b806f1c3344eb3aa2a024b23254b75b3/6d677352a6b2eefec6e345ee2b491521"
		}
	}
`, uid, name)

	httpmock.RegisterResponder("GET", "https://app.cloud66.com/api/3/stacks/"+uid+".json", httpmock.NewStringResponder(200, data))
}

func testAccCloud66SslCertificateLetsEncrypt(stackID string, uid string) {
	sslData := fmt.Sprintf(`
	{
		"uuid": "ssl-%[1]s",
		"name": "my-serv-new",
		"server_group_id": null,
		"server_names": "example.com",
		"sha256_fingerprint": "UXXsUuBNZQhNBBsPjaEATCA8t06O2RvgxuMC16q1XLCCHkIitBvMcDqoUpNO16oK",
		"ca_name": "Let's Encrypt",
		"type": "lets_encrypt",
		"ssl_termination": true,
		"has_intermediate_cert": true,
		"status": 3,
		"created_at": "2019-10-23T14:15:53Z",
		"updated_at": "2020-03-04T12:48:25Z",
		"expires_at": "2020-06-02T11:48:04Z",
		"certificate": null,
		"key": null,
		"intermediate_certificate": null
	}`, uid)

	listSslResponse := fmt.Sprintf(`
	{
		"response": [%[1]s],
		"count": 1,
		"pagination": {
			"previous": null,
			"next": null,
			"current": 1,
			"per_page": 30,
			"count": 1,
			"pages": 1
		}
	}`, sslData)
	deleteSslResponse := fmt.Sprintf(`{"response": %[1]s}`, sslData)
	createSslResponse := fmt.Sprintf(`
	{
		"response": {
			"uuid": null,
			"name": null,
			"server_group_id": null,
			"server_names": "example.com",
			"sha256_fingerprint": "UXXsUuBNZQhNBBsPjaEATCA8t06O2RvgxuMC16q1XLCCHkIitBvMcDqoUpNO16oK",
			"ca_name": "Let's Encrypt",
			"type": "lets_encrypt",
			"ssl_termination": true,
			"has_intermediate_cert": true,
			"status": 1,
			"created_at": "2019-10-23T14:15:53Z",
			"updated_at": "2020-03-04T12:48:25Z",
			"expires_at": null,
			"certificate": null,
			"key": null,
			"intermediate_certificate": null
		}
	}`)

	httpmock.RegisterResponder("POST", "https://app.cloud66.com/api/3/stacks/"+stackID+"/ssl_certificates.json", httpmock.NewStringResponder(200, createSslResponse))
	httpmock.RegisterResponder("GET", "https://app.cloud66.com/api/3/stacks/"+stackID+"/ssl_certificates.json", httpmock.NewStringResponder(200, listSslResponse))
	httpmock.RegisterResponder("DELETE", "https://app.cloud66.com/api/3/stacks/"+stackID+"/ssl_certificates/ssl-"+uid+".json", httpmock.NewStringResponder(200, deleteSslResponse))
}

func testAccCloud66SslCertificateManual(stackID string, uid string) {
	sslData := fmt.Sprintf(`
	{
		"uuid": "ssl-%[1]s",
		"name": "my-serv-new",
		"server_group_id": null,
		"server_names": "example.com",
		"sha256_fingerprint": "f33832c92a78e776c15fed3f9d1f6fb4b7f0f2ce7f126c2495ea62618ef8e195",
		"ca_name": null,
		"type": "manual",
		"ssl_termination": true,
		"has_intermediate_cert": false,
		"status": 3,
		"created_at": "2019-10-23T14:15:53Z",
		"updated_at": "2020-03-04T12:48:25Z",
		"expires_at": "2020-06-02T11:48:04Z",
		"certificate": null,
		"key": null,
		"intermediate_certificate": null
	}`, uid)

	listSslResponse := fmt.Sprintf(`
	{
		"response": [%[1]s],
		"count": 1,
		"pagination": {
			"previous": null,
			"next": null,
			"current": 1,
			"per_page": 30,
			"count": 1,
			"pages": 1
		}
	}`, sslData)
	deleteSslResponse := fmt.Sprintf(`{"response": %[1]s}`, sslData)
	createSslResponse := fmt.Sprintf(`
	{
		"response": {
			"uuid": null,
			"name": null,
			"server_group_id": null,
			"server_names": "example.com",
			"sha256_fingerprint": "f33832c92a78e776c15fed3f9d1f6fb4b7f0f2ce7f126c2495ea62618ef8e195",
			"ca_name": null,
			"type": "manual",
			"ssl_termination": true,
			"has_intermediate_cert": false,
			"status": 1,
			"created_at": "2019-10-23T14:15:53Z",
			"updated_at": "2020-03-04T12:48:25Z",
			"expires_at": null,
			"certificate": null,
			"key": null,
			"intermediate_certificate": null
		}
	}`)

	httpmock.RegisterResponder("POST", "https://app.cloud66.com/api/3/stacks/"+stackID+"/ssl_certificates.json", httpmock.NewStringResponder(200, createSslResponse))
	httpmock.RegisterResponder("GET", "https://app.cloud66.com/api/3/stacks/"+stackID+"/ssl_certificates.json", httpmock.NewStringResponder(200, listSslResponse))
	httpmock.RegisterResponder("DELETE", "https://app.cloud66.com/api/3/stacks/"+stackID+"/ssl_certificates/ssl-"+uid+".json", httpmock.NewStringResponder(200, deleteSslResponse))
}

func testAccCloud66EnvVariable(stackID string, key string, value string) {
	envVarData := fmt.Sprintf(`
	{
		"id": 2426460,
		"key": "%[1]s",
		"value": "%[2]s",
		"readonly": false,
		"created_at": "2019-10-23T14:15:53Z",
		"updated_at": "2020-03-04T12:48:25Z",
		"is_password": false,
		"is_generated": false,
		"history": []
	}`, key, value)

	listEnvVarResponse := fmt.Sprintf(`
	{
		"response": [%[1]s],
		"count": 1,
		"pagination": {
			"previous": null,
			"next": null,
			"current": 1,
			"per_page": 30,
			"count": 1,
			"pages": 1
		}
	}`, envVarData)
	createEnvVarResponse := fmt.Sprint(`
	{
		"response": {
			"id": 3360669,
			"user": "some-user@example.com",
			"resource_type": "stack",
			"action": "env-var-new",
			"resource_id": "66204",
			"started_via": "api",
			"started_at": "2022-04-12T10:12:46Z",
			"finished_at": null,
			"finished_success": null,
			"finished_message": null,
			"finished_result": null
		}
	}`)
	updateEnvVarResponse := createEnvVarResponse
	deleteEnvVarResponse := createEnvVarResponse

	httpmock.RegisterResponder("POST", "https://app.cloud66.com/api/3/stacks/"+stackID+"/environments.json", httpmock.NewStringResponder(200, createEnvVarResponse))
	httpmock.RegisterResponder("GET", "https://app.cloud66.com/api/3/stacks/"+stackID+"/environments.json", httpmock.NewStringResponder(200, listEnvVarResponse))
	httpmock.RegisterResponder("PUT", "https://app.cloud66.com/api/3/stacks/"+stackID+"/environments/"+key+".json", httpmock.NewStringResponder(200, updateEnvVarResponse))
	httpmock.RegisterResponder("DELETE", "https://app.cloud66.com/api/3/stacks/"+stackID+"/environments/"+key+".json", httpmock.NewStringResponder(200, deleteEnvVarResponse))
}

func testAccCloud66FirewallRequest(stackID string) {
	firewallData := fmt.Sprintf(`
	{
		"id": 168806136,
        "from_ip": "0.0.0.0/0",
        "from_group_id": null,
        "from_server_id": null,
        "to_ip": null,
        "to_group_id": 112989,
        "to_server_id": null,
        "protocol": "tcp",
        "port": 5432,
        "rule_type": "user",
        "comments": null,
        "created_at": "2022-04-13T04:21:46Z",
        "updated_at": "2022-04-13T04:21:46Z"
	}`)

	listFirewallResponse := fmt.Sprintf(`
	{
		"response": [%[1]s],
		"count": 1,
		"pagination": {
			"previous": null,
			"next": null,
			"current": 1,
			"per_page": 30,
			"count": 1,
			"pages": 1
		}
	}`, firewallData)
	createFirewallResponse := fmt.Sprint(`
	{
		"response": {
			"id": 3360669,
			"user": "some-user@example.com",
			"resource_type": "stack",
			"action": "update_firewall",
			"resource_id": "66204",
			"started_via": "api",
			"started_at": "2022-04-12T10:12:46Z",
			"finished_at": null,
			"finished_success": null,
			"finished_message": null,
			"finished_result": null
		}
	}`)
	getFirewallResponse := fmt.Sprintf(`{ "response": %[1]s }`, firewallData)

	httpmock.RegisterResponder("POST", "https://app.cloud66.com/api/3/stacks/"+stackID+"/firewalls.json", httpmock.NewStringResponder(200, createFirewallResponse))
	httpmock.RegisterResponder("GET", "https://app.cloud66.com/api/3/stacks/"+stackID+"/firewalls.json", httpmock.NewStringResponder(200, listFirewallResponse))
	httpmock.RegisterResponder("GET", "https://app.cloud66.com/api/3/stacks/"+stackID+"/firewalls/168806136.json", httpmock.NewStringResponder(200, getFirewallResponse))
}
