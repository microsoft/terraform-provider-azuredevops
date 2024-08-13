//go:build (all || resource_environment_resource_kubernetes) && !exclude_resource_environment_resource_kubernetes
// +build all resource_environment_resource_kubernetes
// +build !exclude_resource_environment_resource_kubernetes

package taskagent

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/stretchr/testify/assert"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

var testEnvironmentKubernetesResourceProjectId = uuid.New()
var testEnvironmentKubernetesResourceServiceEndpointId = uuid.New()
var testEnvironmentKubernetesResourceEnvironmentId = rand.Intn(100)
var testEnvironmentKubernetesResourceId = rand.Intn(100)
var testEnvironmentKubernetesResourceTags = []string{"test1", "test2"}

var testEnvironmentKubernetesResource = taskagent.KubernetesResource{
	EnvironmentReference: &taskagent.EnvironmentReference{
		Id: &testEnvironmentKubernetesResourceEnvironmentId,
	},
	Id:                &testEnvironmentKubernetesResourceId,
	Name:              converter.String("Test Kubernetes Resource"),
	Tags:              &testEnvironmentKubernetesResourceTags,
	ClusterName:       converter.String("Test Cluster"),
	Namespace:         converter.String("Test Namespace"),
	ServiceEndpointId: &testEnvironmentKubernetesResourceServiceEndpointId,
}

var testEnvironmentKubernetesResourceProject = taskagent.ProjectReference{Id: &testEnvironmentKubernetesResourceProjectId}

// verifies that the flatten/expand round trip yields the same definition
func TestEnvironmentKubernetesResource_ExpandFlatten_RoundTrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceEnvironmentKubernetes().Schema, nil)

	flattenEnvironmentKubernetesResource(resourceData, &testEnvironmentKubernetesResourceProject, &testEnvironmentKubernetesResource)
	resourceData.SetId(strconv.Itoa(*testEnvironmentKubernetesResource.Id))
	projectAfterRoundTrip, resourceAfterRoundTrip, err := expandEnvironmentKubernetesResource(resourceData)
	resourceAfterRoundTrip.Id = testEnvironmentKubernetesResource.Id
	require.Nil(t, err)
	require.Equal(t, testEnvironmentKubernetesResourceProject, *projectAfterRoundTrip)
	require.Equal(t, testEnvironmentKubernetesResource, *resourceAfterRoundTrip)
}

func TestEnvironmentKubernetesResource_CreateKubernetesResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
	clients := &client.AggregatedClient{
		TaskAgentClient: taskAgentClient,
		Ctx:             context.Background(),
	}

	expectedArgs := taskagent.AddKubernetesResourceArgsExistingEndpoint{
		CreateParameters: &taskagent.KubernetesResourceCreateParametersExistingEndpoint{
			ClusterName:       testEnvironmentKubernetesResource.ClusterName,
			Name:              testEnvironmentKubernetesResource.Name,
			Namespace:         testEnvironmentKubernetesResource.Namespace,
			Tags:              testEnvironmentKubernetesResource.Tags,
			ServiceEndpointId: testEnvironmentKubernetesResource.ServiceEndpointId,
		},
		Project:       converter.String(testEnvironmentKubernetesResourceProject.Id.String()),
		EnvironmentId: testEnvironmentKubernetesResource.EnvironmentReference.Id,
	}

	taskAgentClient.
		EXPECT().
		AddKubernetesResourcExistingEndpoint(clients.Ctx, expectedArgs).
		Return(&testEnvironmentKubernetesResource, nil).
		Times(1)
	taskAgentClient.
		EXPECT().
		GetKubernetesResource(clients.Ctx, taskagent.GetKubernetesResourceArgs{
			Project:       converter.String(testEnvironmentKubernetesResourceProject.Id.String()),
			EnvironmentId: testEnvironmentKubernetesResource.EnvironmentReference.Id,
			ResourceId:    testEnvironmentKubernetesResource.Id,
		}).
		Return(&testEnvironmentKubernetesResource, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceEnvironmentKubernetes().Schema, nil)
	resourceData.SetId(strconv.Itoa(*testEnvironmentKubernetesResource.Id))
	flattenEnvironmentKubernetesResource(resourceData, &testEnvironmentKubernetesResourceProject, &testEnvironmentKubernetesResource)
	err := resourceEnvironmentKubernetesCreate(resourceData, clients)
	require.NoError(t, err)

	project, resource, err := expandEnvironmentKubernetesResource(resourceData)
	resource.Id = testEnvironmentKubernetesResource.Id
	require.NoError(t, err)
	assert.Equal(t, testEnvironmentKubernetesResourceProject, *project)
	assert.Equal(t, testEnvironmentKubernetesResource, *resource)
}

func TestEnvironmentKubernetesResource_CreateKubernetesResourceReturnsErrorOnFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
	clients := &client.AggregatedClient{
		TaskAgentClient: taskAgentClient,
		Ctx:             context.Background(),
	}

	expectedArgs := taskagent.AddKubernetesResourceArgsExistingEndpoint{
		CreateParameters: &taskagent.KubernetesResourceCreateParametersExistingEndpoint{
			ClusterName:       testEnvironmentKubernetesResource.ClusterName,
			Name:              testEnvironmentKubernetesResource.Name,
			Namespace:         testEnvironmentKubernetesResource.Namespace,
			Tags:              testEnvironmentKubernetesResource.Tags,
			ServiceEndpointId: testEnvironmentKubernetesResource.ServiceEndpointId,
		},
		Project:       converter.String(testEnvironmentKubernetesResourceProject.Id.String()),
		EnvironmentId: testEnvironmentKubernetesResource.EnvironmentReference.Id,
	}

	expectedError := errors.New("test error")

	taskAgentClient.
		EXPECT().
		AddKubernetesResourcExistingEndpoint(clients.Ctx, expectedArgs).
		Return(nil, expectedError).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceEnvironmentKubernetes().Schema, nil)
	resourceData.SetId(strconv.Itoa(*testEnvironmentKubernetesResource.Id))
	flattenEnvironmentKubernetesResource(resourceData, &testEnvironmentKubernetesResourceProject, &testEnvironmentKubernetesResource)
	err := resourceEnvironmentKubernetesCreate(resourceData, clients)
	assert.Contains(t, err.Error(), expectedError.Error())
}

