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
					resource.TestCheckResourceAttrSet(tfNode, "id"),
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

// NOTE! This test might be flaky due to eventual consistent reads after update during import/refresh.
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
					resource.TestCheckResourceAttrSet(tfNode, "id"),
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
