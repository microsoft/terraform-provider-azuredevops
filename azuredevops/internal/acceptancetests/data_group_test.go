package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

func TestAccGroupDataSource_ReadersResolvesWithProjectID(t *testing.T) {
	projectName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclGroupDataReadersConfig(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("azuredevops_group.collection_readers", "descriptor"),
					resource.TestCheckResourceAttrSet("data.azuredevops_group.project_readers", "descriptor"),
					testAccCheckCollectionGroupNotInProjectGroup(
						"data.azuredevops_group.project_readers",
						"azuredevops_group.collection_readers",
					),
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

func hclGroupDataReadersConfig(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_group" "collection_readers" {
  display_name = "Readers"
}

resource "azuredevops_project" "test" {
  name       = "%s"
  depends_on = [azuredevops_group.collection_readers]
}

data "azuredevops_group" "project_admins" {
  name       = "Project Administrators"
  project_id = azuredevops_project.test.id
}

resource "azuredevops_group_membership" "make_collection_visible" {
  mode  = "add"
  group = data.azuredevops_group.project_admins.descriptor
  members = [
    azuredevops_group.collection_readers.descriptor
  ]
}

data "azuredevops_group" "project_readers" {
  name       = "Readers"
  project_id = azuredevops_project.test.id
  depends_on = [azuredevops_group_membership.make_collection_visible]
}
`, projectName)
}

func testAccCheckCollectionGroupNotInProjectGroup(projectGroupNode string, collectionGroupNode string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		pg, ok := s.RootModule().Resources[projectGroupNode]
		if !ok {
			return fmt.Errorf("not found: %s", projectGroupNode)
		}

		cg, ok := s.RootModule().Resources[collectionGroupNode]
		if !ok {
			return fmt.Errorf("not found: %s", collectionGroupNode)
		}

		collectionDesc := cg.Primary.Attributes["descriptor"]
		if collectionDesc == "" {
			return fmt.Errorf("collection descriptor missing")
		}

		projectDesc := pg.Primary.Attributes["descriptor"]
		if projectDesc == "" {
			return fmt.Errorf("project descriptor missing")
		}

		if collectionDesc == projectDesc {
			return fmt.Errorf(
				"collection-level group %q present in project-scoped group list",
				collectionDesc,
			)
		}

		return nil
	}
}
