//go:build (all || core || resource_git_repository_branch) && !exclude_resource_git_repository_branch
// +build all core resource_git_repository_branch
// +build !exclude_resource_git_repository_branch

package acceptancetests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/tfhelper"
)

func TestAccGitRepoBranch_fromBranch(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	branchName := testutils.GenerateResourceName()
	resNode := "azuredevops_git_repository_branch.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclGitRepoBranchesFromBranch(projectName, gitRepoName, branchName),
				Check: resource.ComposeTestCheckFunc(
					checkRepositoryBranchExist(branchName),
					resource.TestCheckResourceAttr(resNode, "name", branchName),
					resource.TestCheckResourceAttr(resNode, "ref_branch", "master"),
					resource.TestCheckResourceAttrSet(resNode, "last_commit_id"),
				),
			},
			{
				ResourceName:            resNode,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       hclRepositoryBranchID,
				ImportStateVerifyIgnore: []string{"ref_branch"},
			},
		},
	},
	)
}

func TestAccGitRepoBranch_fromCommit(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	branchName := testutils.GenerateResourceName()
	resNode := "azuredevops_git_repository_branch.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclGitRepoBranchesFromCommit(projectName, gitRepoName, branchName),
				Check: resource.ComposeTestCheckFunc(
					checkRepositoryBranchExist(branchName),
					resource.TestCheckResourceAttr(resNode, "name", branchName),
					resource.TestCheckResourceAttrSet(resNode, "ref_commit_id"),
					resource.TestCheckResourceAttrSet(resNode, "last_commit_id"),
				),
			},
			{
				ResourceName:            resNode,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateIdFunc:       hclRepositoryBranchID,
				ImportStateVerifyIgnore: []string{"ref_commit_id"},
			},
		},
	},
	)
}

func TestAccGitRepoBranch_invalidRef(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	branchName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      testutils.CheckProjectDestroyed,
		Steps: []resource.TestStep{
			{
				Config:      hclGitRepoBranchInvalidRef(projectName, gitRepoName, branchName),
				ExpectError: regexp.MustCompile(`No refs found that match ref "refs/tags/0.0.0"`),
			},
		},
	},
	)
}

func TestAccGitRepoBranch_requireImportError(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	branchName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		CheckDestroy:      testutils.CheckProjectDestroyed,
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      hclGitRepoBranchesImportError(projectName, gitRepoName, branchName),
				ExpectError: regexp.MustCompile(`Update refs failed. Update Status: staleOldObjectId`),
			},
		},
	},
	)
}

func hclRepositoryBranchID(state *terraform.State) (string, error) {
	res := state.RootModule().Resources["azuredevops_git_repository_branch.test"]
	repositoryName := res.Primary.Attributes["repository_id"]
	name := res.Primary.Attributes["name"]
	return fmt.Sprintf("%s:%s", repositoryName, name), nil
}

func checkRepositoryBranchExist(expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		res, ok := s.RootModule().Resources["azuredevops_git_repository_branch.test"]
		if !ok {
			return fmt.Errorf(" Did not find `azuredevops_git_repository_branch` in the TF state")
		}

		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
		repoId, branchName, err := tfhelper.ParseGitRepoBranchID(res.Primary.ID)
		if err != nil {
			return fmt.Errorf(" Parse resource IDs: %w", err)
		}

		branch, err := clients.GitReposClient.GetBranch(clients.Ctx, git.GetBranchArgs{
			RepositoryId: &repoId,
			Name:         &branchName,
		})

		if err != nil {
			return fmt.Errorf(" Repositroy: %s, Branch: %s cannot be found. Error=%v", repoId, branchName, err)
		}

		if *branch.Name != expectedName {
			return fmt.Errorf(" Branch Name=%s, but expected Name=%s", *branch.Name, expectedName)
		}
		return nil
	}
}

func hclGitRepoBranchesFromBranch(projectName, gitRepoName, branchName string) string {
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

resource "azuredevops_git_repository_branch" "test" {
  repository_id = azuredevops_git_repository.test.id
  name          = "%[3]s"
  ref_branch    = "master"
}`, projectName, gitRepoName, branchName)
}

func hclGitRepoBranchesFromCommit(projectName, gitRepoName, branchName string) string {
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

resource "azuredevops_git_repository_branch" "test" {
  repository_id = azuredevops_git_repository.test.id
  name          = "%[3]s"
  ref_commit_id = azuredevops_git_repository_branch.from_master.last_commit_id
}`, projectName, gitRepoName, branchName)
}

func hclGitRepoBranchInvalidRef(projectName, gitRepoName, branchName string) string {
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

resource "azuredevops_git_repository_branch" "from_commit_id" {
  repository_id = azuredevops_git_repository.test.id
  name          = "testbranch2-%[3]s"
  ref_commit_id = azuredevops_git_repository_branch.from_master.last_commit_id
}

resource "azuredevops_git_repository_branch" "from_nonexistent_tag" {
  repository_id = azuredevops_git_repository.test.id
  name          = "testbranch-non-existent-tag"
  ref_tag       = "0.0.0"
}`, projectName, gitRepoName, branchName)
}

func hclGitRepoBranchesImportError(projectName, gitRepoName, branchName string) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_git_repository_branch" "import" {
  repository_id = azuredevops_git_repository_branch.test.repository_id
  name          = azuredevops_git_repository_branch.test.name
  ref_branch    = azuredevops_git_repository_branch.test.ref_branch
}`, hclGitRepoBranchesFromBranch(projectName, gitRepoName, branchName))
}
