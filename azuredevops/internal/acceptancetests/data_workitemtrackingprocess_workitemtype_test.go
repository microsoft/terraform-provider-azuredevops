package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccWorkitemtrackingprocessWorkItemType_DataSource_Get(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfResourceNode := "azuredevops_workitemtrackingprocess_workitemtype.test"
	tfDataNode := "data.azuredevops_workitemtrackingprocess_workitemtype.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclDataSourceWorkItemType(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(tfDataNode, "id", tfResourceNode, "id"),
					resource.TestCheckResourceAttrPair(tfDataNode, "process_id", tfResourceNode, "process_id"),
					resource.TestCheckResourceAttrPair(tfDataNode, "reference_name", tfResourceNode, "id"),
					resource.TestCheckResourceAttrPair(tfDataNode, "name", tfResourceNode, "name"),
					resource.TestCheckResourceAttrPair(tfDataNode, "description", tfResourceNode, "description"),
					resource.TestCheckResourceAttrPair(tfDataNode, "color", tfResourceNode, "color"),
					resource.TestCheckResourceAttrPair(tfDataNode, "icon", tfResourceNode, "icon"),
					resource.TestCheckResourceAttrPair(tfDataNode, "is_enabled", tfResourceNode, "is_enabled"),
					resource.TestCheckResourceAttrPair(tfDataNode, "url", tfResourceNode, "url"),
				),
			},
		},
	})
}

func hclDataSourceWorkItemType(workItemTypeName string, processName string) string {
	process := process(processName)
	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtrackingprocess_workitemtype" "test" {
  name        = "%s"
  process_id  = azuredevops_workitemtrackingprocess_process.test.id
  description = "Test work item type"
}

data "azuredevops_workitemtrackingprocess_workitemtype" "test" {
  process_id     = azuredevops_workitemtrackingprocess_workitemtype.test.process_id
  reference_name = azuredevops_workitemtrackingprocess_workitemtype.test.reference_name
}
`, process, workItemTypeName)
}
