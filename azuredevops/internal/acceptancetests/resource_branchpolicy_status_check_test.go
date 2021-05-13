// +build all resource_branchpolicy_status_check_acceptance_test policy
// +build !exclude_resource_branchpolicy_status_check_acceptance_test !exclude_policy

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccBranchPolicyStatusCheckBasic(t *testing.T) {
	statusCheckTfNode := "azuredevops_branch_policy_status_check.p"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclBranchPolicyStatusCheckResourceBasic(projectName, repoName, "update"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.name", "update")),
			}, {
				ResourceName:      statusCheckTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(statusCheckTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBranchPolicyStatusCheckComplete(t *testing.T) {
	statusCheckTfNode := "azuredevops_branch_policy_status_check.p"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclBranchPolicyStatusCheckResourceComplete(projectName, repoName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(statusCheckTfNode, "settings.0.author_id"),
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.name", "Release"),
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.invalidate_on_update", "true"),
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.applicability", "conditional"),
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.display_name", "PreCheck"),
				),
			}, {
				ResourceName:      statusCheckTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(statusCheckTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBranchPolicyStatusCheckUpdate(t *testing.T) {
	statusCheckTfNode := "azuredevops_branch_policy_status_check.p"
	projectName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclBranchPolicyStatusCheckResourceBasic(projectName, repoName, "update"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.name", "update"),
				),
			}, {
				Config: hclBranchPolicyStatusCheckResourceUpdate(projectName, repoName, "releaseCheck", true, "conditional", "updateName"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(statusCheckTfNode, "settings.0.author_id"),
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.name", "releaseCheck"),
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.invalidate_on_update", "true"),
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.applicability", "conditional"),
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.display_name", "updateName"),
				),
			}, {
				ResourceName:      statusCheckTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(statusCheckTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func hclBranchPolicyStatusCheckResourceTemplate(projectName string, repoName string) string {
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

func hclBranchPolicyStatusCheckResourceBasic(projectName string, repoName string, statusName string) string {
	projectAndRepo := hclBranchPolicyStatusCheckResourceTemplate(projectName, repoName)
	statusCheck := fmt.Sprintf(`
resource "azuredevops_branch_policy_status_check" "p" {
  project_id = azuredevops_project.p.id

  enabled  = true
  blocking = true

  settings {
	name = "%s"
    scope {
      repository_id  = azuredevops_git_repository.r.id
      repository_ref = azuredevops_git_repository.r.default_branch
      match_type     = "Exact"
    }
  }
}
`, statusName)

	return fmt.Sprintf(`%s %s`, projectAndRepo, statusCheck)
}

func hclBranchPolicyStatusCheckResourceComplete(projectName string, repoName string) string {
	return fmt.Sprintf(
		`%s %s`,
		hclBranchPolicyStatusCheckResourceTemplate(projectName, repoName), `

resource "azuredevops_user_entitlement" "user" {
  principal_name       = "mail@email.com"
  account_license_type = "basic"
}

resource "azuredevops_branch_policy_status_check" "p" {
 project_id = azuredevops_project.p.id

 enabled  = true
 blocking = true

 settings {
	name = "Release"
	author_id            = azuredevops_user_entitlement.user.id
	invalidate_on_update = true
	applicability = "conditional"
	display_name = "PreCheck"
	filename_patterns = ["*.go","**.ts"]

   scope {
     repository_id  = azuredevops_git_repository.r.id
     repository_ref = azuredevops_git_repository.r.default_branch
     match_type     = "Exact"
   }
 }
}
`)
}

func hclBranchPolicyStatusCheckResourceUpdate(projectName string, repoName string,
	statusName string, invalid bool, applicability string, displayName string) string {

	statusCheck := fmt.Sprintf(
		`
data "azuredevops_group" "group" {
  project_id = azuredevops_project.p.id
  name       = "Project Administrators"
}

resource "azuredevops_branch_policy_status_check" "p" {
 project_id = azuredevops_project.p.id

 enabled  = true
 blocking = true

 settings {
	name = "%s"
	author_id = data.azuredevops_group.group.origin_id
	invalidate_on_update = %t
	applicability = "%s"
	display_name = "%s"
	filename_patterns = ["*.go","**.ts"]

   scope {
     repository_id  = azuredevops_git_repository.r.id
     repository_ref = azuredevops_git_repository.r.default_branch
     match_type     = "Exact"
   }
 }
}
`, statusName, invalid, applicability, displayName)

	return fmt.Sprintf(
		`%s %s`,
		hclBranchPolicyStatusCheckResourceTemplate(projectName, repoName),
		statusCheck)
}
