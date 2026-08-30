package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccGitRefs_DataSource(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	branchName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclGitRefsDataSource(projectName, gitRepoName, branchName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.azuredevops_git_refs.test", "id"),
					resource.TestCheckResourceAttr("data.azuredevops_git_refs.test", "refs.#", "2"), // Should contain default_branch and the newly created one
				),
			},
		},
	})
}

func hclGitRefsDataSource(projectName, gitRepoName, branchName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_git_repository_branch" "test" {
  repository_id = azuredevops_git_repository.test.id
  name          = "%[3]s"
  ref_branch    = azuredevops_git_repository.test.default_branch
}

data "azuredevops_git_refs" "test" {
  repository_id = azuredevops_git_repository.test.id
  depends_on    = [azuredevops_git_repository_branch.test]
}
`, projectName, gitRepoName, branchName)
}
