package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccWorkitemtrackingprocessControl_Basic(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_control.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicControl(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "process_id"),
					resource.TestCheckResourceAttrSet(tfNode, "work_item_type_reference_name"),
					resource.TestCheckResourceAttrSet(tfNode, "group_id"),
					resource.TestCheckResourceAttr(tfNode, "label", "Test Control"),
					resource.TestCheckResourceAttr(tfNode, "visible", "true"),
					resource.TestCheckResourceAttr(tfNode, "read_only", "false"),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrSet(tfNode, "order"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: controlImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkitemtrackingprocessControl_Update(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_control.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicControl(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "process_id"),
					resource.TestCheckResourceAttrSet(tfNode, "work_item_type_reference_name"),
					resource.TestCheckResourceAttrSet(tfNode, "group_id"),
					resource.TestCheckResourceAttr(tfNode, "label", "Test Control"),
					resource.TestCheckResourceAttr(tfNode, "visible", "true"),
					resource.TestCheckResourceAttr(tfNode, "read_only", "false"),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrSet(tfNode, "order"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: controlImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: updatedControl(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "process_id"),
					resource.TestCheckResourceAttrSet(tfNode, "work_item_type_reference_name"),
					resource.TestCheckResourceAttrSet(tfNode, "group_id"),
					resource.TestCheckResourceAttr(tfNode, "label", "Updated Control"),
					resource.TestCheckResourceAttr(tfNode, "visible", "false"),
					resource.TestCheckResourceAttr(tfNode, "read_only", "true"),
					resource.TestCheckResourceAttr(tfNode, "order", "0"),
					resource.TestCheckResourceAttr(tfNode, "metadata", "test metadata"),
					resource.TestCheckResourceAttr(tfNode, "watermark", "Enter a title"),
					resource.TestCheckResourceAttr(tfNode, "control_type", "FieldControl"),
					resource.TestCheckResourceAttr(tfNode, "inherited", "false"),
					resource.TestCheckResourceAttr(tfNode, "overridden", "false"),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: controlImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkitemtrackingprocessControl_Move(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_control.test"

	var originalGroupId string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicControl(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "process_id"),
					resource.TestCheckResourceAttrSet(tfNode, "work_item_type_reference_name"),
					resource.TestCheckResourceAttrSet(tfNode, "group_id"),
					resource.TestCheckResourceAttr(tfNode, "label", "Test Control"),
					resource.TestCheckResourceAttr(tfNode, "visible", "true"),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrSet(tfNode, "order"),
					resource.TestCheckResourceAttrWith(tfNode, "group_id", func(value string) error {
						originalGroupId = value
						return nil
					}),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: controlImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: movedControl(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "process_id"),
					resource.TestCheckResourceAttrSet(tfNode, "work_item_type_reference_name"),
					resource.TestCheckResourceAttrSet(tfNode, "group_id"),
					resource.TestCheckResourceAttr(tfNode, "label", "Test Control"),
					resource.TestCheckResourceAttr(tfNode, "visible", "true"),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrSet(tfNode, "order"),
					resource.TestCheckResourceAttrWith(tfNode, "group_id", func(value string) error {
						if value == originalGroupId {
							return fmt.Errorf("group_id should have changed, but is still %s", originalGroupId)
						}
						return nil
					}),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: controlImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkitemtrackingprocessControl_Contribution(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_control.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: contributionControl(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "process_id"),
					resource.TestCheckResourceAttrSet(tfNode, "work_item_type_reference_name"),
					resource.TestCheckResourceAttrSet(tfNode, "group_id"),
					resource.TestCheckResourceAttr(tfNode, "is_contribution", "true"),
					resource.TestCheckResourceAttr(tfNode, "contribution.#", "1"),
					resource.TestCheckResourceAttr(tfNode, "contribution.0.contribution_id", "ms-devlabs.vsts-extensions-multivalue-control.multivalue-form-control"),
					resource.TestCheckResourceAttr(tfNode, "contribution.0.height", "50"),
					resource.TestCheckResourceAttr(tfNode, "contribution.0.inputs.FieldName", "System.Tags"),
					resource.TestCheckResourceAttr(tfNode, "contribution.0.inputs.Values", "Option1;Option2;Option3"),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: controlImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func basicControl(workItemTypeName string, processName string) string {
	workItemType := basicWorkItemType(workItemTypeName, processName)
	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtrackingprocess_group" "test" {
  process_id                    = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  page_id                       = azuredevops_workitemtrackingprocess_workitemtype.test.pages[0].id
  section_id                    = azuredevops_workitemtrackingprocess_workitemtype.test.pages[0].sections[0].id
  label                         = "Test Group"
}

resource "azuredevops_workitemtrackingprocess_control" "test" {
  process_id                    = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  group_id                      = azuredevops_workitemtrackingprocess_group.test.id
  control_id                    = "System.Title"
  label                         = "Test Control"
}
`, workItemType)
}

func updatedControl(workItemTypeName string, processName string) string {
	workItemType := basicWorkItemType(workItemTypeName, processName)
	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtrackingprocess_group" "test" {
  process_id                    = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  page_id                       = azuredevops_workitemtrackingprocess_workitemtype.test.pages[0].id
  section_id                    = azuredevops_workitemtrackingprocess_workitemtype.test.pages[0].sections[0].id
  label                         = "Test Group"
}

resource "azuredevops_workitemtrackingprocess_control" "test" {
  process_id                    = azuredevops_workitemtrackingprocess_process.test.id
  control_id                    = "System.Title"
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  group_id                      = azuredevops_workitemtrackingprocess_group.test.id
  label                         = "Updated Control"
  visible                       = false
  read_only                     = true
  order                         = 0
  metadata                      = "test metadata"
  watermark                     = "Enter a title"
}
`, workItemType)
}

func movedControl(workItemTypeName string, processName string) string {
	workItemType := basicWorkItemType(workItemTypeName, processName)
	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtrackingprocess_group" "test" {
  process_id                    = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  page_id                       = azuredevops_workitemtrackingprocess_workitemtype.test.pages[0].id
  section_id                    = azuredevops_workitemtrackingprocess_workitemtype.test.pages[0].sections[0].id
  label                         = "Test Group"
}

resource "azuredevops_workitemtrackingprocess_group" "test2" {
  process_id                    = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  page_id                       = azuredevops_workitemtrackingprocess_workitemtype.test.pages[0].id
  section_id                    = azuredevops_workitemtrackingprocess_workitemtype.test.pages[0].sections[0].id
  label                         = "Test Group 2"
}

resource "azuredevops_workitemtrackingprocess_control" "test" {
  process_id                    = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  group_id                      = azuredevops_workitemtrackingprocess_group.test2.id
  control_id                    = "System.Title"
  label                         = "Test Control"
}
`, workItemType)
}

func contributionControl(workItemTypeName string, processName string) string {
	workItemType := basicWorkItemType(workItemTypeName, processName)
	return fmt.Sprintf(`
%s

resource "azuredevops_extension" "test" {
  publisher_id = "ms-devlabs"
  extension_id = "vsts-extensions-multivalue-control"
}

resource "azuredevops_workitemtrackingprocess_group" "test" {
  process_id                    = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  page_id                       = azuredevops_workitemtrackingprocess_workitemtype.test.pages[0].id
  section_id                    = azuredevops_workitemtrackingprocess_workitemtype.test.pages[0].sections[0].id
  label                         = "Test Group"
}

resource "azuredevops_workitemtrackingprocess_control" "test" {
  depends_on                    = [azuredevops_extension.test]
  process_id                    = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  group_id                      = azuredevops_workitemtrackingprocess_group.test.id
  control_id                    = "ms-devlabs.vsts-extensions-multivalue-control.multivalue-form-control"
  is_contribution               = true
  contribution {
    contribution_id = "ms-devlabs.vsts-extensions-multivalue-control.multivalue-form-control"
    height          = 50
    inputs = {
      FieldName = "System.Tags"
      Values    = "Option1;Option2;Option3"
    }
  }
}
`, workItemType)
}

func controlImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		processId := rs.Primary.Attributes["process_id"]
		witRefName := rs.Primary.Attributes["work_item_type_reference_name"]
		groupId := rs.Primary.Attributes["group_id"]
		controlId := rs.Primary.ID
		return fmt.Sprintf("%s/%s/%s/%s", processId, witRefName, groupId, controlId), nil
	}
}
