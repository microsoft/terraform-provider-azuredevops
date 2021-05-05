// +build all core resource_team
// +build !exclude_resource_team

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccTeam_CreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	teamName := testutils.GenerateResourceName()
	teamDescription := "@@TEAMDESCRIPTION@@1"
	projectResource := testutils.HclProjectResource(projectName)

	config1 := fmt.Sprintf(`

%s

resource "azuredevops_team" "team" {
	project_id = azuredevops_project.project.id
	name = "%s"
	description = "%s"
}

	`, projectResource, teamName, teamDescription)

	teamName2 := testutils.GenerateResourceName()
	teamDescription2 := "@@TEAMDESCRIPTION@@2"
	config2 := fmt.Sprintf(`

	%s

	resource "azuredevops_team" "team" {
		project_id = azuredevops_project.project.id
		name = "%s"
		description = "%s"
	}

		`, projectResource, teamName2, teamDescription2)

	tfNode := "azuredevops_team.team"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "name", teamName),
					resource.TestCheckResourceAttr(tfNode, "description", teamDescription),
					resource.TestCheckResourceAttrSet(tfNode, "administrators.#"),
					resource.TestCheckResourceAttrSet(tfNode, "members.#"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "name", teamName2),
					resource.TestCheckResourceAttr(tfNode, "description", teamDescription2),
					resource.TestCheckResourceAttrSet(tfNode, "administrators.#"),
					resource.TestCheckResourceAttrSet(tfNode, "members.#"),
				),
			},
		},
	})
}
