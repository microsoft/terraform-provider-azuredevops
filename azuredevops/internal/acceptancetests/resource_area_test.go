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
					resource.TestCheckResourceAttr(tfNode, "path", fmt.Sprintf("%s\\%s", projectName, areaName)),
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
					resource.TestCheckResourceAttr(tfNodeParent, "path", fmt.Sprintf("%s\\%s", projectName, parentAreaName)),
					resource.TestCheckResourceAttrSet(tfNodeChild, "id"),
					resource.TestCheckResourceAttr(tfNodeChild, "path", fmt.Sprintf("%s\\%s\\%s", projectName, parentAreaName, childAreaName)),
					resource.TestCheckResourceAttrPair(tfNodeChild, "parent_area_id", tfNodeParent, "id"),
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
					resource.TestCheckResourceAttr(tfNode, "path", fmt.Sprintf("%s\\%s", projectName, areaName)),
					resource.TestCheckResourceAttr(tfNode, "name", areaName),
				),
			},
			{
				Config: hclAreaBasic(projectName, updatedAreaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "path", fmt.Sprintf("%s\\%s", projectName, updatedAreaName)),
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
}
`, testutils.HclProjectResource(projectName), areaName)
}

func hclAreaChild(projectName, parentAreaName, childAreaName string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_area" "parent" {
  project_id = azuredevops_project.project.id
  name       = "%s"
}

resource "azuredevops_area" "child" {
  project_id     = azuredevops_project.project.id
  name           = "%s"
  parent_area_id = azuredevops_area.parent.id
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
		nodeID := rs.Primary.Attributes["id"]

		return fmt.Sprintf("%s/%s", projectID, nodeID), nil
	}
}
