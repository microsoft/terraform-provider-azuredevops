package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccWorkItemDataSource_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	title := testutils.GenerateResourceName()

	tfNode := "data.azuredevops_workitem.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclWorkItemDataSourceBasic(projectName, title),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "url"),
					resource.TestCheckResourceAttrSet(tfNode, "area_path"),
					resource.TestCheckResourceAttrSet(tfNode, "iteration_path"),
					resource.TestCheckResourceAttr(tfNode, "title", title),
					resource.TestCheckResourceAttr(tfNode, "type", "Issue"),
					resource.TestCheckResourceAttr(tfNode, "state", "Active"),
					resource.TestCheckResourceAttr(tfNode, "parent_id", "0"),
					resource.TestCheckResourceAttr(tfNode, "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr(tfNode, "tags.*", "tag1"),
					resource.TestCheckTypeSetElemAttr(tfNode, "tags.*", "tag2"),
					resource.TestCheckResourceAttr(tfNode, "fields.System.Title", title),
					resource.TestCheckResourceAttr(tfNode, "fields.System.WorkItemType", "Issue"),
				),
			},
		},
	})
}

func TestAccWorkItemDataSource_withParent(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	title := testutils.GenerateResourceName()

	tfNode := "data.azuredevops_workitem.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclWorkItemDataSourceWithParent(projectName, title),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "title", title+"-child"),
					resource.TestCheckResourceAttrPair(tfNode, "parent_id", "azuredevops_workitem.parent", "id"),
					resource.TestCheckResourceAttr(tfNode, "relations.#", "1"),
					resource.TestCheckResourceAttr(tfNode, "relations.0.rel", "System.LinkTypes.Hierarchy-Reverse"),
					resource.TestCheckResourceAttrSet(tfNode, "relations.0.url"),
				),
			},
		},
	})
}

func hclWorkItemDataSourceBasic(projectName, title string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_workitem" "test" {
  project_id = azuredevops_project.project.id
  title      = "%s"
  type       = "Issue"
  state      = "Active"
  tags       = ["tag1", "tag2"]
}

data "azuredevops_workitem" "test" {
  id         = azuredevops_workitem.test.id
  project_id = azuredevops_project.project.id
}
`, testutils.HclProjectResource(projectName), title)
}

func hclWorkItemDataSourceWithParent(projectName, title string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_workitem" "parent" {
  project_id = azuredevops_project.project.id
  title      = "%[2]s-parent"
  type       = "Issue"
}

resource "azuredevops_workitem" "test" {
  project_id = azuredevops_project.project.id
  title      = "%[2]s-child"
  type       = "Issue"
  parent_id  = azuredevops_workitem.parent.id
}

data "azuredevops_workitem" "test" {
  id         = azuredevops_workitem.test.id
  project_id = azuredevops_project.project.id
}
`, testutils.HclProjectResource(projectName), title)
}
