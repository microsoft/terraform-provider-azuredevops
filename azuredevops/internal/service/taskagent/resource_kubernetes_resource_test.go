//go:build (all || resource_environment_resource_kubernetes) && !exclude_resource_environment_resource_kubernetes
// +build all resource_environment_resource_kubernetes
// +build !exclude_resource_environment_resource_kubernetes

package taskagent

import (
	"context"
	"errors"
	"math/rand"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdk/taskagentkubernetesresource"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

var testKubernetesResourceProjectId = uuid.New()
var testKubernetesResourceServiceEndpointId = uuid.New()
var testKubernetesResourceEnvironmentId = rand.Intn(100)
var testKubernetesResourceId = rand.Intn(100)
var testKubernetesResourceTags = []string{"test1", "test2"}

var testKubernetesResource = taskagent.KubernetesResource{
	EnvironmentReference: &taskagent.EnvironmentReference{
		Id: &testKubernetesResourceEnvironmentId,
	},
	Id:                &testKubernetesResourceId,
	Name:              converter.String("Test Kubernetes Resource"),
	Tags:              &testKubernetesResourceTags,
	ClusterName:       converter.String("Test Cluster"),
	Namespace:         converter.String("Test Namespace"),
	ServiceEndpointId: &testKubernetesResourceServiceEndpointId,
}

var testKubernetesResourceProject = taskagent.ProjectReference{Id: &testKubernetesResourceProjectId}

// verifies that the flatten/expand round trip yields the same definition
func TestKubernetesResource_ExpandFlatten_RoundTrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceKubernetesResource().Schema, nil)

	flattenKubernetesResource(resourceData, &testKubernetesResourceProject, &testKubernetesResource)
	projectAfterRoundTrip, resourceAfterRoundTrip, err := expandKubernetesResource(resourceData)
	require.Nil(t, err)
	require.Equal(t, testKubernetesResourceProject, *projectAfterRoundTrip)
	require.Equal(t, testKubernetesResource, *resourceAfterRoundTrip)
}

func TestKubernetesResource_CreateKubernetesResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskAgentClient := azdosdkmocks.NewMockTaskagentkubernetesresourceClient(ctrl)
	clients := &client.AggregatedClient{
		TaskAgentKubernetesResourceClient: taskAgentClient,
		Ctx:                               context.Background(),
	}

	expectedArgs := taskagentkubernetesresource.AddKubernetesResourceArgs{
		CreateParameters: &taskagent.KubernetesResourceCreateParametersExistingEndpoint{
			ClusterName:       testKubernetesResource.ClusterName,
			Name:              testKubernetesResource.Name,
			Namespace:         testKubernetesResource.Namespace,
			Tags:              testKubernetesResource.Tags,
			ServiceEndpointId: testKubernetesResource.ServiceEndpointId,
		},
		Project:       converter.String(testKubernetesResourceProject.Id.String()),
		EnvironmentId: testKubernetesResource.EnvironmentReference.Id,
	}

	taskAgentClient.
		EXPECT().
		AddKubernetesResource(clients.Ctx, expectedArgs).
		Return(&testKubernetesResource, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceKubernetesResource().Schema, nil)
	flattenKubernetesResource(resourceData, &testKubernetesResourceProject, &testKubernetesResource)
	err := resourceKubernetesCreate(resourceData, clients)
	require.NoError(t, err)

	project, resource, err := expandKubernetesResource(resourceData)
	require.NoError(t, err)
	assert.Equal(t, testKubernetesResourceProject, *project)
	assert.Equal(t, testKubernetesResource, *resource)
}

func TestKubernetesResource_CreateKubernetesResourceReturnsErrorOnFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskAgentClient := azdosdkmocks.NewMockTaskagentkubernetesresourceClient(ctrl)
	clients := &client.AggregatedClient{
		TaskAgentKubernetesResourceClient: taskAgentClient,
		Ctx:                               context.Background(),
	}

	expectedArgs := taskagentkubernetesresource.AddKubernetesResourceArgs{
		CreateParameters: &taskagent.KubernetesResourceCreateParametersExistingEndpoint{
			ClusterName:       testKubernetesResource.ClusterName,
			Name:              testKubernetesResource.Name,
			Namespace:         testKubernetesResource.Namespace,
			Tags:              testKubernetesResource.Tags,
			ServiceEndpointId: testKubernetesResource.ServiceEndpointId,
		},
		Project:       converter.String(testKubernetesResourceProject.Id.String()),
		EnvironmentId: testKubernetesResource.EnvironmentReference.Id,
	}

	expectedError := errors.New("test error")

	taskAgentClient.
		EXPECT().
		AddKubernetesResource(clients.Ctx, expectedArgs).
		Return(nil, expectedError).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceKubernetesResource().Schema, nil)
	flattenKubernetesResource(resourceData, &testKubernetesResourceProject, &testKubernetesResource)
	err := resourceKubernetesCreate(resourceData, clients)
	assert.Contains(t, err.Error(), expectedError.Error())
}

