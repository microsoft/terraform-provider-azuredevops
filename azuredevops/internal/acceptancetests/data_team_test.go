//go:build (all || core || data_sources || data_team) && (!exclude_data_sources || !exclude_data_team)
// +build all core data_sources data_team
// +build !exclude_data_sources !exclude_data_team

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccTeam_DataSource_Basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_team.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testutils.PreCheck(t, nil) },
		Providers:                 testutils.GetProviders(),
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: hclTeamDataSourceBasic(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttrSet(tfNode, "description"),
					resource.TestCheckResourceAttrSet(tfNode, "administrators.#"),
					resource.TestCheckResourceAttrSet(tfNode, "members.#"),
					resource.TestCheckResourceAttrSet(tfNode, "descriptor"),
				),
			},
		},
	})
}

func hclTeamDataSourceBasic(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  description        = "description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

data "azuredevops_team" "test" {
  project_id = azuredevops_project.test.id
  name       = "${azuredevops_project.test.name} Team"
}


`, name)
}
