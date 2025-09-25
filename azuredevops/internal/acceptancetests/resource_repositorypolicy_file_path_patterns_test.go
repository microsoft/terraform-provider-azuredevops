package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccRepositoryPolicyFilePathPatterns(t *testing.T) {
	testutils.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"RepositoryPolicies": {
			"basic":  testAccRepositoryPolicyFilePathPatternsRepoPolicyBasic,
			"update": testAccRepositoryPolicyFilePathPatternsRepoPolicyUpdate,
		},
		"ProjectPolicies": {
			"basic":  testAccRepositoryPolicyFilePathPatternsProjectPolicyBasic,
			"update": testAccRepositoryPolicyFilePathPatternsProjectPolicyUpdate,
		},
	})
}

func testAccRepositoryPolicyFilePathPatternsRepoPolicyBasic(t *testing.T) {
	authorEmailTfNode := "azuredevops_repository_policy_file_path_pattern.test"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepoPolicyFilePathPatternsResourceRepoPolicyBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(authorEmailTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(authorEmailTfNode, "filepath_patterns.#", "1"),
				),
			}, {
				ResourceName:      authorEmailTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(authorEmailTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRepositoryPolicyFilePathPatternsRepoPolicyUpdate(t *testing.T) {
	authorEmailTfNode := "azuredevops_repository_policy_file_path_pattern.test"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepoPolicyFilePathPatternsResourceRepoPolicyBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(authorEmailTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(authorEmailTfNode, "filepath_patterns.#", "1"),
				),
			}, {
				Config: hclRepoPolicyFilePathPatternsResourceRepoPolicyUpdate(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(authorEmailTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(authorEmailTfNode, "filepath_patterns.#", "2"),
				),
			}, {
				ResourceName:      authorEmailTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(authorEmailTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRepositoryPolicyFilePathPatternsProjectPolicyBasic(t *testing.T) {
	authorEmailTfNode := "azuredevops_repository_policy_file_path_pattern.test"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepoPolicyFilePathPatternsResourceProjectPolicyBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(authorEmailTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(authorEmailTfNode, "filepath_patterns.#", "1"),
				),
			}, {
				ResourceName:      authorEmailTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(authorEmailTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRepositoryPolicyFilePathPatternsProjectPolicyUpdate(t *testing.T) {
	authorEmailTfNode := "azuredevops_repository_policy_file_path_pattern.test"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepoPolicyFilePathPatternsResourceProjectPolicyBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(authorEmailTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(authorEmailTfNode, "filepath_patterns.#", "1"),
				),
			}, {
				Config: hclRepoPolicyFilePathPatternsResourceProjectPolicyUpdate(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(authorEmailTfNode, "filepath_patterns.#", "2"),
					resource.TestCheckResourceAttr(authorEmailTfNode, "enabled", "true"),
				),
			}, {
				ResourceName:      authorEmailTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(authorEmailTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func hclRepoPolicyFilePathPatternsResourceTemplate(projectName string, repoName string) string {
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

func hclRepoPolicyFilePathPatternsResourceRepoPolicyBasic(projectName string, repoName string) string {
	projectAndRepo := hclRepoPolicyFilePathPatternsResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`
%s

resource "azuredevops_repository_policy_file_path_pattern" "test" {
  project_id        = azuredevops_project.test.id
  enabled           = true
  blocking          = true
  filepath_patterns = ["*.go"]
  repository_ids    = [azuredevops_git_repository.test.id]
}`, projectAndRepo)
}

func hclRepoPolicyFilePathPatternsResourceRepoPolicyUpdate(projectName string, repoName string) string {
	projectAndRepo := hclRepoPolicyFilePathPatternsResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`
%s

resource "azuredevops_repository_policy_file_path_pattern" "test" {
  project_id        = azuredevops_project.test.id
  enabled           = true
  blocking          = true
  filepath_patterns = ["*.go", "/home/test/*.ts"]
  repository_ids    = [azuredevops_git_repository.test.id]
}`, projectAndRepo)
}

func hclRepoPolicyFilePathPatternsResourceProjectPolicyBasic(projectName string, repoName string) string {
	projectAndRepo := hclRepoPolicyFilePathPatternsResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`
%s

resource "azuredevops_repository_policy_file_path_pattern" "test" {
  project_id        = azuredevops_project.test.id
  enabled           = true
  blocking          = true
  filepath_patterns = ["*.go"]
  depends_on        = [azuredevops_git_repository.test]
}`, projectAndRepo)
}

func hclRepoPolicyFilePathPatternsResourceProjectPolicyUpdate(projectName string, repoName string) string {
	projectAndRepo := hclRepoPolicyFilePathPatternsResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`
%s

resource "azuredevops_repository_policy_file_path_pattern" "test" {
  project_id        = azuredevops_project.test.id
  enabled           = true
  blocking          = true
  filepath_patterns = ["*.go", "/home/test/*.ts"]
  depends_on        = [azuredevops_git_repository.test]
}`, projectAndRepo)
}
