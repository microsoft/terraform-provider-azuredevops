//go:build (all || resource_branchpolicy_auto_reviewers) && !exclude_resource_branchpolicy_auto_reviewers
// +build all resource_branchpolicy_auto_reviewers
// +build !exclude_resource_branchpolicy_auto_reviewers

package branch

import (
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/policy"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

// verifies that the flatten/expand round trip path produces repeatable results
func TestBranchPolicyAutoReviewers_ExpandFlatten_Roundtrip(t *testing.T) {
	var projectID = uuid.New().String()
	var randomUUID = uuid.New()
	var testPolicy = &policy.PolicyConfiguration{
		Id:         converter.Int(1),
		IsEnabled:  converter.Bool(true),
		IsBlocking: converter.Bool(true),
		Type: &policy.PolicyTypeRef{
			Id: &randomUUID,
		},
		Settings: map[string]interface{}{
			"scope": []map[string]interface{}{
				{
					"repositoryId": "test-repo-id",
					"refName":      "test-ref-name",
					"matchKind":    "test-match-kind",
				},
			},
			"creatorVoteCounts":    false,
			"filenamePatterns":     []string{"*"},
			"requiredReviewerIds":  []string{"some-group"},
			"minimumApproverCount": 1,
			"message":              "",
		},
	}

	resourceData := schema.TestResourceDataRaw(t, ResourceBranchPolicyAutoReviewers().Schema, nil)
	resourceData.SetId(strconv.Itoa(*testPolicy.Id))
	err := autoReviewersFlattenFunc(resourceData, testPolicy, &projectID)
	require.Nil(t, err)
	expandedPolicy, expandedProjectID, err := autoReviewersExpandFunc(resourceData, randomUUID)
	require.Nil(t, err)

	require.Equal(t, testPolicy, expandedPolicy)
	require.Equal(t, projectID, *expandedProjectID)
}
