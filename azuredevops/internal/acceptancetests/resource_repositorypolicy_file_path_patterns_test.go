package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccRepositoryPolicyFilePathPatterns(t *testing.T) {
	testutils.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"RepositoryPolicies": {
			"basic":  testAccRepositoryPolicyFilePathPatternsRepoPolicyBasic,
			"update": testAccRepositoryPolicyFilePathPatternsRepoPolicyUpdate,
		},
		"ProjectPolicies": {
			"basic":  TestAccRepositoryPolicyFilePathPatternsProjectPolicyBasic,
			"update": testAccRepositoryPolicyFilePathPatternsProjectPolicyUpdate,
		},
	})
}

func testAccRepositoryPolicyFilePathPatternsRepoPolicyBasic(t *testing.T) {
	authorEmailTfNode := "azuredevops_repository_policy_file_path_pattern.p"
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
					resource.TestCheckResourceAttr(authorEmailTfNode, "settings.#", "1"),
					resource.TestCheckResourceAttr(authorEmailTfNode, "settings.0.filepath_patterns.#", "1"),
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
	authorEmailTfNode := "azuredevops_repository_policy_file_path_pattern.p"
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
					resource.TestCheckResourceAttr(authorEmailTfNode, "settings.#", "1"),
				),
			}, {
				Config: hclRepoPolicyFilePathPatternsResourceRepoPolicyUpdate(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(authorEmailTfNode, "settings.0.filepath_patterns.#", "2"),
					resource.TestCheckResourceAttr(authorEmailTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(authorEmailTfNode, "settings.#", "1"),
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

func TestAccRepositoryPolicyFilePathPatternsProjectPolicyBasic(t *testing.T) {
	authorEmailTfNode := "azuredevops_repository_policy_file_path_pattern.p"
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
					resource.TestCheckResourceAttr(authorEmailTfNode, "settings.#", "1"),
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
	authorEmailTfNode := "azuredevops_repository_policy_file_path_pattern.p"
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
					resource.TestCheckResourceAttr(authorEmailTfNode, "settings.#", "1"),
				),
			}, {
				Config: hclRepoPolicyFilePathPatternsResourceProjectPolicyUpdate(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(authorEmailTfNode, "settings.0.filepath_patterns.#", "2"),
					resource.TestCheckResourceAttr(authorEmailTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(authorEmailTfNode, "settings.#", "1"),
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
resource "azuredevops_project" "p" {
  name               = "%s"
  description        = "Test Project Description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_git_repository" "r" {
  project_id = azuredevops_project.p.id
  name       = "%s"
  initialization {
    init_type = "Clean"
  }
}
`, projectName, repoName)
}

func hclRepoPolicyFilePathPatternsResourceRepoPolicyBasic(projectName string, repoName string) string {
	projectAndRepo := hclRepoPolicyFilePathPatternsResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`%s %s`, projectAndRepo, `
resource "azuredevops_repository_policy_file_path_pattern" "p" {
  project_id = azuredevops_project.p.id

  enabled  = true
  blocking = true

  settings {
	filepath_patterns = ["*.go"]
    scope {
      repository_id  = azuredevops_git_repository.r.id
    }
  }
}
`)
}

func hclRepoPolicyFilePathPatternsResourceRepoPolicyUpdate(projectName string, repoName string) string {
	projectAndRepo := hclRepoPolicyFilePathPatternsResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`%s %s`, projectAndRepo, `
resource "azuredevops_repository_policy_file_path_pattern" "p" {
  project_id = azuredevops_project.p.id

  enabled  = true
  blocking = true

  settings {
	filepath_patterns = ["*.go", "/home/test/*.ts"]
    scope {
      repository_id  = azuredevops_git_repository.r.id
    }
  }
}
`)
}

func hclRepoPolicyFilePathPatternsResourceProjectPolicyBasic(projectName string, repoName string) string {
	projectAndRepo := hclRepoPolicyFilePathPatternsResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`%s %s`, projectAndRepo, `
resource "azuredevops_repository_policy_file_path_pattern" "p" {
  project_id = azuredevops_project.p.id

  enabled  = true
  blocking = true

  settings {
	filepath_patterns = ["*.go"]
  }
  depends_on = [azuredevops_git_repository.r]
}
`)
}

func hclRepoPolicyFilePathPatternsResourceProjectPolicyUpdate(projectName string, repoName string) string {
	projectAndRepo := hclRepoPolicyFilePathPatternsResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`%s %s`, projectAndRepo, `
resource "azuredevops_repository_policy_file_path_pattern" "p" {
  project_id = azuredevops_project.p.id

  enabled  = true
  blocking = true

  settings {
	filepath_patterns = ["*.go", "/home/test/*.ts"]
  }
  depends_on = [azuredevops_git_repository.r]
}
`)
}
