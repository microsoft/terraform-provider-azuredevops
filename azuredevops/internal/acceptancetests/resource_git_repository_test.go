//go:build (all || core || resource_git_repository) && !exclude_resource_git_repository
// +build all core resource_git_repository
// +build !exclude_resource_git_repository

package acceptancetests

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/git"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// Verifies that the following sequence of events occurrs without error:
//
//	(1) TF apply creates resource
//	(2) TF state values are set
//	(3) resource can be queried by ID and has expected name
//	(4) TF destroy deletes resource
//	(5) resource can no longer be queried by ID
func TestAccGitRepo_CreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoNameFirst := testutils.GenerateResourceName()
	gitRepoNameSecond := testutils.GenerateResourceName()
	tfRepoNode := "azuredevops_git_repository.repository"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkGitRepoDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclGitRepoResource(projectName, gitRepoNameFirst, "Uninitialized"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfRepoNode, "project_id"),
					resource.TestCheckResourceAttr(tfRepoNode, "name", gitRepoNameFirst),
					checkGitRepoExists(gitRepoNameFirst),
					resource.TestCheckResourceAttrSet(tfRepoNode, "is_fork"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "remote_url"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "size"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "ssh_url"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "url"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "web_url"),
				),
			},
			{
				Config: testutils.HclGitRepoResource(projectName, gitRepoNameSecond, "Uninitialized"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfRepoNode, "project_id"),
					resource.TestCheckResourceAttr(tfRepoNode, "name", gitRepoNameSecond),
					checkGitRepoExists(gitRepoNameSecond),
					resource.TestCheckResourceAttrSet(tfRepoNode, "is_fork"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "remote_url"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "size"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "ssh_url"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "url"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "web_url"),
				),
			},
		},
	})
}

// Verifies that the create operation fails if the initialization is
// not specified.
func TestAccGitRepo_Create_IncorrectInitialization(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	azureGitRepoResource := fmt.Sprintf(`
	resource "azuredevops_git_repository" "repository" {
		project_id      = azuredevops_project.project.id
		name            = "%s"
	}`, gitRepoName)
	projectResource := testutils.HclProjectResource(projectName)
	gitRepoResource := fmt.Sprintf("%s\n%s", projectResource, azureGitRepoResource)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config:      gitRepoResource,
				ExpectError: regexp.MustCompile(`config is invalid: "initialization": required field is not set`),
			},
		},
	})

}

func TestAccGitRepo_Create_Import(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	repoImportConfig := testutils.HclProjectGitRepositoryImport(gitRepoName, projectName)
	tfRepoNode := "azuredevops_git_repository.repository"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: repoImportConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfRepoNode, "is_fork"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "remote_url"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "size"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "ssh_url"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "url"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "web_url"),
					resource.TestCheckResourceAttr(tfRepoNode, "initialization.#", "1"),
					checkGitRepoExists(gitRepoName),
				)},
		},
	})

}

// Verifies that a newly created repo with init_type of "Clean" has the expected
// master branch available
func TestAccGitRepo_RepoInitialization_Clean(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	tfRepoNode := "azuredevops_git_repository.repository"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkGitRepoDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclGitRepoResource(projectName, gitRepoName, "Clean"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfRepoNode, "project_id"),
					resource.TestCheckResourceAttr(tfRepoNode, "name", gitRepoName),
					checkGitRepoExists(gitRepoName),
					resource.TestCheckResourceAttr(tfRepoNode, "default_branch", "refs/heads/master"),
				),
			},
		},
	})
}

// Verifies that a newly created repo with init_type of "Uninitialized" does NOT
// have a master branch established
func TestAccGitRepo_RepoInitialization_Uninitialized(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	tfRepoNode := "azuredevops_git_repository.repository"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkGitRepoDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclGitRepoResource(projectName, gitRepoName, "Uninitialized"),
				Check: resource.ComposeTestCheckFunc(
					checkGitRepoExists(gitRepoName),
					resource.TestCheckResourceAttrSet(tfRepoNode, "project_id"),
					resource.TestCheckResourceAttr(tfRepoNode, "name", gitRepoName),
					resource.TestCheckResourceAttr(tfRepoNode, "default_branch", ""),
				),
			},
		},
	})
}

