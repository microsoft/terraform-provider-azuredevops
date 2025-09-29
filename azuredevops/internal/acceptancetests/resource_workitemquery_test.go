//go:build (all || workitemtracking || resource_workitemquery) && (!exclude_workitemtracking || !resource_workitemquery)

package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// Query directly under area (Shared Queries)
func TestAccWorkItemQuery_UnderArea(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	queryName := "tfacc-query-area"
	wiql := "SELECT [System.Id] FROM WorkItems WHERE [System.TeamProject] = @project"

	config := testutils.HclProjectResource(projectName) + `
resource "azuredevops_workitemquery" "query" {
  project_id = azuredevops_project.project.id
  name       = "` + queryName + `"
  area       = "My Queries"
  wiql       = "` + wiql + `"
}
`

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

	config := testutils.HclProjectResource(projectName) + `
resource "azuredevops_workitemquery_folder" "folder" {
  project_id = azuredevops_project.project.id
  name       = "` + folderName + `"
  area       = "My Queries"
}

resource "azuredevops_workitemquery" "query" {
  project_id = azuredevops_project.project.id
  name       = "` + queryName + `"
  parent_id  = azuredevops_workitemquery_folder.folder.id
  wiql       = "` + wiql + `"
}
`

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

	configCreate := testutils.HclProjectResource(projectName) + `
resource "azuredevops_workitemquery" "query" {
  project_id = azuredevops_project.project.id
  name       = "` + initialName + `"
  area       = "My Queries"
  wiql       = "` + initialWiql + `"
}
`

	configUpdate := testutils.HclProjectResource(projectName) + `
resource "azuredevops_workitemquery" "query" {
  project_id = azuredevops_project.project.id
  name       = "` + updatedName + `"
  area       = "My Queries"
  wiql       = "` + updatedWiql + `"
}
`

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