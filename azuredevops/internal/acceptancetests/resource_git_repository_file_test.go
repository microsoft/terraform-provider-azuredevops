//go:build (all || core || resource_git_repository_file) && !exclude_resource_git_repository_file

package acceptancetests

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/git"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

func TestAccGitRepoFile_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	tfRepoFileNode := "azuredevops_git_repository_file.test"

	branch := "refs/heads/master"
	file := "foo.txt"
	contentFirst := "bar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclGitRepositoryFileBasic(projectName, gitRepoName, branch, file, contentFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfRepoFileNode, "file", file),
					resource.TestCheckResourceAttr(tfRepoFileNode, "content", contentFirst),
					resource.TestCheckResourceAttr(tfRepoFileNode, "branch", branch),
					resource.TestCheckResourceAttrSet(tfRepoFileNode, "commit_message"),
					checkGitRepoFileContent(contentFirst),
				),
			},
			{
				ResourceName:      tfRepoFileNode,
				ImportStateIdFunc: repositoryFileIdFunc(tfRepoFileNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGitRepoFile_complete(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	tfRepoFileNode := "azuredevops_git_repository_file.test"

	branch := "refs/heads/master"
	file := "foo.txt"
	contentFirst := "bar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclGitRepositoryFileComplete(projectName, gitRepoName, branch, file, contentFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfRepoFileNode, "file", file),
					resource.TestCheckResourceAttr(tfRepoFileNode, "content", contentFirst),
					resource.TestCheckResourceAttr(tfRepoFileNode, "branch", branch),
					resource.TestCheckResourceAttrSet(tfRepoFileNode, "commit_message"),
					checkGitRepoFileContent(contentFirst),
				),
			},
			{
				ResourceName:      tfRepoFileNode,
				ImportStateIdFunc: repositoryFileIdFunc(tfRepoFileNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGitRepoFile_authorEmailPolicy(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	tfRepoFileNode := "azuredevops_git_repository_file.test"

	branch := "refs/heads/master"
	file := "foo.txt"
	contentFirst := "bar"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclGitRepositoryFileAuthorEmailPolicy(projectName, gitRepoName, branch, file, contentFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfRepoFileNode, "file", file),
					resource.TestCheckResourceAttr(tfRepoFileNode, "content", contentFirst),
					resource.TestCheckResourceAttr(tfRepoFileNode, "branch", branch),
					resource.TestCheckResourceAttrSet(tfRepoFileNode, "commit_message"),
					checkGitRepoFileContent(contentFirst),
				),
			},
			{
				ResourceName:      tfRepoFileNode,
				ImportStateIdFunc: repositoryFileIdFunc(tfRepoFileNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGitRepoFile_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	tfRepoFileNode := "azuredevops_git_repository_file.test"

	branch := "refs/heads/master"
	file := "foo.txt"
	contentFirst := "bar"
	contentSecond := "baz"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclGitRepositoryFileBasic(projectName, gitRepoName, branch, file, contentFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfRepoFileNode, "file", file),
					resource.TestCheckResourceAttr(tfRepoFileNode, "content", contentFirst),
					resource.TestCheckResourceAttr(tfRepoFileNode, "branch", branch),
					resource.TestCheckResourceAttrSet(tfRepoFileNode, "commit_message"),
					checkGitRepoFileContent(contentFirst),
				),
			},
			{
				ResourceName:      tfRepoFileNode,
				ImportStateIdFunc: repositoryFileIdFunc(tfRepoFileNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: hclGitRepositoryFileBasic(projectName, gitRepoName, branch, file, contentSecond),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfRepoFileNode, "file", file),
					resource.TestCheckResourceAttr(tfRepoFileNode, "content", contentSecond),
					resource.TestCheckResourceAttr(tfRepoFileNode, "branch", branch),
					resource.TestCheckResourceAttrSet(tfRepoFileNode, "commit_message"),
					checkGitRepoFileContent(contentSecond),
				),
			},
			{
				ResourceName:      tfRepoFileNode,
				ImportStateIdFunc: repositoryFileIdFunc(tfRepoFileNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: hclGitRepositoryFileWithoutFile(projectName, gitRepoName),
				Check: resource.ComposeTestCheckFunc(
					checkGitRepoFileNotExists(file),
				),
			},
		},
	})
}

