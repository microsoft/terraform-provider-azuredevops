package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccRepositoryPolicyFilePathPatternsBasic(t *testing.T) {
	authorEmailTfNode := "azuredevops_repository_policy_file_path_pattern.p"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepoPolicyFilePathPatternsResourceBasic(projectName, repoName),
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

func TestAccRepositoryPolicyFilePathPatternsComplete(t *testing.T) {
	authorEmailTfNode := "azuredevops_repository_policy_file_path_pattern.p"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepoPolicyFilePathPatternsResourceComplete(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(authorEmailTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(authorEmailTfNode, "settings.#", "1"),
					resource.TestCheckResourceAttr(authorEmailTfNode, "settings.0.filepath_patterns.#", "2"),
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
func TestAccRepositoryPolicyFilePathPatternsUpdate(t *testing.T) {
	authorEmailTfNode := "azuredevops_repository_policy_file_path_pattern.p"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepoPolicyFilePathPatternsResourceBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(authorEmailTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(authorEmailTfNode, "settings.#", "1"),
				),
			}, {
				Config: hclRepoPolicyFilePathPatternsResourceUpdate(projectName, repoName),
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

func hclRepoPolicyFilePathPatternsResourceBasic(projectName string, repoName string) string {
	projectAndRepo := hclRepoPolicyFilePathPatternsResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`%s %s`, projectAndRepo, `
resource "azuredevops_repository_policy_file_path_pattern" "p" {
  project_id = azuredevops_project.p.id

  enabled  = true
  blocking = true

  settings {
    scope {
      repository_id  = azuredevops_git_repository.r.id
    }
  }
}
`)
}

func hclRepoPolicyFilePathPatternsResourceComplete(projectName string, repoName string) string {
	projectAndRepo := hclRepoPolicyFilePathPatternsResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`%s %s`, projectAndRepo, `
resource "azuredevops_repository_policy_file_path_pattern" "p" {
 project_id = azuredevops_project.p.id

 enabled  = true
 blocking = true

 settings {
   filepath_patterns = ["/home/workspace/filefilter.java","*.ts"]
   scope {
     repository_id  = azuredevops_git_repository.r.id
   }
 }
}
`)
}

func hclRepoPolicyFilePathPatternsResourceUpdate(projectName string, repoName string) string {
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
