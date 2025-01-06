//go:build (all || core || resource_git_repository) && !exclude_resource_git_repository
// +build all core resource_git_repository
// +build !exclude_resource_git_repository

package acceptancetests

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func TestAccGitRepository_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoNameFirst := testutils.GenerateResourceName()
	gitRepoNameSecond := testutils.GenerateResourceName()
	tfRepoNode := "azuredevops_git_repository.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkGitRepoDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclGitRepositoryBasic(projectName, gitRepoNameFirst, "Uninitialized"),
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
				Config: hclGitRepositoryBasic(projectName, gitRepoNameSecond, "Uninitialized"),
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

func TestAccGitRepository_disabled(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	tfRepoNode := "azuredevops_git_repository.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkGitRepoDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclGitRepositoryDisable(projectName, gitRepoName, true),
				Check: resource.ComposeTestCheckFunc(
					checkGitRepoExists(gitRepoName),
					resource.TestCheckResourceAttrSet(tfRepoNode, "project_id"),
					resource.TestCheckResourceAttr(tfRepoNode, "name", gitRepoName),
					resource.TestCheckResourceAttr(tfRepoNode, "disabled", "true"),
				),
			},
			{
				Config: hclGitRepositoryDisable(projectName, gitRepoName, false),
				Check: resource.ComposeTestCheckFunc(
					checkGitRepoExists(gitRepoName),
					resource.TestCheckResourceAttrSet(tfRepoNode, "project_id"),
					resource.TestCheckResourceAttr(tfRepoNode, "name", gitRepoName),
					resource.TestCheckResourceAttr(tfRepoNode, "disabled", "false"),
				),
			},
		},
	})
}

func TestAccGitRepository_disabledCannotUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	gitRepoNameUpdate := gitRepoName + "update"
	tfRepoNode := "azuredevops_git_repository.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkGitRepoDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclGitRepositoryDisable(projectName, gitRepoName, true),
				Check: resource.ComposeTestCheckFunc(
					checkGitRepoExists(gitRepoName),
					resource.TestCheckResourceAttrSet(tfRepoNode, "project_id"),
					resource.TestCheckResourceAttr(tfRepoNode, "name", gitRepoName),
					resource.TestCheckResourceAttr(tfRepoNode, "disabled", "true"),
				),
			},
			{
				Config: hclGitRepositoryDisable(projectName, gitRepoNameUpdate, true),
				Check: resource.ComposeTestCheckFunc(
					checkGitRepoExists(gitRepoName),
					resource.TestCheckResourceAttrSet(tfRepoNode, "project_id"),
					resource.TestCheckResourceAttr(tfRepoNode, "name", gitRepoNameUpdate),
					resource.TestCheckResourceAttr(tfRepoNode, "disabled", "true"),
				),
				ExpectError: regexp.MustCompile(`A disabled repository cannot be updated, please enable the repository before attempting to update`),
			},
			{
				Config: hclGitRepositoryDisable(projectName, gitRepoName, false),
				Check: resource.ComposeTestCheckFunc(
					checkGitRepoExists(gitRepoName),
					resource.TestCheckResourceAttrSet(tfRepoNode, "project_id"),
					resource.TestCheckResourceAttr(tfRepoNode, "name", gitRepoName),
					resource.TestCheckResourceAttr(tfRepoNode, "disabled", "false"),
				),
			},
		},
	})
}

func TestAccGitRepository_incorrectInitialization(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config:      hclGitRepositoryIncorrectInitialization(projectName, gitRepoName),
				ExpectError: regexp.MustCompile(`Insufficient initialization blocks`),
			},
		},
	})

}

func TestAccGitRepository_import(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()

	tfRepoNode := "azuredevops_git_repository.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclGitRepositoryImport(projectName, gitRepoName),
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

func TestAccGitRepository_initializationClean(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	tfRepoNode := "azuredevops_git_repository.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkGitRepoDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclGitRepositoryBasic(projectName, gitRepoName, "Clean"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfRepoNode, "project_id"),
					resource.TestCheckResourceAttr(tfRepoNode, "name", gitRepoName),
					checkGitRepoExists(gitRepoName),
					resource.TestCheckResourceAttr(tfRepoNode, "default_branch", "refs/heads/master"),
					resource.TestCheckResourceAttrSet(tfRepoNode, "disabled"),
				),
			},
		},
	})
}

