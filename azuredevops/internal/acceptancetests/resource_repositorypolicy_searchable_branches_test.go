package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccRepositoryPolicySearchableBranches(t *testing.T) {
	testutils.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"RepositoryPolicies": {
			"basic":  testAccRepositoryPolicySearchableBranchesBasic,
			"update": testAccRepositoryPolicySearchableBranchesUpdate,
		},
	})
}

func testAccRepositoryPolicySearchableBranchesBasic(t *testing.T) {
	searchableBranchesTfNode := "azuredevops_repository_searchable_branches.test"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepoPolicySearchableBranchesResourceBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(searchableBranchesTfNode, "searchable_branches.#", "1"),
				),
			}, {
				ResourceName:      searchableBranchesTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(searchableBranchesTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRepositoryPolicySearchableBranchesUpdate(t *testing.T) {
	searchableBranchesTfNode := "azuredevops_repository_searchable_branches.test"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepoPolicySearchableBranchesResourceBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(searchableBranchesTfNode, "searchable_branches.#", "1"),
				),
			}, {
				Config: hclRepoPolicySearchableBranchesResourceUpdate(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(searchableBranchesTfNode, "searchable_branches.#", "2"),
				),
			}, {
				ResourceName:      searchableBranchesTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(searchableBranchesTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func hclRepoPolicySearchableBranchesResourceTemplate(projectName string, repoName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%s"
  description        = "Test Project Description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_git_repository" "test" {
  project_id = azuredevops_project.test.id
  name       = "%s"
  initialization {
    init_type = "Clean"
  }
}
`, projectName, repoName)
}

func hclRepoPolicySearchableBranchesResourceBasic(projectName string, repoName string) string {
	projectAndRepo := hclRepoPolicySearchableBranchesResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`
%s

resource "azuredevops_repository_policy_searchable_branches" "test" {
  project_id = azuredevops_project.test.id

  searchable_branches     = ["testbranch"]
  repository_ids          = [azuredevops_git_repository.test.id]
}`, projectAndRepo)
}

func hclRepoPolicySearchableBranchesResourceUpdate(projectName string, repoName string) string {
	projectAndRepo := hclRepoPolicySearchableBranchesResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`
%s

resource "azuredevops_repository_policy_searchable_branches" "test" {
  project_id = azuredevops_project.test.id

  searchable_branches     = ["testbranch2"]
  repository_ids          = [azuredevops_git_repository.test.id]
}`, projectAndRepo)
}
