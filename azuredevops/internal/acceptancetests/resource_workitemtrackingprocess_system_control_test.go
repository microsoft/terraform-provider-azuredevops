package acceptancetests

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func TestAccWorkitemtrackingprocessSystemControl_Basic(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_system_control.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicSystemControl(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getSystemControlImportIdFunc(tfNode),
			},
		},
	})
}

func TestAccWorkitemtrackingprocessSystemControl_Update(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_system_control.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicSystemControl(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getSystemControlImportIdFunc(tfNode),
			},
			{
				Config: updatedSystemControl(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getSystemControlImportIdFunc(tfNode),
			},
		},
	})
}

func TestAccWorkitemtrackingprocessSystemControl_Revert(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_system_control.test"

	var processId, witRefName, controlId string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicSystemControl(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrWith(tfNode, "process_id", func(value string) error {
						processId = value
						return nil
					}),
					resource.TestCheckResourceAttrWith(tfNode, "work_item_type_id", func(value string) error {
						witRefName = value
						return nil
					}),
					resource.TestCheckResourceAttrWith(tfNode, "id", func(value string) error {
						controlId = value
						return nil
					}),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getSystemControlImportIdFunc(tfNode),
			},
			{
				Config: removedSystemControl(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					testCheckSystemControlReverted(&processId, &witRefName, &controlId),
				),
			},
		},
	})
}

func testCheckSystemControlReverted(processId, witRefName, controlId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

		controls, err := clients.WorkItemTrackingProcessClient.GetSystemControls(context.Background(),
			workitemtrackingprocess.GetSystemControlsArgs{
				ProcessId:  converter.UUID(*processId),
				WitRefName: witRefName,
			})
		if err != nil {
			return fmt.Errorf("error getting system controls: %+v", err)
		}

		if controls != nil {
			for _, c := range *controls {
				if c.Id != nil && *c.Id == *controlId {
					return fmt.Errorf("system control %s should have been reverted but is still modified", *controlId)
				}
			}
		}
		return nil
	}
}

func basicSystemControl(workItemTypeName string, processName string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtrackingprocess_process" "test" {
  name                   = "%s"
  parent_process_type_id = "%s"
}

resource "azuredevops_workitemtrackingprocess_workitemtype" "test" {
  name       = "%s"
  process_id = azuredevops_workitemtrackingprocess_process.test.id
}

resource "azuredevops_workitemtrackingprocess_system_control" "test" {
  process_id        = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_id = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  control_id        = "System.AreaPath"
  visible           = false
}
`, processName, agileSystemProcessTypeId, workItemTypeName)
}

func updatedSystemControl(workItemTypeName string, processName string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtrackingprocess_process" "test" {
  name                   = "%s"
  parent_process_type_id = "%s"
}

resource "azuredevops_workitemtrackingprocess_workitemtype" "test" {
  name       = "%s"
  process_id = azuredevops_workitemtrackingprocess_process.test.id
}

resource "azuredevops_workitemtrackingprocess_system_control" "test" {
  process_id        = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_id = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  control_id        = "System.AreaPath"
  visible           = true
  label             = "Custom Area"
}
`, processName, agileSystemProcessTypeId, workItemTypeName)
}

func removedSystemControl(workItemTypeName string, processName string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtrackingprocess_process" "test" {
  name                   = "%s"
  parent_process_type_id = "%s"
}

resource "azuredevops_workitemtrackingprocess_workitemtype" "test" {
  name       = "%s"
  process_id = azuredevops_workitemtrackingprocess_process.test.id
}
`, processName, agileSystemProcessTypeId, workItemTypeName)
}

func getSystemControlImportIdFunc(tfNode string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		res := state.RootModule().Resources[tfNode]
		id := res.Primary.Attributes["id"]
		processId := res.Primary.Attributes["process_id"]
		witRefName := res.Primary.Attributes["work_item_type_id"]
		return fmt.Sprintf("%s/%s/%s", processId, witRefName, id), nil
	}
}
