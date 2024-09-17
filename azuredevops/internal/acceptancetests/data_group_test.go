//go:build (all || core || data_sources || data_group) && (!exclude_data_sources || !exclude_data_group)
// +build all core data_sources data_group
// +build !exclude_data_sources !exclude_data_group

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// Validates that a configuration containing a project group lookup is able to read the resource correctly.
// Because this is a data source, there are no resources to inspect in AzDO
func TestAccGroupDataSource_Read_HappyPath(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	group := "Build Administrators"
	tfBuildDefNode := "data.azuredevops_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclGroupDataBasic(projectName, group),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "name"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "descriptor"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "origin"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "origin_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "group_id"),
				),
			},
		},
	})
}

func TestAccGroupDataSource_Read_ProjectCollectionAdministrators(t *testing.T) {
	group := "Project Collection Administrators"
	tfBuildDefNode := "data.azuredevops_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclGroupDataAllGroups(group),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "name"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "descriptor"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "origin"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "origin_id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "group_id"),
				),
			},
		},
	})
}

func hclGroupDataBasic(projectName, groupName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  description        = "description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

data "azuredevops_group" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s"
}`, projectName, groupName)
}

func hclGroupDataAllGroups(groupName string) string {
	return fmt.Sprintf(`
data "azuredevops_group" "test" {
  name = "%s"
}`, groupName)
}