func TestAccGitRepository_uninitialized(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	tfRepoNode := "azuredevops_git_repository.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkGitRepoDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclGitRepositoryBasic(projectName, gitRepoName, "Uninitialized"),
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

func TestAccGitRepository_forkBranchNotEmpty(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	gitForkedRepoName := testutils.GenerateResourceName()
	tfRepoNode := "azuredevops_git_repository.test"
	tfForkedRepoNode := "azuredevops_git_repository.fork"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkGitRepoDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclGitRepositoryForkBranchNotEmpty(projectName, gitRepoName, gitForkedRepoName),
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

func TestAccGitRssepository_privateImportBranchNotEmpty(t *testing.T) {
	if os.Getenv("AZDO_GENERIC_GIT_SERVICE_CONNECTION_USERNAME") == "" ||
		os.Getenv("AZDO_GENERIC_GIT_SERVICE_CONNECTION_PASSWORD") == "" {
		t.Skip("Skipping as AZDO_GENERIC_GIT_SERVICE_CONNECTION_USERNAME or AZDO_GENERIC_GIT_SERVICE_CONNECTION_PASSWORD is not specified")
	}
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
				Config: hclGitRepositoryImportPrivate(projectName, gitRepoName, gitImportRepoName, serviceEndpointName),
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

		gitRepo, ok := s.RootModule().Resources["azuredevops_git_repository.test"]
		if !ok {
			return fmt.Errorf(" Did not find a repo definition in the TF state")
		}

		repoID := gitRepo.Primary.ID
		projectID := gitRepo.Primary.Attributes["project_id"]

		repo, err := readGitRepo(clients, repoID, projectID)
		if err != nil {
			return err
		}

		if *repo.Name != expectedName {
			return fmt.Errorf(" AzDO Git Repository has Name=%s, but expected Name=%s", *repo.Name, expectedName)
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
	repo, err := clients.GitReposClient.GetRepository(clients.Ctx, git.GetRepositoryArgs{
		RepositoryId: converter.String(repoID),
		Project:      converter.String(projectID),
	})

	// If the repository is disabled, the repository cannot be obtained through the GET API
	if utils.ResponseWasNotFound(err) {
		var allRepo *[]git.GitRepository
		allRepo, err = clients.GitReposClient.GetRepositories(clients.Ctx, git.GetRepositoriesArgs{
			Project: converter.String(projectID),
			// This flag is used to include disabled repos
			IncludeHidden: converter.Bool(true),
		})
		if err != nil {
			return nil, err
		}
		for _, gitRepo := range *allRepo {
			if strings.EqualFold((*gitRepo.Id).String(), repoID) ||
				strings.EqualFold(*gitRepo.Name, repoID) {
				repo = &gitRepo
				break
			}
		}
	}
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func hclGitRepositoryBasic(projectName, repoName, initType string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%s"
  initialization {
    init_type = "%s"
  }
}
`, projectName, repoName, initType)
}

func hclGitRepositoryDisable(projectName, repoName string, disabled bool) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%s"
  disabled   = %t
  initialization {
    init_type = "Clean"
  }
}
`, projectName, repoName, disabled)
}

func hclGitRepositoryIncorrectInitialization(projectName, repoName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%s"
}
`, projectName, repoName)
}

func hclGitRepositoryImport(projectName, repoName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%s"
  initialization {
    init_type   = "Import"
    source_type = "Git"
    source_url  = "https://github.com/microsoft/terraform-provider-azuredevops.git"
  }
}
`, projectName, repoName)
}

func hclGitRepositoryForkBranchNotEmpty(projectName, repoName, forkRepoName string) string {
	repoInit := hclGitRepositoryBasic(projectName, repoName, "Clean")
	return fmt.Sprintf(`
%s

resource "azuredevops_git_repository" "fork" {
  project_id           = azuredevops_project.test.id
  parent_repository_id = azuredevops_git_repository.test.id
  name                 = "%s"
  initialization {
    init_type = "Fork"
  }
}`, repoInit, forkRepoName)
}

func hclGitRepositoryImportPrivate(projectName, repoName, importRepoName, serviceEndpointName string) string {
	repoInit := hclGitRepositoryBasic(projectName, repoName, "Clean")
	return fmt.Sprintf(`
%s

resource "azuredevops_serviceendpoint_generic_git" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  repository_url        = azuredevops_git_repository.test.remote_url
}

resource "azuredevops_git_repository" "import" {
  project_id = azuredevops_project.test.id
  name       = "%s"
  initialization {
    init_type             = "Import"
    source_type           = "Git"
    source_url            = azuredevops_git_repository.test.remote_url
    service_connection_id = azuredevops_serviceendpoint_generic_git.test.id
  }
}
`, repoInit, serviceEndpointName, importRepoName)
}
