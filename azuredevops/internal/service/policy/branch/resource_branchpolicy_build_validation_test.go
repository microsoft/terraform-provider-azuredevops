//go:build (all || resource_branchpolicy_build_validation) && !exclude_resource_branchpolicy_build_validation
// +build all resource_branchpolicy_build_validation
// +build !exclude_resource_branchpolicy_build_validation

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
func TestBranchPolicyBuildValidation_ExpandFlatten_Roundtrip(t *testing.T) {
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
			"buildDefinitionId":       77,
			"displayName":             "test policy",
			"manualQueueOnly":         true,
			"queueOnSourceUpdateOnly": true,
			"validDuration":           700,
			"filenamePatterns":        &([]string{"*md"}),
		},
	}

	resourceData := schema.TestResourceDataRaw(t, ResourceBranchPolicyBuildValidation().Schema, nil)
	resourceData.SetId(strconv.Itoa(*testPolicy.Id))
	err := buildValidationFlattenFunc(resourceData, testPolicy, &projectID)
	require.Nil(t, err)
	expandedPolicy, expandedProjectID, err := buildValidationExpandFunc(resourceData, randomUUID)
	require.Nil(t, err)

	require.Equal(t, testPolicy, expandedPolicy)
	require.Equal(t, projectID, *expandedProjectID)
}
