package cloud66

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	// Provider name for single configuration testing
	ProviderNameCloud66 = "cloud66"
)

var (
	testAccProviders map[string]*schema.Provider
	testAccProvider  *schema.Provider

	// Integration test access token
	testAccCloud66AccessToken string = "bEVxqlsN800QT0UqVnDBRPXKbaPbkpOAQOMCbZTVV9uvEF2gVEcR2sHhhxMPQaMx"
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

func testAccCloud66Stacks(uid string, name string) string {
	return fmt.Sprintf(`
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
}

func testAccCloud66Stack(uid string, name string) string {
	return fmt.Sprintf(`
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
}
