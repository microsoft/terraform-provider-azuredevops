//go:build (all || core || data_sources || data_teams) && (!exclude_data_sources || !exclude_data_teams)
// +build all core data_sources data_teams
// +build !exclude_data_sources !exclude_data_teams

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccTeams_DataSource(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	projectResource := testutils.HclProjectResource(projectName)

	config := fmt.Sprintf(`

%s

data "azuredevops_teams" "all_teams" {
	project_id = azuredevops_project.project.id
}

	`, projectResource)

	tfNode := "data.azuredevops_teams.all_teams"
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testutils.PreCheck(t, nil) },
		Providers:                 testutils.GetProviders(),
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: config,
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
