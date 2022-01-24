//go:build (all || core || resource_team) && !exclude_resource_team
// +build all core resource_team
// +build !exclude_resource_team

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccTeam_CreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	teamName := testutils.GenerateResourceName()
	teamDescription := "@@TEAMDESCRIPTION@@1"
	config1 := testutils.HclTeamConfiguration(projectName, teamName, teamDescription, nil, nil)

	teamName2 := testutils.GenerateResourceName()
	teamDescription2 := "@@TEAMDESCRIPTION@@2"
	config2 := testutils.HclTeamConfiguration(projectName, teamName2, teamDescription2, nil, nil)

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

func TestAccTeam_CreateAndUpdateAdministrators(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	teamName := testutils.GenerateResourceName()
	teamDescription := "@@TEAMDESCRIPTION@@1"
	config1 := fmt.Sprintf(`
%s

data "azuredevops_group" "builtin_project_contributors" {
	project_id = azuredevops_project.project.id
	name       = "Contributors"
}
`, testutils.HclTeamConfiguration(projectName, teamName, teamDescription, &[]string{
		"data.azuredevops_group.builtin_project_contributors.descriptor",
	}, nil))

	config2 := fmt.Sprintf(`
%s

data "azuredevops_group" "builtin_project_readers" {
	project_id = azuredevops_project.project.id
	name       = "Readers"
}

data "azuredevops_group" "builtin_project_contributors" {
	project_id = azuredevops_project.project.id
	name       = "Contributors"
}
`, testutils.HclTeamConfiguration(projectName, teamName, teamDescription, &[]string{
		"data.azuredevops_group.builtin_project_readers.descriptor",
		"data.azuredevops_group.builtin_project_contributors.descriptor",
	}, nil))

	config3 := fmt.Sprintf(`
%s

data "azuredevops_group" "builtin_project_readers" {
	project_id = azuredevops_project.project.id
	name       = "Readers"
}

`, testutils.HclTeamConfiguration(projectName, teamName, teamDescription, &[]string{
		"data.azuredevops_group.builtin_project_readers.descriptor",
	}, nil))

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
					resource.TestCheckResourceAttr(tfNode, "administrators.#", "1"),
					resource.TestCheckResourceAttrSet(tfNode, "members.#"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "name", teamName),
					resource.TestCheckResourceAttr(tfNode, "description", teamDescription),
					resource.TestCheckResourceAttr(tfNode, "administrators.#", "2"),
					resource.TestCheckResourceAttrSet(tfNode, "members.#"),
				),
			},
			{
				Config: config3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "name", teamName),
					resource.TestCheckResourceAttr(tfNode, "description", teamDescription),
					resource.TestCheckResourceAttr(tfNode, "administrators.#", "1"),
					resource.TestCheckResourceAttrSet(tfNode, "members.#"),
				),
			},
		},
	})
}

func TestAccTeam_CreateAndUpdateMembers(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	teamName := testutils.GenerateResourceName()
	teamDescription := "@@TEAMDESCRIPTION@@1"
	config1 := fmt.Sprintf(`
%s
data "azuredevops_group" "builtin_project_contributors" {
	project_id = azuredevops_project.project.id
	name       = "Contributors"
}
`, testutils.HclTeamConfiguration(projectName, teamName, teamDescription, nil, &[]string{
		"data.azuredevops_group.builtin_project_contributors.descriptor",
	}))

	config2 := fmt.Sprintf(`
%s

data "azuredevops_group" "builtin_project_readers" {
	project_id = azuredevops_project.project.id
	name       = "Readers"
}

data "azuredevops_group" "builtin_project_contributors" {
	project_id = azuredevops_project.project.id
	name       = "Contributors"
}
`, testutils.HclTeamConfiguration(projectName, teamName, teamDescription, nil, &[]string{
		"data.azuredevops_group.builtin_project_readers.descriptor",
		"data.azuredevops_group.builtin_project_contributors.descriptor",
	}))

	config3 := fmt.Sprintf(`
%s

data "azuredevops_group" "builtin_project_readers" {
	project_id = azuredevops_project.project.id
	name       = "Readers"
}

`, testutils.HclTeamConfiguration(projectName, teamName, teamDescription, nil, &[]string{
		"data.azuredevops_group.builtin_project_readers.descriptor",
	}))

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
					resource.TestCheckResourceAttr(tfNode, "members.#", "1"),
				),
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "name", teamName),
					resource.TestCheckResourceAttr(tfNode, "description", teamDescription),
					resource.TestCheckResourceAttrSet(tfNode, "administrators.#"),
					resource.TestCheckResourceAttr(tfNode, "members.#", "2"),
				),
			},
			{
				Config: config3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "name", teamName),
					resource.TestCheckResourceAttr(tfNode, "description", teamDescription),
					resource.TestCheckResourceAttrSet(tfNode, "administrators.#"),
					resource.TestCheckResourceAttr(tfNode, "members.#", "1"),
				),
			},
		},
	})
}
