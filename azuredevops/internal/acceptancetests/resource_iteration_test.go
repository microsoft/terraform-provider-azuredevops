package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccIteration_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	iterationName := "TestIteration"
	tfNode := "azuredevops_iteration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclIterationBasic(projectName, iterationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "path", fmt.Sprintf("%s\\%s", projectName, iterationName)),
					resource.TestCheckResourceAttr(tfNode, "start_date", ""),
					resource.TestCheckResourceAttr(tfNode, "finish_date", ""),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateIdFunc: testAccIterationImportStateIdFunc(tfNode),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIteration_child(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	parentIterationName := "ParentIteration"
	childIterationName := "ChildIteration"
	tfNodeParent := "azuredevops_iteration.parent"
	tfNodeChild := "azuredevops_iteration.child"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclIterationChild(projectName, parentIterationName, childIterationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNodeParent, "id"),
					resource.TestCheckResourceAttr(tfNodeParent, "path", fmt.Sprintf("%s\\%s", projectName, parentIterationName)),
					resource.TestCheckResourceAttrSet(tfNodeChild, "id"),
					resource.TestCheckResourceAttr(tfNodeChild, "path", fmt.Sprintf("%s\\%s\\%s", projectName, parentIterationName, childIterationName)),
					resource.TestCheckResourceAttrPair(tfNodeChild, "parent_iteration_id", tfNodeParent, "id"),
				),
			},
			{
				ResourceName:      tfNodeParent,
				ImportState:       true,
				ImportStateIdFunc: testAccIterationImportStateIdFunc(tfNodeParent),
				ImportStateVerify: true,
			},
			{
				ResourceName:      tfNodeChild,
				ImportState:       true,
				ImportStateIdFunc: testAccIterationImportStateIdFunc(tfNodeChild),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccIteration_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	iterationName := "OriginalIteration"
	updatedIterationName := "UpdatedIteration"
	tfNode := "azuredevops_iteration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclIterationBasic(projectName, iterationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "path", fmt.Sprintf("%s\\%s", projectName, iterationName)),
					resource.TestCheckResourceAttr(tfNode, "name", iterationName),
				),
			},
			{
				Config: hclIterationWithDates(projectName, updatedIterationName, "2026-01-05", "2026-01-16"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "path", fmt.Sprintf("%s\\%s", projectName, updatedIterationName)),
					resource.TestCheckResourceAttr(tfNode, "name", updatedIterationName),
					resource.TestCheckResourceAttr(tfNode, "start_date", "2026-01-05"),
					resource.TestCheckResourceAttr(tfNode, "finish_date", "2026-01-16"),
				),
			},
			{
				Config: hclIterationWithDates(projectName, updatedIterationName, "2026-01-19", "2026-01-30"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "start_date", "2026-01-19"),
					resource.TestCheckResourceAttr(tfNode, "finish_date", "2026-01-30"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateIdFunc: testAccIterationImportStateIdFunc(tfNode),
				ImportStateVerify: true,
			},
			{
				Config: hclIterationBasic(projectName, updatedIterationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "start_date", ""),
					resource.TestCheckResourceAttr(tfNode, "finish_date", ""),
				),
			},
		},
	})
}

func hclIterationBasic(projectName, iterationName string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_iteration" "test" {
  project_id = azuredevops_project.project.id
  name       = "%s"
}
`, testutils.HclProjectResource(projectName), iterationName)
}

func hclIterationWithDates(projectName, iterationName, startDate, finishDate string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_iteration" "test" {
  project_id  = azuredevops_project.project.id
  name        = "%s"
  start_date  = "%s"
  finish_date = "%s"
}
`, testutils.HclProjectResource(projectName), iterationName, startDate, finishDate)
}

func hclIterationChild(projectName, parentIterationName, childIterationName string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_iteration" "parent" {
  project_id = azuredevops_project.project.id
  name       = "%s"
}

resource "azuredevops_iteration" "child" {
  project_id          = azuredevops_project.project.id
  name                = "%s"
  parent_iteration_id = azuredevops_iteration.parent.id
}
`, testutils.HclProjectResource(projectName), parentIterationName, childIterationName)
}

func testAccIterationImportStateIdFunc(resourceName string) func(s *terraform.State) (string, error) {
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
