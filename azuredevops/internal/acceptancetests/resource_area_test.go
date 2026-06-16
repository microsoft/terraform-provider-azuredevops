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
					resource.TestCheckResourceAttrSet(tfNode, "area_id"),
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
					resource.TestCheckResourceAttrSet(tfNodeParent, "area_id"),
					resource.TestCheckResourceAttrSet(tfNodeChild, "id"),
					resource.TestCheckResourceAttr(tfNodeChild, "name", childAreaName),
					resource.TestCheckResourceAttrPair(tfNodeChild, "parent_area_id", tfNodeParent, "area_id"),
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
  project_id = azuredevops_project.project.id
  name       = "%s"
  parent_area_id  = azuredevops_area.parent.area_id
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
		nodeID := rs.Primary.Attributes["area_id"]

		return fmt.Sprintf("%s/%s", projectID, nodeID), nil
	}
}
