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
	tfNode := "azuredevops_workitemtrackingprocess_field.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicField(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "reference_name", "Custom.TestField"),
					resource.TestCheckResourceAttrSet(tfNode, "process_id"),
					resource.TestCheckResourceAttrSet(tfNode, "work_item_type_ref_name"),
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttrSet(tfNode, "type"),
					resource.TestCheckResourceAttr(tfNode, "read_only", "false"),
					resource.TestCheckResourceAttr(tfNode, "required", "false"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getFieldStateIdFunc(tfNode),
			},
		},
	})
}

func TestAccWorkitemtrackingprocessField_Update(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_field.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicField(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "reference_name", "Custom.TestField"),
					resource.TestCheckResourceAttr(tfNode, "read_only", "false"),
					resource.TestCheckResourceAttr(tfNode, "required", "false"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getFieldStateIdFunc(tfNode),
			},
			{
				Config: updatedField(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "reference_name", "Custom.TestField"),
					resource.TestCheckResourceAttr(tfNode, "required", "true"),
					resource.TestCheckResourceAttr(tfNode, "default_value_json", "\"default\""),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getFieldStateIdFunc(tfNode),
			},
		},
	})
}

func basicField(workItemTypeName string, processName string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtrackingprocess_process" "test" {
  name = "%s"
}

resource "azuredevops_workitemtrackingprocess_workitemtype" "test" {
  name       = "%s"
  process_id = azuredevops_workitemtrackingprocess_process.test.id
}

resource "azuredevops_workitemtracking_field" "test" {
  name           = "TestField"
  reference_name = "Custom.TestField"
  type           = "string"
}

resource "azuredevops_workitemtrackingprocess_field" "test" {
  process_id              = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_ref_name = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  reference_name          = azuredevops_workitemtracking_field.test.reference_name
}
`, processName, workItemTypeName)
}

func updatedField(workItemTypeName string, processName string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtrackingprocess_process" "test" {
  name = "%s"
}

resource "azuredevops_workitemtrackingprocess_workitemtype" "test" {
  name       = "%s"
  process_id = azuredevops_workitemtrackingprocess_process.test.id
}

resource "azuredevops_workitemtracking_field" "test" {
  name           = "TestField"
  reference_name = "Custom.TestField"
  type           = "string"
}

resource "azuredevops_workitemtrackingprocess_field" "test" {
  process_id              = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_ref_name = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
  reference_name          = azuredevops_workitemtracking_field.test.reference_name
  required                = true
  default_value_json      = "\"default\""
}
`, processName, workItemTypeName)
}

func getFieldStateIdFunc(tfNode string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		res := state.RootModule().Resources[tfNode]
		id := res.Primary.ID
		processId := res.Primary.Attributes["process_id"]
		witRefName := res.Primary.Attributes["work_item_type_ref_name"]
		return fmt.Sprintf("%s/%s/%s", processId, witRefName, id), nil
	}
}
