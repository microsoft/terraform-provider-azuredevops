//go:build (all || core || resource_git_repository_branch_lock) && !exclude_resource_git_repository_branch_lock
// +build all core resource_git_repository_branch_lock
// +build !exclude_resource_git_repository_branch_lock

package acceptancetests

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccGitRepoBranchLock_CreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	branchName := testutils.GenerateResourceName()
	branchNameChanged := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			// Test branch lock
			{
				Config: hclGitRepoBranchLock(projectName, gitRepoName, branchName),
				Check: resource.ComposeTestCheckFunc(
					// test-branch
					resource.TestCheckResourceAttr("azuredevops_git_repository_branch.from_master", "name", fmt.Sprintf("testbranch-%s", branchName)),
					resource.TestCheckResourceAttr("azuredevops_git_repository_branch.from_master", "ref_branch", "master"),
					resource.TestCheckResourceAttr("azuredevops_git_repository_branch_lock.master", "is_locked", "true"),
				),
			},
			// Test branch unlock
			{
				Config: hclGitRepoBranchUnlock(projectName, gitRepoName, branchNameChanged),
				Check: resource.ComposeTestCheckFunc(
					// test-branch
					resource.TestCheckResourceAttr("azuredevops_git_repository_branch.from_master", "name", fmt.Sprintf("testbranch-%s", branchNameChanged)),
					resource.TestCheckResourceAttr("azuredevops_git_repository_branch.from_master", "ref_branch", "master"),
					resource.TestCheckResourceAttr("azuredevops_git_repository_branch_lock.master", "is_locked", "false"),
				),
			},
		},
	},
	)
}

func hclGitRepoBranchLock(projectName, gitRepoName, branchName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  description        = "description"
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

resource "azuredevops_git_repository_branch" "from_master" {
  repository_id = azuredevops_git_repository.test.id
  name          = "testbranch-%[3]s"
  ref_branch    = "master"
}

resource "azuredevops_git_repository_branch_lock" "master" {
  repository_id = azuredevops_git_repository.test.id
  branch   = azuredevops_git_repository_branch.from_master.name
  is_locked = true
}
  `, projectName, gitRepoName, branchName)
}

func hclGitRepoBranchUnlock(projectName, gitRepoName, branchName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  description        = "description"
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

resource "azuredevops_git_repository_branch" "from_master" {
  repository_id = azuredevops_git_repository.test.id
  name          = "testbranch-%[3]s"
  ref_branch    = "master"
}

resource "azuredevops_git_repository_branch_lock" "master" {
  repository_id = azuredevops_git_repository.test.id
  branch   = azuredevops_git_repository_branch.from_master.name
  is_locked = false
}
  `, projectName, gitRepoName, branchName)
}
