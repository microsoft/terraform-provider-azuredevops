package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccRepositoryPolicyAuthorEmailPatterns(t *testing.T) {
	testutils.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"RepositoryPolicies": {
			"basic":  testAccRepositoryPolicyAuthorEmailPatternsRepoPolicyBasic,
			"update": testAccRepositoryPolicyAuthorEmailPatternsRepoPolicyUpdate,
		},
		"ProjectPolicies": {
			"basic":  testAccRepositoryPolicyAuthorEmailPatternsProjectPolicyBasic,
			"update": testAccRepositoryPolicyAuthorEmailPatternsProjectPolicyUpdate,
		},
	})
}

func testAccRepositoryPolicyAuthorEmailPatternsRepoPolicyBasic(t *testing.T) {
	authorEmailTfNode := "azuredevops_repository_policy_author_email_pattern.test"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepositoryPolicyAuthorEmailPatternsResourceRepoPolicyBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(authorEmailTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(authorEmailTfNode, "author_email_patterns.0", "test1@test.com"),
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

func testAccRepositoryPolicyAuthorEmailPatternsRepoPolicyUpdate(t *testing.T) {
	authorEmailTfNode := "azuredevops_repository_policy_author_email_pattern.test"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepositoryPolicyAuthorEmailPatternsResourceRepoPolicyBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(authorEmailTfNode, "enabled", "true"),
				),
			}, {
				Config: hclRepositoryPolicyAuthorEmailPatternsResourceRepoPolicyUpdate(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(authorEmailTfNode, "author_email_patterns.0", "test2@test.com"),
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

func testAccRepositoryPolicyAuthorEmailPatternsProjectPolicyBasic(t *testing.T) {
	authorEmailTfNode := "azuredevops_repository_policy_author_email_pattern.test"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepositoryPolicyAuthorEmailPatternsResourceProjectPolicyBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(authorEmailTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(authorEmailTfNode, "author_email_patterns.#", "1"),
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

func testAccRepositoryPolicyAuthorEmailPatternsProjectPolicyUpdate(t *testing.T) {
	authorEmailTfNode := "azuredevops_repository_policy_author_email_pattern.test"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepositoryPolicyAuthorEmailPatternsResourceProjectPolicyBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(authorEmailTfNode, "enabled", "true"),
				),
			}, {
				Config: hclRepositoryPolicyAuthorEmailPatternsResourceProjectPolicyUpdate(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(authorEmailTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(authorEmailTfNode, "author_email_patterns.#", "2"),
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

func hclRepositoryPolicyAuthorEmailPatternsResourceTemplate(projectName string, repoName string) string {
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

func hclRepositoryPolicyAuthorEmailPatternsResourceRepoPolicyBasic(projectName string, repoName string) string {
	projectAndRepo := hclRepositoryPolicyAuthorEmailPatternsResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`
%s

resource "azuredevops_repository_policy_author_email_pattern" "test" {
  project_id = azuredevops_project.test.id

  enabled  = true
  blocking = true

  author_email_patterns = ["test1@test.com"]
  repository_ids        = [azuredevops_git_repository.test.id]
}`, projectAndRepo)
}

func hclRepositoryPolicyAuthorEmailPatternsResourceRepoPolicyUpdate(projectName string, repoName string) string {
	projectAndRepo := hclRepositoryPolicyAuthorEmailPatternsResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`
%s

resource "azuredevops_repository_policy_author_email_pattern" "test" {
  project_id = azuredevops_project.test.id

  enabled  = true
  blocking = true

  author_email_patterns = ["test2@test.com"]
  repository_ids        = [azuredevops_git_repository.test.id]
}`, projectAndRepo)
}

func hclRepositoryPolicyAuthorEmailPatternsResourceProjectPolicyBasic(projectName string, repoName string) string {
	projectAndRepo := hclRepositoryPolicyAuthorEmailPatternsResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`
%s

resource "azuredevops_repository_policy_author_email_pattern" "test" {
  project_id = azuredevops_project.test.id

  enabled               = true
  blocking              = true
  author_email_patterns = ["test1@test.com"]
  depends_on            = [azuredevops_git_repository.test]
}`, projectAndRepo)
}

func hclRepositoryPolicyAuthorEmailPatternsResourceProjectPolicyUpdate(projectName string, repoName string) string {
	projectAndRepo := hclRepositoryPolicyAuthorEmailPatternsResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`
%s

resource "azuredevops_repository_policy_author_email_pattern" "test" {
  project_id = azuredevops_project.test.id

  enabled  = true
  blocking = true

  author_email_patterns = ["test1@test.com", "test2@test.com"]
  depends_on            = [azuredevops_git_repository.test]
}`, projectAndRepo)
}
