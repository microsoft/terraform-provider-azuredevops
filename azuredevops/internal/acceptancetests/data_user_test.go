//go:build (all || data_sources || data_user) && (!exclude_data_sources || !exclude_data_user)

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccUser_dataSource(t *testing.T) {
	userName := "foo@email.com"
	tfNode := "data.azuredevops_user.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclDataUserBasic(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "subject_kind"),
					resource.TestCheckResourceAttrSet(tfNode, "principal_name"),
					resource.TestCheckResourceAttrSet(tfNode, "mail_address"),
					resource.TestCheckResourceAttrSet(tfNode, "origin"),
					resource.TestCheckResourceAttrSet(tfNode, "origin_id"),
					resource.TestCheckResourceAttrSet(tfNode, "display_name"),
					resource.TestCheckResourceAttrSet(tfNode, "domain"),
				),
			},
		},
	})
}

func hclDataUserBasic(uname string) string {
	return fmt.Sprintf(`
resource "azuredevops_user_entitlement" "test" {
  principal_name       = "%[1]s"
  account_license_type = "basic"
}

data "azuredevops_user" "test" {
  descriptor = azuredevops_user_entitlement.test.descriptor
  depends_on = [azuredevops_user_entitlement.test]
}`, uname)
}
