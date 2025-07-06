//go:build (all || core || data_sources || data_teams) && (!exclude_data_sources || !exclude_data_teams)

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccTeams_DataSource_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_teams.test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclTeamsDataSourceBasic(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttr(tfNode, "teams.#", "1"),
					resource.TestCheckResourceAttrSet(tfNode, "teams.0.project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "teams.0.id"),
					resource.TestCheckResourceAttrSet(tfNode, "teams.0.name"),
					resource.TestCheckResourceAttrSet(tfNode, "teams.0.description"),
					resource.TestCheckResourceAttrSet(tfNode, "teams.0.administrators.#"),
					resource.TestCheckResourceAttrSet(tfNode, "teams.0.members.#"),
				),
			},
		},
	})
}

func hclTeamsDataSourceBasic(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  description        = "description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

data "azuredevops_teams" "test" {
  project_id = azuredevops_project.test.id
}
`, name)
}
