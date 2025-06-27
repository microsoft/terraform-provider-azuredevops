//go:build (all || resource_branchpolicy_merge_types) && !exclude_resource_branchpolicy_merge_types
// +build all resource_branchpolicy_merge_types
// +build !exclude_resource_branchpolicy_merge_types

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
func TestBranchPolicyMergeTypes_ExpandFlatten_Roundtrip(t *testing.T) {
	projectID := uuid.New().String()
	randomUUID := uuid.New()
	testPolicy := &policy.PolicyConfiguration{
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
			"allowSquash":        true,
			"allowRebase":        true,
			"allowNoFastForward": true,
			"allowRebaseMerge":   true,
		},
	}

	resourceData := schema.TestResourceDataRaw(t, ResourceBranchPolicyMergeTypes().Schema, nil)
	resourceData.SetId(strconv.Itoa(*testPolicy.Id))
	err := mergeTypesFlattenFunc(resourceData, testPolicy, &projectID)
	require.Nil(t, err)
	expandedPolicy, expandedProjectID, err := mergeTypesExpandFunc(resourceData, randomUUID)
	require.Nil(t, err)

	require.Equal(t, testPolicy, expandedPolicy)
	require.Equal(t, projectID, *expandedProjectID)
}
