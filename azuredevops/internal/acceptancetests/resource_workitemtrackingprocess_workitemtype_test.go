package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccWorkitemtrackingprocessWorkItemType_Basic(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_workitemtype.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicWorkItemType(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", workItemTypeName),
					resource.TestCheckResourceAttrSet(tfNode, "process_id"),
					resource.TestCheckResourceAttr(tfNode, "is_disabled", "false"),
					resource.TestCheckResourceAttrSet(tfNode, "url"),
					resource.TestCheckResourceAttrSet(tfNode, "color"),
					resource.TestCheckResourceAttrSet(tfNode, "icon"),
					resource.TestCheckResourceAttr(tfNode, "inherits_from", ""),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getWorkItemTypeStateIdFunc(tfNode),
			},
		},
	})
}

func TestAccWorkitemtrackingprocessWorkItemType_CreateAndUpdate(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()

	tfNode := "azuredevops_workitemtrackingprocess_workitemtype.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicWorkItemType(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", workItemTypeName),
					resource.TestCheckResourceAttrSet(tfNode, "process_id"),
					resource.TestCheckResourceAttr(tfNode, "is_disabled", "false"),
					resource.TestCheckResourceAttrSet(tfNode, "url"),
					resource.TestCheckResourceAttrSet(tfNode, "color"),
					resource.TestCheckResourceAttrSet(tfNode, "icon"),
					resource.TestCheckResourceAttr(tfNode, "inherits_from", ""),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getWorkItemTypeStateIdFunc(tfNode),
			},
			{
				Config: disabledWorkItemType(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", workItemTypeName),
					resource.TestCheckResourceAttrSet(tfNode, "process_id"),
					resource.TestCheckResourceAttr(tfNode, "is_disabled", "true"),
					resource.TestCheckResourceAttrSet(tfNode, "url"),
					resource.TestCheckResourceAttrSet(tfNode, "color"),
					resource.TestCheckResourceAttrSet(tfNode, "icon"),
					resource.TestCheckResourceAttr(tfNode, "inherits_from", ""),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getWorkItemTypeStateIdFunc(tfNode),
			},
		},
	})
}

func basicWorkItemType(name string, processName string) string {
	process := process(processName)
	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtrackingprocess_workitemtype" "test" {
  name       = "%s"
  process_id = azuredevops_workitemtrackingprocess_process.test.id
}
`, process, name)
}

func disabledWorkItemType(name string, processName string) string {
	process := process(processName)
	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtrackingprocess_workitemtype" "test" {
  name        = "%s"
  process_id  = azuredevops_workitemtrackingprocess_process.test.id
  is_disabled = true
}
`, process, name)
}

func getWorkItemTypeStateIdFunc(tfNode string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		res := state.RootModule().Resources[tfNode]
		id := res.Primary.Attributes["id"]
		processId := res.Primary.Attributes["process_id"]
		return fmt.Sprintf("%s/%s", processId, id), nil
	}
}
