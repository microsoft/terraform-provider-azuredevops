package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccGroupMembership_overwrite(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_group_membership.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: overwriteEmpty(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "members.#", "0"),
				),
			},
			{
				Config: overwriteWithMember(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "members.#", "1"),
				),
			},
			{
				Config: overwriteEmpty(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "members.#", "0"),
				),
			},
		},
	})
}

func overwriteEmpty(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "acctest-%[1]s"
}

resource "azuredevops_group" "test" {
  display_name = "acctest-%[1]s"
  scope        = azuredevops_project.test.id
}

resource "azuredevops_group_membership" "test" {
  group   = azuredevops_group.test.id
  mode    = "overwrite"
  members = []
}
`, name)
}

func overwriteWithMember(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "acctest-%[1]s"
}

resource "azuredevops_group" "test" {
  display_name = "acctest-%[1]s"
  scope        = azuredevops_project.test.id
}

resource "azuredevops_group" "member" {
  display_name = "acctest-member-%[1]s"
  scope        = azuredevops_project.test.id
}

resource "azuredevops_group_membership" "test" {
  group   = azuredevops_group.test.id
  mode    = "overwrite"
  members = [azuredevops_group.member.id]
}
`, name)
}
