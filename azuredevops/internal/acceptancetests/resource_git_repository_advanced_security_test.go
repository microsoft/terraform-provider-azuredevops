package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccGitRepositoryAdvancedSecurity_EnableUpdateImport(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()
	tfNode := "azuredevops_git_repository_advanced_security.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclGitRepositoryAdvancedSecurity(projectName, repoName, true, true, true, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(tfNode, "repository_id", "azuredevops_git_repository.repository", "id"),
					resource.TestCheckResourceAttr(tfNode, "secret_protection_enabled", "true"),
					resource.TestCheckResourceAttr(tfNode, "block_pushes", "true"),
					resource.TestCheckResourceAttr(tfNode, "code_security_enabled", "true"),
					resource.TestCheckResourceAttr(tfNode, "dependency_scanning_injection_enabled", "true"),
				),
			},
			{
				Config: hclGitRepositoryAdvancedSecurity(projectName, repoName, true, false, false, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "secret_protection_enabled", "true"),
					resource.TestCheckResourceAttr(tfNode, "block_pushes", "false"),
					resource.TestCheckResourceAttr(tfNode, "code_security_enabled", "false"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func hclGitRepositoryAdvancedSecurity(projectName, repoName string, secretProtection, blockPushes, codeSecurity, dependencyScanningInjection bool) string {
	return fmt.Sprintf(`
%s

resource "azuredevops_git_repository_advanced_security" "test" {
  project_id                            = azuredevops_project.project.id
  repository_id                         = azuredevops_git_repository.repository.id
  secret_protection_enabled             = %t
  block_pushes                          = %t
  code_security_enabled                 = %t
  dependency_scanning_injection_enabled = %t
}
`, testutils.HclGitRepoResource(projectName, repoName, "Clean"), secretProtection, blockPushes, codeSecurity, dependencyScanningInjection)
}
