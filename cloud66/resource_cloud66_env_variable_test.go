package cloud66

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCloud66EnvVariable_Deployment(t *testing.T) {
	t.Parallel()

	rnd := generateRandomResourceName()
	stackID := generateRandomUid()
	key := generateRandomEnvKey()
	value := generateRandomEnvValue()

	testAccCloud66EnvVariable(stackID, key, value)

	resourceName := "cloud66_env_variable." + rnd
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloud66EnvVariable_Deployment(stackID, rnd, key, value),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "stack_id", stackID),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "value", value),
					resource.TestCheckResourceAttr(resourceName, "readonly", "false"),
				),
			},
		},
	})
}

func testAccCloud66EnvVariable_Deployment(stactID string, rnd string, key string, value string) string {
	return fmt.Sprintf(`
provider "cloud66" {
	access_token = "%[1]s"
}

resource "cloud66_env_variable" "%[3]s" {
	stack_id = "%[2]s"
	key = "%[4]s"
	value = "%[5]s"
	apply_strategy = "deployment"
}
`, testAccCloud66AccessToken, stactID, rnd, key, value)
}

func TestAccCloud66EnvVariable_Immediately(t *testing.T) {
	t.Parallel()

	rnd := generateRandomResourceName()
	stackID := generateRandomUid()
	key := generateRandomEnvKey()
	value := generateRandomEnvValue()

	testAccCloud66EnvVariable(stackID, key, value)

	resourceName := "cloud66_env_variable." + rnd
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCloud66EnvVariable_Immediately(stackID, rnd, key, value),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "stack_id", stackID),
					resource.TestCheckResourceAttr(resourceName, "key", key),
					resource.TestCheckResourceAttr(resourceName, "value", value),
					resource.TestCheckResourceAttr(resourceName, "readonly", "false"),
				),
			},
		},
	})
}

func testAccCloud66EnvVariable_Immediately(stactID string, rnd string, key string, value string) string {
	return fmt.Sprintf(`
provider "cloud66" {
	access_token = "%[1]s"
}

resource "cloud66_env_variable" "%[3]s" {
	stack_id = "%[2]s"
	key = "%[4]s"
	value = "%[5]s"
	apply_strategy = "immediately"
}
`, testAccCloud66AccessToken, stactID, rnd, key, value)
}
