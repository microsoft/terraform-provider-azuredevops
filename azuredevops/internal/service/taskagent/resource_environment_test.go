//go:build (all || resource_environment) && !exclude_resource_environment
// +build all resource_environment
// +build !exclude_resource_environment

package taskagent

import (
	"context"
	"errors"
	"math/rand"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

var testEnvironmentProjectId = uuid.New()
var testEnvironmentId = rand.Intn(100)

var testEnvironment = taskagent.EnvironmentInstance{
	Id:          &testEnvironmentId,
	Name:        converter.String("EnvironmentName"),
	Description: converter.String(""),
	Project: &taskagent.ProjectReference{
		Id: &testEnvironmentProjectId,
	},
}

// verifies that the flatten/expand round trip yields the same definition
func TestEnvironment_ExpandFlatten_Roundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceEnvironment().Schema, nil)
	flattenEnvironment(resourceData, &testEnvironment)
	environmentAfterRoundTrip, err := expandEnvironment(resourceData, true)
	require.Nil(t, err)
	require.Equal(t, testEnvironment, *environmentAfterRoundTrip)
}

// verifies that the create operation is considered failed if the API call fails.
func TestEnvironment_CreateEnvironment_DoesNotSwallowErrorFromFailedAddAgentCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
	clients := &client.AggregatedClient{
		TaskAgentClient: taskAgentClient,
		Ctx:             context.Background(),
	}

	projectId := testEnvironment.Project.Id.String()
	expectedArgs := taskagent.AddEnvironmentArgs{
		Project: &projectId,
		EnvironmentCreateParameter: &taskagent.EnvironmentCreateParameter{
			Name:        testEnvironment.Name,
			Description: testEnvironment.Description,
		},
	}

	taskAgentClient.
		EXPECT().
		AddEnvironment(clients.Ctx, expectedArgs).
		Return(nil, errors.New("AddEnvironment() Failed")).
		Times(1)

	newEnvironment, err := createEnvironment(clients, &testEnvironment)
	require.Nil(t, newEnvironment)
	require.Equal(t, "AddEnvironment() Failed", err.Error())
}

func TestEnvironment_DeleteEnvironment_ReturnsErrorIfIdReadFails(t *testing.T) {
	client := &client.AggregatedClient{}

	resourceData := schema.TestResourceDataRaw(t, ResourceEnvironment().Schema, nil)
	flattenEnvironment(resourceData, &testEnvironment)
	resourceData.SetId("")

	err := resourceEnvironmentDelete(resourceData, client)
	require.Equal(t, "Error getting environment id: strconv.Atoi: parsing \"\": invalid syntax", err.Error())
}

func TestEnvironment_UpdateEnvironment_ReturnsErrorIfIdReadFails(t *testing.T) {
	client := &client.AggregatedClient{}

	resourceData := schema.TestResourceDataRaw(t, ResourceEnvironment().Schema, nil)
	flattenEnvironment(resourceData, &testEnvironment)
	resourceData.SetId("")

	err := resourceEnvironmentUpdate(resourceData, client)
	require.Equal(t, "Error converting terraform data model to AzDO environment reference: Error getting environment id: strconv.Atoi: parsing \"\": invalid syntax", err.Error())
}

func TestEnvironment_UpdateEnvironment_UpdateAndRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskAgentClient := azdosdkmocks.NewMockTaskagentClient(ctrl)
	clients := &client.AggregatedClient{
		TaskAgentClient: taskAgentClient,
		Ctx:             context.Background(),
	}

	environmentToUpdate := taskagent.EnvironmentInstance{
		Id:          &testEnvironmentId,
		Name:        converter.String("Updated Name"),
		Description: converter.String("Some Description"),
		Project: &taskagent.ProjectReference{
			Id: &testEnvironmentProjectId,
		},
	}

	resourceData := schema.TestResourceDataRaw(t, ResourceEnvironment().Schema, nil)
	flattenEnvironment(resourceData, &environmentToUpdate)

	projectIdString := testEnvironmentProjectId.String()
	taskAgentClient.
		EXPECT().
		UpdateEnvironment(clients.Ctx, taskagent.UpdateEnvironmentArgs{
			Project:       &projectIdString,
			EnvironmentId: &testEnvironmentId,
			EnvironmentUpdateParameter: &taskagent.EnvironmentUpdateParameter{
				Name:        environmentToUpdate.Name,
				Description: environmentToUpdate.Description,
			},
		}).
		Return(&environmentToUpdate, nil).
		Times(1)

	taskAgentClient.
		EXPECT().
		GetEnvironmentById(clients.Ctx, taskagent.GetEnvironmentByIdArgs{
			Project:       &projectIdString,
			EnvironmentId: &testEnvironmentId,
		}).
		Return(&environmentToUpdate, nil).
		Times(1)

	err := resourceEnvironmentUpdate(resourceData, clients)
	require.Nil(t, err)

	updatedEnvironment, _ := expandEnvironment(resourceData, false)
	require.Equal(t, environmentToUpdate.Id, updatedEnvironment.Id)
	require.Equal(t, environmentToUpdate.Name, updatedEnvironment.Name)
	require.Equal(t, environmentToUpdate.Description, updatedEnvironment.Description)
	require.Equal(t, environmentToUpdate.Project.Id, updatedEnvironment.Project.Id)
}
