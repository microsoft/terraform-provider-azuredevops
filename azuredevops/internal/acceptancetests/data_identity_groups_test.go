//go:build (all || core || data_sources || data_groups) && (!exclude_data_sources || !exclude_data_groups)
// +build all core data_sources data_groups
// +build !exclude_data_sources !exclude_data_groups

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func generateIdentityGroupsDataSourceConfig(projectName string) string {
	if projectName == "" {
		return `
data "azuredevops_identity_groups" "groups" {
}`
	}

	dataSource := `
data "azuredevops_project" "project" {
	name = "Default"
}

data "azuredevops_identity_groups" "groups" {
	project_id = azuredevops_project.project.id
}`

	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s", projectResource, dataSource)
}

func testIdentityGroupsDataSource(t *testing.T, projectName string) {
	tfNode := "data.azuredevops_identity_groups.groups"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: generateIdentityGroupsDataSourceConfig(projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "groups.#"),
				),
			},
		},
	})
}

func TestAccIdentityGroupsDataSource_ReadProject(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	testIdentityGroupsDataSource(t, projectName)
}

func TestAccIdentityGroupsDataSource_ReadNoProject(t *testing.T) {
	testIdentityGroupsDataSource(t, "")
}
