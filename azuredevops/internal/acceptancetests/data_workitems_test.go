package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccWorkItemsDataSource_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	title := testutils.GenerateResourceName()

	tfNode := "data.azuredevops_workitems.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclWorkItemsDataSourceBasic(projectName, title),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "work_items.#", "1"),
					resource.TestCheckResourceAttrPair(tfNode, "work_items.0.id", "azuredevops_workitem.test", "id"),
					resource.TestCheckResourceAttr(tfNode, "work_items.0.title", title),
					resource.TestCheckResourceAttr(tfNode, "work_items.0.type", "Issue"),
					resource.TestCheckResourceAttrSet(tfNode, "work_items.0.url"),
				),
			},
		},
	})
}

func TestAccWorkItemsDataSource_recursiveTree(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	title := testutils.GenerateResourceName()

	tfNode := "data.azuredevops_workitems.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclWorkItemsDataSourceRecursiveTree(projectName, title),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "work_items.#", "3"),
					resource.TestCheckResourceAttr(tfNode, "work_items.0.title", title+"-root"),
					resource.TestCheckResourceAttr(tfNode, "work_items.1.title", title+"-child"),
					resource.TestCheckResourceAttr(tfNode, "work_items.2.title", title+"-grandchild"),
					resource.TestCheckResourceAttrPair(tfNode, "work_items.1.parent_id", "azuredevops_workitem.root", "id"),
					resource.TestCheckResourceAttrPair(tfNode, "work_items.2.parent_id", "azuredevops_workitem.child", "id"),
				),
			},
		},
	})
}

func TestAccWorkItemsDataSource_savedQuery(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	title := testutils.GenerateResourceName()

	tfNode := "data.azuredevops_workitems.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclWorkItemsDataSourceSavedQuery(projectName, title),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "work_items.#", "1"),
					resource.TestCheckResourceAttrPair(tfNode, "work_items.0.id", "azuredevops_workitem.test", "id"),
					resource.TestCheckResourceAttr(tfNode, "work_items.0.title", title),
				),
			},
		},
	})
}

func hclWorkItemsDataSourceBasic(projectName, title string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_workitem" "test" {
  project_id = azuredevops_project.project.id
  title      = "%s"
  type       = "Issue"
}

data "azuredevops_workitems" "test" {
  project_id = azuredevops_project.project.id
  wiql       = "SELECT [System.Id] FROM WorkItems WHERE [System.Id] = ${azuredevops_workitem.test.id}"
}
`, testutils.HclProjectResource(projectName), title)
}

func hclWorkItemsDataSourceRecursiveTree(projectName, title string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_workitem" "root" {
  project_id = azuredevops_project.project.id
  title      = "%[2]s-root"
  type       = "Issue"
}

resource "azuredevops_workitem" "child" {
  project_id = azuredevops_project.project.id
  title      = "%[2]s-child"
  type       = "Issue"
  parent_id  = azuredevops_workitem.root.id
}

resource "azuredevops_workitem" "grandchild" {
  project_id = azuredevops_project.project.id
  title      = "%[2]s-grandchild"
  type       = "Issue"
  parent_id  = azuredevops_workitem.child.id
}

data "azuredevops_workitems" "test" {
  project_id = azuredevops_project.project.id
  wiql       = <<-EOT
    SELECT [System.Id]
    FROM WorkItemLinks
    WHERE [Source].[System.Id] = ${azuredevops_workitem.root.id}
      AND [System.Links.LinkType] = 'System.LinkTypes.Hierarchy-Forward'
    MODE (Recursive)
  EOT

  depends_on = [azuredevops_workitem.grandchild]
}
`, testutils.HclProjectResource(projectName), title)
}

func hclWorkItemsDataSourceSavedQuery(projectName, title string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_workitem" "test" {
  project_id = azuredevops_project.project.id
  title      = "%s"
  type       = "Issue"
}

resource "azuredevops_workitemquery" "test" {
  project_id = azuredevops_project.project.id
  area       = "Shared Queries"
  name       = "%[2]s"
  wiql       = "SELECT [System.Id] FROM WorkItems WHERE [System.Id] = ${azuredevops_workitem.test.id}"
}

data "azuredevops_workitems" "test" {
  project_id = azuredevops_project.project.id
  query_id   = azuredevops_workitemquery.test.id
}
`, testutils.HclProjectResource(projectName), title)
}
