package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtrackingprocess"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

func TestAccWorkitemtrackingprocessInheritedState_Basic(t *testing.T) {
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_inherited_state.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: inheritedStateConfig(processName, true),
				Check:  resource.TestCheckResourceAttrSet(tfNode, "id"),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: inheritedStateImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkitemtrackingprocessInheritedState_Update(t *testing.T) {
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_inherited_state.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: inheritedStateConfig(processName, true),
				Check:  resource.TestCheckResourceAttrSet(tfNode, "id"),
			},
			{
				Config: inheritedStateConfig(processName, false),
				Check:  resource.TestCheckResourceAttrSet(tfNode, "id"),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: inheritedStateImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkitemtrackingprocessInheritedState_RemoveFromState(t *testing.T) {
	processName := testutils.GenerateResourceName()
	inheritedStateNode := "azuredevops_workitemtrackingprocess_inherited_state.test"

	var stateId string
	var processId string
	var witRefName string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: inheritedStateConfig(processName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrWith(inheritedStateNode, "id", func(value string) error {
						stateId = value
						return nil
					}),
					resource.TestCheckResourceAttrWith(inheritedStateNode, "process_id", func(value string) error {
						processId = value
						return nil
					}),
					resource.TestCheckResourceAttrWith(inheritedStateNode, "work_item_type_reference_name", func(value string) error {
						witRefName = value
						return nil
					}),
				),
			},
			{
				// Remove the inherited_state resource from config - state should still exist in Azure DevOps
				Config: removedInheritedState(processName),
				Check: resource.ComposeTestCheckFunc(
					checkInheritedStateStillExists(&processId, &witRefName, &stateId),
				),
			},
		},
	})
}

func inheritedStateConfig(processName string, hidden bool) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtrackingprocess_process" "test" {
  name                   = "%s"
  parent_process_type_id = "adcc42ab-9882-485e-a3ed-7678f01f66bc"
}

resource "azuredevops_workitemtrackingprocess_workitemtype" "test" {
  process_id                      = azuredevops_workitemtrackingprocess_process.test.id
  name                            = "Bug"
  parent_work_item_reference_name = "Microsoft.VSTS.WorkItemTypes.Bug"
}

resource "azuredevops_workitemtrackingprocess_inherited_state" "test" {
  process_id                    = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  name                          = "New"
  hidden                        = %t
}
`, processName, hidden)
}

func removedInheritedState(processName string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtrackingprocess_process" "test" {
  name                   = "%s"
  parent_process_type_id = "adcc42ab-9882-485e-a3ed-7678f01f66bc"
}

resource "azuredevops_workitemtrackingprocess_workitemtype" "test" {
  process_id                      = azuredevops_workitemtrackingprocess_process.test.id
  name                            = "Bug"
  parent_work_item_reference_name = "Microsoft.VSTS.WorkItemTypes.Bug"
}
`, processName)
}

func checkInheritedStateStillExists(processIdStr *string, witRefName *string, stateId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

		processId, err := uuid.Parse(*processIdStr)
		if err != nil {
			return fmt.Errorf("invalid process_id: %w", err)
		}

		stateUUID, err := uuid.Parse(*stateId)
		if err != nil {
			return fmt.Errorf("invalid state_id: %w", err)
		}

		state, err := clients.WorkItemTrackingProcessClient.GetStateDefinition(clients.Ctx, workitemtrackingprocess.GetStateDefinitionArgs{
			ProcessId:  &processId,
			WitRefName: witRefName,
			StateId:    &stateUUID,
		})
		if err != nil {
			return fmt.Errorf("getting state definition: %w", err)
		}

		if state == nil {
			return fmt.Errorf("inherited state should still exist after removing from Terraform but it was not found")
		}

		return nil
	}
}

func inheritedStateImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		processId := rs.Primary.Attributes["process_id"]
		witRefName := rs.Primary.Attributes["work_item_type_reference_name"]
		stateName := rs.Primary.Attributes["name"]
		return fmt.Sprintf("%s/%s/%s", processId, witRefName, stateName), nil
	}
}
