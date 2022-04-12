package cloud66

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCloud66EnvVariableDataSource(t *testing.T) {

	t.Parallel()

	rnd := generateRandomResourceName()
	stackID := generateRandomUid()
	key := generateRandomEnvKey()
	value := generateRandomEnvValue()

	name := fmt.Sprintf("data.cloud66_env_variable.%s", rnd)

	testAccCloud66EnvVariable(stackID, key, value)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloud66EnvVariableDataSource(stackID, rnd, key),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloud66EnvVariableDataSourceID(name, key),
					resource.TestCheckResourceAttr(name, "stack_id", stackID),
					resource.TestCheckResourceAttr(name, "key", key),
					resource.TestCheckResourceAttr(name, "value", value),
				),
			},
		},
	})
}

func testAccCloud66EnvVariableDataSource(stactID string, rnd string, key string) string {
	return fmt.Sprintf(`
provider "cloud66" {
  access_token = "%[1]s"
}

data "cloud66_env_variable" "%[3]s" {
	stack_id = "%[2]s"
	key = "%[4]s"
}
`, testAccCloud66AccessToken, stactID, rnd, key)
}

func testAccCheckCloud66EnvVariableDataSourceID(n string, key string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		all := s.RootModule().Resources
		rs, ok := all[n]
		if !ok {
			return fmt.Errorf("can't find env variable data source: %s", n)
		}

		if rs.Primary.ID != key {
			return fmt.Errorf("Snapshot stacks source ID not set")
		}
		return nil
	}
}
