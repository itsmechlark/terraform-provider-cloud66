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
