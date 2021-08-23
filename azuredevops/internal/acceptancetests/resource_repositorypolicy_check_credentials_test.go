// +build all resource_policy_check_credentials
// +build !resource_policy_check_credentials

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

const checkCredentialsTfNode = "azuredevops_repository_policy_check_credentials.p"

func TestAccRepoPolicyCheckCredentials(t *testing.T) {
	testutils.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"RepositoryPolicies": {
			"basic":  testAccRepoPolicyCheckCredentialsBasic,
			"update": testAccRepoPolicyCheckCredentialsUpdate,
		},
		"ProjectPolicies": {
			"basic":  testAccProjectPolicyCheckCredentialsBasic,
			"update": testAccProjectPolicyCheckCredentialsUpdate,
		},
	})
}

func testAccRepoPolicyCheckCredentialsBasic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepoPolicyCheckCredentialsBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(checkCredentialsTfNode, "enabled", "true"),
				),
			}, {
				ResourceName:      checkCredentialsTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(checkCredentialsTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRepoPolicyCheckCredentialsUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepoPolicyCheckCredentialsBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(checkCredentialsTfNode, "enabled", "true"),
				),
			}, {
				Config: hclRepoPolicyCheckCredentialsUpdate(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(checkCredentialsTfNode, "enabled", "false"),
				),
			}, {
				ResourceName:      checkCredentialsTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(checkCredentialsTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProjectPolicyCheckCredentialsBasic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclProjectPolicyCheckCredentialsBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(checkCredentialsTfNode, "enabled", "true"),
				),
			}, {
				ResourceName:      checkCredentialsTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(checkCredentialsTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProjectPolicyCheckCredentialsUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclProjectPolicyCheckCredentialsBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(checkCredentialsTfNode, "enabled", "true"),
				),
			}, {
				Config: hclProjectPolicyCheckCredentialsUpdate(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(checkCredentialsTfNode, "enabled", "false"),
				),
			}, {
				ResourceName:      checkCredentialsTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(checkCredentialsTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func hclPolicyCheckCredentialsResourceTemplate(projectName string, repoName string) string {
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

func hclRepoPolicyCheckCredentialsBasic(projectName string, repoName string) string {
	projectAndRepo := hclPolicyCheckCredentialsResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`%s %s`, projectAndRepo, `
resource "azuredevops_repository_policy_check_credentials" "p" {
  project_id = azuredevops_project.p.id

  enabled  = true
  blocking = true
  repository_ids  = [azuredevops_git_repository.r.id]
}
`)
}

func hclRepoPolicyCheckCredentialsUpdate(projectName string, repoName string) string {
	projectAndRepo := hclPolicyCheckCredentialsResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`%s %s`, projectAndRepo, `
resource "azuredevops_repository_policy_check_credentials" "p" {
  project_id = azuredevops_project.p.id
  enabled  = false
  blocking = true
  repository_ids  = [azuredevops_git_repository.r.id]
}
`)
}

func hclProjectPolicyCheckCredentialsBasic(projectName string, repoName string) string {
	projectAndRepo := hclPolicyCheckCredentialsResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`%s %s`, projectAndRepo, `
resource "azuredevops_repository_policy_check_credentials" "p" {
  project_id = azuredevops_project.p.id
  enabled  = true
  blocking = true
  depends_on = [azuredevops_git_repository.r]
}
`)
}

func hclProjectPolicyCheckCredentialsUpdate(projectName string, repoName string) string {
	projectAndRepo := hclPolicyCheckCredentialsResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`%s %s`, projectAndRepo, `
resource "azuredevops_repository_policy_check_credentials" "p" {
  project_id = azuredevops_project.p.id
  enabled  = false
  blocking = true
  depends_on = [azuredevops_git_repository.r]
}
`)
}
