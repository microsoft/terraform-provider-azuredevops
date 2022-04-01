//go:build (all || core || resource_project || resource_project_features) && !exclude_resource_project_features
// +build all core resource_project resource_project_features
// +build !exclude_resource_project_features

package core

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/featuremanagement"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/testhelper"
	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
)

/**
 * Begin unit tests
 */

const invalidFeatureID = "invalid.projectFeature"

func TestProjectFeatures_ReadProjectFeatureStates_DontSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	featureClient := azdosdkmocks.NewMockFeaturemanagementClient(ctrl)
	clients := &client.AggregatedClient{
		FeatureManagementClient: featureClient,
		Ctx:                     context.Background(),
	}

	const projectID = "2efa4ede-32a3-4373-9703-1c3ed30faeab"
	const errMsg = "QueryFeatureStates() Failed"
	featureClient.
		EXPECT().
		QueryFeatureStates(clients.Ctx, gomock.Any()).
		Return(nil, errors.New(errMsg)).
		Times(1)

	ret, err := readProjectFeatureStates(clients.Ctx, clients.FeatureManagementClient, projectID)
	require.Nil(t, ret)
	require.NotNil(t, err)
	require.Equal(t, errMsg, err.Error())
}

func TestProjectFeatures_ReadProjectFeatureStates_OnlyReturnProjectFeatures(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	featureClient := azdosdkmocks.NewMockFeaturemanagementClient(ctrl)
	clients := &client.AggregatedClient{
		FeatureManagementClient: featureClient,
		Ctx:                     context.Background(),
	}

	const projectID = "2123a037-b1e8-4917-864f-0d025d166a4e"
	const errMsg = "QueryFeatureStates() Failed"
	featureStates := map[string]featuremanagement.ContributedFeatureState{}
	for _, featureName := range *getProjectFeatureIDs() {
		featureStates[featureName] = featuremanagement.ContributedFeatureState{
			FeatureId: converter.String(featureName),
			State:     &featuremanagement.ContributedFeatureEnabledValueValues.Enabled,
		}
	}
	featureStates[invalidFeatureID] = featuremanagement.ContributedFeatureState{
		FeatureId: converter.String(invalidFeatureID),
		State:     &featuremanagement.ContributedFeatureEnabledValueValues.Enabled,
	}

	featureClient.
		EXPECT().
		QueryFeatureStates(clients.Ctx, gomock.Any()).
		Return(&featuremanagement.ContributedFeatureStateQuery{
			FeatureIds:    nil,
			FeatureStates: &featureStates,
			ScopeValues:   nil,
		}, nil).
		Times(1)

	ret, err := readProjectFeatureStates(clients.Ctx, clients.FeatureManagementClient, projectID)
	require.NotNil(t, ret)
	require.Nil(t, err)
	require.True(t, len(featureStates) > len(*ret))
	if _, ok := (*ret)[invalidFeatureID]; ok {
		require.Fail(t, fmt.Sprintf("readProjectFeatureStates must not return state for feature %s", invalidFeatureID))
	}

	featureNameMap := *getProjectFeatureNameMap()
	for k, v := range featureStates {
		if k == invalidFeatureID {
			continue
		}
		state, ok := (*ret)[featureNameMap[k]]
		if !ok {
			require.Fail(t, fmt.Sprintf("readProjectFeatureStates must return state for feature %s", k))
		}
		if state != *v.State {
			require.Fail(t, fmt.Sprintf("readProjectFeatureStates must return state %s for feature %s", featuremanagement.ContributedFeatureEnabledValueValues.Enabled, k))
		}
	}
}

func TestProjectFeatures_GetConfiguredProjectFeatureStates_DontSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	featureClient := azdosdkmocks.NewMockFeaturemanagementClient(ctrl)
	clients := &client.AggregatedClient{
		FeatureManagementClient: featureClient,
		Ctx:                     context.Background(),
	}

	const projectID = "f011627f-b25b-40ae-9f1f-309168f661b2"
	const errMsg = "QueryFeatureStates() Failed"
	featureClient.
		EXPECT().
		QueryFeatureStates(clients.Ctx, gomock.Any()).
		Return(nil, errors.New(errMsg)).
		Times(1)

	ret, err := getConfiguredProjectFeatureStates(clients.Ctx, clients.FeatureManagementClient, &map[string]interface{}{}, projectID)
	require.Nil(t, ret)
	require.NotNil(t, err)
	require.Equal(t, errMsg, err.Error())
}

