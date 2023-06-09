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

func TestAccTeam_DataSource(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	projectResource := testutils.HclProjectResource(projectName)

	config := fmt.Sprintf(`

%s

data "azuredevops_team" "team" {
	project_id = azuredevops_project.project.id
	name = "${azuredevops_project.project.name} Team"
}

	`, projectResource)

	tfNode := "data.azuredevops_team.team"
	resource.Test(t, resource.TestCase{
		PreCheck:                  func() { testutils.PreCheck(t, nil) },
		Providers:                 testutils.GetProviders(),
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: config,
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
