package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccTeam_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	teamName := testutils.GenerateResourceName()

	tfNode := "azuredevops_team.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclTeamBasic(projectName, teamName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", teamName),
				),
			},
		},
	})
}

func TestAccTeam_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	teamName := testutils.GenerateResourceName()
	teamDescription := "description"
	teamName2 := testutils.GenerateResourceName()
	teamDescription2 := "@@descriptionUpdate"

	tfNode := "azuredevops_team.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclTeamUpdate(projectName, teamName, teamDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", teamName),
					resource.TestCheckResourceAttr(tfNode, "description", teamDescription),
					resource.TestCheckResourceAttrSet(tfNode, "administrators.#"),
					resource.TestCheckResourceAttrSet(tfNode, "members.#"),
					resource.TestCheckResourceAttrSet(tfNode, "descriptor"),
				),
			},
			{
				Config: hclTeamUpdate(projectName, teamName2, teamDescription2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", teamName2),
					resource.TestCheckResourceAttr(tfNode, "description", teamDescription2),
					resource.TestCheckResourceAttrSet(tfNode, "administrators.#"),
					resource.TestCheckResourceAttrSet(tfNode, "descriptor"),
				),
			},
		},
	})
}

func TestAccTeam_membersUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	teamName := testutils.GenerateResourceName()
	tfNode := "azuredevops_team.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclTeamMembersBasic(projectName, teamName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", teamName),
					resource.TestCheckResourceAttr(tfNode, "members.#", "1"),
				),
			},
			{
				Config: hclTeamMembersUpdate(projectName, teamName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", teamName),
					resource.TestCheckResourceAttr(tfNode, "members.#", "2"),
				),
			},
			{
				Config: hclTeamMembersBasic(projectName, teamName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", teamName),
					resource.TestCheckResourceAttr(tfNode, "members.#", "1"),
				),
			},
		},
	})
}

func TestAccTeam_administratorsUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	teamName := testutils.GenerateResourceName()

	tfNode := "azuredevops_team.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclTeamAdministratorsBasic(projectName, teamName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", teamName),
					resource.TestCheckResourceAttrSet(tfNode, "administrators.#"),
					resource.TestCheckResourceAttrSet(tfNode, "descriptor"),
				),
			},
			{
				Config: hclTeamAdministratorsUpdate(projectName, teamName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", teamName),
					resource.TestCheckResourceAttr(tfNode, "administrators.#", "2"),
				),
			},
			{
				Config: hclTeamAdministratorsBasic(projectName, teamName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", teamName),
					resource.TestCheckResourceAttr(tfNode, "administrators.#", "1"),
				),
			},
		},
	})
}

func TestAccTeam_complete(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	teamName := testutils.GenerateResourceName()

	tfNode := "azuredevops_team.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclTeamComplete(projectName, teamName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", teamName),
					resource.TestCheckResourceAttr(tfNode, "administrators.#", "1"),
					resource.TestCheckResourceAttr(tfNode, "members.#", "1"),
					resource.TestCheckResourceAttrSet(tfNode, "descriptor"),
				),
			},
		},
	})
}

func hclTeamBasic(projectName, teamName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_team" "test" {
  project_id = azuredevops_project.test.id
  name       = "%s"
}
`, projectName, teamName)
}

func hclTeamUpdate(projectName, teamName, description string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_team" "test" {
  project_id  = azuredevops_project.test.id
  name        = "%s"
  description = "%s"
}
`, projectName, teamName, description)
}

func hclTeamMembersBasic(projectName, teamName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

data "azuredevops_group" "test" {
  project_id = azuredevops_project.test.id
  name       = "Contributors"
}

resource "azuredevops_team" "test" {
  project_id  = azuredevops_project.test.id
  name        = "%s"
  description = "Test sTeam"

  members = [
    data.azuredevops_group.test.descriptor
  ]
}
`, projectName, teamName)
}

func hclTeamMembersUpdate(projectName, teamName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

data "azuredevops_group" "test" {
  project_id = azuredevops_project.test.id
  name       = "Contributors"
}

data "azuredevops_group" "test2" {
  project_id = azuredevops_project.test.id
  name       = "Readers"
}

resource "azuredevops_team" "test" {
  project_id  = azuredevops_project.test.id
  name        = "%s"
  description = "Test sTeam"

  members = [
    data.azuredevops_group.test.descriptor,
    data.azuredevops_group.test2.descriptor
  ]
}
`, projectName, teamName)
}

func hclTeamAdministratorsBasic(projectName, teamName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

data "azuredevops_group" "test" {
  project_id = azuredevops_project.test.id
  name       = "Contributors"
}

resource "azuredevops_team" "test" {
  project_id  = azuredevops_project.test.id
  name        = "%s"
  description = "Test sTeam"

  administrators = [
    data.azuredevops_group.test.descriptor
  ]
}
`, projectName, teamName)
}

func hclTeamAdministratorsUpdate(projectName, teamName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

data "azuredevops_group" "test" {
  project_id = azuredevops_project.test.id
  name       = "Contributors"
}

data "azuredevops_group" "test2" {
  project_id = azuredevops_project.test.id
  name       = "Readers"
}

resource "azuredevops_team" "test" {
  project_id  = azuredevops_project.test.id
  name        = "%s"
  description = "Test sTeam"

  administrators = [
    data.azuredevops_group.test.descriptor,
    data.azuredevops_group.test2.descriptor
  ]
}
`, projectName, teamName)
}

func hclTeamComplete(projectName, teamName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

data "azuredevops_group" "test" {
  project_id = azuredevops_project.test.id
  name       = "Contributors"
}

data "azuredevops_group" "test2" {
  project_id = azuredevops_project.test.id
  name       = "Readers"
}

resource "azuredevops_team" "test" {
  project_id  = azuredevops_project.test.id
  name        = "%s"
  description = "Test sTeam"

  administrators = [
    data.azuredevops_group.test.descriptor,
  ]

  members = [
    data.azuredevops_group.test2.descriptor
  ]
}`, projectName, teamName)
}
