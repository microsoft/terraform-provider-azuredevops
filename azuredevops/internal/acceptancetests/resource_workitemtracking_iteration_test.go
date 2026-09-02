package acceptancetests

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func TestAccWorkItemTrackingIteration_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	iterationName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtracking_iteration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkIterationDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclIterationBasic(projectName, iterationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", iterationName),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttr(tfNode, "path", "\\"), // Root path
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkItemTrackingIteration_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	iterationName := testutils.GenerateResourceName()
	iterationNameUpdated := iterationName + "_updated"
	tfNode := "azuredevops_workitemtracking_iteration.test"

	startDate := "2023-01-01T00:00:00Z"
	finishDate := "2023-01-31T00:00:00Z"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkIterationDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclIterationBasic(projectName, iterationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", iterationName),
				),
			},
			{
				Config: hclIterationUpdate(projectName, iterationNameUpdated, startDate, finishDate),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", iterationNameUpdated),
					resource.TestCheckResourceAttr(tfNode, "attributes.0.start_date", startDate),
					resource.TestCheckResourceAttr(tfNode, "attributes.0.finish_date", finishDate),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func checkIterationDestroyed(s *terraform.State) error {
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_workitemtracking_iteration" {
			continue
		}

		id, err := strconv.Atoi(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("Iteration ID=%d cannot be parsed!. Error=%v", id, err)
		}

		projectID := resource.Primary.Attributes["project_id"]
		path := resource.Primary.Attributes["path"]
		if path == "" {
			path = "\\" + resource.Primary.Attributes["name"]
		} else {
			path = path + "\\" + resource.Primary.Attributes["name"]
		}

		structureGroup := workitemtracking.TreeStructureGroupValues.Iterations

		_, err = clients.WorkItemTrackingClient.GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
			Project:        &projectID,
			StructureGroup: &structureGroup,
			Path:           converter.String(path),
			Depth:          converter.Int(0),
		})

		if err == nil {
			return fmt.Errorf("Iteration '%s' (ID: %d) still exists", path, id)
		}
	}

	return nil
}

func hclIterationBasic(projectName, iterationName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%s"
  work_item_template = "Agile"
}

resource "azuredevops_workitemtracking_iteration" "test" {
  project_id = azuredevops_project.test.id
  name       = "%s"
  path       = "\\"
}`, projectName, iterationName)
}

func hclIterationUpdate(projectName, iterationName, start, finish string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%s"
  work_item_template = "Agile"
}

resource "azuredevops_workitemtracking_iteration" "test" {
  project_id = azuredevops_project.test.id
  name       = "%s"
  path       = "\\"

  attributes {
    start_date  = "%s"
    finish_date = "%s"
  }
}`, projectName, iterationName, start, finish)
}

func TestAccWorkItemTrackingIteration_child(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	iterationName := testutils.GenerateResourceName()
	childIterationName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtracking_iteration.child"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkIterationDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclIterationChild(projectName, iterationName, childIterationName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", childIterationName),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttr(tfNode, "path", "\\\\"+iterationName),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func hclIterationChild(projectName, iterationName, childIterationName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%s"
  work_item_template = "Agile"
}

resource "azuredevops_workitemtracking_iteration" "test" {
  project_id = azuredevops_project.test.id
  name       = "%s"
  path       = "\\"
}

resource "azuredevops_workitemtracking_iteration" "child" {
  project_id = azuredevops_project.test.id
  name       = "%s"
  path       = "\\\\${azuredevops_workitemtracking_iteration.test.name}"
}
`, projectName, iterationName, childIterationName)
}
