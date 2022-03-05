package cloud66

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	mockHttp "github.com/moemoe89/go-helpers"
)

func TestAccCloud66Stack_NameLookup(t *testing.T) {

	t.Parallel()
	rnd := generateRandomResourceName()
	uid := generateRandomUid()
	name := fmt.Sprintf("data.cloud66_stack.%s", rnd)

	srv := mockHttp.HttpMock("/stacks.json", http.StatusOK, testAccCloud66Stacks(uid, "awesome-app"))
	defer srv.Close()

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloud66StackName(srv.URL, rnd),
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

func testAccCloud66StackName(url string, rnd string) string {
	return fmt.Sprintf(`
provider "cloud66" {
  url = "%[1]s"	
  access_token = "%[2]s"
}

data "cloud66_stack" "%[3]s" {
  name = "awesome-app"
}
`, url, testAccCloud66AccessToken, rnd)
}

func TestAccCloud66Stack_UidLookup(t *testing.T) {
	t.Parallel()

	rnd := generateRandomResourceName()
	uid := generateRandomUid()
	name := fmt.Sprintf("data.cloud66_stack.%s", rnd)

	srv := mockHttp.HttpMock("/stacks/"+uid+".json", http.StatusOK, testAccCloud66Stack(uid, "awesome-app"))
	defer srv.Close()

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloud66StackUid(srv.URL, rnd, uid),
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

func testAccCloud66StackUid(url string, rnd string, uid string) string {
	return fmt.Sprintf(`
provider "cloud66" {
  url = "%[1]s"	
  access_token = "%[2]s"
}

data "cloud66_stack" "%[3]s" {
  uid = "%[4]s"
}
`, url, testAccCloud66AccessToken, rnd, uid)
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
