package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccUsers_DataSource(t *testing.T) {
	userName := "foo@email.com"
	tfNode := "data.azuredevops_users.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: dataUser_basic(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "users.#", "1"),
					resource.TestCheckResourceAttr(tfNode, "principal_name", "foo@email.com"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func dataUser_basic(uname string) string {
	return fmt.Sprintf(
		`
resource "azuredevops_user_entitlement" "test" {
  principal_name       = "%[1]s"
  account_license_type = "basic"
}

data "azuredevops_users" "test" {
  principal_name = "%[1]s"
  depends_on = [azuredevops_user_entitlement.test]
}`, uname)
}
