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

func testAccCloud66SslCertificateLetsEncrypt(uid string) string {
	return fmt.Sprintf(`
	{
		"response": {
			"uid": "ssl-%[1]s",
			"name": "my-serv-new",
			"server_group_id": null,
			"server_names": "master.my-serv-new.c66.me",
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
		}
	}
`, uid)
}

func testAccCloud66SslCertificateManual(uid string) string {
	return fmt.Sprintf(`
	{
		"response": {
			"uid": "ssl-%[1]s",
			"name": "my-serv-new",
			"server_group_id": null,
			"server_names": "master.my-serv-new.c66.me",
			"sha256_fingerprint": "UXXsUuBNZQhNBBsPjaEATCA8t06O2RvgxuMC16q1XLCCHkIitBvMcDqoUpNO16oK",
			"ca_name": "Manual",
			"type": "manual",
			"ssl_termination": true,
			"has_intermediate_cert": false,
			"status": 3,
			"created_at": "2019-10-23T14:15:53Z",
			"updated_at": "2020-03-04T12:48:25Z",
			"expires_at": "2020-06-02T11:48:04Z",
			"certificate": "-----BEGIN CERTIFICATE-----\nMIIFQDCCBCigAwIBAgISBITqGEnFOTnEKy0WPSby659TMA0GCSqGSIb3DQEBCwUA\nMDIxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MQswCQYDVQQD\nEwJSMzAeFw0yMTA5MTIyMjE0MjdaFw0yMTEyMTEyMjE0MjZaMB4xHDAaBgNVBAMT\nE3RlcnJhZm9ybS5jZmFwaS5uZXQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK\nAoIBAQDVHcbuasE2kqFqPagrdUN2OddOkZsujnMe+GVDV65hwK8OFQGRdeiLuXhM\nc4yyAt4eEUNxP+H51HssdKPKPur9lWvBkciHGNvVsoVsWY1QKzhctcZi/TXGi89p\nqnynyMbLSEosr7QXLoVih0i6EgHIZhqT3Iz9MQd5ZymuPnyZBN/DCv32Dhdlueav\n0Q2Dqd7PThmtRBYs5odlF2MNWfwOyxRmJXfI66zTGtdgUTq8Fxk9d/RLt+kIWO7y\nBpMdIUPRVmLwkPO07tFYiG6VtcmTdPMtZsmJwcDABc0qU+U8NpRmigwnLIzsjfwb\nH06wwRhO8N1pQfBPDGtpYp4T3/LDAgMBAAGjggJiMIICXjAOBgNVHQ8BAf8EBAMC\nBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAw\nHQYDVR0OBBYEFK7MWl1qlU2RrL+WlL+QWJjY8poCMB8GA1UdIwQYMBaAFBQusxe3\nWFbLrlAJQOYfr52LFMLGMFUGCCsGAQUFBwEBBEkwRzAhBggrBgEFBQcwAYYVaHR0\ncDovL3IzLm8ubGVuY3Iub3JnMCIGCCsGAQUFBzAChhZodHRwOi8vcjMuaS5sZW5j\nci5vcmcvMDQGA1UdEQQtMCuCE3RlcnJhZm9ybS5jZmFwaS5uZXSCFHRlcnJhZm9y\nbTIuY2ZhcGkubmV0MEwGA1UdIARFMEMwCAYGZ4EMAQIBMDcGCysGAQQBgt8TAQEB\nMCgwJgYIKwYBBQUHAgEWGmh0dHA6Ly9jcHMubGV0c2VuY3J5cHQub3JnMIIBAgYK\nKwYBBAHWeQIEAgSB8wSB8ADuAHUAXNxDkv7mq0VEsV6a1FbmEDf71fpH3KFzlLJe\n5vbHDsoAAAF73EpjEQAABAMARjBEAiB7JXTWsVOKjbRJUhh8nD7BTpo4kYavQ88V\n6AdiTJJTGgIgI9gdMaF0NLpV3SO6J7LvH8ruQ+aTdgmQRoG5o89xVt0AdQD2XJQv\n0XcwIhRUGAgwlFaO400TGTO/3wwvIAvMTvFk4wAAAXvcSmMFAAAEAwBGMEQCIDqN\nolVOMaRyX57A952HltGv7kHvbpP1Cq1Hlx6wtBHvAiBpF6WhzPklj4omAmALxcHR\nmunqNwK1RTZWi0GVAVRQsjANBgkqhkiG9w0BAQsFAAOCAQEAeUhP+bGbtpwREWn6\nbDbGGg5lIBZ1zgzrotM16YcrpzS/BHOpQps7uqmt8aP68aGAyJl3lB2fF2TM8klv\nEoXvG4rGHlRtHZhllCtD1T5f9APKH88F+LoqYyp/m049LZCY/9WCgkXrqNtSbLut\nAr7b1LqvDpyS4m7cW/uG1mk3dsHjmJuwYhk3W/xWyBa6FFxHowbxDSXRGkSJ6SWC\nEXD0YagNvfpm+kNB58pJSIbBpbNL0mJA7gy2BWN58Sb0DMK+gam79QSLrZKdIlq/\nYQWun8yGsH8gHJFWyGcHtnQsGYvMd0Dr7Xf1uIOn/eQujFjF6i9/D5FTxnR5Stbb\nPwneVQ==\n-----END CERTIFICATE-----\n",
			"key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDVHcbuasE2kqFq\nPagrdUN2OddOkZsujnMe+GVDV65hwK8OFQGRdeiLuXhMc4yyAt4eEUNxP+H51Hss\ndKPKPur9lWvBkciHGNvVsoVsWY1QKzhctcZi/TXGi89pqnynyMbLSEosr7QXLoVi\nh0i6EgHIZhqT3Iz9MQd5ZymuPnyZBN/DCv32Dhdlueav0Q2Dqd7PThmtRBYs5odl\nF2MNWfwOyxRmJXfI66zTGtdgUTq8Fxk9d/RLt+kIWO7yBpMdIUPRVmLwkPO07tFY\niG6VtcmTdPMtZsmJwcDABc0qU+U8NpRmigwnLIzsjfwbH06wwRhO8N1pQfBPDGtp\nYp4T3/LDAgMBAAECggEAZZty0w0W3Xv/dXW8Diw0Y9Oj8ZO+Vu4XuPZY4UiWnYiO\nbbpaKw36N1PQJTMaK2zulYtJil8Y0FIb/9AEn1JsG0b4PyvQXYjelv4sWsI/e69/\nicQot91dnDHgS9K66Avzq8vlgXSr+jl14sn5RK19KBx2I3UNy1Fq7NjgqHCmWxV2\nfkIb4BIte6sFzhPr1uImtIY9Q6h9wwKngIxzTXrcFBGzb4HB4MZ2IcrHQHc1SeET\ncxY3/OKpoj+E+gePKIUGyEFOx8T7+XLGDlpjFPACiMsXI0oikcdG3ytw2mMN9mXt\n4fX5ZNy62dqcGeYSnJDrElUkpBMu7zblB6xpk+j6UQKBgQD7WOE8IndQ9TrIwk3l\nSnAUtQ5iuQgwFjDSB1rMGC2pkFfXiZHQVKPNe9zMcMTr9qP66EKwvU3UUxFvO61x\n5ZRm/kKt/KRnt/yKgLI7ZdjS2LvIyvyG8ZttNVhF8Y442pdLipoBHEmfp2O9Bds/\nGEy96znBn26aWK2k27jMwMkNhQKBgQDZD7kqy1zanKtY6KGrMObUfFxw9BaobpGS\nqrw/LY5t7txPOHXR8GLiuwro16t1reqP77LriZug0XQ03ULAxI0U0lrx8w62xMTL\nAQ99iUXryxRnZzStwpdJaVRfC2IoYSGqnpbb4GE/5oX55ZS1iETWHFXuVvoQcwVK\nA647TqTtpwKBgFcey5NIbwsEtUd48f8T+h1zVHUrpYblai6ilfpANzOa8Jeo+322\nmMBUuoeyXs9bQiNp9hPEygFaeaSQjuH3raS1ZO9hrqq0vzhSu3STLMCIly5WDYnI\nnRMRdnNn8uAKBH8On6ra3zoTjyKpsQEBrzf1HKPcWz3sluOZtUhjWkzxAoGAeMHS\nlghFRCnc+b2SE5dFE/mLxBtHb7Tzr9DkoZFKp8Y3MquKgJ1nphPA4gD6FqIG2MTV\nmUwZFMLyD2b4+B1hD7BngCtkiDG3+ehBIen4yFFWrKAyImkbmW/LzIScuzIudKl9\n7B1MfSxWZMxgiw2gni1tcQdaX0ReMOsTR1NdVgkCgYEAz6dwXQIz41hVEx5Nykip\n8OKvizl9pgvcmbkJJFblnSFvZbv3dfMhLbe428MMptWeaxSQvp4abE9lZnG9DdoW\nrj5hLF0qD5IsbnSKjwlu9fQpc0gmx7aYXsTTHrIVZq3Upojhqa1abcbij6NNZZyL\nhjMX9yA6Kco9zdioQGtSRuY=\n-----END PRIVATE KEY-----\n",
			"intermediate_certificate": null
		}
	}
`, uid)
}
