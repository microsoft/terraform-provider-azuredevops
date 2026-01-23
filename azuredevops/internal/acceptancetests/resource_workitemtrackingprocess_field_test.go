package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccWorkitemtrackingprocessField_Basic(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	fieldName := generateFieldName()
	tfNode := "azuredevops_workitemtrackingprocess_field.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkProcessAndFieldDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicField(workItemTypeName, processName, fieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
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

func TestAccWorkitemtrackingprocessField_Identity(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	fieldName := generateFieldName()
	tfNode := "azuredevops_workitemtrackingprocess_field.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkProcessAndFieldDestroyed,
		Steps: []resource.TestStep{
			{
				Config: identityField(workItemTypeName, processName, fieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
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

func TestAccWorkitemtrackingprocessField_Update(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	fieldName := generateFieldName()
	tfNode := "azuredevops_workitemtrackingprocess_field.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkProcessAndFieldDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicField(workItemTypeName, processName, fieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: updatedField(workItemTypeName, processName, fieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
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

func basicField(workItemTypeName string, processName string, fieldName string) string {
	testProcessAndWit := basicWorkItemType(workItemTypeName, processName)
	testField := fieldBasic(fieldName)
	return fmt.Sprintf(`
%s

%s

resource "azuredevops_workitemtrackingprocess_field" "test" {
  process_id        = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_id = azuredevops_workitemtrackingprocess_workitemtype.test.id
  field_id          = azuredevops_workitemtracking_field.test.id
}
`, testProcessAndWit, testField)
}

func updatedField(workItemTypeName string, processName string, fieldName string) string {
	testProcessAndWit := basicWorkItemType(workItemTypeName, processName)
	testField := fieldBasic(fieldName)
	return fmt.Sprintf(`
%s

%s

resource "azuredevops_workitemtrackingprocess_field" "test" {
  process_id        = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_id = azuredevops_workitemtrackingprocess_workitemtype.test.id
  field_id          = azuredevops_workitemtracking_field.test.id
  read_only         = true
  required          = true
  default_value     = "default"
}
`, testProcessAndWit, testField)
}

func identityField(workItemTypeName string, processName string, fieldName string) string {
	testProcessAndWit := basicWorkItemType(workItemTypeName, processName)
	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtracking_field" "test" {
  name           = "%s"
  reference_name = "Custom.%s"
  type           = "identity"
}

resource "azuredevops_workitemtrackingprocess_field" "test" {
  process_id        = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_id = azuredevops_workitemtrackingprocess_workitemtype.test.id
  field_id          = azuredevops_workitemtracking_field.test.id
  allow_groups      = true
}
`, testProcessAndWit, fieldName, fieldName)
}

func checkProcessAndFieldDestroyed(s *terraform.State) error {
	if err := testutils.CheckProcessDestroyed(s); err != nil {
		return err
	}
	return checkFieldDestroyed(s)
}