func TestKubernetesResource_ReadKubernetesResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
	clients := &client.AggregatedClient{
		TaskAgentClient: taskAgentClient,
		Ctx:             context.Background(),
	}

	expectedArgs := taskagent.GetKubernetesResourceArgs{
		Project:       converter.String(testKubernetesResourceProject.Id.String()),
		EnvironmentId: testKubernetesResource.EnvironmentReference.Id,
		ResourceId:    testKubernetesResource.Id,
	}

	taskAgentClient.
		EXPECT().
		GetKubernetesResource(clients.Ctx, expectedArgs).
		Return(&testKubernetesResource, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceKubernetesResource().Schema, nil)
	flattenKubernetesResource(resourceData, &testKubernetesResourceProject, &testKubernetesResource)
	err := resourceKubernetesRead(resourceData, clients)
	require.NoError(t, err)

	project, resource, err := expandKubernetesResource(resourceData)
	require.NoError(t, err)
	assert.Equal(t, testKubernetesResourceProject, *project)
	assert.Equal(t, testKubernetesResource, *resource)
}

func TestKubernetesResource_ReadKubernetesResourceReturnsErrorOnFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
	clients := &client.AggregatedClient{
		TaskAgentClient: taskAgentClient,
		Ctx:             context.Background(),
	}

	expectedArgs := taskagent.GetKubernetesResourceArgs{
		Project:       converter.String(testKubernetesResourceProject.Id.String()),
		EnvironmentId: testKubernetesResource.EnvironmentReference.Id,
		ResourceId:    testKubernetesResource.Id,
	}

	expectedError := errors.New("test error")

	taskAgentClient.
		EXPECT().
		GetKubernetesResource(clients.Ctx, expectedArgs).
		Return(nil, expectedError).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceKubernetesResource().Schema, nil)
	flattenKubernetesResource(resourceData, &testKubernetesResourceProject, &testKubernetesResource)
	err := resourceKubernetesRead(resourceData, clients)
	assert.Contains(t, err.Error(), expectedError.Error())
}

func TestKubernetesResource_DeleteKubernetesResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
	clients := &client.AggregatedClient{
		TaskAgentClient: taskAgentClient,
		Ctx:             context.Background(),
	}

	expectedArgs := taskagent.DeleteKubernetesResourceArgs{
		Project:       converter.String(testKubernetesResourceProject.Id.String()),
		EnvironmentId: testKubernetesResource.EnvironmentReference.Id,
		ResourceId:    testKubernetesResource.Id,
	}

	taskAgentClient.
		EXPECT().
		DeleteKubernetesResource(clients.Ctx, expectedArgs).
		Return(nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceKubernetesResource().Schema, nil)
	flattenKubernetesResource(resourceData, &testKubernetesResourceProject, &testKubernetesResource)
	err := resourceKubernetesDelete(resourceData, clients)
	require.NoError(t, err)

	_, resource, err := expandKubernetesResource(resourceData)
	require.NoError(t, err)
	assert.Nil(t, resource.Id)
}

func TestKubernetesResource_DeleteKubernetesResourceReturnsErrorOnFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
	clients := &client.AggregatedClient{
		TaskAgentClient: taskAgentClient,
		Ctx:             context.Background(),
	}

	expectedArgs := taskagent.DeleteKubernetesResourceArgs{
		Project:       converter.String(testKubernetesResourceProject.Id.String()),
		EnvironmentId: testKubernetesResource.EnvironmentReference.Id,
		ResourceId:    testKubernetesResource.Id,
	}

	expectedError := errors.New("test error")

	taskAgentClient.
		EXPECT().
		DeleteKubernetesResource(clients.Ctx, expectedArgs).
		Return(expectedError).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceKubernetesResource().Schema, nil)
	flattenKubernetesResource(resourceData, &testKubernetesResourceProject, &testKubernetesResource)
	err := resourceKubernetesDelete(resourceData, clients)
	assert.Contains(t, err.Error(), expectedError.Error())
}
