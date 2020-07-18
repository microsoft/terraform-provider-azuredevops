// +build all resource_branchpolicy_build_validation
// +build !exclude_resource_branchpolicy_build_validation

package policy

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/policy"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
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
	err := buildValidationFlattenFunc(resourceData, testPolicy, &projectID)
	require.Nil(t, err)
	expandedPolicy, expandedProjectID, err := buildValidationExpandFunc(resourceData, randomUUID)
	require.Nil(t, err)

	require.Equal(t, testPolicy, expandedPolicy)
	require.Equal(t, projectID, *expandedProjectID)
}
