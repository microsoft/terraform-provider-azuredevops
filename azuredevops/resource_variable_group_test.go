// +build all resource_variable_group
// +build !exclude_resource_variable_group

package azuredevops

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/build"
	"github.com/microsoft/azure-devops-go-api/azuredevops/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/stretchr/testify/require"
)

func TestVariableGroupAllowAccess_ExpandFlatten_Roundtrip(t *testing.T) {
	testVariableGroup := taskagent.VariableGroup{
		Id:   converter.Int(100),
		Name: converter.String("Name"),
	}
	resourceRefType := "variablegroup"
	testDefinitionResource := build.DefinitionResourceReference{
		Type:       &resourceRefType,
		Authorized: converter.Bool(true),
		Name:       testVariableGroup.Name,
		Id:         converter.String("100"),
	}
	resourceData := schema.TestResourceDataRaw(t, resourceVariableGroup().Schema, nil)

	testArrayDefinitionResourceReference := []build.DefinitionResourceReference{testDefinitionResource}
	flattenAllowAccess(resourceData, &testArrayDefinitionResourceReference)

	definitionResourceReferenceArgs := expandAllowAccess(resourceData, &testVariableGroup)
	require.Equal(t, testDefinitionResource.Authorized, definitionResourceReferenceArgs[0].Authorized)
	require.Equal(t, testDefinitionResource.Id, definitionResourceReferenceArgs[0].Id)
}

func TestVariableGroup_ExpandFlatten_Roundtrip(t *testing.T) {
	testVariableGroup := taskagent.VariableGroup{
		Id:          converter.Int(100),
		Name:        converter.String("Name"),
		Description: converter.String("This is a test variable group."),
		Variables: &map[string]interface{}{
			"var1": map[string]interface{}{
				"value":    converter.String("value1"),
				"isSecret": converter.Bool(false),
			},
		},
	}
	resourceData := schema.TestResourceDataRaw(t, resourceVariableGroup().Schema, nil)
	testVarGroupProjectID := uuid.New().String()

	err := flattenVariableGroup(resourceData, &testVariableGroup, &testVarGroupProjectID)
	require.Equal(t, nil, err)

	variableGroupParams, projectID, _ := expandVariableGroupParameters(resourceData)
	require.Equal(t, *testVariableGroup.Name, *variableGroupParams.Name)
	require.Equal(t, *testVariableGroup.Description, *variableGroupParams.Description)
	require.Equal(t, testVarGroupProjectID, *projectID)

	variablesExpected, _ := json.Marshal(testVariableGroup.Variables)
	variableActual, _ := json.Marshal(variableGroupParams.Variables)
	require.Equal(t, variablesExpected, variableActual)
}

func TestVariableGroupKeyVault_ExpandFlatten_Roundtrip(t *testing.T) {
	testVariableGroupKeyvault := taskagent.VariableGroup{
		Id:          converter.Int(100),
		Name:        converter.String("Name"),
		Description: converter.String("This is a test variable group."),
		Variables: &map[string]interface{}{
			"var1": map[string]interface{}{
				"isSecret": converter.Bool(false),
			},
		},
		ProviderData: map[string]interface{}{
			"serviceEndpointId": converter.String(uuid.New().String()),
			"vault":             converter.String("VaultName"),
		},
		Type: converter.String(azureKeyVaultType),
	}
	resourceData := schema.TestResourceDataRaw(t, resourceVariableGroup().Schema, nil)
	testVarGroupProjectID := uuid.New().String()

	err := flattenVariableGroup(resourceData, &testVariableGroupKeyvault, &testVarGroupProjectID)
	require.Equal(t, nil, err)

	variableGroupParams, projectID, _ := expandVariableGroupParameters(resourceData)
	require.Equal(t, *testVariableGroupKeyvault.Name, *variableGroupParams.Name)
	require.Equal(t, *testVariableGroupKeyvault.Description, *variableGroupParams.Description)
	require.Equal(t, testVarGroupProjectID, *projectID)

	variablesExpected, _ := json.Marshal(testVariableGroupKeyvault.Variables)
	variableActual, _ := json.Marshal(variableGroupParams.Variables)
	require.Equal(t, variablesExpected, variableActual)

	providerDataExpected, _ := json.Marshal(testVariableGroupKeyvault.ProviderData)
	providerDataActual, _ := json.Marshal(variableGroupParams.ProviderData)
	require.Equal(t, providerDataExpected, providerDataActual)
}
