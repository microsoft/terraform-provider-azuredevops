package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccWorkitemtrackingprocessPage_Basic(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_page.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicPage(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					testutils.TestCheckAttrGreaterThan(tfNode, "section.#", 0),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: pageImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWorkitemtrackingprocessPage_Update(t *testing.T) {
	workItemTypeName := testutils.GenerateWorkItemTypeName()
	processName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_page.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProcessDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicPage(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					testutils.TestCheckAttrGreaterThan(tfNode, "section.#", 0),
				),
			},
			{
				Config: updatedPage(workItemTypeName, processName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					testutils.TestCheckAttrGreaterThan(tfNode, "section.#", 0),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: pageImportStateIdFunc(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func basicPage(workItemTypeName string, processName string) string {
	workItemType := basicWorkItemType(workItemTypeName, processName)
	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtrackingprocess_page" "test" {
  process_id        = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_id = azuredevops_workitemtrackingprocess_workitemtype.test.id
  label             = "Test Page"
}
`, workItemType)
}

func updatedPage(workItemTypeName string, processName string) string {
	workItemType := basicWorkItemType(workItemTypeName, processName)
	return fmt.Sprintf(`
%s

resource "azuredevops_workitemtrackingprocess_page" "test" {
  process_id        = azuredevops_workitemtrackingprocess_process.test.id
  work_item_type_id = azuredevops_workitemtrackingprocess_workitemtype.test.id
  label             = "Updated Page"
  visible           = false
  order             = 4
}
`, workItemType)
}

func pageImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}
		processId := rs.Primary.Attributes["process_id"]
		witRefName := rs.Primary.Attributes["work_item_type_id"]
		pageId := rs.Primary.ID
		return fmt.Sprintf("%s/%s/%s", processId, witRefName, pageId), nil
	}
}
