//go:build (all || resource_branchpolicy_min_reviewers) && !exclude_resource_branchpolicy_min_reviewers
// +build all resource_branchpolicy_min_reviewers
// +build !exclude_resource_branchpolicy_min_reviewers

package branch

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/policy"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// verifies that the flatten/expand round trip path produces repeatable results
func TestBranchPolicyMinReviewers_ExpandFlatten_Roundtrip(t *testing.T) {
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
			"minimumApproverCount":        2,
			"creatorVoteCounts":           true,
			"allowDownvotes":              true,
			"resetOnSourcePush":           true,
			"requireVoteOnLastIteration":  true,
			"resetRejectionsOnSourcePush": true,
			"blockLastPusherVote":         true,
		},
	}

	resourceData := schema.TestResourceDataRaw(t, ResourceBranchPolicyMinReviewers().Schema, nil)
	err := minReviewersFlattenFunc(resourceData, testPolicy, &projectID)
	require.Nil(t, err)
	expandedPolicy, expandedProjectID, err := minReviewersExpandFunc(resourceData, randomUUID)
	require.Nil(t, err)

	require.Equal(t, testPolicy, expandedPolicy)
	require.Equal(t, projectID, *expandedProjectID)
}
