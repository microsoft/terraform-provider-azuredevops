package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccIdentityUsersDataSource(t *testing.T) {
	userName := "tfsdev install"
	tfNode := "data.azuredevops_identity_user.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclIdentityUsersDataSourceConfig(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "name", userName),
				),
			},
		},
	})
}

func hclIdentityUsersDataSourceConfig(userName string) string {
	return fmt.Sprintf(`
data "azuredevops_identity_user" "test" {
  name = "%[1]s"
}`, userName)
}
