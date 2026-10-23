package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// TestAccAreaTree_basic verifies that a simple, single-level tree creates
// exactly the area(s) described and populates area_ids/area_paths.
func TestAccAreaTree_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_area_tree.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclAreaTreeBasic(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "area_ids.%", "1"),
					resource.TestCheckResourceAttr(tfNode, "area_paths.%", "1"),
					resource.TestCheckResourceAttrSet(tfNode, "area_ids.Team A"),
					resource.TestCheckResourceAttr(tfNode, "area_paths.Team A", fmt.Sprintf("%s\\Area\\Team A", projectName)),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateIdFunc: testAccAreaTreeImportStateIdFunc(tfNode),
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccAreaTree_nested verifies that a deeply nested tree creates every
// implied ancestor node automatically.
func TestAccAreaTree_nested(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_area_tree.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclAreaTreeNested(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "area_ids.%", "4"),
					resource.TestCheckResourceAttrSet(tfNode, "area_ids.Team A"),
					resource.TestCheckResourceAttrSet(tfNode, "area_ids.Team A/Sub Area"),
					resource.TestCheckResourceAttrSet(tfNode, "area_ids.Team A/Sub Area/Grandchild"),
					resource.TestCheckResourceAttrSet(tfNode, "area_ids.Team B"),
					resource.TestCheckResourceAttr(tfNode, "area_paths.Team A/Sub Area/Grandchild",
						fmt.Sprintf("%s\\Area\\Team A\\Sub Area\\Grandchild", projectName)),
				),
			},
		},
	})
}

// TestAccAreaTree_update verifies that growing/shrinking the tree between
// applies creates newly added nodes and prunes removed ones, including
// orphaned ancestors.
func TestAccAreaTree_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "azuredevops_area_tree.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclAreaTreeBasic(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "area_ids.%", "1"),
					resource.TestCheckResourceAttrSet(tfNode, "area_ids.Team A"),
				),
			},
			{
				Config: hclAreaTreeNested(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "area_ids.%", "4"),
					resource.TestCheckResourceAttrSet(tfNode, "area_ids.Team A/Sub Area/Grandchild"),
					resource.TestCheckResourceAttrSet(tfNode, "area_ids.Team B"),
				),
			},
			{
				Config: hclAreaTreeBasic(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "area_ids.%", "1"),
					resource.TestCheckResourceAttrSet(tfNode, "area_ids.Team A"),
				),
			},
		},
	})
}

func hclAreaTreeBasic(projectName string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_area_tree" "test" {
  project_id = azuredevops_project.project.id
  paths = jsonencode({
    "Team A" = {}
  })
}
`, testutils.HclProjectResource(projectName))
}

func hclAreaTreeNested(projectName string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_area_tree" "test" {
  project_id = azuredevops_project.project.id
  paths = jsonencode({
    "Team A" = {
      "Sub Area" = {
        "Grandchild" = {}
      }
    }
    "Team B" = {}
  })
}
`, testutils.HclProjectResource(projectName))
}

func testAccAreaTreeImportStateIdFunc(resourceName string) func(s *terraform.State) (string, error) {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("not found: %s", resourceName)
		}
		return rs.Primary.Attributes["project_id"], nil
	}
}
