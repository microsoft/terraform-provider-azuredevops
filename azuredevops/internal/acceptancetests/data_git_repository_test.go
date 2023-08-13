//go:build (all || data_sources || git || data_git_repository) && (!exclude_data_sources || !exclude_git || !data_git_repository)
// +build all data_sources git data_git_repository
// +build !exclude_data_sources !exclude_git !data_git_repository

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// Verifies that the following sequence of events occurrs without error:
//
//	(1) TF can create a project
//	(2) A data source is added to the configuration, and that data source can find the created project
func TestAccGitRepository_DataSource(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	gitRepoName := testutils.GenerateResourceName()
	tfConfigStep1 := testutils.HclGitRepoResource(projectName, gitRepoName, "Clean")
	tfConfigStep2 := fmt.Sprintf("%s\n%s", tfConfigStep1, testutils.HclProjectGitRepository(projectName, gitRepoName))

	tfNode := "data.azuredevops_git_repository.repository"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testutils.PreCheck(t, nil) },
		Providers:                 testutils.GetProviders(),
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: tfConfigStep1,
			}, {
				Config: tfConfigStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttr(tfNode, "name", gitRepoName),
					resource.TestCheckResourceAttr(tfNode, "default_branch", "refs/heads/master"),
				),
			},
		},
	})
}
