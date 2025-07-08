package taskagent

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.
// var azProjectRef = &core.TeamProjectReference{
//	Id:   testhelper.CreateUUID(),
//	Name: converter.String("project-01"),
//}

//
// func TestDataSourceVariableGroup_Read_VariableGroup(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
//	clients := &client.AggregatedClient{
//		Ctx:             context.Background(),
//		TaskAgentClient: taskAgentClient,
//	}
//
//	name := "vgName"
//	variableGroup := taskagent.VariableGroup{
//		Id:   converter.Int(100),
//		Name: converter.String(name),
//		Type: converter.String("Vsts"),
//		Variables: &map[string]interface{}{
//			"var1": map[string]interface{}{
//				"value":    converter.String("value1"),
//				"isSecret": converter.Bool(false),
//			},
//		},
//	}
//
//	taskAgentClient.
//		EXPECT().
//		GetVariableGroups(clients.Ctx,
//			taskagent.GetVariableGroupsArgs{
//				Project:   converter.String(azProjectRef.Id.String()),
//				GroupName: variableGroup.Name,
//				Top:       converter.Int(1),
//			}).
//		Return(&[]taskagent.VariableGroup{
//			variableGroup,
//		}, nil)
//
//	resourceData := schema.TestResourceDataRaw(t, DataVariableGroup().Schema, nil)
//	resourceData.Set(vgName, variableGroup.Name)
//	resourceData.Set(vgProjectID, azProjectRef.Id.String())
//
//	err := dataSourceVariableGroupRead(resourceData, clients)
//	require.Nil(t, err)
//	require.Equal(t, resourceData.Id(), fmt.Sprintf("%d", *variableGroup.Id))
//	require.Equal(t, resourceData.Get(vgName), name)
//	require.Equal(t, resourceData.Get(vgProjectID), azProjectRef.Id.String())
//}
//
// func TestDataSourceVariableGroup_Read_FindAllVariables(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
//	clients := &client.AggregatedClient{
//		Ctx:             context.Background(),
//		TaskAgentClient: taskAgentClient,
//	}
//
//	name := "vgName"
//	variableGroup := taskagent.VariableGroup{
//		Id:   converter.Int(100),
//		Name: converter.String(name),
//		Type: converter.String("Vsts"),
//		Variables: &map[string]interface{}{
//			"var1": map[string]interface{}{
//				"value":    converter.String("value1"),
//				"isSecret": converter.Bool(false),
//			},
//			"var2": map[string]interface{}{
//				"value":    converter.String("value2"),
//				"isSecret": converter.Bool(false),
//			},
//		},
//	}
//
//	taskAgentClient.
//		EXPECT().
//		GetVariableGroups(clients.Ctx,
//			taskagent.GetVariableGroupsArgs{
//				Project:   converter.String(azProjectRef.Id.String()),
//				GroupName: variableGroup.Name,
//				Top:       converter.Int(1),
//			}).
//		Return(&[]taskagent.VariableGroup{
//			variableGroup,
//		}, nil)
//
//	resourceData := schema.TestResourceDataRaw(t, DataVariableGroup().Schema, nil)
//	resourceData.Set(vgName, variableGroup.Name)
//	resourceData.Set(vgProjectID, azProjectRef.Id.String())
//
//	err := dataSourceVariableGroupRead(resourceData, clients)
//	require.Nil(t, err)
//	variables, ok := resourceData.GetOk(vgVariable)
//	require.True(t, ok)
//	require.NotNil(t, variables)
//	variablesSet, ok := variables.(*schema.Set)
//	require.True(t, ok)
//	require.NotNil(t, variablesSet)
//	require.Equal(t, 2, variablesSet.Len())
//}
//
// func TestDataSourceVariableGroup_Read_VariableGroupNotFound(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
//	clients := &client.AggregatedClient{
//		Ctx:             context.Background(),
//		TaskAgentClient: taskAgentClient,
//	}
//
//	name := "nonexistent"
//	taskAgentClient.
//		EXPECT().
//		GetVariableGroups(clients.Ctx, taskagent.GetVariableGroupsArgs{
//			Project:   converter.String(azProjectRef.Id.String()),
//			GroupName: &name,
//			Top:       converter.Int(1),
//		}).
//		Return(&[]taskagent.VariableGroup{}, nil)
//
//	resourceData := schema.TestResourceDataRaw(t, DataVariableGroup().Schema, nil)
//	resourceData.Set(vgName, &name)
//	resourceData.Set(vgProjectID, azProjectRef.Id.String())
//
//	err := dataSourceVariableGroupRead(resourceData, clients)
//	require.Contains(t, err.Error(), "Unable to find variable group")
//}
