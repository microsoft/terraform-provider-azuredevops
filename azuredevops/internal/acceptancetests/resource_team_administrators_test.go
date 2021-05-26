// +build all core resource_team_administrators
// +build !exclude_resource_team_administrators

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccTeamAdministrators_CreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	teamName := testutils.GenerateResourceName()

	config1 := fmt.Sprintf(`

%s

data "azuredevops_group" "builtin_project_contributors" {
	project_id = azuredevops_project.project.id
	name       = "Contributors"
}

resource "azuredevops_team_administrators" "team_administrators" {
	project_id = azuredevops_team.team.project_id
	team_id = azuredevops_team.team.id
	administrators = [
	  data.azuredevops_group.builtin_project_contributors.descriptor
	]
}


	`, testutils.HclTeamConfiguration(projectName, teamName, "", nil, nil))

	config2 := fmt.Sprintf(`

%s

data "azuredevops_group" "builtin_project_contributors" {
	project_id = azuredevops_project.project.id
	name       = "Contributors"
}

data "azuredevops_group" "builtin_project_readers" {
	project_id = azuredevops_project.project.id
	name       = "Readers"
}

resource "azuredevops_team_administrators" "team_administrators" {
	project_id = azuredevops_team.team.project_id
	team_id = azuredevops_team.team.id
	administrators = [
	  data.azuredevops_group.builtin_project_contributors.descriptor,
	  data.azuredevops_group.builtin_project_readers.descriptor
	]
}

		`, testutils.HclTeamConfiguration(projectName, teamName, "", nil, nil))

	tfNode := "azuredevops_team_administrators.team_administrators"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "team_id"),
					resource.TestCheckResourceAttr(tfNode, "administrators.#", "1"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "team_id"),
					resource.TestCheckResourceAttr(tfNode, "administrators.#", "2"),
				),
			},
		},
	})
}

func TestAccTeamAdministrators_CreateAndUpdate_Overwrite(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	teamName := testutils.GenerateResourceName()

	config1 := fmt.Sprintf(`

%s

data "azuredevops_group" "builtin_project_contributors" {
	project_id = azuredevops_project.project.id
	name       = "Contributors"
}

resource "azuredevops_team_administrators" "team_administrators" {
	project_id = azuredevops_team.team.project_id
	team_id = azuredevops_team.team.id
	mode = "overwrite"
	administrators = [
	  data.azuredevops_group.builtin_project_contributors.descriptor
	]
}


	`, testutils.HclTeamConfiguration(projectName, teamName, "", nil, nil))

	config2 := fmt.Sprintf(`

%s

data "azuredevops_group" "builtin_project_contributors" {
	project_id = azuredevops_project.project.id
	name       = "Contributors"
}

data "azuredevops_group" "builtin_project_readers" {
	project_id = azuredevops_project.project.id
	name       = "Readers"
}

resource "azuredevops_team_administrators" "team_administrators" {
	project_id = azuredevops_team.team.project_id
	team_id = azuredevops_team.team.id
	mode = "overwrite"
	administrators = [
	  data.azuredevops_group.builtin_project_contributors.descriptor,
	  data.azuredevops_group.builtin_project_readers.descriptor
	]
}

		`, testutils.HclTeamConfiguration(projectName, teamName, "", nil, nil))

	tfNode := "azuredevops_team_administrators.team_administrators"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: config1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "team_id"),
					resource.TestCheckResourceAttr(tfNode, "administrators.#", "1"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "team_id"),
					resource.TestCheckResourceAttr(tfNode, "administrators.#", "2"),
				),
			},
		},
	})
}
