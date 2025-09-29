//go:build (all || workitemtracking || resource_workitemquery_folder) && (!exclude_workitemtracking || !resource_workitemquery_folder)

package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// Folder under an area (Shared Queries)
func TestAccWorkItemQueryFolder_UnderArea(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	folderName := "tfacc-folder-area"

	config := testutils.HclProjectResource(projectName) + `
resource "azuredevops_workitemquery_folder" "folder" {
  project_id = azuredevops_project.project.id
  name       = "` + folderName + `"
  area       = "My Queries"
}
`

	res := "azuredevops_workitemquery_folder.folder"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{{
			Config: config,
			Check: resource.ComposeTestCheckFunc(
				testutils.CheckProjectExists(projectName),
				resource.TestCheckResourceAttr(res, "name", folderName),
				resource.TestCheckResourceAttrSet(res, "project_id"),
			),
		}},
	})
}

// Folder under a folder (child folder has parent_id referencing parent folder)
func TestAccWorkItemQueryFolder_UnderFolder(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	parentFolderName := "tfacc-folder-parent"
	childFolderName := "tfacc-folder-child"

	config := testutils.HclProjectResource(projectName) + `
resource "azuredevops_workitemquery_folder" "parent" {
  project_id = azuredevops_project.project.id
  name       = "` + parentFolderName + `"
  area       = "My Queries"
}

resource "azuredevops_workitemquery_folder" "child" {
  project_id = azuredevops_project.project.id
  name       = "` + childFolderName + `"
  parent_id  = azuredevops_workitemquery_folder.parent.id
}
`

	parentRes := "azuredevops_workitemquery_folder.parent"
	childRes := "azuredevops_workitemquery_folder.child"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{{
			Config: config,
			Check: resource.ComposeTestCheckFunc(
				testutils.CheckProjectExists(projectName),
				resource.TestCheckResourceAttr(parentRes, "name", parentFolderName),
				resource.TestCheckResourceAttr(childRes, "name", childFolderName),
				resource.TestCheckResourceAttrSet(childRes, "parent_id"),
			),
		}},
	})
}
