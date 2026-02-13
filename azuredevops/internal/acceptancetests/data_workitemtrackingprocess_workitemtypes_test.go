package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccWorkitemtrackingprocessWorkItemTypes_DataSource_List(t *testing.T) {
	workItemTypeName1 := testutils.GenerateWorkItemTypeName()
	workItemTypeName2 := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfDataNode := "data.azuredevops_workitemtrackingprocess_workitemtypes.test"
	tfResourceNode1 := "azuredevops_workitemtrackingprocess_workitemtype.test1"
	tfResourceNode2 := "azuredevops_workitemtrackingprocess_workitemtype.test2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclDataSourceWorkItemTypes(workItemTypeName1, workItemTypeName2, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfDataNode, "id"),
					resource.TestCheckResourceAttrSet(tfDataNode, "process_id"),
					// It might contain an unknown amount of inherited work item types
					testutils.TestCheckAttrGreaterThan(tfDataNode, "work_item_types.#", 1),
					resource.TestCheckTypeSetElemNestedAttrs(tfDataNode, "work_item_types.*", map[string]string{
						"name":        workItemTypeName1,
						"description": "Test work item type 1",
					}),
					resource.TestCheckTypeSetElemAttrPair(tfDataNode, "work_item_types.*.reference_name", tfResourceNode1, "reference_name"),
					resource.TestCheckTypeSetElemNestedAttrs(tfDataNode, "work_item_types.*", map[string]string{
						"name":        workItemTypeName2,
						"description": "Test work item type 2",
					}),
					resource.TestCheckTypeSetElemAttrPair(tfDataNode, "work_item_types.*.reference_name", tfResourceNode2, "reference_name"),
				),
			},
		},
	})
}

func hclDataSourceWorkItemTypes(workItemTypeName1 string, workItemTypeName2 string, processName string) string {
	process := process(processName)
	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtrackingprocess_workitemtype" "test1" {
  name        = "%s"
  process_id  = azuredevops_workitemtrackingprocess_process.test.id
  description = "Test work item type 1"
}

resource "azuredevops_workitemtrackingprocess_workitemtype" "test2" {
  name        = "%s"
  process_id  = azuredevops_workitemtrackingprocess_process.test.id
  description = "Test work item type 2"
}

data "azuredevops_workitemtrackingprocess_workitemtypes" "test" {
  process_id = azuredevops_workitemtrackingprocess_process.test.id
  depends_on = [
    azuredevops_workitemtrackingprocess_workitemtype.test1,
    azuredevops_workitemtrackingprocess_workitemtype.test2
  ]
}
`, process, workItemTypeName1, workItemTypeName2)
}
