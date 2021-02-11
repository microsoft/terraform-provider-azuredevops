// +build all core resource_git_repository_file
// +build !exclude_resource_git_repository_file

package acceptancetests

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/git"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// TestAccGitRepoFile_CreateUpdateDelete verifies that a file can
// be added to a repository and the contents can be updated
func TestAccGitRepoFile_CreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	tfRepoFileNode := "azuredevops_git_repository_file.file"

	branch := "refs/heads/main"
	file := "foo.txt"
	contentFirst := "bar"
	contentSecond := "baz"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclGitRepoFileResource(projectName, gitRepoName, "Clean", branch, file, contentFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfRepoFileNode, "file", file),
					resource.TestCheckResourceAttr(tfRepoFileNode, "content", contentFirst),
					resource.TestCheckResourceAttr(tfRepoFileNode, "branch", branch),
					resource.TestCheckResourceAttrSet(tfRepoFileNode, "commit_message"),
					checkGitRepoFileContent(contentFirst),
				),
			},
			{
				Config: testutils.HclGitRepoFileResource(projectName, gitRepoName, "Clean", branch, file, contentSecond),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfRepoFileNode, "file", file),
					resource.TestCheckResourceAttr(tfRepoFileNode, "content", contentSecond),
					resource.TestCheckResourceAttr(tfRepoFileNode, "branch", branch),
					resource.TestCheckResourceAttrSet(tfRepoFileNode, "commit_message"),
					checkGitRepoFileContent(contentSecond),
				),
			},
			{
				Config: testutils.HclGitRepoResource(projectName, gitRepoName, "Clean"),
				Check: resource.ComposeTestCheckFunc(
					checkGitRepoFileNotExists(file),
				),
			},
		},
	})
}

// TestAccGitRepo_Create_IncorrectBranch verifies a file
// can't be added to a non existant branch
func TestAccGitRepoFile_Create_IncorrectBranch(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config:      testutils.HclGitRepoFileResource(projectName, gitRepoName, "Clean", "foobar", "foo", "bar"),
				ExpectError: regexp.MustCompile(`errors during apply: Branch "foobar" does not exist`),
			},
		},
	})
}

func checkGitRepoFileNotExists(fileName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

		repo, ok := s.RootModule().Resources["azuredevops_git_repository.repository"]
		if !ok {
			return fmt.Errorf("Did not find a repo definition in the TF state")
		}

		ctx := context.Background()
		_, err := clients.GitReposClient.GetItem(ctx, git.GetItemArgs{
			RepositoryId: &repo.Primary.ID,
			Path:         &fileName,
		})
		if err != nil && strings.Contains(err.Error(), "could not be found in the repository") {
			return err
		}

		return nil
	}
}

func checkGitRepoFileContent(expectedContent string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

		gitFile, ok := s.RootModule().Resources["azuredevops_git_repository_file.file"]
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
