// +build all resource_branchpolicy_acceptance_test
// +build !exclude_resource_branchpolicy_acceptance_test

package acceptancetests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccBranchPolicyMinReviewers_CreateAndUpdate(t *testing.T) {
	minReviewerTfNode := "azuredevops_branch_policy_min_reviewers.p"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: getMinReviewersHcl(true, true, 1, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(minReviewerTfNode, "id"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "blocking", "true"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "enabled", "true"),
				),
			}, {
				Config: getMinReviewersHcl(false, false, 2, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(minReviewerTfNode, "id"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "blocking", "false"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "enabled", "false"),
				),
			}, {
				ResourceName:      minReviewerTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(minReviewerTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func getMinReviewersHcl(enabled bool, blocking bool, reviewers int, submitterCanVote bool) string {
	minReviewCountPolicy := fmt.Sprintf(`
	resource "azuredevops_branch_policy_min_reviewers" "p" {
		project_id = azuredevops_project.project.id
		enabled  = %t
		blocking = %t
		settings {
			reviewer_count     = %d
			submitter_can_vote = %t
			scope {
				repository_id  = azuredevops_git_repository.repository.id
				repository_ref = azuredevops_git_repository.repository.default_branch
				match_type     = "exact"
			}
		}
	}
	`, enabled, blocking, reviewers, submitterCanVote)

	return strings.Join(
		[]string{
			getProjectRepoBuildUserEntitlementResource(),
			minReviewCountPolicy,
		},
		"\n",
	)
}

func TestAccBranchPolicyAutoReviewers_CreateAndUpdate(t *testing.T) {
	autoReviewerTfNode := "azuredevops_branch_policy_auto_reviewers.p"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: getAutoReviewersHcl(true, true, false, "auto reviewer", fmt.Sprintf("\"%s\",\"%s\"", "*/API*.cs", "README.md")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(autoReviewerTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(autoReviewerTfNode, "blocking", "true"),
				),
			}, {
				Config: getAutoReviewersHcl(false, false, true, "new auto reviewer", fmt.Sprintf("\"%s\",\"%s\"", "*/API*.cs", "README.md")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(autoReviewerTfNode, "enabled", "false"),
					resource.TestCheckResourceAttr(autoReviewerTfNode, "blocking", "false"),
				),
			}, {
				ResourceName:      autoReviewerTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(autoReviewerTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func getAutoReviewersHcl(enabled bool, blocking bool, submitterCanVote bool, message string, pathFilters string) string {
	autoReviewerPolicy := fmt.Sprintf(`
	resource "azuredevops_branch_policy_auto_reviewers" "p" {
		project_id = azuredevops_project.project.id
		enabled  = %t
		blocking = %t
		settings {
			auto_reviewer_ids     = [azuredevops_user_entitlement.user.id]
			submitter_can_vote = %t
			message = "%s"
			path_filters = [%s]
			scope {
				repository_id  = azuredevops_git_repository.repository.id
				repository_ref = azuredevops_git_repository.repository.default_branch
				match_type     = "exact"
			}
		}
	}
	`, enabled, blocking, submitterCanVote, message, pathFilters)

	return strings.Join(
		[]string{
			getProjectRepoBuildUserEntitlementResource(),
			autoReviewerPolicy,
		},
		"\n",
	)
}

func TestAccBranchPolicyBuildValidation_CreateAndUpdate(t *testing.T) {
	buildValidationTfNode := "azuredevops_branch_policy_build_validation.p"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: getBuildValidationHcl(true, true, "build validation", 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(buildValidationTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(buildValidationTfNode, "settings.0.filename_patterns.#", "3"),
				),
			}, {
				Config: getBuildValidationHcl(false, false, "build validation rename", 720),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(buildValidationTfNode, "enabled", "false"),
					resource.TestCheckResourceAttr(buildValidationTfNode, "settings.0.filename_patterns.#", "3"),
				),
			}, {
				ResourceName:      buildValidationTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(buildValidationTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func getBuildValidationHcl(enabled bool, blocking bool, displayName string, validDuration int) string {
	buildValidationPolicy := fmt.Sprintf(`
	resource "azuredevops_branch_policy_build_validation" "p" {
		project_id = azuredevops_project.project.id
		enabled  = %t
		blocking = %t
		settings {
			display_name = "%s"
			valid_duration = %d
			build_definition_id = azuredevops_build_definition.build.id
			filename_patterns =  [
				"/WebApp/*",
				"!/WebApp/Tests/*",
				"*.cs"
			]
			scope {
				repository_id  = azuredevops_git_repository.repository.id
				repository_ref = azuredevops_git_repository.repository.default_branch
				match_type     = "exact"
			}
		}
	}
	`, enabled, blocking, displayName, validDuration)

	return strings.Join(
		[]string{
			getProjectRepoBuildUserEntitlementResource(),
			buildValidationPolicy,
		},
		"\n",
	)
}

func getProjectRepoBuildUserEntitlementResource() string {
	projectAndRepo := testutils.HclGitRepoResource(testutils.GenerateResourceName(), testutils.GenerateResourceName(), "Clean")
	userEntitlement := testutils.HclUserEntitlementResource("acc@test.com")
	buildDef := testutils.HclBuildDefinitionResource(
		"Sample Build Definition",
		`\`,
		"TfsGit",
		"${azuredevops_git_repository.repository.id}",
		"master",
		"path/to/yaml",
		"")

	return strings.Join(
		[]string{
			projectAndRepo,
			userEntitlement,
			buildDef,
		},
		"\n",
	)
}
