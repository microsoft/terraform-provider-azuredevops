//go:build (all || resource_branchpolicy_acceptance_test || policy) && (!exclude_resource_branchpolicy_acceptance_test || !exclude_policy)
// +build all resource_branchpolicy_acceptance_test policy
// +build !exclude_resource_branchpolicy_acceptance_test !exclude_policy

package acceptancetests

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// TestAccBranchPolicyMinReviewers_CreateAndUpdate - acceptance test for min reviewers branch policy attributes
func TestAccBranchPolicyMinReviewers_CreateAndUpdate(t *testing.T) {
	minReviewerTfNode := "azuredevops_branch_policy_min_reviewers.p"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: getMinReviewersHcl(true, true, 1, false, "\"refs/heads/release\"", "Exact"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(minReviewerTfNode, "id"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "blocking", "true"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "settings.0.submitter_can_vote", "false"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "settings.0.allow_completion_with_rejects_or_waits", "false"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "settings.0.last_pusher_cannot_approve", "false"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "settings.0.on_last_iteration_require_vote", "false"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "settings.0.on_last_iteration_require_vote", "false"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "settings.0.on_push_reset_approved_votes", "true"),
				),
			}, {
				Config: getMinReviewersHcl(false, false, 2, true, "\"refs/heads/release\"", "Exact"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(minReviewerTfNode, "id"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "blocking", "false"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "enabled", "false"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "settings.0.submitter_can_vote", "true"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "settings.0.allow_completion_with_rejects_or_waits", "true"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "settings.0.last_pusher_cannot_approve", "true"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "settings.0.on_last_iteration_require_vote", "true"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "settings.0.on_last_iteration_require_vote", "true"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "settings.0.on_push_reset_all_votes", "true"),
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

func getMinReviewersHcl(enabled bool, blocking bool, reviewers int, flag bool, repositoryRef string, matchType string) string {
	votes := "all"
	if !flag {
		votes = "approved"
	}
	settings := fmt.Sprintf(
		`
		reviewer_count     = %[1]d
		submitter_can_vote = %[2]t
		allow_completion_with_rejects_or_waits =%[2]t
		last_pusher_cannot_approve = %[2]t
		on_last_iteration_require_vote = %[2]t
		on_push_reset_%[3]s_votes = true
		`, reviewers, flag, votes)

	return getBranchPolicyHcl("azuredevops_branch_policy_min_reviewers", enabled, blocking, settings, "azuredevops_git_repository.repository.id", repositoryRef, matchType)
}

func TestAccBranchPolicyAutoReviewers_CreateAndUpdate(t *testing.T) {
	autoReviewerTfNode := "azuredevops_branch_policy_auto_reviewers.p"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, &[]string{"AZDO_TEST_AAD_USER_EMAIL"}) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: getAutoReviewersHcl(true, true, false, "auto reviewer", fmt.Sprintf("\"%s\",\"%s\"", "*/API*.cs", "README.md"), "\"refs/heads/release\"", "Exact"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(autoReviewerTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(autoReviewerTfNode, "blocking", "true"),
				),
			}, {
				Config: getAutoReviewersHcl(false, false, true, "new auto reviewer", fmt.Sprintf("\"%s\",\"%s\"", "*/API*.cs", "README.md"), "\"refs/heads/release\"", "Exact"),
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

func getAutoReviewersHcl(enabled bool, blocking bool, submitterCanVote bool, message string, pathFilters string, repositoryRef string, matchType string) string {
	settings := fmt.Sprintf(
		`
		auto_reviewer_ids  = [azuredevops_user_entitlement.user.id]
		submitter_can_vote = %t
		message 		   = "%s"
		path_filters       = [%s]
		`, submitterCanVote, message, pathFilters,
	)
	userPrincipalName := os.Getenv("AZDO_TEST_AAD_USER_EMAIL")
	userEntitlement := testutils.HclUserEntitlementResource(userPrincipalName)

	return strings.Join(
		[]string{
			userEntitlement,
			getBranchPolicyHcl("azuredevops_branch_policy_auto_reviewers", enabled, blocking, settings, "azuredevops_git_repository.repository.id", repositoryRef, matchType),
		},
		"\n",
	)
}

func TestAccBranchPolicyBuildValidation_CreateAndUpdate(t *testing.T) {
	buildValidationTfNode := "azuredevops_branch_policy_build_validation.p"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: getBuildValidationHcl(true, true, "build validation", 0, "\"refs/heads/release\"", "Exact"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(buildValidationTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(buildValidationTfNode, "settings.0.filename_patterns.#", "3"),
				),
			}, {
Config: getBuildValidationHcl(false, false, "build validation rename", 720, "\"refs/heads/release\"", "Exact"),
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

func getBuildValidationHcl(enabled bool, blocking bool, displayName string, validDuration int, repositoryRef string, matchType string) string {
	settings := fmt.Sprintf(
		`
		display_name = "%s"
		valid_duration = %d
		build_definition_id = azuredevops_build_definition.build.id
		filename_patterns =  [
			"/WebApp/*",
			"!/WebApp/Tests/*",
			"*.cs"
		]
		`, displayName, validDuration,
	)

	return getBranchPolicyHcl("azuredevops_branch_policy_build_validation", enabled, blocking, settings, "azuredevops_git_repository.repository.id", repositoryRef, matchType)
}

func TestAccBranchPolicyWorkItemLinking_CreateAndUpdate(t *testing.T) {
	resourceName := "azuredevops_branch_policy_work_item_linking"
	workItemLinkingTfNode := fmt.Sprintf("%s.p", resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: getBranchPolicyHcl(resourceName, true, true, "", "azuredevops_git_repository.repository.id", "\"refs/heads/release\"", "Exact"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(workItemLinkingTfNode, "enabled", "true"),
				),
			}, {
				Config: getBranchPolicyHcl(resourceName, false, false, "", "azuredevops_git_repository.repository.id", "\"refs/heads/release\"", "Exact"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(workItemLinkingTfNode, "enabled", "false"),
				),
			}, {
				ResourceName:      workItemLinkingTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(workItemLinkingTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBranchPolicyCommentResolution_CreateAndUpdate(t *testing.T) {
	resourceName := "azuredevops_branch_policy_comment_resolution"
	workItemLinkingTfNode := fmt.Sprintf("%s.p", resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: getBranchPolicyHcl(resourceName, true, true, "", "azuredevops_git_repository.repository.id", "\"refs/heads/release\"", "Exact"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(workItemLinkingTfNode, "enabled", "true"),
				),
			}, {
				Config: getBranchPolicyHcl(resourceName, false, false, "", "azuredevops_git_repository.repository.id", "\"refs/heads/release\"", "Exact"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(workItemLinkingTfNode, "enabled", "false"),
				),
			}, {
				ResourceName:      workItemLinkingTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(workItemLinkingTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBranchPolicyMergeTypes_CreateAndUpdate(t *testing.T) {
	buildValidationTfNode := "azuredevops_branch_policy_merge_types.p"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: getMergeTypesHcl(true, true, true, true, true, true, "\"refs/heads/release\"", "Exact"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(buildValidationTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(buildValidationTfNode, "settings.0.allow_squash", "true"),
					resource.TestCheckResourceAttr(buildValidationTfNode, "settings.0.allow_rebase_and_fast_forward", "true"),
					resource.TestCheckResourceAttr(buildValidationTfNode, "settings.0.allow_basic_no_fast_forward", "true"),
					resource.TestCheckResourceAttr(buildValidationTfNode, "settings.0.allow_rebase_with_merge", "true"),
				),
			}, {
				Config: getMergeTypesHcl(false, false, false, false, false, false, "\"refs/heads/release\"", "Exact"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(buildValidationTfNode, "enabled", "false"),
					resource.TestCheckResourceAttr(buildValidationTfNode, "settings.0.allow_squash", "false"),
					resource.TestCheckResourceAttr(buildValidationTfNode, "settings.0.allow_rebase_and_fast_forward", "false"),
					resource.TestCheckResourceAttr(buildValidationTfNode, "settings.0.allow_basic_no_fast_forward", "false"),
					resource.TestCheckResourceAttr(buildValidationTfNode, "settings.0.allow_rebase_with_merge", "false"),
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

func getMergeTypesHcl(enabled bool, blocking bool, allowSquash bool, allowRebase bool, allowNoFastForward bool, allowRebaseMerge bool, repositoryRef string, matchType string) string {
	settings := fmt.Sprintf(
		`
		allow_squash = %t
		allow_rebase_and_fast_forward = %t
		allow_basic_no_fast_forward = %t
		allow_rebase_with_merge = %t
		`, allowSquash, allowRebase, allowNoFastForward, allowRebaseMerge,
	)

	return getBranchPolicyHcl("azuredevops_branch_policy_merge_types", enabled, blocking, settings, "azuredevops_git_repository.repository.id", repositoryRef, matchType)
}

func getBranchPolicyHcl(resourceName string, enabled bool, blocking bool, settings string, repositoryId string, repositoryRef string, matchType string) string {
	branchPolicy := fmt.Sprintf(`
	resource "%s" "p" {
		project_id = azuredevops_project.project.id
		enabled  = %t
		blocking = %t
		settings {
			%s
			scope {
				repository_id  = %s
				repository_ref = %s
				match_type     = "%s"
			}
		}
	}
	`, resourceName, enabled, blocking, settings, repositoryId, repositoryRef, matchType)
	projectAndRepo := testutils.HclGitRepoResource(testutils.GenerateResourceName(), testutils.GenerateResourceName(), "Clean")
	buildDef := testutils.HclBuildDefinitionResource(
		"Sample Build Definition",
		`\\`,
		"TfsGit",
		"${azuredevops_git_repository.repository.id}",
		"master",
		"path/to/yaml",
		"")

	return strings.Join(
		[]string{
			branchPolicy,
			projectAndRepo,
			buildDef,
		},
		"\n",
	)
}

func getStatusCheckHcl(enabled bool, blocking bool, name string, invalidateOnUpdate bool, applicability string, repositoryId string, repositoryRef string, matchType string) string {
	settings := fmt.Sprintf(
		`
		name = "%s"
		invalidate_on_update = %t
		applicability = "%s"
		filename_patterns =  [
			"/WebApp/*",
			"!/WebApp/Tests/*",
			"*.cs"
		]
		`, name, invalidateOnUpdate, applicability,
	)

	return getBranchPolicyHcl("azuredevops_branch_policy_status_check", enabled, blocking, settings, repositoryId, repositoryRef, matchType)
}

func TestAccBranchPolicyStatusCheck_CreateAndUpdate(t *testing.T) {
	statusCheckTfNode := "azuredevops_branch_policy_status_check.p"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: getStatusCheckHcl(true, true, "abc-1", true, "default", "null", "null", "defaultBranch"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(statusCheckTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(statusCheckTfNode, "blocking", "true"),
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.name", "abc-1"),
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.invalidate_on_update", "true"),
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.applicability", "default"),
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.filename_patterns.#", "3"),
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.scope.0.repository_id", ""),
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.scope.0.repository_ref", ""),
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.scope.0.match_type", "DefaultBranch"),
				),
			}, {
				Config: getStatusCheckHcl(false, false, "abc-2", false, "conditional", "null", "\"refs/heads/release\"", "Prefix"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(statusCheckTfNode, "enabled", "false"),
					resource.TestCheckResourceAttr(statusCheckTfNode, "blocking", "false"),
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.name", "abc-2"),
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.invalidate_on_update", "false"),
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.applicability", "conditional"),
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.scope.0.repository_id", ""),
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.scope.0.repository_ref", "refs/heads/release"),
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.scope.0.match_type", "Prefix"),
				),
			}, {
				Config: getStatusCheckHcl(false, false, "abc-3", false, "conditional", "null", "\"refs/heads/release\"", "Exact"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.name", "abc-3"),
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.scope.0.match_type", "Exact"),
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.scope.0.repository_id", ""),
					resource.TestCheckResourceAttr(statusCheckTfNode, "settings.0.scope.0.repository_ref", "refs/heads/release"),
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
