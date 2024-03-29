package cloud66

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	api "github.com/itsmechlark/cloud66"
)

func TestAccCloud66SslCertificate_LetsEncrypt(t *testing.T) {
	t.Parallel()

	var ssl api.SslCertificate
	rnd := generateRandomResourceName()
	stackID := generateRandomUid()
	uid := generateRandomUid()

	testAccCloud66SslCertificateLetsEncrypt(stackID, uid)

	resourceName := "cloud66_ssl_certificate." + rnd
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloud66SslCertificate_LetsEncrypt(stackID, rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCloud66SslCertificateAttributes(stackID, &ssl),
					resource.TestCheckResourceAttr(resourceName, "ca_name", "Let's Encrypt"),
					resource.TestCheckResourceAttr(resourceName, "type", "lets_encrypt"),
					resource.TestCheckResourceAttr(resourceName, "ssl_termination", "true"),
					resource.TestCheckResourceAttr(resourceName, "server_group_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "server_names.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "has_intermediate_cert", "true"),
					resource.TestCheckResourceAttr(resourceName, "sha256_fingerprint", "UXXsUuBNZQhNBBsPjaEATCA8t06O2RvgxuMC16q1XLCCHkIitBvMcDqoUpNO16oK"),
				),
			},
		},
	})
}

func testAccCloud66SslCertificate_LetsEncrypt(stactID string, rnd string) string {
	return fmt.Sprintf(`
provider "cloud66" {
	access_token = "%[1]s"
}

resource "cloud66_ssl_certificate" "%[3]s" {
	stack_id = "%[2]s"
	type = "lets_encrypt"
	server_names = ["example.com"]
}
`, testAccCloud66AccessToken, stactID, rnd)
}

func TestAccCloud66SslCertificate_Manual(t *testing.T) {
	t.Parallel()

	var ssl api.SslCertificate
	rnd := generateRandomResourceName()
	stackID := generateRandomUid()
	uid := generateRandomUid()

	testAccCloud66SslCertificateManual(stackID, uid)

	resourceName := "cloud66_ssl_certificate." + rnd
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloud66SslCertificate_Manual(stackID, rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCloud66SslCertificateAttributes(stackID, &ssl),
					resource.TestCheckResourceAttr(resourceName, "ca_name", ""),
					resource.TestCheckResourceAttr(resourceName, "type", "manual"),
					resource.TestCheckResourceAttr(resourceName, "ssl_termination", "true"),
					resource.TestCheckResourceAttr(resourceName, "server_group_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "server_names.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "has_intermediate_cert", "false"),
					resource.TestCheckResourceAttr(resourceName, "sha256_fingerprint", "f33832c92a78e776c15fed3f9d1f6fb4b7f0f2ce7f126c2495ea62618ef8e195"),
				),
			},
		},
	})
}

