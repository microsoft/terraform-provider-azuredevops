package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
					resource.TestCheckResourceAttr(tfNode, "is_enabled", "true"),
					resource.TestCheckResourceAttrSet(tfNode, "url"),
					resource.TestCheckResourceAttrSet(tfNode, "color"),
					resource.TestCheckResourceAttrSet(tfNode, "icon"),
					resource.TestCheckResourceAttrSet(tfNode, "reference_name"),
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
					resource.TestCheckResourceAttr(tfNode, "is_enabled", "true"),
					resource.TestCheckResourceAttrSet(tfNode, "url"),
					resource.TestCheckResourceAttrSet(tfNode, "color"),
					resource.TestCheckResourceAttrSet(tfNode, "icon"),
					resource.TestCheckResourceAttrSet(tfNode, "reference_name"),
					resource.TestCheckResourceAttrSet(tfNode, "pages.#"),
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
					resource.TestCheckResourceAttr(tfNode, "is_enabled", "false"),
					resource.TestCheckResourceAttrSet(tfNode, "url"),
					resource.TestCheckResourceAttrSet(tfNode, "color"),
					resource.TestCheckResourceAttrSet(tfNode, "icon"),
					resource.TestCheckResourceAttrSet(tfNode, "reference_name"),
					resource.TestCheckResourceAttrSet(tfNode, "pages.#"),
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

func TestAccWorkitemtrackingprocessWorkItemType_States(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_workitemtype.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: workItemTypeWithStates(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "state.#", "2"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getWorkItemTypeStateIdFunc(tfNode),
			},
			{
				Config: workItemTypeWithStatesUpdated(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "state.#", "3"),
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

// Azure DevOps rejects any Update against a Completed state (VS403093) even
// when the payload is identical, so the sync must skip those calls.
func TestAccWorkitemtrackingprocessWorkItemType_StatesWithNoChanges(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_workitemtype.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: workItemTypeStatesWithNoChanges(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "state.#", "3"),
				),
			},
		},
	})
}

func workItemTypeStatesWithNoChanges(name, processName string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtrackingprocess_workitemtype" "test" {
  name       = "%s"
  process_id = azuredevops_workitemtrackingprocess_process.test.id

  state {
    name           = "New"
    color          = "#3544ca"
    state_category = "Proposed"
  }

  state {
    name           = "Active"
    color          = "#ff9d00"
    state_category = "InProgress"
  }

  state {
    name           = "Closed"
    color          = "#339933"
    state_category = "Completed"
  }
}
`, process(processName), name)
}

func workItemTypeWithStates(name, processName string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtrackingprocess_workitemtype" "test" {
  name       = "%s"
  process_id = azuredevops_workitemtrackingprocess_process.test.id

  state {
    name           = "Active 2"
    color          = "#ff9d00"
    state_category = "InProgress"
  }

  state {
    name           = "Closed 2"
    color          = "#339933"
    state_category = "Completed"
  }
}
`, process(processName), name)
}

func workItemTypeWithStatesUpdated(name, processName string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtrackingprocess_workitemtype" "test" {
  name       = "%s"
  process_id = azuredevops_workitemtrackingprocess_process.test.id

  state {
    name           = "New"
    color          = "#3544ca"
    state_category = "Proposed"
  }

  state {
    name           = "Active 2"
    color          = "#020100"
    state_category = "InProgress"
  }

  state {
    name           = "Closed 3"
    color          = "#339933"
    state_category = "Completed"
  }
}
`, process(processName), name)
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
  name       = "%s"
  process_id = azuredevops_workitemtrackingprocess_process.test.id
  is_enabled = false
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
