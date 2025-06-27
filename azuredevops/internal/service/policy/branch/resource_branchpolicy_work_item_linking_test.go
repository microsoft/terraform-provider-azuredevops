//go:build (all || resource_branchpolicy_work_item_linking) && !exclude_resource_branchpolicy_work_item_linking
// +build all resource_branchpolicy_work_item_linking
// +build !exclude_resource_branchpolicy_work_item_linking

package branch

import (
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/policy"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// verifies that the flatten/expand round trip path produces repeatable results
func TestBranchPolicyWorkItemLinking_ExpandFlatten_Roundtrip(t *testing.T) {
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
		},
	}

	resourceData := schema.TestResourceDataRaw(t, ResourceBranchPolicyWorkItemLinking().Schema, nil)
	resourceData.SetId(strconv.Itoa(*testPolicy.Id))
	err := workItemLinkingFlattenFunc(resourceData, testPolicy, &projectID)
	require.Nil(t, err)
	expandedPolicy, expandedProjectID, err := workItemLinkingExpandFunc(resourceData, randomUUID)
	require.Nil(t, err)

	require.Equal(t, testPolicy, expandedPolicy)
	require.Equal(t, projectID, *expandedProjectID)
}
