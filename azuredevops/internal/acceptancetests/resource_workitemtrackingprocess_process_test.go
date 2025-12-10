package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccWorkitemtrackingprocessProcess_Basic(t *testing.T) {
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_process.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: process(processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", processName),
					resource.TestCheckResourceAttrSet(tfNode, "reference_name"),
					resource.TestCheckResourceAttr(tfNode, "is_default", "false"),
					resource.TestCheckResourceAttr(tfNode, "is_enabled", "true"),
					resource.TestCheckResourceAttrSet(tfNode, "customization_type"),
					resource.TestCheckResourceAttr(tfNode, "parent_process_type_id", "adcc42ab-9882-485e-a3ed-7678f01f66bc"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getProcessStateIdFunc(tfNode),
			},
		},
	})
}

func TestAccWorkitemtrackingprocessProcess_CreateDisabled(t *testing.T) {
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_process.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: disabledProcess(processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", processName),
					resource.TestCheckResourceAttr(tfNode, "is_enabled", "false"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getProcessStateIdFunc(tfNode),
			},
		},
	})
}

func TestAccWorkitemtrackingprocessProcess_CreateAndUpdate(t *testing.T) {
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_process.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: process(processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", processName),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getProcessStateIdFunc(tfNode),
			},
			{
				Config: disabledProcess(processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "is_enabled", "false"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getProcessStateIdFunc(tfNode),
			},
		},
	})
}

func process(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtrackingprocess_process" "test" {
  name                   = "%s"
  parent_process_type_id = "%s"
}
`, name, agileSystemProcessTypeId)
}

func disabledProcess(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtrackingprocess_process" "test" {
  name                   = "%s"
  parent_process_type_id = "%s"
  is_enabled             = false
}
`, name, agileSystemProcessTypeId)
}

func getProcessStateIdFunc(tfNode string) resource.ImportStateIdFunc {
	return func(state *terraform.State) (string, error) {
		res := state.RootModule().Resources[tfNode]
		return res.Primary.Attributes["id"], nil
	}
}

// Sourced from https://learn.microsoft.com/en-us/rest/api/azure/devops/processes/processes/list?view=azure-devops-rest-7.1&tabs=HTTP#get-the-list-of-processes
const (
	agileSystemProcessTypeId string = "adcc42ab-9882-485e-a3ed-7678f01f66bc"
)