func TestProjectFeatures_GetConfiguredProjectFeatureStates_ReturnsOnlyRelevantFeatures(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	featureClient := azdosdkmocks.NewMockFeaturemanagementClient(ctrl)
	clients := &client.AggregatedClient{
		FeatureManagementClient: featureClient,
		Ctx:                     context.Background(),
	}

	const projectID = "f011627f-b25b-40ae-9f1f-309168f661b2"
	featureIDs := *getProjectFeatureIDs()
	featureStates := map[string]featuremanagement.ContributedFeatureState{}
	for _, featureName := range featureIDs {
		featureStates[featureName] = featuremanagement.ContributedFeatureState{
			FeatureId: converter.String(featureName),
			State:     &featuremanagement.ContributedFeatureEnabledValueValues.Enabled,
		}
	}
	featureClient.
		EXPECT().
		QueryFeatureStates(clients.Ctx, gomock.Any()).
		Return(&featuremanagement.ContributedFeatureStateQuery{
			FeatureIds:    nil,
			FeatureStates: &featureStates,
			ScopeValues:   nil,
		}, nil).
		Times(1)

	idx := testhelper.RandIntSlice(0, len(featureIDs)-1, len(featureIDs)/2)
	featureNameMap := *getProjectFeatureNameMap()
	relevantFeatures := make(map[string]interface{}, len(idx))
	for i := range idx {
		relevantFeatures[string(featureNameMap[featureIDs[i]])] = string(*featureStates[featureIDs[i]].State)
	}
	ret, err := getConfiguredProjectFeatureStates(clients.Ctx, clients.FeatureManagementClient, &relevantFeatures, projectID)
	require.Nil(t, err)
	require.NotNil(t, ret)
	for k, v := range relevantFeatures {
		retVal, ok := (*ret)[ProjectFeatureType(k)]
		if !ok {
			require.Fail(t, fmt.Sprintf("getConfiguredProjectFeatureStates must return state for feature %s", k))
		}
		if v.(string) != string(retVal) {
			require.Fail(t, fmt.Sprintf("getConfiguredProjectFeatureStates must return state % for feature %s", v, k))
		}
	}
}

func TestProjectFeatures_SetProjectFeatureStates_HandleError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	featureClient := azdosdkmocks.NewMockFeaturemanagementClient(ctrl)
	clients := &client.AggregatedClient{
		FeatureManagementClient: featureClient,
		Ctx:                     context.Background(),
	}

	const projectID = "ab3abf46-a5ed-4f61-9bde-4fd8427ec54e"
	featureIDs := *getProjectFeatureIDs()
	featureMap := *getProjectFeatureNameMap()
	idx := testhelper.RandInt(0, len(featureIDs)-1)

	featureStates := make(map[string]interface{}, 1)
	featureStates[string(featureMap[featureIDs[idx]])] = string(featuremanagement.ContributedFeatureEnabledValueValues.Enabled)

	expectedSetFeatureStateForScopeArgs := featuremanagement.SetFeatureStateForScopeArgs{
		Feature: &featuremanagement.ContributedFeatureState{
			FeatureId: converter.String(featureIDs[idx]),
			State:     &featuremanagement.ContributedFeatureEnabledValueValues.Enabled,
			Scope: &featuremanagement.ContributedFeatureSettingScope{
				SettingScope: converter.String("project"),
				UserScoped:   converter.Bool(false),
			},
		},
		FeatureId:  converter.String(featureIDs[idx]),
		UserScope:  converter.String("host"),
		ScopeName:  converter.String("project"),
		ScopeValue: converter.String(projectID),
	}

	const errMsg = "SetFeatureStateForScope() Failed"
	featureClient.
		EXPECT().
		SetFeatureStateForScope(clients.Ctx, expectedSetFeatureStateForScopeArgs).
		Return(nil, errors.New(errMsg)).
		Times(1)

	err := setProjectFeatureStates(clients.Ctx, clients.FeatureManagementClient, projectID, &featureStates)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), errMsg)
}

func TestProjectFeatures_SetProjectFeatureStates_OnlyValidProjectFeatures(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	featureClient := azdosdkmocks.NewMockFeaturemanagementClient(ctrl)
	clients := &client.AggregatedClient{
		FeatureManagementClient: featureClient,
		Ctx:                     context.Background(),
	}

	const projectID = "ab3abf46-a5ed-4f61-9bde-4fd8427ec54e"
	featureStates := map[string]interface{}{
		invalidFeatureID: string(featuremanagement.ContributedFeatureEnabledValueValues.Enabled),
	}
	errMsg := fmt.Sprintf("unknown feature: %s", invalidFeatureID)

	err := setProjectFeatureStates(clients.Ctx, clients.FeatureManagementClient, projectID, &featureStates)
	require.NotNil(t, err)
	require.Equal(t, errMsg, err.Error())
}
