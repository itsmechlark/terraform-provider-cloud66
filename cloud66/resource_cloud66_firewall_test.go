package cloud66

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCloud66Firewall(t *testing.T) {
	t.Parallel()

	rnd := generateRandomResourceName()
	stackID := generateRandomUid()

	testAccCloud66FirewallRequest(stackID)

	resourceName := "cloud66_firewall." + rnd
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloud66Firewall(stackID, rnd),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "stack_id", stackID),
					resource.TestCheckResourceAttr(resourceName, "from_ip", "0.0.0.0/0"),
					resource.TestCheckResourceAttr(resourceName, "to_group_id", "112989"),
					resource.TestCheckResourceAttr(resourceName, "protocol", "tcp"),
					resource.TestCheckResourceAttr(resourceName, "port", "5432"),
				),
			},
		},
	})
}

func testAccCloud66Firewall(stactID string, rnd string) string {
	return fmt.Sprintf(`
provider "cloud66" {
	access_token = "%[1]s"
}

resource "cloud66_firewall" "%[3]s" {
	stack_id = "%[2]s"
	from_ip = "0.0.0.0/0"
	to_group_id = 112989
	protocol = "tcp"
	port = 5432
}
`, testAccCloud66AccessToken, stactID, rnd)
}
