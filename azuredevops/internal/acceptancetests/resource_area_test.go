package acceptancetests

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func TestAccArea_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	areaName := testutils.GenerateResourceName()
	tfNode := "azuredevops_area.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkAreaDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclAreaBasic(projectName, areaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", areaName),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttr(tfNode, "path", "\\"),
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

func TestAccArea_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	areaName := testutils.GenerateResourceName()
	areaNameUpdated := areaName + "_updated"
	tfNode := "azuredevops_area.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkAreaDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclAreaBasic(projectName, areaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", areaName),
				),
			},
			{
				Config: hclAreaBasic(projectName, areaNameUpdated), // Update name
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", areaNameUpdated),
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

func checkAreaDestroyed(s *terraform.State) error {
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_area" {
			continue
		}

		id, err := strconv.Atoi(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("Area ID=%d cannot be parsed!. Error=%v", id, err)
		}

		projectID := resource.Primary.Attributes["project_id"]
		path := resource.Primary.Attributes["path"]
		name := resource.Primary.Attributes["name"]

		fullPath := path
		if fullPath == "" || fullPath == "\\" {
			fullPath = "\\" + name
		} else {
			fullPath = fullPath + "\\" + name
		}

		structureGroup := workitemtracking.TreeStructureGroupValues.Areas

		_, err = clients.WorkItemTrackingClient.GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
			Project:        &projectID,
			StructureGroup: &structureGroup,
			Path:           converter.String(fullPath),
			Depth:          converter.Int(0),
		})

		if err == nil {
			return fmt.Errorf("Area '%s' (ID: %d) still exists", fullPath, id)
		}
	}

	return nil
}

func hclAreaBasic(projectName, areaName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%s"
  work_item_template = "Agile"
}

resource "azuredevops_area" "test" {
  project_id = azuredevops_project.test.id
  name       = "%s"
  path       = "\\"
}`, projectName, areaName)
}
