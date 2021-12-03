//go:build (all || core || resource_project) && !exclude_resource_project
// +build all core resource_project
// +build !exclude_resource_project

package core

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/featuremanagement"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/operations"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/testhelper"
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
func TestProject_CreateProject_DoesNotSwallowErrorFromFailedCreateCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	clients := &client.AggregatedClient{
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
func TestProject_CreateProject_DoesNotSwallowErrorFromFailedAsyncStatusCheckCall(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	operationsClient := azdosdkmocks.NewMockOperationsClient(ctrl)
	clients := &client.AggregatedClient{
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

	err := createProject(clients, &testProject, 10*time.Minute)
	require.Equal(t, " waiting for project ready. GetOperation() failed ", err.Error())
}

// verifies that polling is done to validate the status of the asynchronous
// testProject create operation.
func TestProject_CreateProject_PollsUntilOperationIsSuccessful(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	operationsClient := azdosdkmocks.NewMockOperationsClient(ctrl)
	clients := &client.AggregatedClient{
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

	err := createProject(clients, &testProject, 10*time.Minute)
	require.Equal(t, nil, err)
}

// verifies that if a project takes too long to create, an error is returned
func TestProject_CreateProject_ReportsErrorIfNoSuccessForLongTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	operationsClient := azdosdkmocks.NewMockOperationsClient(ctrl)
	clients := &client.AggregatedClient{
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

	err := createProject(clients, &testProject, 20*time.Second)
	require.NotNil(t, err, "Expected error indicating timeout")
}

func TestProject_FlattenExpand_RoundTrip(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	clients := &client.AggregatedClient{
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

	resourceData := schema.TestResourceDataRaw(t, ResourceProject().Schema, nil)
	err := flattenProject(clients, resourceData, &testProject)
	require.Nil(t, err)

	projectAfterRoundTrip, err := expandProject(clients, resourceData, true)
	require.Nil(t, err)
	require.Equal(t, testProject, *projectAfterRoundTrip)
}

// verifies that the project ID is used for reads if the ID is set
func TestProject_ProjectRead_UsesIdIfSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	clients := &client.AggregatedClient{
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

	_, _ = projectRead(clients, id, name)
}

// verifies that the project name is used for reads if the ID is not set
func TestProject_ProjectRead_UsesNameIfIdNotSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	clients := &client.AggregatedClient{
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

	_, _ = projectRead(clients, id, name)
}

// creates an operation given a status
func operationWithStatus(status operations.OperationStatus) operations.Operation {
	return operations.Operation{Status: &status}
}

func TestAzureDevOpsProject_ConfigureProjectFeatures_HandleErrorCorrectly(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	featureClient := azdosdkmocks.NewMockFeaturemanagementClient(ctrl)
	operationsClient := azdosdkmocks.NewMockOperationsClient(ctrl)
	clients := &client.AggregatedClient{
		CoreClient:              coreClient,
		OperationsClient:        operationsClient,
		FeatureManagementClient: featureClient,
		Ctx:                     context.Background(),
	}

	const projectName = "test.project"
	const projectID = "925a1c6a-49f6-4a29-b3ae-f467a345545a"
	const errMsg = "GOMOCK: SetFeatureStateForScope failed"

	expectedGetProjectArgs := core.GetProjectArgs{
		ProjectId:           converter.String(projectName),
		IncludeCapabilities: converter.Bool(true),
		IncludeHistory:      converter.Bool(false),
	}

	coreClient.
		EXPECT().
		GetProject(clients.Ctx, expectedGetProjectArgs).
		Return(&core.TeamProject{
			Id:   converter.UUID(projectID),
			Name: converter.String(projectName),
		}, nil).
		Times(1)

	featureClient.
		EXPECT().
		SetFeatureStateForScope(clients.Ctx, gomock.Any()).
		Return(nil, errors.New(errMsg)).
		Times(1)

	expectedQueueDeleteProjectArgs := core.QueueDeleteProjectArgs{
		ProjectId: converter.UUID(projectID),
	}

	const operationID = "83a87383-807e-46d0-b9c3-35cd5030017a"
	const pluginID = "cbbbd22d-353a-474b-adb2-8e95945d58c3"

	coreClient.
		EXPECT().
		QueueDeleteProject(clients.Ctx, expectedQueueDeleteProjectArgs).
		Return(&operations.OperationReference{
			Id:       converter.UUID(operationID),
			PluginId: converter.UUID(pluginID),
		}, nil).
		Times(1)

	expectedGetOperationArgs := operations.GetOperationArgs{
		OperationId: converter.UUID(operationID),
		PluginId:    converter.UUID(pluginID),
	}

	operationsClient.
		EXPECT().
		GetOperation(clients.Ctx, expectedGetOperationArgs).
		Return(&operations.Operation{
			Id:       converter.UUID(operationID),
			PluginId: converter.UUID(pluginID),
			Status:   &operations.OperationStatusValues.Succeeded,
		}, nil).
		Times(1)

	featureIDs := *getProjectFeatureIDs()
	featureMap := *getProjectFeatureNameMap()
	idx := testhelper.RandInt(0, len(featureIDs)-1)

	featureStatesMap := make(map[string]interface{}, 1)
	featureStatesMap[string(featureMap[featureIDs[idx]])] = string(featuremanagement.ContributedFeatureEnabledValueValues.Enabled)
	featureStates := interface{}(featureStatesMap)

	err := configureProjectFeatures(clients, "", projectName, &featureStates, 10*time.Minute)
	require.NotNil(t, err)
}
