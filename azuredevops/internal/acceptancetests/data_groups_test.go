package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// Validates that a configuration containing a project group lookup is able to read the resource correctly.
// Because this is a data source, there are no resources to inspect in AzDO
func TestAccGroupsDataSource_Read_Project(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_groups.groups"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclGroupsDataSourceBasic(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "groups.#"),
				),
			},
		},
	})
}

func TestAccGroupsDataSource_Read_NoProject(t *testing.T) {
	tfNode := "data.azuredevops_groups.groups"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclGroupsDataSourceAllGroups(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "groups.#"),
				),
			},
		},
	})
}

func hclGroupsDataSourceBasic(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
  name               = "%[1]s"
  description        = "description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

data "azuredevops_groups" "groups" {
  project_id = azuredevops_project.project.id
}
`, projectName)
}

func hclGroupsDataSourceAllGroups() string {
	return fmt.Sprintf(`data "azuredevops_groups" "groups" {}`)
}
