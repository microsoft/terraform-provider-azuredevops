//go:build (all || resource_policy_file_size) && !resource_policy_file_size

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccRepositoryPolicyFileSize(t *testing.T) {
	testutils.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"RepositoryPolicies": {
			"basic":  testAccRepoPolicyFileSizeBasic,
			"update": testAccRepoPolicyFileSizeUpdate,
		},
		"ProjectPolicies": {
			"basic":  testAccProjectPolicyFileSizeBasic,
			"update": testAccProjectPolicyFileSizeUpdate,
		},
	})
}

func testAccRepoPolicyFileSizeBasic(t *testing.T) {
	fileSizeTfNode := "azuredevops_repository_policy_max_file_size.test"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepoPolicyFileSizeBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fileSizeTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(fileSizeTfNode, "max_file_size", "1"),
				),
			}, {
				ResourceName:      fileSizeTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(fileSizeTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRepoPolicyFileSizeUpdate(t *testing.T) {
	fileSizeTfNode := "azuredevops_repository_policy_max_file_size.test"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepoPolicyFileSizeBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fileSizeTfNode, "enabled", "true"),
				),
			}, {
				Config: hclRepoPolicyFileSizeUpdate(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fileSizeTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(fileSizeTfNode, "max_file_size", "5"),
				),
			}, {
				ResourceName:      fileSizeTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(fileSizeTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProjectPolicyFileSizeBasic(t *testing.T) {
	fileSizeTfNode := "azuredevops_repository_policy_max_file_size.test"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclProjectPolicyFileSizeBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fileSizeTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(fileSizeTfNode, "max_file_size", "1"),
				),
			}, {
				ResourceName:      fileSizeTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(fileSizeTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProjectPolicyFileSizeUpdate(t *testing.T) {
	fileSizeTfNode := "azuredevops_repository_policy_max_file_size.test"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclProjectPolicyFileSizeBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fileSizeTfNode, "enabled", "true"),
				),
			}, {
				Config: hclProjectPolicyFileSizeUpdate(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fileSizeTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(fileSizeTfNode, "max_file_size", "5"),
				),
			}, {
				ResourceName:      fileSizeTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(fileSizeTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func hclPolicyFileSizeResourceTemplate(projectName string, repoName string) string {
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

func hclRepoPolicyFileSizeBasic(projectName string, repoName string) string {
	projectAndRepo := hclPolicyFileSizeResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`
%s

resource "azuredevops_repository_policy_max_file_size" "test" {
  project_id     = azuredevops_project.test.id
  enabled        = true
  blocking       = true
  max_file_size  = 1
  repository_ids = [azuredevops_git_repository.test.id]
}`, projectAndRepo)
}

func hclRepoPolicyFileSizeUpdate(projectName string, repoName string) string {
	projectAndRepo := hclPolicyFileSizeResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`
%s

resource "azuredevops_repository_policy_max_file_size" "test" {
  project_id     = azuredevops_project.test.id
  enabled        = true
  blocking       = true
  max_file_size  = 5
  repository_ids = [azuredevops_git_repository.test.id]
}`, projectAndRepo)
}

func hclProjectPolicyFileSizeBasic(projectName string, repoName string) string {
	projectAndRepo := hclPolicyFileSizeResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`
%s

resource "azuredevops_repository_policy_max_file_size" "test" {
  project_id    = azuredevops_project.test.id
  enabled       = true
  blocking      = true
  max_file_size = 1
  depends_on    = [azuredevops_git_repository.test]
}`, projectAndRepo)
}

func hclProjectPolicyFileSizeUpdate(projectName string, repoName string) string {
	projectAndRepo := hclPolicyFileSizeResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`
%s

resource "azuredevops_repository_policy_max_file_size" "test" {
  project_id    = azuredevops_project.test.id
  enabled       = true
  blocking      = true
  max_file_size = 5
  depends_on    = [azuredevops_git_repository.test]
}`, projectAndRepo)
}
