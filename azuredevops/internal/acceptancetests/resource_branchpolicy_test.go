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

// Verifies that the following sequence of events occurrs without error:
//	(1) Branch policies can be created with no errors
//	(2) Branch policies can be updated with no errors
//	(3) Branch policies can be deleted with no errors
func TestAccBranchPolicy_CreateAndUpdate(t *testing.T) {
	projName := testutils.GenerateResourceName()
	repoName := testutils.GenerateResourceName()
	opts1 := hclOptions{
		projectName:            projName,
		repoName:               repoName,
		minReviewerOptions:     minReviewPolicyOpts{true, true, 1, false},
		buildValidationOptions: buildValidationPolicyOpts{true, true, "build validation", 0},
	}

	opts2 := hclOptions{
		projectName:            projName,
		repoName:               repoName,
		minReviewerOptions:     minReviewPolicyOpts{false, false, 2, true},
		buildValidationOptions: buildValidationPolicyOpts{false, false, "build validation rename", 720},
	}

	minReviewerTfNode := "azuredevops_branch_policy_min_reviewers.p"
	buildVlidationTfNode := "azuredevops_branch_policy_build_validation.p"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: getHCL(opts1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(minReviewerTfNode, "id"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "blocking", "true"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(buildVlidationTfNode, "enabled", "true"),
					resource.TestCheckResourceAttr(buildVlidationTfNode, "enabled", "true"),
					resource.TestCheckResourceAttrSet(buildVlidationTfNode, "filename_patterns"),
				),
			}, {
				Config: getHCL(opts2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(minReviewerTfNode, "id"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "blocking", "false"),
					resource.TestCheckResourceAttr(minReviewerTfNode, "enabled", "false"),
					resource.TestCheckResourceAttr(buildVlidationTfNode, "enabled", "false"),
					resource.TestCheckResourceAttrSet(buildVlidationTfNode, "filename_patterns"),
				),
			}, {
				ResourceName:      minReviewerTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(minReviewerTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			}, {
				ResourceName:      buildVlidationTfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(buildVlidationTfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

type minReviewPolicyOpts struct {
	enabled          bool
	blocking         bool
	reviewers        int
	submitterCanVote bool
}

type buildValidationPolicyOpts struct {
	enabled       bool
	blocking      bool
	displayName   string
	validDuration int
}

type hclOptions struct {
	projectName            string
	repoName               string
	minReviewerOptions     minReviewPolicyOpts
	buildValidationOptions buildValidationPolicyOpts
}

func getHCL(opts hclOptions) string {
	projectAndRepo := testutils.HclGitRepoResource(opts.projectName, opts.repoName, "Clean")
	buildDef := `
	resource "azuredevops_build_definition" "build" {
		project_id = azuredevops_project.project.id
		name       = "Sample Build Definition"

		repository {
			repo_type   = "TfsGit"
			repo_id     = azuredevops_git_repository.gitrepo.id
			yml_path    = "azure-pipelines.yml"
		}
	}
`
	minReviewCountPolicyFmt := `
	resource "azuredevops_branch_policy_min_reviewers" "p" {
		project_id = azuredevops_project.project.id
		enabled  = %t
		blocking = %t
		settings {
			reviewer_count     = %d
			submitter_can_vote = %t
			scope {
				repository_id  = azuredevops_git_repository.gitrepo.id
				repository_ref = azuredevops_git_repository.gitrepo.default_branch
				match_type     = "exact"
			}
		}
	}
`

	minReviewCountPolicy := fmt.Sprintf(
		minReviewCountPolicyFmt,
		opts.minReviewerOptions.enabled,
		opts.minReviewerOptions.blocking,
		opts.minReviewerOptions.reviewers,
		opts.minReviewerOptions.submitterCanVote)

	buildValidationPolicyFmt := `
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
				repository_id  = azuredevops_git_repository.gitrepo.id
				repository_ref = azuredevops_git_repository.gitrepo.default_branch
				match_type     = "exact"
			}
		}
	}
`
	buildValidationPolicyFmt = fmt.Sprintf(
		buildValidationPolicyFmt,
		opts.buildValidationOptions.enabled,
		opts.buildValidationOptions.blocking,
		opts.buildValidationOptions.displayName,
		opts.buildValidationOptions.validDuration)

	return strings.Join(
		[]string{
			projectAndRepo,
			buildDef,
			minReviewCountPolicy,
			buildValidationPolicyFmt,
		},
		"\n",
	)

}