func TestAccGitRepoFile_incorrectBranch(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config:      hclGitRepositoryFileBasic(projectName, gitRepoName, "foobar", "foo", "bar"),
				ExpectError: regexp.MustCompile(`Creating Git file. Branch not found. Name: foobar`),
			},
		},
	})
}

func repositoryFileIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Resource node not found: %s", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["repository_id"], rs.Primary.Attributes["file"]), nil
	}
}

func checkGitRepoFileNotExists(fileName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

		repo, ok := s.RootModule().Resources["azuredevops_git_repository.test"]
		if !ok {
			return fmt.Errorf("Did not find a repo definition in the TF state")
		}

		ctx := context.Background()
		_, err := clients.GitReposClient.GetItem(ctx, git.GetItemArgs{
			RepositoryId: &repo.Primary.ID,
			Path:         &fileName,
		})
		if err != nil && !strings.Contains(err.Error(), "could not be found in the repository") {
			return err
		}

		return nil
	}
}

func checkGitRepoFileContent(expectedContent string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

		gitFile, ok := s.RootModule().Resources["azuredevops_git_repository_file.test"]
		if !ok {
			return fmt.Errorf("Did not find a repo definition in the TF state")
		}

		fileID := gitFile.Primary.ID
		comps := strings.Split(fileID, "/")
		repoID := comps[0]
		file := comps[1]

		ctx := context.Background()
		r, err := clients.GitReposClient.GetItemContent(ctx, git.GetItemContentArgs{
			RepositoryId: &repoID,
			Path:         &file,
		})
		if err != nil {
			return err
		}

		buf := new(bytes.Buffer)
		if _, err = buf.ReadFrom(r); err != nil {
			return err
		}

		if buf.String() != expectedContent {
			return fmt.Errorf("Unexpected git file content: %v", buf.String())
		}

		return nil
	}
}

func hclGitRepositoryFileBasic(name, repoName, branch, file, content string) string {
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

resource "azuredevops_git_repository_file" "test" {
  repository_id = azuredevops_git_repository.test.id
  branch        = "%[3]s"
  file          = "%[4]s"
  content       = "%[5]s"
}
`, name, repoName, branch, file, content)
}

func hclGitRepositoryFileComplete(name, repoName, branch, file, content string) string {
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

resource "azuredevops_git_repository_file" "test" {
  repository_id   = azuredevops_git_repository.test.id
  branch          = "%[3]s"
  file            = "%[4]s"
  content         = "%[5]s"
  author_name     = "author"
  author_email    = "auhtor@test.com"
  committer_name  = "comitter"
  committer_email = "committer@test.com"
}
`, name, repoName, branch, file, content)
}

func hclGitRepositoryFileAuthorEmailPolicy(name, repoName, branch, file, content string) string {
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

resource "azuredevops_repository_policy_author_email_pattern" "test" {
  project_id            = azuredevops_project.test.id
  enabled               = true
  blocking              = true
  author_email_patterns = ["auhtor@test.com"]
  repository_ids        = [azuredevops_git_repository.test.id]
}

resource "azuredevops_git_repository_file" "test" {
  repository_id   = azuredevops_git_repository.test.id
  branch          = "%[3]s"
  file            = "%[4]s"
  content         = "%[5]s"
  author_name     = "author"
  author_email    = "auhtor@test.com"
  committer_name  = "comitter"
  committer_email = "committer@test.com"
  depends_on      = [azuredevops_repository_policy_author_email_pattern.test]
}
`, name, repoName, branch, file, content)
}

func hclGitRepositoryFileWithoutFile(name, repoName string) string {
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
}`, name, repoName)
}
