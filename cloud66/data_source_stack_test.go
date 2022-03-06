package cloud66

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCloud66Stack_NameLookup(t *testing.T) {

	t.Parallel()
	rnd := generateRandomResourceName()
	uid := generateRandomUid()
	name := fmt.Sprintf("data.cloud66_stack.%s", rnd)

	testAccCloud66Stacks(uid, "awesome-app")

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloud66StackName(rnd),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloud66StackDataSourceID(name),
					resource.TestCheckResourceAttr(name, "id", uid),
					resource.TestCheckResourceAttr(name, "name", "awesome-app"),
					resource.TestCheckResourceAttr(name, "environment", "production"),
					resource.TestCheckResourceAttr(name, "fqdn", "awesome-app.dev.c66.me"),
					resource.TestCheckResourceAttr(name, "language", "ruby"),
					resource.TestCheckResourceAttr(name, "framework", "rails"),
				),
			},
		},
	})
}

func testAccCloud66StackName(rnd string) string {
	return fmt.Sprintf(`
provider "cloud66" {
  access_token = "%[1]s"
}

data "cloud66_stack" "%[2]s" {
  name = "awesome-app"
}
`, testAccCloud66AccessToken, rnd)
}

func TestAccCloud66Stack_UidLookup(t *testing.T) {
	t.Parallel()

	rnd := generateRandomResourceName()
	uid := generateRandomUid()
	name := fmt.Sprintf("data.cloud66_stack.%s", rnd)

	testAccCloud66Stack(uid, "awesome-app")

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloud66StackUid(rnd, uid),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloud66StackDataSourceID(name),
					resource.TestCheckResourceAttr(name, "id", uid),
					resource.TestCheckResourceAttr(name, "name", "awesome-app"),
					resource.TestCheckResourceAttr(name, "environment", "production"),
					resource.TestCheckResourceAttr(name, "fqdn", "awesome-app.dev.c66.me"),
					resource.TestCheckResourceAttr(name, "language", "ruby"),
					resource.TestCheckResourceAttr(name, "framework", "rails"),
				),
			},
		},
	})
}

func testAccCloud66StackUid(rnd string, uid string) string {
	return fmt.Sprintf(`
provider "cloud66" {
  access_token = "%[1]s"
}

data "cloud66_stack" "%[2]s" {
  uid = "%[3]s"
}
`, testAccCloud66AccessToken, rnd, uid)
}

func testAccCheckCloud66StackDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		all := s.RootModule().Resources
		rs, ok := all[n]
		if !ok {
			return fmt.Errorf("can't find stacks data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Snapshot stacks source ID not set")
		}
		return nil
	}
}
