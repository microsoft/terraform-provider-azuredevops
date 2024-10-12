//go:build (all || resource_policy_reserved_names) && !resource_policy_reserved_names
// +build all resource_policy_reserved_names
// +build !resource_policy_reserved_names

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccRepositoryPolicyReservedNames(t *testing.T) {
	testutils.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"RepositoryPolicies": {
			"basic":  testAccRepoPolicyReservedNamesBasic,
			"update": testAccRepoPolicyReservedNamesUpdate,
		},
		"ProjectPolicies": {
			"basic":  testAccProjectPolicyReservedNamesBasic,
			"update": testAccProjectPolicyReservedNamesUpdate,
		},
	})
}

func testAccRepoPolicyReservedNamesBasic(t *testing.T) {
	reservedNameTfNode := "azuredevops_repository_policy_reserved_names.test"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepoPolicyReservedNamesBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(reservedNameTfNode, "enabled", "true"),
				),
			}, {
				ResourceName:      reservedNameTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(reservedNameTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRepoPolicyReservedNamesUpdate(t *testing.T) {
	reservedNameTfNode := "azuredevops_repository_policy_reserved_names.test"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepoPolicyReservedNamesBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(reservedNameTfNode, "enabled", "true"),
				),
			}, {
				Config: hclRepoPolicyReservedNamesUpdate(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(reservedNameTfNode, "enabled", "false"),
				),
			}, {
				ResourceName:      reservedNameTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(reservedNameTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProjectPolicyReservedNamesBasic(t *testing.T) {
	reservedNameTfNode := "azuredevops_repository_policy_reserved_names.test"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclProjectPolicyReservedNamesBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(reservedNameTfNode, "enabled", "true"),
				),
			}, {
				ResourceName:      reservedNameTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(reservedNameTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProjectPolicyReservedNamesUpdate(t *testing.T) {
	reservedNameTfNode := "azuredevops_repository_policy_reserved_names.test"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclProjectPolicyReservedNamesBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(reservedNameTfNode, "enabled", "true"),
				),
			}, {
				Config: hclProjectPolicyReservedNamesUpdate(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(reservedNameTfNode, "enabled", "false"),
				),
			}, {
				ResourceName:      reservedNameTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(reservedNameTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func hclPolicyReservedNamesResourceTemplate(projectName string, repoName string) string {
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

func hclRepoPolicyReservedNamesBasic(projectName string, repoName string) string {
	projectAndRepo := hclPolicyReservedNamesResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`
%s

resource "azuredevops_repository_policy_reserved_names" "test" {
  project_id     = azuredevops_project.test.id
  enabled        = true
  blocking       = true
  repository_ids = [azuredevops_git_repository.test.id]
}`, projectAndRepo)
}

func hclRepoPolicyReservedNamesUpdate(projectName string, repoName string) string {
	projectAndRepo := hclPolicyReservedNamesResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`
%s

resource "azuredevops_repository_policy_reserved_names" "test" {
  project_id     = azuredevops_project.test.id
  enabled        = false
  blocking       = true
  repository_ids = [azuredevops_git_repository.test.id]
}`, projectAndRepo)
}

func hclProjectPolicyReservedNamesBasic(projectName string, repoName string) string {
	projectAndRepo := hclPolicyReservedNamesResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`
%s

resource "azuredevops_repository_policy_reserved_names" "test" {
  project_id = azuredevops_project.test.id
  enabled    = true
  blocking   = true
  depends_on = [azuredevops_git_repository.test]
}`, projectAndRepo)
}

func hclProjectPolicyReservedNamesUpdate(projectName string, repoName string) string {
	projectAndRepo := hclPolicyReservedNamesResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`
%s

resource "azuredevops_repository_policy_reserved_names" "test" {
  project_id = azuredevops_project.test.id
  enabled    = false
  blocking   = true
  depends_on = [azuredevops_git_repository.test]
}`, projectAndRepo)
}