func testAccCloud66SslCertificate_Manual(stactID string, rnd string) string {
	return fmt.Sprintf(`
provider "cloud66" {
	access_token = "%[1]s"
}

resource "cloud66_ssl_certificate" "%[3]s" {
	stack_id = "%[2]s"
	type = "manual"
	server_names = ["example.com"]
    certificate = "-----BEGIN CERTIFICATE-----\nMIIFQDCCBCigAwIBAgISBITqGEnFOTnEKy0WPSby659TMA0GCSqGSIb3DQEBCwUA\nMDIxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MQswCQYDVQQD\nEwJSMzAeFw0yMTA5MTIyMjE0MjdaFw0yMTEyMTEyMjE0MjZaMB4xHDAaBgNVBAMT\nE3RlcnJhZm9ybS5jZmFwaS5uZXQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK\nAoIBAQDVHcbuasE2kqFqPagrdUN2OddOkZsujnMe+GVDV65hwK8OFQGRdeiLuXhM\nc4yyAt4eEUNxP+H51HssdKPKPur9lWvBkciHGNvVsoVsWY1QKzhctcZi/TXGi89p\nqnynyMbLSEosr7QXLoVih0i6EgHIZhqT3Iz9MQd5ZymuPnyZBN/DCv32Dhdlueav\n0Q2Dqd7PThmtRBYs5odlF2MNWfwOyxRmJXfI66zTGtdgUTq8Fxk9d/RLt+kIWO7y\nBpMdIUPRVmLwkPO07tFYiG6VtcmTdPMtZsmJwcDABc0qU+U8NpRmigwnLIzsjfwb\nH06wwRhO8N1pQfBPDGtpYp4T3/LDAgMBAAGjggJiMIICXjAOBgNVHQ8BAf8EBAMC\nBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAw\nHQYDVR0OBBYEFK7MWl1qlU2RrL+WlL+QWJjY8poCMB8GA1UdIwQYMBaAFBQusxe3\nWFbLrlAJQOYfr52LFMLGMFUGCCsGAQUFBwEBBEkwRzAhBggrBgEFBQcwAYYVaHR0\ncDovL3IzLm8ubGVuY3Iub3JnMCIGCCsGAQUFBzAChhZodHRwOi8vcjMuaS5sZW5j\nci5vcmcvMDQGA1UdEQQtMCuCE3RlcnJhZm9ybS5jZmFwaS5uZXSCFHRlcnJhZm9y\nbTIuY2ZhcGkubmV0MEwGA1UdIARFMEMwCAYGZ4EMAQIBMDcGCysGAQQBgt8TAQEB\nMCgwJgYIKwYBBQUHAgEWGmh0dHA6Ly9jcHMubGV0c2VuY3J5cHQub3JnMIIBAgYK\nKwYBBAHWeQIEAgSB8wSB8ADuAHUAXNxDkv7mq0VEsV6a1FbmEDf71fpH3KFzlLJe\n5vbHDsoAAAF73EpjEQAABAMARjBEAiB7JXTWsVOKjbRJUhh8nD7BTpo4kYavQ88V\n6AdiTJJTGgIgI9gdMaF0NLpV3SO6J7LvH8ruQ+aTdgmQRoG5o89xVt0AdQD2XJQv\n0XcwIhRUGAgwlFaO400TGTO/3wwvIAvMTvFk4wAAAXvcSmMFAAAEAwBGMEQCIDqN\nolVOMaRyX57A952HltGv7kHvbpP1Cq1Hlx6wtBHvAiBpF6WhzPklj4omAmALxcHR\nmunqNwK1RTZWi0GVAVRQsjANBgkqhkiG9w0BAQsFAAOCAQEAeUhP+bGbtpwREWn6\nbDbGGg5lIBZ1zgzrotM16YcrpzS/BHOpQps7uqmt8aP68aGAyJl3lB2fF2TM8klv\nEoXvG4rGHlRtHZhllCtD1T5f9APKH88F+LoqYyp/m049LZCY/9WCgkXrqNtSbLut\nAr7b1LqvDpyS4m7cW/uG1mk3dsHjmJuwYhk3W/xWyBa6FFxHowbxDSXRGkSJ6SWC\nEXD0YagNvfpm+kNB58pJSIbBpbNL0mJA7gy2BWN58Sb0DMK+gam79QSLrZKdIlq/\nYQWun8yGsH8gHJFWyGcHtnQsGYvMd0Dr7Xf1uIOn/eQujFjF6i9/D5FTxnR5Stbb\nPwneVQ==\n-----END CERTIFICATE-----\n"
    key = "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDVHcbuasE2kqFq\nPagrdUN2OddOkZsujnMe+GVDV65hwK8OFQGRdeiLuXhMc4yyAt4eEUNxP+H51Hss\ndKPKPur9lWvBkciHGNvVsoVsWY1QKzhctcZi/TXGi89pqnynyMbLSEosr7QXLoVi\nh0i6EgHIZhqT3Iz9MQd5ZymuPnyZBN/DCv32Dhdlueav0Q2Dqd7PThmtRBYs5odl\nF2MNWfwOyxRmJXfI66zTGtdgUTq8Fxk9d/RLt+kIWO7yBpMdIUPRVmLwkPO07tFY\niG6VtcmTdPMtZsmJwcDABc0qU+U8NpRmigwnLIzsjfwbH06wwRhO8N1pQfBPDGtp\nYp4T3/LDAgMBAAECggEAZZty0w0W3Xv/dXW8Diw0Y9Oj8ZO+Vu4XuPZY4UiWnYiO\nbbpaKw36N1PQJTMaK2zulYtJil8Y0FIb/9AEn1JsG0b4PyvQXYjelv4sWsI/e69/\nicQot91dnDHgS9K66Avzq8vlgXSr+jl14sn5RK19KBx2I3UNy1Fq7NjgqHCmWxV2\nfkIb4BIte6sFzhPr1uImtIY9Q6h9wwKngIxzTXrcFBGzb4HB4MZ2IcrHQHc1SeET\ncxY3/OKpoj+E+gePKIUGyEFOx8T7+XLGDlpjFPACiMsXI0oikcdG3ytw2mMN9mXt\n4fX5ZNy62dqcGeYSnJDrElUkpBMu7zblB6xpk+j6UQKBgQD7WOE8IndQ9TrIwk3l\nSnAUtQ5iuQgwFjDSB1rMGC2pkFfXiZHQVKPNe9zMcMTr9qP66EKwvU3UUxFvO61x\n5ZRm/kKt/KRnt/yKgLI7ZdjS2LvIyvyG8ZttNVhF8Y442pdLipoBHEmfp2O9Bds/\nGEy96znBn26aWK2k27jMwMkNhQKBgQDZD7kqy1zanKtY6KGrMObUfFxw9BaobpGS\nqrw/LY5t7txPOHXR8GLiuwro16t1reqP77LriZug0XQ03ULAxI0U0lrx8w62xMTL\nAQ99iUXryxRnZzStwpdJaVRfC2IoYSGqnpbb4GE/5oX55ZS1iETWHFXuVvoQcwVK\nA647TqTtpwKBgFcey5NIbwsEtUd48f8T+h1zVHUrpYblai6ilfpANzOa8Jeo+322\nmMBUuoeyXs9bQiNp9hPEygFaeaSQjuH3raS1ZO9hrqq0vzhSu3STLMCIly5WDYnI\nnRMRdnNn8uAKBH8On6ra3zoTjyKpsQEBrzf1HKPcWz3sluOZtUhjWkzxAoGAeMHS\nlghFRCnc+b2SE5dFE/mLxBtHb7Tzr9DkoZFKp8Y3MquKgJ1nphPA4gD6FqIG2MTV\nmUwZFMLyD2b4+B1hD7BngCtkiDG3+ehBIen4yFFWrKAyImkbmW/LzIScuzIudKl9\n7B1MfSxWZMxgiw2gni1tcQdaX0ReMOsTR1NdVgkCgYEAz6dwXQIz41hVEx5Nykip\n8OKvizl9pgvcmbkJJFblnSFvZbv3dfMhLbe428MMptWeaxSQvp4abE9lZnG9DdoW\nrj5hLF0qD5IsbnSKjwlu9fQpc0gmx7aYXsTTHrIVZq3Upojhqa1abcbij6NNZZyL\nhjMX9yA6Kco9zdioQGtSRuY=\n-----END PRIVATE KEY-----\n"
}`, testAccCloud66AccessToken, stactID, rnd)
}

