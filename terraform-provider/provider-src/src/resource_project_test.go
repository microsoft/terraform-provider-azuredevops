package main

import (
	"context"
	"errors"
	"testing"

	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/operations"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCoreClient struct {
	mock.Mock
}

func (m *MockCoreClient) GetProcesses(ctx context.Context, args core.GetProcessesArgs) (*[]core.Process, error) {
	return nil, errors.New("Whoops")
}
	
func (m *MockCoreClient) GetProjects(ctx context.Context, args core.GetProjectsArgs) (*core.GetProjectsResponseValue, error) {
	return nil, errors.New("Whoops")
}

func (m *MockCoreClient) QueueCreateProject(ctx context.Context, args core.QueueCreateProjectArgs) (*operations.OperationReference, error) {
	return nil, errors.New("Whoops")
}

func TestProjectCreate_MapsTfState(t *testing.T) {

	testValues := &projectValues{
		projectName:      "Test Project",
		description:      "A description",
		visibility:       "public",
		versionControl:   "git",
		workItemTemplate: "Agile",
	}

	clients := &AggregatedClient{CoreClient: &MockCoreClient{}}
	//ctx, _ := context.WithDeadline(context.Background(), time.Now())
	projectCreate(context.Background(), clients, testValues) //map[string]string, error) 
}