//go:build (all || core || resource_git_repository_branch) && !exclude_resource_git_repository_branch
// +build all core resource_git_repository_branch
// +build !exclude_resource_git_repository_branch

package acceptancetests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// TestAccGitRepoBranch_CreateUpdateDelete verifies that a branch can
// be added to a repository and that it can be replaced
func TestAccGitRepoBranch_CreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	branchName := testutils.GenerateResourceName()
	branchNameChanged := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclGitRepoBranches(projectName, gitRepoName, "Clean", branchName),
				Check: resource.ComposeTestCheckFunc(
					// test-branch
					resource.TestCheckResourceAttr("azuredevops_git_repository_branch.from_master", "name", fmt.Sprintf("testbranch-%s", branchName)),
					resource.TestCheckResourceAttr("azuredevops_git_repository_branch.from_master", "ref", "master"),
					resource.TestCheckResourceAttr("azuredevops_git_repository_branch.from_master", "branch_reference", fmt.Sprintf("refs/heads/testbranch-%s", branchName)),
					resource.TestCheckResourceAttrSet("azuredevops_git_repository_branch.from_master", "branch_head"),
					// test-branch2
					resource.TestCheckResourceAttr("azuredevops_git_repository_branch.from_commit_id", "name", fmt.Sprintf("testbranch2-%s", branchName)),
					resource.TestCheckResourceAttrSet("azuredevops_git_repository_branch.from_commit_id", "commit_id"),
					resource.TestCheckResourceAttr("azuredevops_git_repository_branch.from_commit_id", "branch_reference", fmt.Sprintf("refs/heads/testbranch2-%s", branchName)),
					resource.TestCheckResourceAttrSet("azuredevops_git_repository_branch.from_commit_id", "branch_head"),
				),
			},
			// Test import branch created from ref, ignore fields set only on create
			{
				ResourceName:            "azuredevops_git_repository_branch.from_master",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ref", "tag", "commit_id"},
			},
			// Test replace/update branch when name changes
			{
				Config: hclGitRepoBranches(projectName, gitRepoName, "Clean", branchNameChanged),
				Check: resource.ComposeTestCheckFunc(
					// test-branch
					resource.TestCheckResourceAttr("azuredevops_git_repository_branch.from_master", "name", fmt.Sprintf("testbranch-%s", branchNameChanged)),
					resource.TestCheckResourceAttr("azuredevops_git_repository_branch.from_master", "ref", "master"),
					resource.TestCheckResourceAttr("azuredevops_git_repository_branch.from_master", "branch_reference", fmt.Sprintf("refs/heads/testbranch-%s", branchNameChanged)),
					resource.TestCheckResourceAttrSet("azuredevops_git_repository_branch.from_master", "branch_head"),
					// test-branch2
					resource.TestCheckResourceAttr("azuredevops_git_repository_branch.from_commit_id", "name", fmt.Sprintf("testbranch2-%s", branchNameChanged)),
					resource.TestCheckResourceAttrSet("azuredevops_git_repository_branch.from_commit_id", "commit_id"),
					resource.TestCheckResourceAttr("azuredevops_git_repository_branch.from_commit_id", "branch_reference", fmt.Sprintf("refs/heads/testbranch2-%s", branchNameChanged)),
					resource.TestCheckResourceAttrSet("azuredevops_git_repository_branch.from_commit_id", "branch_head"),
				),
			},
			// Test invalid ref
			{
				Config: fmt.Sprintf(`
%s

resource "azuredevops_git_repository_branch" "from_nonexistent_tag" {
	repository_id = azuredevops_git_repository.repository.id
    name = "testbranch-non-existent-tag"
	tag = "0.0.0"
}
`, hclGitRepoBranches(projectName, gitRepoName, "Clean", branchNameChanged)),
				ExpectError: regexp.MustCompile(`No refs found that match ref "refs/tags/0.0.0"`),
			},
		},
	},
	)
}

func hclGitRepoBranches(projectName, gitRepoName, initType, branchName string) string {
	gitRepoResource := testutils.HclGitRepoResource(projectName, gitRepoName, initType)
	return fmt.Sprintf(`
%[1]s

resource "azuredevops_git_repository_branch" "from_master" {
	repository_id = azuredevops_git_repository.repository.id
	name = "testbranch-%[2]s"
    ref = "master"
}
resource "azuredevops_git_repository_branch" "from_commit_id" {
	repository_id = azuredevops_git_repository.repository.id
    name = "testbranch2-%[2]s"
	commit_id = azuredevops_git_repository_branch.from_master.branch_head
}
  `, gitRepoResource, branchName)
}
