package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccWorkitemtrackingprocessState_Basic(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_state.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicState(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrSet(tfNode, "url"),
					resource.TestCheckResourceAttrSet(tfNode, "order"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getStateImportIdFunc(tfNode),
			},
		},
	})
}

func TestAccWorkitemtrackingprocessState_Update(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_state.test"

	var stateId string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicState(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					captureStateId(tfNode, &stateId),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getStateImportIdFunc(tfNode),
			},
			{
				Config: updatedState(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(tfNode, "id", &stateId),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getStateImportIdFunc(tfNode),
			},
			{
				Config: basicState(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(tfNode, "id", &stateId),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getStateImportIdFunc(tfNode),
			},
		},
	})
}

func captureStateId(tfNode string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		res := s.RootModule().Resources[tfNode]
		*id = res.Primary.Attributes["id"]
		return nil
	}
}

func basicState(workItemTypeName string, processName string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtrackingprocess_process" "test" {
  name                   = "%s"
  parent_process_type_id = "%s"
}

resource "azuredevops_workitemtrackingprocess_workitemtype" "test" {
  name       = "%s"
  process_id = azuredevops_workitemtrackingprocess_process.test.id
}

resource "azuredevops_workitemtrackingprocess_state" "test" {
  process_id        = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_id = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  name              = "Ready"
  color             = "#b2b2b2"
  state_category    = "Proposed"
  order             = 2
}
`, processName, agileSystemProcessTypeId, workItemTypeName)
}

func updatedState(workItemTypeName string, processName string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtrackingprocess_process" "test" {
  name                   = "%s"
  parent_process_type_id = "%s"
}

resource "azuredevops_workitemtrackingprocess_workitemtype" "test" {
  name       = "%s"
  process_id = azuredevops_workitemtrackingprocess_process.test.id
}

resource "azuredevops_workitemtrackingprocess_state" "test" {
  process_id        = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_id = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  name              = "Ready"
  color             = "#5688E0"
  state_category    = "InProgress"
  order             = 3
}
`, processName, agileSystemProcessTypeId, workItemTypeName)
}

func getStateImportIdFunc(tfNode string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		res := state.RootModule().Resources[tfNode]
		id := res.Primary.Attributes["id"]
		processId := res.Primary.Attributes["process_id"]
		witRefName := res.Primary.Attributes["work_item_type_id"]
		return fmt.Sprintf("%s/%s/%s", processId, witRefName, id), nil
	}
}
