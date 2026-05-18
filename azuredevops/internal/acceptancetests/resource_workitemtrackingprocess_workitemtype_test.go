package acceptancetests

import (
	"fmt"
	"regexp"
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
					resource.TestCheckResourceAttr(tfNode, "state.0.name", "New Active"),
					resource.TestCheckResourceAttr(tfNode, "state.0.order", "1"),
					resource.TestCheckResourceAttr(tfNode, "state.1.name", "New Closed"),
					resource.TestCheckResourceAttr(tfNode, "state.1.order", "2"),
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
					resource.TestCheckResourceAttr(tfNode, "state.#", "4"),
					resource.TestCheckResourceAttr(tfNode, "state.0.name", "New"),
					resource.TestCheckResourceAttr(tfNode, "state.0.order", "1"),
					resource.TestCheckResourceAttr(tfNode, "state.1.name", "New Active"),
					resource.TestCheckResourceAttr(tfNode, "state.1.order", "2"),
					resource.TestCheckResourceAttr(tfNode, "state.2.name", "Active Order"),
					resource.TestCheckResourceAttr(tfNode, "state.2.order", "3"),
					resource.TestCheckResourceAttr(tfNode, "state.3.name", "New Closed"),
					resource.TestCheckResourceAttr(tfNode, "state.3.order", "4"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getWorkItemTypeStateIdFunc(tfNode),
			},
			{
				Config: workItemTypeWithStatesChangingOrder(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "state.#", "4"),
					resource.TestCheckResourceAttr(tfNode, "state.0.name", "New"),
					resource.TestCheckResourceAttr(tfNode, "state.0.order", "1"),
					resource.TestCheckResourceAttr(tfNode, "state.1.name", "Active Order"),
					resource.TestCheckResourceAttr(tfNode, "state.1.order", "2"),
					resource.TestCheckResourceAttr(tfNode, "state.2.name", "New Active"),
					resource.TestCheckResourceAttr(tfNode, "state.2.order", "3"),
					resource.TestCheckResourceAttr(tfNode, "state.3.name", "New Closed"),
					resource.TestCheckResourceAttr(tfNode, "state.3.order", "4"),
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

func TestAccWorkitemtrackingprocessWorkItemType_StatesForbiddenOnInherited(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config:      workItemTypeWithStatesAndParent(workItemTypeName, processName),
				ExpectError: regexp.MustCompile(`state.*blocks are only valid on non-inherited work item types`),
			},
		},
	})
}

func workItemTypeWithStatesAndParent(name, processName string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtrackingprocess_workitemtype" "test" {
  name                            = "%s"
  process_id                      = azuredevops_workitemtrackingprocess_process.test.id
  parent_work_item_reference_name = "Microsoft.VSTS.WorkItemTypes.Bug"

  state {
    name           = "Custom"
    color          = "#3544ca"
    state_category = "Proposed"
  }
}
`, process(processName), name)
}

func TestAccWorkitemtrackingprocessWorkItemType_StatesRemovedFromConfig(t *testing.T) {
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
			{
				Config: basicWorkItemType(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "reference_name"),
				),
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
    name           = "New Active"
    color          = "#ff9d01"
    state_category = "InProgress"
  }

  state {
    name           = "New Closed"
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
    name           = "New Active"
    color          = "#ff9d01"
    state_category = "InProgress"
  }

  state {
    name           = "Active Order"
    color          = "#020100"
    state_category = "InProgress"
  }

  state {
    name           = "New Closed"
    color          = "#339933"
    state_category = "Completed"
  }
}
`, process(processName), name)
}

func workItemTypeWithStatesChangingOrder(name, processName string) string {
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
    name           = "Active Order"
    color          = "#020100"
    state_category = "InProgress"
  }

  state {
    name           = "New Active"
    color          = "#ff9d01"
    state_category = "InProgress"
  }

  state {
    name           = "New Closed"
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
