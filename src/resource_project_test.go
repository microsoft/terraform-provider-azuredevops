package main

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestProjectCreate_MapsTfState(t *testing.T) {

	testValues := &projectValues{
		projectName:      "Test Project",
		description:      "A description",
		visibility:       "public",
		versionControl:   "git",
		workItemTemplate: "Agile",
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := NewMockCoreClient(ctrl)
	coreClient.EXPECT().QueueCreateProject(gomock.Not(gomock.Nil()), gomock.Not(gomock.Nil())).Return(nil, errors.New("Whoops")).MinTimes(1).MaxTimes(1)
	clients := &aggregatedClient{
		CoreClient: coreClient,
		ctx:        context.Background(),
	}

	err := projectCreate(clients, testValues)
	require.Equal(t, "Whoops", err.Error())
}