func testAccCloud66SslCertificateAttributes(stackID string, ssl *api.SslCertificate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		actual := schema.NewSet(schema.HashString, []interface{}{})
		for _, h := range ssl.ServerNames {
			actual.Add(h)
		}
		expected := schema.NewSet(schema.HashString, []interface{}{stackID, fmt.Sprintf("*.%s", stackID)})
		if actual.Difference(expected).Len() > 0 {
			return fmt.Errorf("incorrect servernames: expected %v, got %v", expected, actual)
		}

		return nil
	}
}

func TestAccCloud66SslCertificate_Overwrite(t *testing.T) {
	t.Parallel()

	var ssl api.SslCertificate
	rnd := generateRandomResourceName()
	stackID := generateRandomUid()
	uid := generateRandomUid()

	testAccCloud66SslCertificateLetsEncrypt(stackID, uid)

	resourceName := "cloud66_ssl_certificate." + rnd
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloud66SslCertificate_Overwrite(stackID, rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCloud66SslCertificateAttributes(stackID, &ssl),
					resource.TestCheckResourceAttr(resourceName, "ca_name", "Let's Encrypt"),
					resource.TestCheckResourceAttr(resourceName, "type", "lets_encrypt"),
					resource.TestCheckResourceAttr(resourceName, "ssl_termination", "true"),
					resource.TestCheckResourceAttr(resourceName, "server_group_id", "0"),
					resource.TestCheckResourceAttr(resourceName, "server_names.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "has_intermediate_cert", "true"),
					resource.TestCheckResourceAttr(resourceName, "sha256_fingerprint", "UXXsUuBNZQhNBBsPjaEATCA8t06O2RvgxuMC16q1XLCCHkIitBvMcDqoUpNO16oK"),
				),
			},
		},
	})
}

func testAccCloud66SslCertificate_Overwrite(stactID string, rnd string) string {
	return fmt.Sprintf(`
provider "cloud66" {
	access_token = "%[1]s"
}

resource "cloud66_ssl_certificate" "%[3]s" {
	stack_id = "%[2]s"
	type = "lets_encrypt"
	server_names = ["example.com"]
	allow_overwrite = true
}
`, testAccCloud66AccessToken, stactID, rnd)
}
