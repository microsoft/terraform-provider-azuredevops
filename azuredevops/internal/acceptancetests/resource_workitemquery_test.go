package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// Query directly under area (Shared Queries)
func TestAccWorkItemQuery_UnderArea(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	queryName := "tfacc-query-area"
	wiql := "SELECT [System.Id] FROM WorkItems WHERE [System.TeamProject] = @project"

	config := hclWorkItemQueryResource(projectName, queryName, wiql)

	res := "azuredevops_workitemquery.query"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{{
			Config: config,
			Check: resource.ComposeTestCheckFunc(
				testutils.CheckProjectExists(projectName),
				resource.TestCheckResourceAttr(res, "name", queryName),
				resource.TestCheckResourceAttr(res, "wiql", wiql),
			),
		}},
	})
}

// Query under a folder
func TestAccWorkItemQuery_UnderFolder(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	folderName := "tfacc-folder-for-query"
	queryName := "tfacc-query-under-folder"
	wiql := "SELECT [System.Id] FROM WorkItems WHERE [System.TeamProject] = @project"

	config := hclWorkItemQueryUnderFolderResource(projectName, folderName, queryName, wiql)

	folderRes := "azuredevops_workitemquery_folder.folder"
	queryRes := "azuredevops_workitemquery.query"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{{
			Config: config,
			Check: resource.ComposeTestCheckFunc(
				testutils.CheckProjectExists(projectName),
				resource.TestCheckResourceAttr(folderRes, "name", folderName),
				resource.TestCheckResourceAttr(queryRes, "name", queryName),
				resource.TestCheckResourceAttr(queryRes, "wiql", wiql),
			),
		}},
	})
}

// Update existing query (name + WIQL)
func TestAccWorkItemQuery_Update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	initialName := "tfacc-query-update-initial"
	updatedName := "tfacc-query-update-updated"
	initialWiql := "SELECT [System.Id] FROM WorkItems WHERE [System.TeamProject] = @project"
	updatedWiql := "SELECT [System.Id] FROM WorkItems WHERE [System.TeamProject] = @project ORDER BY [System.Id] DESC"

	configCreate := hclWorkItemQueryResource(projectName, initialName, initialWiql)
	configUpdate := hclWorkItemQueryResource(projectName, updatedName, updatedWiql)

	res := "azuredevops_workitemquery.query"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: configCreate,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttr(res, "name", initialName),
					resource.TestCheckResourceAttr(res, "wiql", initialWiql),
				),
			},
			{
				Config: configUpdate,
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttr(res, "name", updatedName),
					resource.TestCheckResourceAttr(res, "wiql", updatedWiql),
				),
			},
		},
	})
}

// Helpers

func hclWorkItemQueryResource(projectName, queryName, wiql string) string {
	return fmt.Sprintf(`%s
resource "azuredevops_workitemquery" "query" {
  project_id = azuredevops_project.project.id
  name       = %q
  area       = "My Queries"
  wiql       = %q
}
`, testutils.HclProjectResource(projectName), queryName, wiql)
}

func hclWorkItemQueryUnderFolderResource(projectName, folderName, queryName, wiql string) string {
	return fmt.Sprintf(`%s
resource "azuredevops_workitemquery_folder" "folder" {
  project_id = azuredevops_project.project.id
  name       = %q
  area       = "My Queries"
}

resource "azuredevops_workitemquery" "query" {
  project_id = azuredevops_project.project.id
  name       = %q
  parent_id  = azuredevops_workitemquery_folder.folder.id
  wiql       = %q
}
`, testutils.HclProjectResource(projectName), folderName, queryName, wiql)
}
