package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccWorkitemtrackingprocessList_Basic(t *testing.T) {
	listName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_list.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckListDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicList(listName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", listName),
					resource.TestCheckResourceAttr(tfNode, "type", "string"),
					resource.TestCheckResourceAttr(tfNode, "is_suggested", "false"),
					resource.TestCheckResourceAttr(tfNode, "items.#", "3"),
					resource.TestCheckResourceAttr(tfNode, "items.0", "Red"),
					resource.TestCheckResourceAttr(tfNode, "items.1", "Green"),
					resource.TestCheckResourceAttr(tfNode, "items.2", "Blue"),
					resource.TestCheckResourceAttrSet(tfNode, "url"),
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

func TestAccWorkitemtrackingprocessList_Update(t *testing.T) {
	listName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_list.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckListDestroyed,
		Steps: []resource.TestStep{
			{
				Config: basicList(listName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", listName),
					resource.TestCheckResourceAttr(tfNode, "items.#", "3"),
					resource.TestCheckResourceAttr(tfNode, "is_suggested", "false"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: updatedList(listName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", listName),
					resource.TestCheckResourceAttr(tfNode, "items.#", "4"),
					resource.TestCheckResourceAttr(tfNode, "items.0", "Red"),
					resource.TestCheckResourceAttr(tfNode, "items.1", "Green"),
					resource.TestCheckResourceAttr(tfNode, "items.2", "Blue"),
					resource.TestCheckResourceAttr(tfNode, "items.3", "Yellow"),
					resource.TestCheckResourceAttr(tfNode, "is_suggested", "true"),
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

func TestAccWorkitemtrackingprocessList_Integer(t *testing.T) {
	listName := testutils.GenerateResourceName()
	tfNode := "azuredevops_workitemtrackingprocess_list.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckListDestroyed,
		Steps: []resource.TestStep{
			{
				Config: integerList(listName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", listName),
					resource.TestCheckResourceAttr(tfNode, "type", "integer"),
					resource.TestCheckResourceAttr(tfNode, "items.#", "3"),
					resource.TestCheckResourceAttr(tfNode, "items.0", "1"),
					resource.TestCheckResourceAttr(tfNode, "items.1", "2"),
					resource.TestCheckResourceAttr(tfNode, "items.2", "3"),
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

func basicList(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtrackingprocess_list" "test" {
  name  = "%s"
  items = ["Red", "Green", "Blue"]
}
`, name)
}

func updatedList(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtrackingprocess_list" "test" {
  name         = "%s"
  items        = ["Red", "Green", "Blue", "Yellow"]
  is_suggested = true
}
`, name)
}

func integerList(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_workitemtrackingprocess_list" "test" {
  name  = "%s"
  type  = "integer"
  items = ["1", "2", "3"]
}
`, name)
}
