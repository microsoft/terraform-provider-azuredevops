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

func TestAccWorkitemtrackingprocessInheritedControl_Basic(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_inherited_control.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicInheritedControl(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getInheritedControlImportIdFunc(tfNode),
			},
		},
	})
}

func TestAccWorkitemtrackingprocessInheritedControl_Update(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_inherited_control.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicInheritedControl(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				Config: updatedInheritedControl(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getInheritedControlImportIdFunc(tfNode),
			},
		},
	})
}

func TestAccWorkitemtrackingprocessInheritedControl_Revert(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_inherited_control.test"

	var processId, witRefName, groupId, controlId string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicInheritedControl(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrWith(tfNode, "process_id", func(value string) error {
						processId = value
						return nil
					}),
					resource.TestCheckResourceAttrWith(tfNode, "work_item_type_reference_name", func(value string) error {
						witRefName = value
						return nil
					}),
					resource.TestCheckResourceAttrWith(tfNode, "group_id", func(value string) error {
						groupId = value
						return nil
					}),
					resource.TestCheckResourceAttrWith(tfNode, "id", func(value string) error {
						controlId = value
						return nil
					}),
				),
			},
			{
				// Remove the inherited control customization, keeping only process and work item type
				Config: inheritedControlRevertConfig(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					checkInheritedControlRevertedFunc(&processId, &witRefName, &groupId, &controlId),
				),
			},
		},
	})
}

func checkInheritedControlRevertedFunc(processId, witRefName, groupId, controlId *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

		// Get the work item type layout to check if control is still overridden
		args := workitemtrackingprocess.GetProcessWorkItemTypeArgs{
			ProcessId:  converter.UUID(*processId),
			WitRefName: witRefName,
			Expand:     &workitemtrackingprocess.GetWorkItemTypeExpandValues.Layout,
		}
		workItemType, err := clients.WorkItemTrackingProcessClient.GetProcessWorkItemType(context.Background(), args)
		if err != nil {
			return fmt.Errorf("error getting work item type: %+v", err)
		}

		if workItemType == nil || workItemType.Layout == nil {
			return fmt.Errorf("work item type or layout is nil")
		}

		// Find the group and control
		group := findGroupById(workItemType.Layout, *groupId)
		if group == nil {
			// Group may have been removed, which is fine
			return nil
		}

		control := findControlInGroup(group, *controlId)
		if control == nil {
			// Control not found in this group, which is expected after revert
			return nil
		}

		// If control exists but is marked as inherited (not overridden), that's the expected state
		if control.Inherited != nil && *control.Inherited && (control.Overridden == nil || !*control.Overridden) {
			return nil
		}

		return fmt.Errorf("inherited control %s should have been reverted but is still customized", *controlId)
	}
}

func basicInheritedControl(workItemTypeName string, processName string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtrackingprocess_process" "test" {
  name                   = "%s"
  parent_process_type_id = "%s"
}

resource "azuredevops_workitemtrackingprocess_workitemtype" "test" {
  name       = "%s"
  process_id = azuredevops_workitemtrackingprocess_process.test.id
}

resource "azuredevops_workitemtrackingprocess_inherited_control" "test" {
  process_id                    = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  group_id                      = azuredevops_workitemtrackingprocess_workitemtype.test.pages[0].sections[0].groups[0].id
  control_id                    = azuredevops_workitemtrackingprocess_workitemtype.test.pages[0].sections[0].groups[0].controls[0].id
  visible                       = false
}
`, processName, agileSystemProcessTypeId, workItemTypeName)
}

func updatedInheritedControl(workItemTypeName string, processName string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtrackingprocess_process" "test" {
  name                   = "%s"
  parent_process_type_id = "%s"
}

resource "azuredevops_workitemtrackingprocess_workitemtype" "test" {
  name       = "%s"
  process_id = azuredevops_workitemtrackingprocess_process.test.id
}

resource "azuredevops_workitemtrackingprocess_inherited_control" "test" {
  process_id                    = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  group_id                      = azuredevops_workitemtrackingprocess_workitemtype.test.pages[0].sections[0].groups[0].id
  control_id                    = azuredevops_workitemtrackingprocess_workitemtype.test.pages[0].sections[0].groups[0].controls[0].id
  visible                       = true
  label                         = "Custom Label"
}
`, processName, agileSystemProcessTypeId, workItemTypeName)
}

func inheritedControlRevertConfig(workItemTypeName string, processName string) string {
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

func getInheritedControlImportIdFunc(tfNode string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		res := state.RootModule().Resources[tfNode]
		id := res.Primary.Attributes["id"]
		processId := res.Primary.Attributes["process_id"]
		witRefName := res.Primary.Attributes["work_item_type_reference_name"]
		groupId := res.Primary.Attributes["group_id"]
		return fmt.Sprintf("%s/%s/%s/%s", processId, witRefName, groupId, id), nil
	}
}

func findGroupById(layout *workitemtrackingprocess.FormLayout, groupId string) *workitemtrackingprocess.Group {
	if layout == nil || layout.Pages == nil {
		return nil
	}
	for _, page := range *layout.Pages {
		if page.Sections == nil {
			continue
		}
		for _, section := range *page.Sections {
			if section.Groups == nil {
				continue
			}
			for _, group := range *section.Groups {
				if group.Id != nil && *group.Id == groupId {
					return &group
				}
			}
		}
	}
	return nil
}

func findControlInGroup(group *workitemtrackingprocess.Group, controlId string) *workitemtrackingprocess.Control {
	if group == nil || group.Controls == nil {
		return nil
	}
	for _, control := range *group.Controls {
		if control.Id != nil && *control.Id == controlId {
			return &control
		}
	}
	return nil
}
