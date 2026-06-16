package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccArea_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	areaName := "TestArea"
	tfNode := "azuredevops_area.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclAreaBasic(projectName, areaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "name", areaName),
					resource.TestCheckResourceAttr(tfNode, "path", "/"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateIdFunc: testAccAreaImportStateIdFunc(tfNode),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccArea_child(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	parentAreaName := "ParentArea"
	childAreaName := "ChildArea"
	tfNodeParent := "azuredevops_area.parent"
	tfNodeChild := "azuredevops_area.child"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclAreaChild(projectName, parentAreaName, childAreaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNodeParent, "id"),
					resource.TestCheckResourceAttr(tfNodeParent, "name", parentAreaName),
					resource.TestCheckResourceAttr(tfNodeParent, "path", "/"),
					resource.TestCheckResourceAttrSet(tfNodeChild, "id"),
					resource.TestCheckResourceAttr(tfNodeChild, "name", childAreaName),
					resource.TestCheckResourceAttr(tfNodeChild, "path", "/"+parentAreaName),
				),
			},
			{
				ResourceName:      tfNodeParent,
				ImportState:       true,
				ImportStateIdFunc: testAccAreaImportStateIdFunc(tfNodeParent),
				ImportStateVerify: true,
			},
			{
				ResourceName:      tfNodeChild,
				ImportState:       true,
				ImportStateIdFunc: testAccAreaImportStateIdFunc(tfNodeChild),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccArea_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	areaName := "OriginalArea"
	updatedAreaName := "UpdatedArea"
	tfNode := "azuredevops_area.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclAreaBasic(projectName, areaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "name", areaName),
				),
			},
			{
				Config: hclAreaBasic(projectName, updatedAreaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "name", updatedAreaName),
				),
			},
		},
	})
}

func hclAreaBasic(projectName, areaName string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_area" "test" {
  project_id = azuredevops_project.project.id
  name       = "%s"
  path       = "/"
}
`, testutils.HclProjectResource(projectName), areaName)
}

func hclAreaChild(projectName, parentAreaName, childAreaName string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_area" "parent" {
  project_id = azuredevops_project.project.id
  name       = "%s"
  path       = "/"
}

resource "azuredevops_area" "child" {
  project_id = azuredevops_project.project.id
  name       = "%s"
  path       = "/${azuredevops_area.parent.name}"

  depends_on = [azuredevops_area.parent]
}
`, testutils.HclProjectResource(projectName), parentAreaName, childAreaName)
}

func testAccAreaImportStateIdFunc(resourceName string) func(s *terraform.State) (string, error) {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		projectID := rs.Primary.Attributes["project_id"]
		name := rs.Primary.Attributes["name"]
		path := rs.Primary.Attributes["path"]

		var importID string
		if path == "/" || path == "" {
			importID = fmt.Sprintf("%s/%s", projectID, name)
		} else {
			importID = fmt.Sprintf("%s%s/%s", projectID, path, name)
		}
		return importID, nil
	}
}
