package acceptancetests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccBuildFolder_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()

	tfNode := "azuredevops_build_folder.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: buildFolderBasic(projectName, "\\\\test folder", "Acceptance Test Folder"),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "path"),
					resource.TestCheckResourceAttrSet(tfNode, "description"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateIdFunc: importBuildFolderID("azuredevops_build_folder.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBuildFolder_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()

	tfNode := "azuredevops_build_folder.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: buildFolderBasic(projectName, "\\\\test folder", "Acceptance Test Folder"),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttr(tfNode, "path", `\test folder`),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateIdFunc: importBuildFolderID("azuredevops_build_folder.test"),
				ImportStateVerify: true,
			},
			{
				Config: buildFolderBasic(projectName, "\\\\test folderupdate", "Acceptance Test Folder"),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttr(tfNode, "path", `\test folderupdate`),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateIdFunc: importBuildFolderID("azuredevops_build_folder.test"),
				ImportStateVerify: true,
			},
			{
				Config: buildFolderBasic(projectName, "\\\\test folder", "Acceptance Test Folder"),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
					resource.TestCheckResourceAttr(tfNode, "path", `\test folder`),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateIdFunc: importBuildFolderID("azuredevops_build_folder.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBuildFolder_requiresImportErrorStep(t *testing.T) {
	projectName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: buildFolderBasic(projectName, "\\\\test folder", "Acceptance Test Folder"),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckProjectExists(projectName),
				),
			},
			{
				ResourceName:      "azuredevops_build_folder.test",
				ImportState:       true,
				ImportStateIdFunc: importBuildFolderID("azuredevops_build_folder.test"),
				ImportStateVerify: true,
			},
			{
				Config:      buildFolderRequiresImport(projectName, "\\\\test folder", "Acceptance Test Folder"),
				ExpectError: buildFolderRequiresImportError(`\\test folder\\`, projectName),
			},
		},
	})
}

func buildFolderRequiresImportError(resourceName, projectName string) *regexp.Regexp {
	message := `failed creating resource Build Folder, Folder %[1]s already exists for project %[2]s.`
	return regexp.MustCompile(fmt.Sprintf(message, resourceName, projectName))
}

func importBuildFolderID(resName string) func(state *terraform.State) (string, error) {
	return func(state *terraform.State) (string, error) {
		res := state.RootModule().Resources[resName]
		projectID := res.Primary.Attributes["project_id"]
		path := res.Primary.Attributes["path"]
		return fmt.Sprintf("%s/%s", projectID, path), nil
	}
}

func buildFolderBasic(projectName, path, description string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
  name               = "%[1]s"
  description        = "%[1]s-description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_build_folder" "test" {
  project_id  = azuredevops_project.project.id
  path        = "%[2]s"
  description = "%[3]s"
}
`, projectName, path, description)
}

func buildFolderRequiresImport(projectName, path, description string) string {
	basicConfig := buildFolderBasic(projectName, path, description)
	return fmt.Sprintf(`


%s

resource "azuredevops_build_folder" "import" {
  project_id  = azuredevops_build_folder.test.project_id
  path        = azuredevops_build_folder.test.path
  description = azuredevops_build_folder.test.description
}
`, basicConfig)
}
