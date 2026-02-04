package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccWorkitemtrackingprocessRule_Basic(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_rule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicRule(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: ruleImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkitemtrackingprocessRule_Update(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_rule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicRule(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: ruleImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: updatedRule(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: ruleImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkitemtrackingprocessRule_MultipleConditionsAndActions(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_rule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: multipleConditionsRule(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: ruleImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func basicRule(workItemTypeName string, processName string) string {
	workItemType := basicWorkItemType(workItemTypeName, processName)
	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtrackingprocess_rule" "test" {
  process_id        = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_id = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  name              = "Test Rule"

  condition {
    condition_type = "when"
    field          = "System.State"
    value          = "New"
  }

  action {
    action_type  = "makeRequired"
    target_field = "System.Title"
  }
}
`, workItemType)
}

func updatedRule(workItemTypeName string, processName string) string {
	workItemType := basicWorkItemType(workItemTypeName, processName)
	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtrackingprocess_rule" "test" {
  process_id        = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_id = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  name              = "Updated Rule"
  is_disabled       = true

  condition {
    condition_type = "when"
    field          = "System.State"
    value          = "Active"
  }

  action {
    action_type  = "makeReadOnly"
    target_field = "System.Title"
  }
}
`, workItemType)
}

func multipleConditionsRule(workItemTypeName string, processName string) string {
	workItemType := basicWorkItemType(workItemTypeName, processName)
	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtrackingprocess_rule" "test" {
  process_id        = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_id = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  name              = "Multiple Conditions Rule"

  condition {
    condition_type = "whenWas"
    field          = "System.State"
    value          = "New"
  }

  condition {
    condition_type = "when"
    field          = "System.State"
    value          = "Active"
  }

  condition {
    condition_type = "whenChanged"
    field          = "System.Title"
  }

  action {
    action_type  = "makeRequired"
    target_field = "System.Title"
  }

  action {
    action_type  = "makeReadOnly"
    target_field = "System.Description"
  }
}
`, workItemType)
}

func ruleImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		processId := rs.Primary.Attributes["process_id"]
		witRefName := rs.Primary.Attributes["work_item_type_id"]
		ruleId := rs.Primary.ID
		return fmt.Sprintf("%s/%s/%s", processId, witRefName, ruleId), nil
	}
}
