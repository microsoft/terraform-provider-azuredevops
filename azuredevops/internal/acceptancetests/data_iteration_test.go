//go:build (all || core || data_sources || data_iteration) && (!exclude_data_sources || !exclude_data_iteration)
// +build all core data_sources data_iteration
// +build !exclude_data_sources !exclude_data_iteration

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccIterationDataSource_Read(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config := fmt.Sprintf(`
%s

data "azuredevops_iteration" "root-iteration" {
  project_id = azuredevops_project.project.id
}


`, testutils.HclProjectResource(projectName))

	tfNode := "data.azuredevops_iteration.root-iteration"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "path"),
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttr(tfNode, "has_children", "true"),
					resource.TestCheckResourceAttr(tfNode, "children.#", "3"),
				),
			},
		},
	})
}

func TestAccIterationDataSource_ReadNoChildren(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config := fmt.Sprintf(`
%s

data "azuredevops_iteration" "root-iteration" {
  project_id     = azuredevops_project.project.id
  fetch_children = false
}


`, testutils.HclProjectResource(projectName))

	tfNode := "data.azuredevops_iteration.root-iteration"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "path"),
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttr(tfNode, "has_children", "true"),
					resource.TestCheckResourceAttr(tfNode, "children.#", "0"),
				),
			},
		},
	})
}
