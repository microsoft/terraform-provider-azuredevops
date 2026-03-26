package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
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
				Config: inheritedStateConfig(processName, false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(tfNode, tfjsonpath.New("id"), knownvalue.NotNull()),
				},
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
				Config: inheritedStateConfig(processName, false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(tfNode, tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: inheritedStateImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: inheritedStateConfig(processName, true),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(tfNode, tfjsonpath.New("id"), knownvalue.NotNull()),
				},
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
	tfNode := "azuredevops_workitemtrackingprocess_inherited_state.test"

	var stateId string
	var processId string
	var witRefName string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: inheritedStateConfig(processName, false),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(tfNode, tfjsonpath.New("id"),
						knownvalue.StringFunc(func(value string) error {
							stateId = value
							return nil
						})),
					statecheck.ExpectKnownValue(tfNode, tfjsonpath.New("process_id"),
						knownvalue.StringFunc(func(value string) error {
							processId = value
							return nil
						})),
					statecheck.ExpectKnownValue(tfNode, tfjsonpath.New("work_item_type_id"),
						knownvalue.StringFunc(func(value string) error {
							witRefName = value
							return nil
						})),
				},
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: inheritedStateImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Remove the inherited_state resource from config - state should still exist in Azure DevOps
				Config: removedInheritedState(processName),
				Check: resource.ComposeTestCheckFunc(
					checkInheritedStateReverted(&processId, &witRefName, &stateId),
				),
			},
		},
	})
}

func inheritedStateConfig(processName string, visible bool) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtrackingprocess_process" "test" {
  name                   = "%s"
  parent_process_type_id = "%s"
}

resource "azuredevops_workitemtrackingprocess_workitemtype" "test" {
  process_id                      = azuredevops_workitemtrackingprocess_process.test.id
  name                            = "Bug"
  parent_work_item_reference_name = "Microsoft.VSTS.WorkItemTypes.Bug"
}

resource "azuredevops_workitemtrackingprocess_inherited_state" "test" {
  process_id        = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_id = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  name              = "New"
  visible           = %t
}
`, processName, agileSystemProcessTypeId, visible)
}

func removedInheritedState(processName string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtrackingprocess_process" "test" {
  name                   = "%s"
  parent_process_type_id = "%s"
}

resource "azuredevops_workitemtrackingprocess_workitemtype" "test" {
  process_id                      = azuredevops_workitemtrackingprocess_process.test.id
  name                            = "Bug"
  parent_work_item_reference_name = "Microsoft.VSTS.WorkItemTypes.Bug"
}
`, processName, agileSystemProcessTypeId)
}

func checkInheritedStateReverted(processIdStr *string, witRefName *string, stateId *string) resource.TestCheckFunc {
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

		if state.Hidden != nil && *state.Hidden == true {
			return fmt.Errorf("inherited state should have reverted to visible")
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
		witRefName := rs.Primary.Attributes["work_item_type_id"]
		stateName := rs.Primary.Attributes["name"]
		return fmt.Sprintf("%s/%s/%s", processId, witRefName, stateName), nil
	}
}
