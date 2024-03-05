package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccIdentityUsers_DataSource(t *testing.T) {
	userName := "dummy_user"
	tfNode := "data.azuredevops_identity_user.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: datadentityUsers_basic(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "name", "dummy_user"),
				),
			},
		},
	})
}

func datadentityUsers_basic(uname string) string {
	return fmt.Sprintf(
		`
data "azuredevops_identity_user" "test" {
  name       = "%[1]s"
}`, uname)
}
