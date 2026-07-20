package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccGitRef_Basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	branchName := "refs/heads/" + testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclGitRefResource(projectName, gitRepoName, branchName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("azuredevops_git_ref.test", "id"),
					resource.TestCheckResourceAttr("azuredevops_git_ref.test", "name", branchName),
					resource.TestCheckResourceAttrSet("azuredevops_git_ref.test", "object_id"),
				),
			},
		},
	})
}

func hclGitRefResource(projectName, gitRepoName, branchName string) string {
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

resource "azuredevops_git_ref" "test" {
  repository_id = azuredevops_git_repository.test.id
  name          = "%[3]s"
  ref_branch    = azuredevops_git_repository.test.default_branch
}
`, projectName, gitRepoName, branchName)
}