// Verifies that a newly forked repo does NOT return an empty branch_name
func TestAccGitRepo_RepoFork_BranchNotEmpty(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	gitForkedRepoName := testutils.GenerateResourceName()
	tfRepoNode := "azuredevops_git_repository.repository"
	tfForkedRepoNode := "azuredevops_git_repository.gitforkedrepo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkGitRepoDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclForkedGitRepoResource(projectName, gitRepoName, gitForkedRepoName, "Clean", "Uninitialized"),
				Check: resource.ComposeTestCheckFunc(
					checkGitRepoExists(gitRepoName),
					resource.TestCheckResourceAttrSet(tfRepoNode, "project_id"),
					resource.TestCheckResourceAttr(tfRepoNode, "name", gitRepoName),
					resource.TestCheckResourceAttr(tfRepoNode, "default_branch", "refs/heads/master"),
					resource.TestCheckResourceAttr(tfForkedRepoNode, "name", gitForkedRepoName),
					resource.TestCheckResourceAttr(tfForkedRepoNode, "default_branch", "refs/heads/master"),
				),
			},
		},
	})
}

func TestAccGitRepo_PrivateImport_BranchNotEmpty(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	gitImportRepoName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()

	tfRepoNode := "azuredevops_git_repository.repository"
	tfImportRepoNode := "azuredevops_git_repository.import"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, &[]string{
				"AZDO_GENERIC_GIT_SERVICE_CONNECTION_USERNAME",
				"AZDO_GENERIC_GIT_SERVICE_CONNECTION_PASSWORD",
			})
		},
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkGitRepoDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclProjectGitRepoImportPrivate(projectName, gitRepoName, gitImportRepoName, serviceEndpointName),
				Check: resource.ComposeTestCheckFunc(
					checkGitRepoExists(gitRepoName),
					resource.TestCheckResourceAttrSet(tfRepoNode, "project_id"),
					resource.TestCheckResourceAttr(tfRepoNode, "name", gitRepoName),
					resource.TestCheckResourceAttr(tfRepoNode, "default_branch", "refs/heads/master"),
					resource.TestCheckResourceAttrSet(tfImportRepoNode, "project_id"),
					resource.TestCheckResourceAttr(tfImportRepoNode, "name", gitImportRepoName),
					resource.TestCheckResourceAttr(tfImportRepoNode, "default_branch", "refs/heads/master"),
				),
			},
		},
	})
}

// or not the definition (1) exists in the state and (2) exist in AzDO and (3) has the correct name
func checkGitRepoExists(expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

		gitRepo, ok := s.RootModule().Resources["azuredevops_git_repository.repository"]
		if !ok {
			return fmt.Errorf("Did not find a repo definition in the TF state")
		}

		repoID := gitRepo.Primary.ID
		projectID := gitRepo.Primary.Attributes["project_id"]

		repo, err := readGitRepo(clients, repoID, projectID)
		if err != nil {
			return err
		}

		if *repo.Name != expectedName {
			return fmt.Errorf("AzDO Git Repository has Name=%s, but expected Name=%s", *repo.Name, expectedName)
		}

		return nil
	}
}

func checkGitRepoDestroyed(s *terraform.State) error {
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

	// verify that every repository referenced in the state does not exist in AzDO
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_git_repository" {
			continue
		}

		repoID := resource.Primary.ID
		projectID := resource.Primary.Attributes["project_id"]

		// indicates the git repository still exists - this should fail the test
		if _, err := readGitRepo(clients, repoID, projectID); err == nil {
			return fmt.Errorf("repository with ID %s should not exist", repoID)
		}
	}

	return nil
}

// Lookup an Azure Git Repository using the ID, or name if the ID is not set.
func readGitRepo(clients *client.AggregatedClient, repoID string, projectID string) (*git.GitRepository, error) {
	return clients.GitReposClient.GetRepository(clients.Ctx, git.GetRepositoryArgs{
		RepositoryId: converter.String(repoID),
		Project:      converter.String(projectID),
	})
}
