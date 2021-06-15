// +build all resource_repositorypolicy_author_email_patterns
// +build !exclude_serviceendpoints

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccRepositoryPolicyAuthorEmailPatternsBasic(t *testing.T) {
	authorEmailTfNode := "azuredevops_repository_policy_author_email_pattern.p"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepositoryPolicyAuthorEmailPatternsResourceBasic(projectName, repoName),
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

func TestAccRepositoryPolicyAuthorEmailPatternsComplete(t *testing.T) {
	authorEmailTfNode := "azuredevops_repository_policy_author_email_pattern.p"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepositoryPolicyAuthorEmailPatternsResourceComplete(projectName, repoName),
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

func hclRepositoryPolicyAuthorEmailPatternsResourceTemplate(projectName string, repoName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "p" {
  name               = "%s"
  description        = "Test Project Description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  features = {
    "testplans" = "disabled"
    "artifacts" = "disabled"
  }
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

func hclRepositoryPolicyAuthorEmailPatternsResourceBasic(projectName string, repoName string) string {
	projectAndRepo := hclRepositoryPolicyAuthorEmailPatternsResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`%s %s`, projectAndRepo, `
resource "azuredevops_repository_policy_author_email_pattern" "p" {
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

func hclRepositoryPolicyAuthorEmailPatternsResourceComplete(projectName string, repoName string) string {
	projectAndRepo := hclRepositoryPolicyAuthorEmailPatternsResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`%s %s`, projectAndRepo, `
resource "azuredevops_repository_policy_author_email_pattern" "p" {
 project_id = azuredevops_project.p.id

 enabled  = true
 blocking = true

 settings {
   author_email_patterns = ["test1@test.com", "test2@test.com"]
   scope {
     repository_id  = azuredevops_git_repository.r.id
   }
 }
}
`)
}
