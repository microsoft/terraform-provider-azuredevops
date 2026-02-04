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

func TestAccWorkitemtrackingprocessRule_ConditionTypes(t *testing.T) {
	testCases := []struct {
		conditionType string
		field         string
		value         string
	}{
		{"when", "System.State", "New"},
		{"whenNot", "System.State", "Closed"},
		{"whenChanged", "System.Title", ""},
		{"whenNotChanged", "System.Title", ""},
		{"whenWas", "System.State", "New"},
	}

	for _, tc := range testCases {
		t.Run(tc.conditionType, func(t *testing.T) {
			workItemTypeName := testutils.GenerateWorkItemTypeName()
			processName := testutils.GenerateResourceName()
			tfNode := "azuredevops_workitemtrackingprocess_rule.test"

			resource.ParallelTest(t, resource.TestCase{
				PreCheck:          func() { testutils.PreCheck(t, nil) },
				ProviderFactories: testutils.GetProviderFactories(),
				CheckDestroy:      testutils.CheckProcessDestroyed,
				Steps: []resource.TestStep{
					{
						Config: ruleWithConditionType(workItemTypeName, processName, tc.conditionType, tc.field, tc.value),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttrSet(tfNode, "id"),
						),
					},
				},
			})
		})
	}
}

func TestAccWorkitemtrackingprocessRule_ActionTypes(t *testing.T) {
	testCases := []struct {
		actionType  string
		targetField string
		value       string
	}{
		{"makeRequired", "System.Title", ""},
		{"makeReadOnly", "System.Title", ""},
		{"setDefaultValue", "System.Title", "Default Title"},
		{"setDefaultFromClock", "System.ChangedDate", ""},
		{"setDefaultFromCurrentUser", "System.AssignedTo", ""},
		{"setDefaultFromField", "System.Description", "System.Title"},
		{"copyValue", "System.Title", "Copied Value"},
		{"copyFromClock", "System.ChangedDate", ""},
		{"copyFromCurrentUser", "System.AssignedTo", ""},
		{"copyFromField", "System.Description", "System.Title"},
		{"setValueToEmpty", "System.Description", ""},
		{"copyFromServerClock", "System.ChangedDate", ""},
		{"copyFromServerCurrentUser", "System.AssignedTo", ""},
		{"hideTargetField", "System.Description", ""},
		{"disallowValue", "System.State", "Closed"},
	}

	for _, tc := range testCases {
		t.Run(tc.actionType, func(t *testing.T) {
			workItemTypeName := testutils.GenerateWorkItemTypeName()
			processName := testutils.GenerateResourceName()
			tfNode := "azuredevops_workitemtrackingprocess_rule.test"

			resource.ParallelTest(t, resource.TestCase{
				PreCheck:          func() { testutils.PreCheck(t, nil) },
				ProviderFactories: testutils.GetProviderFactories(),
				CheckDestroy:      testutils.CheckProcessDestroyed,
				Steps: []resource.TestStep{
					{
						Config: ruleWithActionType(workItemTypeName, processName, tc.actionType, tc.targetField, tc.value),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttrSet(tfNode, "id"),
						),
					},
				},
			})
		})
	}
}

func TestAccWorkitemtrackingprocessRule_ConditionWhenCurrentUserIsMemberOfGroup(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	groupName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_rule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: ruleWithCurrentUserIsMemberOfGroup(workItemTypeName, processName, groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
		},
	})
}

func TestAccWorkitemtrackingprocessRule_ConditionWhenCurrentUserIsNotMemberOfGroup(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	groupName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_rule.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: ruleWithCurrentUserIsNotMemberOfGroup(workItemTypeName, processName, groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
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

func ruleWithCurrentUserIsMemberOfGroup(workItemTypeName, processName, groupName string) string {
	workItemType := basicWorkItemType(workItemTypeName, processName)
	return fmt.Sprintf(`
%s

resource "azuredevops_project" "group_test" {
  name = "%s-proj"
}

resource "azuredevops_group" "test" {
  scope        = azuredevops_project.group_test.id
  display_name = "%s"
}

resource "azuredevops_workitemtrackingprocess_rule" "test" {
  process_id        = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_id = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  name              = "Test whenCurrentUserIsMemberOfGroup Rule"

  condition {
    condition_type = "whenCurrentUserIsMemberOfGroup"
    value          = azuredevops_group.test.descriptor
  }

  action {
    action_type  = "makeRequired"
    target_field = "System.Title"
  }
}
`, workItemType, groupName, groupName)
}

func ruleWithCurrentUserIsNotMemberOfGroup(workItemTypeName, processName, groupName string) string {
	workItemType := basicWorkItemType(workItemTypeName, processName)
	return fmt.Sprintf(`
%s

resource "azuredevops_project" "group_test" {
  name = "%s-proj"
}

resource "azuredevops_group" "test" {
  scope        = azuredevops_project.group_test.id
  display_name = "%s"
}

resource "azuredevops_workitemtrackingprocess_rule" "test" {
  process_id        = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_id = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  name              = "Test whenCurrentUserIsNotMemberOfGroup Rule"

  condition {
    condition_type = "whenCurrentUserIsNotMemberOfGroup"
    value          = azuredevops_group.test.descriptor
  }

  action {
    action_type  = "makeRequired"
    target_field = "System.Title"
  }
}
`, workItemType, groupName, groupName)
}

func ruleWithConditionType(workItemTypeName, processName, conditionType, field, value string) string {
	workItemType := basicWorkItemType(workItemTypeName, processName)

	fieldAttr := ""
	if field != "" {
		fieldAttr = fmt.Sprintf(`field = "%s"`, field)
	}
	valueAttr := ""
	if value != "" {
		valueAttr = fmt.Sprintf(`value = "%s"`, value)
	}

	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtrackingprocess_rule" "test" {
  process_id        = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_id = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  name              = "Test %s Rule"

  condition {
    condition_type = "%s"
    %s
    %s
  }

  action {
    action_type  = "makeRequired"
    target_field = "System.Title"
  }
}
`, workItemType, conditionType, conditionType, fieldAttr, valueAttr)
}

func ruleWithActionType(workItemTypeName, processName, actionType, targetField, value string) string {
	workItemType := basicWorkItemType(workItemTypeName, processName)

	valueAttr := ""
	if value != "" {
		valueAttr = fmt.Sprintf(`value = "%s"`, value)
	}

	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtrackingprocess_rule" "test" {
  process_id        = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_id = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  name              = "Test %s Rule"

  condition {
    condition_type = "when"
    field          = "System.State"
    value          = "New"
  }

  action {
    action_type  = "%s"
    target_field = "%s"
    %s
  }
}
`, workItemType, actionType, actionType, targetField, valueAttr)
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
