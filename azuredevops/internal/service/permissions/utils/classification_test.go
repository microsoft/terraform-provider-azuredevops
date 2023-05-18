//go:build all || utils || classification_helper
// +build all utils classification_helper

package utils

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/assert"
)

var iterationProjectID = "a417ffff-fb0d-4cd4-8aac-54d8878b60f0"
var iterationRootID = "0b401c26-b0da-4655-995a-ab62f0b05187"

func TestClassificationNode_CreateIterationToken_RootIteration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workitemtrackingClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{
		WorkItemTrackingClient: workitemtrackingClient,
		Ctx:                    context.Background(),
	}

	for _, path := range []string{"", "/", "    ", "    /", "/   "} {
		workitemtrackingClient.
			EXPECT().
			GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
				Project:        &iterationProjectID,
				Path:           nil,
				StructureGroup: &workitemtracking.TreeStructureGroupValues.Iterations,
				Depth:          converter.Int(1),
			}).
			Return(&workitemtracking.WorkItemClassificationNode{
				Identifier: converter.UUID(iterationRootID),
			}, nil).
			Times(1)

		token, err := CreateClassificationNodeSecurityToken(clients.Ctx, clients.WorkItemTrackingClient, workitemtracking.TreeStructureGroupValues.Iterations, iterationProjectID, path)
		assert.Nil(t, err)
		ref := fmt.Sprintf("%s%s", aclClassificationNodeTokenPrefix, iterationRootID)
		assert.Equal(t, ref, token)
	}
}

func TestClassificationNode_CreateIterationToken_HandleError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workitemtrackingClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{
		WorkItemTrackingClient: workitemtrackingClient,
		Ctx:                    context.Background(),
	}

	const errMsg = "@@GetClassificationNode@@failed"
	workitemtrackingClient.
		EXPECT().
		GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
			Project:        &iterationProjectID,
			Path:           nil,
			StructureGroup: &workitemtracking.TreeStructureGroupValues.Iterations,
			Depth:          converter.Int(1),
		}).
		Return(nil, fmt.Errorf(errMsg)).
		Times(1)

	token, err := CreateClassificationNodeSecurityToken(clients.Ctx, clients.WorkItemTrackingClient, workitemtracking.TreeStructureGroupValues.Iterations, iterationProjectID, "/")
	assert.Empty(t, token)
	assert.NotNil(t, err)
}

func TestClassificationNode_CreateIterationToken_HandleErrorInPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workitemtrackingClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{
		WorkItemTrackingClient: workitemtrackingClient,
		Ctx:                    context.Background(),
	}

	var errMsg = "@@GetClassificationNode@@failed"

	workitemtrackingClient.
		EXPECT().
		GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
			Project:        &iterationProjectID,
			Path:           nil,
			StructureGroup: &workitemtracking.TreeStructureGroupValues.Iterations,
			Depth:          converter.Int(1),
		}).
		Return(&workitemtracking.WorkItemClassificationNode{
			Identifier:  converter.UUID(iterationRootID),
			HasChildren: converter.Bool(true),
		}, nil).
		Times(1)

	workitemtrackingClient.
		EXPECT().
		GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
			Project:        &iterationProjectID,
			Path:           converter.String("iteration"),
			StructureGroup: &workitemtracking.TreeStructureGroupValues.Iterations,
			Depth:          converter.Int(1),
		}).
		Return(nil, fmt.Errorf(errMsg)).
		Times(1)

	token, err := CreateClassificationNodeSecurityToken(clients.Ctx, clients.WorkItemTrackingClient, workitemtracking.TreeStructureGroupValues.Iterations, iterationProjectID, "/iteration")
	assert.Empty(t, token)
	assert.NotNil(t, err)
}

func TestClassificationNode_CreateIterationToken_HandleNoChildren(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workitemtrackingClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{
		WorkItemTrackingClient: workitemtrackingClient,
		Ctx:                    context.Background(),
	}

	workitemtrackingClient.
		EXPECT().
		GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
			Project:        &iterationProjectID,
			Path:           nil,
			StructureGroup: &workitemtracking.TreeStructureGroupValues.Iterations,
			Depth:          converter.Int(1),
		}).
		Return(&workitemtracking.WorkItemClassificationNode{
			Identifier: converter.UUID(iterationRootID),
		}, nil).
		Times(1)

	token, err := CreateClassificationNodeSecurityToken(clients.Ctx, clients.WorkItemTrackingClient, workitemtracking.TreeStructureGroupValues.Iterations, iterationProjectID, "/iteration")
	assert.Empty(t, token)
	assert.NotNil(t, err)
}

func TestClassificationNode_CreateIterationToken_ValidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workitemtrackingClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{
		WorkItemTrackingClient: workitemtrackingClient,
		Ctx:                    context.Background(),
	}

	workitemtrackingClient.
		EXPECT().
		GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
			Project:        &iterationProjectID,
			Path:           nil,
			StructureGroup: &workitemtracking.TreeStructureGroupValues.Iterations,
			Depth:          converter.Int(1),
		}).
		Return(&workitemtracking.WorkItemClassificationNode{
			Identifier:  converter.UUID(iterationRootID),
			HasChildren: converter.Bool(true),
		}, nil).
		Times(1)

	const count = 3
	path := "/"
	idList := make([]string, count)

	for i := 0; i < count; i++ {
		pathItem := fmt.Sprintf("iteration_%d", i)
		if i == 0 {
			path = pathItem
		} else {
			path = path + "/" + pathItem
		}
		idList[i] = uuid.New().String()
		workitemtrackingClient.
			EXPECT().
			GetClassificationNode(clients.Ctx, workitemtracking.GetClassificationNodeArgs{
				Project:        &iterationProjectID,
				Path:           converter.String(path),
				StructureGroup: &workitemtracking.TreeStructureGroupValues.Iterations,
				Depth:          converter.Int(1),
			}).
			Return(&workitemtracking.WorkItemClassificationNode{
				Identifier:  converter.UUID(idList[i]),
				HasChildren: converter.Bool(i+1 < count),
			}, nil).
			Times(1)

		idList[i] = aclClassificationNodeTokenPrefix + idList[i]
	}

	token, err := CreateClassificationNodeSecurityToken(clients.Ctx, clients.WorkItemTrackingClient, workitemtracking.TreeStructureGroupValues.Iterations, iterationProjectID, path)
	assert.Nil(t, err)
	ref := fmt.Sprintf("%s%s:%s", aclClassificationNodeTokenPrefix, iterationRootID, strings.Join(idList, ":"))
	assert.Equal(t, ref, token)
}
