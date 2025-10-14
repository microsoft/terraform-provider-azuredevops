package acceptancetests

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccGitRepositoryFile_DataSource(t *testing.T) {
	tfNode := "data.azuredevops_git_repository_file.test"

	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()
	branch := "refs/heads/master"
	file := "foo.txt"
	content := "bar"
	commitMessage := "first commit"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testutils.PreCheck(t, nil) },
		Providers:                 testutils.GetProviders(),
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: hclDataRepositoryFile(projectName, repoName, branch, file, content, commitMessage, file),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "content", content),
					resource.TestCheckResourceAttr(tfNode, "last_commit_message", commitMessage),
				),
			},
		},
	})
}

func TestAccGitRepositoryFile_DataSource_notExist(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()
	branch := "refs/heads/master"
	file := "foo.txt"
	content := "bar"
	commitMessage := "first commit"
	not_exists_file := "not_exists.txt"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testutils.PreCheck(t, nil) },
		Providers:                 testutils.GetProviders(),
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config:      hclDataRepositoryFile(projectName, repoName, branch, file, content, commitMessage, not_exists_file),
				ExpectError: regexp.MustCompile(fmt.Sprintf("Error: Item not found, repositoryID: [A-Za-z0-9-]+, branch: %s, file: %s", regexp.QuoteMeta(strings.Split(branch, "/")[2]), regexp.QuoteMeta(not_exists_file))),
			},
		},
	})
}

func hclDataRepositoryFile(projectName, repoName, branch, rfile, content, commitMessage, dfile string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%[1]s"
}

resource "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%[2]s"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_git_repository_file" "test" {
  repository_id  = azuredevops_git_repository.test.id
  branch         = "%[3]s"
  file           = "%[4]s"
  content        = "%[5]s"
  commit_message = "%[6]s"
}

data "azuredevops_git_repository_file" "test" {
  repository_id = azuredevops_git_repository.test.id
  branch        = "%[3]s"
  file          = "%[7]s"
  depends_on    = [azuredevops_git_repository_file.test]
}


`, projectName, repoName, branch, rfile, content, commitMessage, dfile)
}
