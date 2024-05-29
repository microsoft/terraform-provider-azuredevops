//go:build (all || data_sources || git || data_git_repository) && (!exclude_data_sources || !exclude_git || !data_git_repository)
// +build all data_sources git data_git_repository
// +build !exclude_data_sources !exclude_git !data_git_repository

package acceptancetests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccGitRepository_DataSource(t *testing.T) {
	name := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_git_repository.repository"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testutils.PreCheck(t, nil) },
		Providers:                 testutils.GetProviders(),
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: hclDataRepository(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttr(tfNode, "name", name),
					resource.TestCheckResourceAttrSet(tfNode, "disabled"),
				),
			},
		},
	})
}

func TestAccGitRepository_DataSource_notExist(t *testing.T) {
	name := testutils.GenerateResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testutils.PreCheck(t, nil) },
		Providers:                 testutils.GetProviders(),
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config:      hclDataRepositoryNotExist(name),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`Repository with name notExist does not exist`)),
			},
		},
	})
}
func hclDataRepository(projectName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%[1]s"
}

data "azuredevops_git_repository" "repository" {
  project_id = azuredevops_project.test.id
  name       = "%[1]s"
}
`, projectName)

}

func hclDataRepositoryNotExist(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%[1]s"
}

data "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "notExist"
}
`, name)

}
