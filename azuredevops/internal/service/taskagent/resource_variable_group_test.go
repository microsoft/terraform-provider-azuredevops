//go:build (all || resource_variable_group) && !exclude_resource_variable_group
// +build all resource_variable_group
// +build !exclude_resource_variable_group

package taskagent

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

//var serviceEndpointResult = &serviceendpoint.ServiceEndpointRequestResult{
//	ErrorMessage: converter.String(""),
//	Result: []interface{}{
//		"{\"value\": [{\"contentType\": \"type\",\"id\": \"https://mock.vault.azure.net/secrets/var1\",\"attributes\": {\"enabled\": true,\"exp\": 1675424159,\"created\": 1611734011,\"updated\": 1612353657,\"recoveryLevel\": \"Recoverable+Purgeable\"},\"tags\": {}}],\"nextLink\": null}",
//	},
//	StatusCode: converter.String("ok"),
//}
//
//func TestVariableGroupAllowAccess_ExpandFlatten_Roundtrip(t *testing.T) {
//	testVariableGroup := taskagent.VariableGroup{
//		Id:          converter.Int(100),
//		Name:        converter.String("Name"),
//		Description: converter.String("This is a test variable group."),
//		Variables: &map[string]interface{}{
//			"var1": map[string]interface{}{
//				"value":    converter.String("value1"),
//				"isSecret": converter.Bool(false),
//			},
//		},
//	}
//	resourceRefType := "variablegroup"
//	testDefinitionResource := build.DefinitionResourceReference{
//		Type:       &resourceRefType,
//		Authorized: converter.Bool(true),
//		Name:       testVariableGroup.Name,
//		Id:         converter.String("100"),
//	}
//	resourceData := schema.TestResourceDataRaw(t, ResourceVariableGroup().Schema, nil)
//	testVarGroupProjectID := uuid.New().String()
//
//	err := flattenVariableGroup(resourceData, &testVariableGroup, &testVarGroupProjectID)
//	require.Equal(t, nil, err)
//
//	testArrayDefinitionResourceReference := []build.DefinitionResourceReference{testDefinitionResource}
//	flattenAllowAccess(resourceData, &testArrayDefinitionResourceReference)
//
//	definitionResourceReferenceArgs := expandAllowAccess(resourceData, &testVariableGroup)
//
//	var definitionRes build.DefinitionResourceReference
//	for _, authResource := range definitionResourceReferenceArgs {
//		if *testDefinitionResource.Id == *authResource.Id {
//			definitionRes = authResource
//		}
//	}
//	require.Equal(t, testDefinitionResource.Authorized, definitionRes.Authorized)
//	require.Equal(t, testDefinitionResource.Id, definitionRes.Id)
//}
//
//func TestVariableGroup_ExpandFlatten_Roundtrip(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	serviceEndpointClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
//	clients := &client.AggregatedClient{
//		ServiceEndpointClient: serviceEndpointClient,
//		Ctx:                   context.Background(),
//	}
//
//	testVariableGroup := taskagent.VariableGroup{
//		Id:          converter.Int(100),
//		Name:        converter.String("Name"),
//		Description: converter.String("This is a test variable group."),
//		Variables: &map[string]interface{}{
//			"var1": map[string]interface{}{
//				"value":    converter.String("value1"),
//				"isSecret": converter.Bool(false),
//			},
//		},
//	}
//	resourceData := schema.TestResourceDataRaw(t, ResourceVariableGroup().Schema, nil)
//	testVarGroupProjectID := uuid.New().String()
//
//	err := flattenVariableGroup(resourceData, &testVariableGroup, &testVarGroupProjectID)
//	require.Equal(t, nil, err)
//
//	variableGroupParams, projectID, _ := expandVariableGroupParameters(clients, resourceData)
//	require.Equal(t, *testVariableGroup.Name, *variableGroupParams.Name)
//	require.Equal(t, *testVariableGroup.Description, *variableGroupParams.Description)
//	require.Equal(t, testVarGroupProjectID, *projectID)
//
//	variablesExpected, _ := json.Marshal(testVariableGroup.Variables)
//	variableActual, _ := json.Marshal(variableGroupParams.Variables)
//	require.Equal(t, variablesExpected, variableActual)
//}
//
//func TestVariableGroupKeyVault_ExpandFlatten_Roundtrip(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	serviceEndpointClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
//	clients := &client.AggregatedClient{
//		ServiceEndpointClient: serviceEndpointClient,
//		Ctx:                   context.Background(),
//	}
//
//	serviceEndpointClient.
//		EXPECT().
//		ExecuteServiceEndpointRequest(clients.Ctx, gomock.Any()).
//		Return(serviceEndpointResult, nil).Times(1)
//
//	testVariableGroupKeyvault := taskagent.VariableGroup{
//		Id:          converter.Int(100),
//		Name:        converter.String("Name"),
//		Description: converter.String("This is a test variable group."),
//		Variables: &map[string]interface{}{
//			"var1": taskagent.AzureKeyVaultVariableValue{
//				IsSecret:    converter.Bool(true),
//				Value:       nil,
//				ContentType: converter.String("type"),
//				Enabled:     converter.Bool(true),
//				Expires: &azuredevops.Time{
//					Time: time.Unix(1675424159, 0),
//				},
//			},
//		},
//		ProviderData: map[string]interface{}{
//			"serviceEndpointId": converter.String(uuid.New().String()),
//			"vault":             converter.String("VaultName"),
//		},
//		Type: converter.String(azureKeyVaultType),
//	}
//	resourceData := schema.TestResourceDataRaw(t, ResourceVariableGroup().Schema, nil)
//	testVarGroupProjectID := uuid.New().String()
//
//	err := flattenVariableGroup(resourceData, &testVariableGroupKeyvault, &testVarGroupProjectID)
//	require.Equal(t, nil, err)
//
//	variableGroupParams, projectID, _ := expandVariableGroupParameters(clients, resourceData)
//	require.Equal(t, *testVariableGroupKeyvault.Name, *variableGroupParams.Name)
//	require.Equal(t, *testVariableGroupKeyvault.Description, *variableGroupParams.Description)
//	require.Equal(t, testVarGroupProjectID, *projectID)
//
//	variablesExpected, _ := json.Marshal(testVariableGroupKeyvault.Variables)
//	variableActual, _ := json.Marshal(variableGroupParams.Variables)
//	require.Equal(t, variablesExpected, variableActual)
//
//	providerDataExpected, _ := json.Marshal(testVariableGroupKeyvault.ProviderData)
//	providerDataActual, _ := json.Marshal(variableGroupParams.ProviderData)
//	require.Equal(t, providerDataExpected, providerDataActual)
//}
