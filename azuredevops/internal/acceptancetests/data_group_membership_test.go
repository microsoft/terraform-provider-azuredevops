//go:build (all || data_sources || data_group_membership) && (!exclude_data_sources || !exclude_data_group_membership)

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccGroupMembershipData_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	groupName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_group_membership.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclDataMemberShipBasic(projectName, groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "members.#", "2"),
				),
			},
		},
	})
}

func hclDataMemberShipBasic(projectName, groupName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

data "azuredevops_group" "admin" {
  project_id = azuredevops_project.test.id
  name       = "Build Administrators"
}

data "azuredevops_group" "contributors" {
  project_id = azuredevops_project.test.id
  name       = "Contributors"
}

resource "azuredevops_group" "test" {
  scope        = azuredevops_project.test.id
  display_name = "%s"

  members = [
    data.azuredevops_group.admin.descriptor,
    data.azuredevops_group.contributors.descriptor
  ]
}
data "azuredevops_group_membership" "test" {
  group_descriptor = azuredevops_group.test.descriptor
}`, projectName, groupName)
}