func TestEnvironmentKubernetesResource_ReadKubernetesResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
	clients := &client.AggregatedClient{
		TaskAgentClient: taskAgentClient,
		Ctx:             context.Background(),
	}

	expectedArgs := taskagent.GetKubernetesResourceArgs{
		Project:       converter.String(testEnvironmentKubernetesResourceProject.Id.String()),
		EnvironmentId: testEnvironmentKubernetesResource.EnvironmentReference.Id,
		ResourceId:    testEnvironmentKubernetesResource.Id,
	}

	taskAgentClient.
		EXPECT().
		GetKubernetesResource(clients.Ctx, expectedArgs).
		Return(&testEnvironmentKubernetesResource, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceEnvironmentKubernetes().Schema, nil)
	resourceData.SetId(strconv.Itoa(*testEnvironmentKubernetesResource.Id))
	flattenEnvironmentKubernetesResource(resourceData, &testEnvironmentKubernetesResourceProject, &testEnvironmentKubernetesResource)
	err := resourceEnvironmentKubernetesRead(resourceData, clients)
	require.NoError(t, err)

	project, resource, err := expandEnvironmentKubernetesResource(resourceData)
	resource.Id = testEnvironmentKubernetesResource.Id
	require.NoError(t, err)
	assert.Equal(t, testEnvironmentKubernetesResourceProject, *project)
	assert.Equal(t, testEnvironmentKubernetesResource, *resource)
}

func TestEnvironmentKubernetesResource_ReadKubernetesResourceReturnsErrorOnFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
	clients := &client.AggregatedClient{
		TaskAgentClient: taskAgentClient,
		Ctx:             context.Background(),
	}

	expectedArgs := taskagent.GetKubernetesResourceArgs{
		Project:       converter.String(testEnvironmentKubernetesResourceProject.Id.String()),
		EnvironmentId: testEnvironmentKubernetesResource.EnvironmentReference.Id,
		ResourceId:    testEnvironmentKubernetesResource.Id,
	}

	expectedError := errors.New("test error")

	taskAgentClient.
		EXPECT().
		GetKubernetesResource(clients.Ctx, expectedArgs).
		Return(nil, expectedError).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceEnvironmentKubernetes().Schema, nil)
	resourceData.SetId(strconv.Itoa(*testEnvironmentKubernetesResource.Id))
	flattenEnvironmentKubernetesResource(resourceData, &testEnvironmentKubernetesResourceProject, &testEnvironmentKubernetesResource)
	err := resourceEnvironmentKubernetesRead(resourceData, clients)
	assert.Contains(t, err.Error(), expectedError.Error())
}

func TestEnvironmentKubernetesResource_DeleteKubernetesResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
	clients := &client.AggregatedClient{
		TaskAgentClient: taskAgentClient,
		Ctx:             context.Background(),
	}

	expectedArgs := taskagent.DeleteKubernetesResourceArgs{
		Project:       converter.String(testEnvironmentKubernetesResourceProject.Id.String()),
		EnvironmentId: testEnvironmentKubernetesResource.EnvironmentReference.Id,
		ResourceId:    testEnvironmentKubernetesResource.Id,
	}

	taskAgentClient.
		EXPECT().
		DeleteKubernetesResource(clients.Ctx, expectedArgs).
		Return(nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceEnvironmentKubernetes().Schema, nil)
	resourceData.SetId(strconv.Itoa(*testEnvironmentKubernetesResource.Id))
	flattenEnvironmentKubernetesResource(resourceData, &testEnvironmentKubernetesResourceProject, &testEnvironmentKubernetesResource)
	err := resourceEnvironmentKubernetesDelete(resourceData, clients)
	require.NoError(t, err)

	_, resource, err := expandEnvironmentKubernetesResource(resourceData)
	require.NoError(t, err)
	assert.Nil(t, resource.Id)
}

func TestEnvironmentKubernetesResource_DeleteKubernetesResourceReturnsErrorOnFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
	clients := &client.AggregatedClient{
		TaskAgentClient: taskAgentClient,
		Ctx:             context.Background(),
	}

	expectedArgs := taskagent.DeleteKubernetesResourceArgs{
		Project:       converter.String(testEnvironmentKubernetesResourceProject.Id.String()),
		EnvironmentId: testEnvironmentKubernetesResource.EnvironmentReference.Id,
		ResourceId:    testEnvironmentKubernetesResource.Id,
	}

	expectedError := errors.New("test error")

	taskAgentClient.
		EXPECT().
		DeleteKubernetesResource(clients.Ctx, expectedArgs).
		Return(expectedError).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceEnvironmentKubernetes().Schema, nil)
	resourceData.SetId(strconv.Itoa(*testEnvironmentKubernetesResource.Id))
	flattenEnvironmentKubernetesResource(resourceData, &testEnvironmentKubernetesResourceProject, &testEnvironmentKubernetesResource)
	err := resourceEnvironmentKubernetesDelete(resourceData, clients)
	assert.Contains(t, err.Error(), expectedError.Error())
}
