//go:build (all || resource_policy_path_lenght) && !resource_policy_path_lenght

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccRepositoryPolicyPathLength(t *testing.T) {
	testutils.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"RepositoryPolicies": {
			"basic":  testAccRepoPolicyPathLengthBasic,
			"update": testAccRepoPolicyPathLengthUpdate,
		},
		"ProjectPolicies": {
			"basic":  testAccProjectPolicyPathLengthBasic,
			"update": testAccProjectPolicyPathLengthUpdate,
		},
	})
}

func testAccRepoPolicyPathLengthBasic(t *testing.T) {
	pathLengthTfNode := "azuredevops_repository_policy_max_path_length.test"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepoPolicyPathLengthBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(pathLengthTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(pathLengthTfNode, "max_path_length", "500"),
				),
			}, {
				ResourceName:      pathLengthTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(pathLengthTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRepoPolicyPathLengthUpdate(t *testing.T) {
	pathLengthTfNode := "azuredevops_repository_policy_max_path_length.test"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepoPolicyPathLengthBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(pathLengthTfNode, "enabled", "true"),
				),
			}, {
				Config: hclRepoPolicyPathLengthUpdate(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(pathLengthTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(pathLengthTfNode, "max_path_length", "1000"),
				),
			}, {
				ResourceName:      pathLengthTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(pathLengthTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProjectPolicyPathLengthBasic(t *testing.T) {
	pathLengthTfNode := "azuredevops_repository_policy_max_path_length.test"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclProjectPolicyPathLengthBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(pathLengthTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(pathLengthTfNode, "max_path_length", "500"),
				),
			}, {
				ResourceName:      pathLengthTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(pathLengthTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProjectPolicyPathLengthUpdate(t *testing.T) {
	pathLengthTfNode := "azuredevops_repository_policy_max_path_length.test"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclProjectPolicyPathLengthBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(pathLengthTfNode, "enabled", "true"),
				),
			}, {
				Config: hclProjectPolicyPathLengthUpdate(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(pathLengthTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(pathLengthTfNode, "max_path_length", "1000"),
				),
			}, {
				ResourceName:      pathLengthTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(pathLengthTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func hclPolicyPathLengthResourceTemplate(projectName string, repoName string) string {
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

func hclRepoPolicyPathLengthBasic(projectName string, repoName string) string {
	projectAndRepo := hclPolicyPathLengthResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`
%s

resource "azuredevops_repository_policy_max_path_length" "test" {
  project_id      = azuredevops_project.test.id
  enabled         = true
  blocking        = true
  max_path_length = 500
  repository_ids  = [azuredevops_git_repository.test.id]
}`, projectAndRepo)
}

func hclRepoPolicyPathLengthUpdate(projectName string, repoName string) string {
	projectAndRepo := hclPolicyPathLengthResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`
%s

resource "azuredevops_repository_policy_max_path_length" "test" {
  project_id      = azuredevops_project.test.id
  enabled         = true
  blocking        = true
  max_path_length = 1000
  repository_ids  = [azuredevops_git_repository.test.id]
}`, projectAndRepo)
}

func hclProjectPolicyPathLengthBasic(projectName string, repoName string) string {
	projectAndRepo := hclPolicyPathLengthResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`
%s

resource "azuredevops_repository_policy_max_path_length" "test" {
  project_id      = azuredevops_project.test.id
  enabled         = true
  blocking        = true
  max_path_length = 500
  depends_on      = [azuredevops_git_repository.test]
}`, projectAndRepo)
}

func hclProjectPolicyPathLengthUpdate(projectName string, repoName string) string {
	projectAndRepo := hclPolicyPathLengthResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`
%s

resource "azuredevops_repository_policy_max_path_length" "test" {
  project_id      = azuredevops_project.test.id
  enabled         = true
  blocking        = true
  max_path_length = 1000
  depends_on      = [azuredevops_git_repository.test]
}`, projectAndRepo)
}
