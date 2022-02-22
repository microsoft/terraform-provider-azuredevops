//go:build (all || core || data_sources || data_area) && (!exclude_data_sources || !exclude_data_area)
// +build all core data_sources data_area
// +build !exclude_data_sources !exclude_data_area

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccAreaDataSource_Read(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config := fmt.Sprintf(`
%s

data "azuredevops_area" "root-area" {
	project_id     = azuredevops_project.project.id
}

`, testutils.HclProjectResource(projectName))

	tfNode := "data.azuredevops_area.root-area"
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
					resource.TestCheckResourceAttr(tfNode, "has_children", "false"),
					resource.TestCheckResourceAttr(tfNode, "children.#", "0"),
				),
			},
		},
	})
}

func TestAccAreaDataSource_ReadNoChildren(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	config := fmt.Sprintf(`
%s

data "azuredevops_area" "root-area" {
	project_id     = azuredevops_project.project.id
	fetch_children = false
}

`, testutils.HclProjectResource(projectName))

	tfNode := "data.azuredevops_area.root-area"
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
					resource.TestCheckResourceAttr(tfNode, "has_children", "false"),
					resource.TestCheckResourceAttr(tfNode, "children.#", "0"),
				),
			},
		},
	})
}
