package cloud66

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCloud66Servers_UidLookup(t *testing.T) {
	t.Parallel()

	rnd := generateRandomResourceName()
	uid := generateRandomUid()
	name := fmt.Sprintf("data.cloud66_servers.%s", rnd)

	testAccCloud66Stack(uid, "awesome-app")
	testAccCloud66Servers(uid)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloud66ServersUid(rnd, uid),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloud66ServersDataSourceID(name),
					resource.TestCheckResourceAttr(name, "id", uid),
				),
			},
		},
	})
}

func testAccCloud66ServersUid(rnd string, uid string) string {
	return fmt.Sprintf(`
provider "cloud66" {
  access_token = "%[1]s"
}

data "cloud66_servers" "%[2]s" {
  uid = "%[3]s"
}
`, testAccCloud66AccessToken, rnd, uid)
}

func testAccCheckCloud66ServersDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		all := s.RootModule().Resources
		rs, ok := all[n]
		if !ok {
			return fmt.Errorf("can't find servers data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Snapshot servers source ID not set")
		}
		return nil
	}
}
