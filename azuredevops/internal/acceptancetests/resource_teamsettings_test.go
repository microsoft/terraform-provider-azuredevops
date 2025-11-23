package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccTeamSettingsResource_projectTeamSettings(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	teamName := testutils.GenerateResourceName()

	tf := "azuredevops_work_team_settings.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclTeamSettingsResourceProjectTeamSettings(projectName, teamName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tf, "project_id"),
					resource.TestCheckResourceAttrSet(tf, "team_id"),
					resource.TestCheckResourceAttrSet(tf, "backlog_visibilities.#"),
					resource.TestCheckResourceAttrSet(tf, "working_days.#"),
					resource.TestCheckResourceAttrSet(tf, "bugs_behavior"),
				),
			},
			{
				ResourceName: tf,
			},
		},
	})
}

func hclTeamSettingsResourceProjectTeamSettings(projectName string, teamName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_team" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s]"
}

resource "azuredevops_work_team_settings" "example" {
  project_id = azuredevops_project.test.id
  team_id    = azuredevops_team.test.id
  
  backlog_visibilities = [
    "Microsoft.EpicCategory",
  ]

  working_days =  [
    "monday",
    "tuesday",
    "wednesday",
    "thursday",
    "friday",
  ]

  "bugs_behavior" = "asRequirements"
}
`, projectName, teamName)
}
