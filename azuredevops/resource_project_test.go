// +build all core resource_project
// +build !exclude_resource_project

package azuredevops

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/operations"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/stretchr/testify/require"
)

var testID = uuid.New()
var testProject = core.TeamProject{
	Id:          &testID,
	Name:        converter.String("Name"),
	Visibility:  &core.ProjectVisibilityValues.Public,
	Description: converter.String("Description"),
	Capabilities: &map[string]map[string]string{
		"versioncontrol":  {"sourceControlType": "SouceControlType"},
		"processTemplate": {"templateTypeId": testID.String()},
	},
}

// verifies that the create operation is considered failed if the initial API
// call fails.
func TestAzureDevOpsProject_CreateProject_DoesNotSwallowErrorFromFailedCreateCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	clients := &config.AggregatedClient{
		CoreClient: coreClient,
		Ctx:        context.Background(),
	}

	expectedProjectCreateArgs := core.QueueCreateProjectArgs{ProjectToCreate: &testProject}

	coreClient.
		EXPECT().
		QueueCreateProject(clients.Ctx, expectedProjectCreateArgs).
		Return(nil, errors.New("QueueCreateProject() Failed")).
		Times(1)

	err := createProject(clients, &testProject, 5)
	require.Equal(t, "QueueCreateProject() Failed", err.Error())
}

// verifies that the create operation is considered failed if there is an issue
// verifying via the async polling operation API that it has completed successfully.
func TestAzureDevOpsProject_CreateProject_DoesNotSwallowErrorFromFailedAsyncStatusCheckCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	operationsClient := azdosdkmocks.NewMockOperationsClient(ctrl)
	clients := &config.AggregatedClient{
		CoreClient:       coreClient,
		OperationsClient: operationsClient,
		Ctx:              context.Background(),
	}

	expectedProjectCreateArgs := core.QueueCreateProjectArgs{ProjectToCreate: &testProject}
	mockedOperationReference := operations.OperationReference{Id: &testID}
	expectedOperationArgs := operations.GetOperationArgs{OperationId: &testID}

	coreClient.
		EXPECT().
		QueueCreateProject(clients.Ctx, expectedProjectCreateArgs).
		Return(&mockedOperationReference, nil).
		Times(1)

	operationsClient.
		EXPECT().
		GetOperation(clients.Ctx, expectedOperationArgs).
		Return(nil, errors.New("GetOperation() failed")).
		Times(1)

	err := createProject(clients, &testProject, 5)
	require.Equal(t, "GetOperation() failed", err.Error())
}

// verifies that polling is done to validate the status of the asynchronous
// testProject create operation.
func TestAzureDevOpsProject_CreateProject_PollsUntilOperationIsSuccessful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	operationsClient := azdosdkmocks.NewMockOperationsClient(ctrl)
	clients := &config.AggregatedClient{
		CoreClient:       coreClient,
		OperationsClient: operationsClient,
		Ctx:              context.Background(),
	}

	expectedProjectCreateArgs := core.QueueCreateProjectArgs{ProjectToCreate: &testProject}
	mockedOperationReference := operations.OperationReference{Id: &testID}
	expectedOperationArgs := operations.GetOperationArgs{OperationId: &testID}

	coreClient.
		EXPECT().
		QueueCreateProject(clients.Ctx, expectedProjectCreateArgs).
		Return(&mockedOperationReference, nil).
		Times(1)

	firstStatus := operationWithStatus(operations.OperationStatusValues.InProgress)
	firstPoll := operationsClient.
		EXPECT().
		GetOperation(clients.Ctx, expectedOperationArgs).
		Return(&firstStatus, nil)

	secondStatus := operationWithStatus(operations.OperationStatusValues.Succeeded)
	secondPoll := operationsClient.
		EXPECT().
		GetOperation(clients.Ctx, expectedOperationArgs).
		Return(&secondStatus, nil)

	gomock.InOrder(firstPoll, secondPoll)

	err := createProject(clients, &testProject, 5)
	require.Equal(t, nil, err)
}

// verifies that if a project takes too long to create, an error is returned
func TestAzureDevOpsProject_CreateProject_ReportsErrorIfNoSuccessForLongTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	operationsClient := azdosdkmocks.NewMockOperationsClient(ctrl)
	clients := &config.AggregatedClient{
		CoreClient:       coreClient,
		OperationsClient: operationsClient,
		Ctx:              context.Background(),
	}

	expectedProjectCreateArgs := core.QueueCreateProjectArgs{ProjectToCreate: &testProject}
	mockedOperationReference := operations.OperationReference{Id: &testID}
	expectedOperationArgs := operations.GetOperationArgs{OperationId: &testID}

	coreClient.
		EXPECT().
		QueueCreateProject(clients.Ctx, expectedProjectCreateArgs).
		Return(&mockedOperationReference, nil).
		Times(1)

	// the operation will forever be "in progress"
	status := operationWithStatus(operations.OperationStatusValues.InProgress)
	operationsClient.
		EXPECT().
		GetOperation(clients.Ctx, expectedOperationArgs).
		Return(&status, nil).
		MinTimes(1)

	err := createProject(clients, &testProject, 5)
	require.NotNil(t, err, "Expected error indicating timeout")
}

func TestAzureDevOpsProject_FlattenExpand_RoundTrip(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	clients := &config.AggregatedClient{
		CoreClient: coreClient,
		Ctx:        context.Background(),
	}

	expectedProcesses := []core.Process{
		{
			Name: converter.String("TemplateName"),
			Id:   &testID,
		},
	}

	// mock the list of all process IDs. This is needed for the call to flattenProject()
	coreClient.
		EXPECT().
		GetProcesses(clients.Ctx, core.GetProcessesArgs{}).
		Return(&expectedProcesses, nil).
		Times(1)

	// mock the lookup of a specific process. This is needed for the call to expandProject()
	coreClient.
		EXPECT().
		GetProcessById(clients.Ctx, core.GetProcessByIdArgs{ProcessId: &testID}).
		Return(&expectedProcesses[0], nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, resourceProject().Schema, nil)
	err := flattenProject(clients, resourceData, &testProject)
	require.Nil(t, err)

	projectAfterRoundTrip, err := expandProject(clients, resourceData, true)
	require.Nil(t, err)
	require.Equal(t, testProject, *projectAfterRoundTrip)
}

// verifies that the project ID is used for reads if the ID is set
func TestAzureDevOpsProject_ProjectRead_UsesIdIfSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	clients := &config.AggregatedClient{
		CoreClient: coreClient,
		Ctx:        context.Background(),
	}

	id := "id"
	name := "name"

	coreClient.
		EXPECT().
		GetProject(clients.Ctx, core.GetProjectArgs{
			ProjectId:           &id,
			IncludeCapabilities: converter.Bool(true),
			IncludeHistory:      converter.Bool(false),
		}).
		Times(1)

	_, _ = ProjectRead(clients, id, name)
}

// verifies that the project name is used for reads if the ID is not set
func TestAzureDevOpsProject_ProjectRead_UsesNameIfIdNotSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	clients := &config.AggregatedClient{
		CoreClient: coreClient,
		Ctx:        context.Background(),
	}

	id := ""
	name := "name"

	coreClient.
		EXPECT().
		GetProject(clients.Ctx, core.GetProjectArgs{
			ProjectId:           &name,
			IncludeCapabilities: converter.Bool(true),
			IncludeHistory:      converter.Bool(false),
		}).
		Times(1)

	_, _ = ProjectRead(clients, id, name)
}

// creates an operation given a status
func operationWithStatus(status operations.OperationStatus) operations.Operation {
	return operations.Operation{Status: &status}
}
