//go:build (all || resource_policy_case_enforcement) && !resource_policy_case_enforcement
// +build all resource_policy_case_enforcement
// +build !resource_policy_case_enforcement

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

const caseEnforceTfNode = "azuredevops_repository_policy_case_enforcement.p"

func TestAccPolicyCaseEnforcement(t *testing.T) {
	testutils.RunTestsInSequence(t, map[string]map[string]func(t *testing.T){
		"RepositoryPolicies": {
			"basic":  testAccRepoPolicyEnforceConsistentCaseBasic,
			"update": testAccRepoPolicyEnforceConsistentCaseUpdate,
		},
		"ProjectPolicies": {
			"basic":  testAccProjectPolicyEnforceConsistentCaseBasic,
			"update": testAccProjectPolicyEnforceConsistentCaseUpdate,
		},
	})
}

func testAccRepoPolicyEnforceConsistentCaseBasic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepoPolicyEnforceConsistentCaseBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(caseEnforceTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(caseEnforceTfNode, "enforce_consistent_case", "true"),
				),
			}, {
				ResourceName:      caseEnforceTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(caseEnforceTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccRepoPolicyEnforceConsistentCaseUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclRepoPolicyEnforceConsistentCaseBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(caseEnforceTfNode, "enabled", "true"),
				),
			}, {
				Config: hclRepoPolicyEnforceConsistentCaseUpdate(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(caseEnforceTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(caseEnforceTfNode, "enforce_consistent_case", "false"),
				),
			}, {
				ResourceName:      caseEnforceTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(caseEnforceTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProjectPolicyEnforceConsistentCaseBasic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclProjectPolicyEnforceConsistentCaseBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(caseEnforceTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(caseEnforceTfNode, "enforce_consistent_case", "true"),
				),
			}, {
				ResourceName:      caseEnforceTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(caseEnforceTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProjectPolicyEnforceConsistentCaseUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclProjectPolicyEnforceConsistentCaseBasic(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(caseEnforceTfNode, "enabled", "true"),
				),
			}, {
				Config: hclProjectPolicyEnforceConsistentCaseUpdate(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(caseEnforceTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(caseEnforceTfNode, "enforce_consistent_case", "false"),
				),
			}, {
				ResourceName:      caseEnforceTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(caseEnforceTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func hclPolicyEnforceConsistentCaseResourceTemplate(projectName string, repoName string) string {
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

func hclRepoPolicyEnforceConsistentCaseBasic(projectName string, repoName string) string {
	projectAndRepo := hclPolicyEnforceConsistentCaseResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`%s %s`, projectAndRepo, `
resource "azuredevops_repository_policy_case_enforcement" "p" {
  project_id = azuredevops_project.p.id

  enabled  = true
  blocking = true
  enforce_consistent_case = true
  repository_ids  = [azuredevops_git_repository.r.id]
}
`)
}

func hclRepoPolicyEnforceConsistentCaseUpdate(projectName string, repoName string) string {
	projectAndRepo := hclPolicyEnforceConsistentCaseResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`%s %s`, projectAndRepo, `
resource "azuredevops_repository_policy_case_enforcement" "p" {
 project_id = azuredevops_project.p.id

 enabled  = true
 blocking = true
 enforce_consistent_case = false
 repository_ids  = [azuredevops_git_repository.r.id]
}
`)
}

func hclProjectPolicyEnforceConsistentCaseBasic(projectName string, repoName string) string {
	projectAndRepo := hclPolicyEnforceConsistentCaseResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`%s %s`, projectAndRepo, `
resource "azuredevops_repository_policy_case_enforcement" "p" {
  project_id = azuredevops_project.p.id
  enabled  = true
  blocking = true
  enforce_consistent_case = true
  depends_on = [azuredevops_git_repository.r]
}
`)
}

func hclProjectPolicyEnforceConsistentCaseUpdate(projectName string, repoName string) string {
	projectAndRepo := hclPolicyEnforceConsistentCaseResourceTemplate(projectName, repoName)
	return fmt.Sprintf(`%s %s`, projectAndRepo, `
resource "azuredevops_repository_policy_case_enforcement" "p" {
 project_id = azuredevops_project.p.id

 enabled  = true
 blocking = true
 enforce_consistent_case = false
 depends_on = [azuredevops_git_repository.r]
}
`)
}
