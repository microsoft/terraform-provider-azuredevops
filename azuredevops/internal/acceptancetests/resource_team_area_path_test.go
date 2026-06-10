package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccTeamAreaPath_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	teamName := testutils.GenerateResourceName()

	tfNode := "azuredevops_team_area_path.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclTeamAreaPathBasic(projectName, teamName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "team_id"),
					resource.TestCheckResourceAttr(tfNode, "include_children", "true"),
				),
			},
		},
	})
}

func TestAccTeamAreaPath_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	teamName := testutils.GenerateResourceName()

	tfNode := "azuredevops_team_area_path.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclTeamAreaPathWithIncludeChildren(projectName, teamName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "include_children", "true"),
				),
			},
			{
				Config: hclTeamAreaPathWithIncludeChildren(projectName, teamName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "include_children", "false"),
				),
			},
		},
	})
}

func TestAccTeamAreaPath_import(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	teamName := testutils.GenerateResourceName()

	tfNode := "azuredevops_team_area_path.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclTeamAreaPathBasic(projectName, teamName),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func hclTeamAreaPathBasic(projectName, teamName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_team" "test" {
  project_id = azuredevops_project.test.id
  name       = "%s"
}

resource "azuredevops_team_area_path" "test" {
  project_id       = azuredevops_project.test.id
  team_id          = azuredevops_team.test.id
  area_path        = azuredevops_project.test.name
  include_children = true
}
`, projectName, teamName)
}

func hclTeamAreaPathWithIncludeChildren(projectName, teamName string, includeChildren bool) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_team" "test" {
  project_id = azuredevops_project.test.id
  name       = "%s"
}

resource "azuredevops_team_area_path" "test" {
  project_id       = azuredevops_project.test.id
  team_id          = azuredevops_team.test.id
  area_path        = azuredevops_project.test.name
  include_children = %t
}
`, projectName, teamName, includeChildren)
}
