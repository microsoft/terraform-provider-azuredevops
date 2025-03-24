//go:build (all || permissions || resource_workitemquery_permissions) && (!exclude_permissions || !resource_workitemquery_permissions)
// +build all permissions resource_workitemquery_permissions
// +build !exclude_permissions !resource_workitemquery_permissions

package permissions

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/workitemtracking"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var wiqProjectID = "f454422e-57b3-442a-8dde-b1b6b7c40b95"
var wiqSharedQueryID = uuid.MustParse("ebfd6f15-1411-4b7b-86ea-41a1b1a9d38d")
var wiqSharedQueryName = "Shared Queries"
var wiqFldrID = uuid.MustParse("0d5eea5f-6d96-4802-82fc-867df52d2014")
var wiqFldrName = "folder"

func TestWorkItemQueryPermissions_CreateWorkItemQueryToken_ProjectGlobal(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workitemtrackingClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{
		WorkItemTrackingClient: workitemtrackingClient,
		Ctx:                    context.Background(),
	}

	d := getWorkItemQueryPermissionsResource(t, wiqProjectID, "")
	token, err := createWorkItemQueryToken(d, clients)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
	ref := fmt.Sprintf("$/%s", wiqProjectID)
	assert.Equal(t, ref, token)
}

func TestWorkItemQueryPermissions_CreateWorkItemQueryToken_SharedQueries(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workitemtrackingClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{
		WorkItemTrackingClient: workitemtrackingClient,
		Ctx:                    context.Background(),
	}

	workitemtrackingClient.
		EXPECT().
		GetQuery(clients.Ctx, workitemtracking.GetQueryArgs{
			Project: &wiqProjectID,
			Query:   converter.String(wiqSharedQueryName),
			Depth:   converter.Int(1),
		}).
		Return(&workitemtracking.QueryHierarchyItem{
			Id:   &wiqSharedQueryID,
			Name: converter.String(wiqSharedQueryName),
		}, nil).
		Times(1)

	d := getWorkItemQueryPermissionsResource(t, wiqProjectID, "/")
	token, err := createWorkItemQueryToken(d, clients)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
	ref := fmt.Sprintf("$/%s/%s", wiqProjectID, wiqSharedQueryID)
	assert.Equal(t, ref, token)
}

func TestWorkItemQueryPermissions_CreateWorkItemQueryToken_HandleError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workitemtrackingClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{
		WorkItemTrackingClient: workitemtrackingClient,
		Ctx:                    context.Background(),
	}

	var errMsg = "@@GetQuery@@failed"
	workitemtrackingClient.
		EXPECT().
		GetQuery(clients.Ctx, workitemtracking.GetQueryArgs{
			Project: &wiqProjectID,
			Query:   converter.String(wiqSharedQueryName),
			Depth:   converter.Int(1),
		}).
		Return(nil, fmt.Errorf(errMsg)).
		Times(1)

	d := getWorkItemQueryPermissionsResource(t, wiqProjectID, "/")
	token, err := createWorkItemQueryToken(d, clients)
	assert.Empty(t, token)
	assert.NotNil(t, err)
	assert.EqualError(t, err, errMsg)
}

func TestWorkItemQueryPermissions_CreateWorkItemQueryToken_HandleErrorInPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workitemtrackingClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{
		WorkItemTrackingClient: workitemtrackingClient,
		Ctx:                    context.Background(),
	}

	var errMsg = "@@GetQuery@@failed"

	workitemtrackingClient.
		EXPECT().
		GetQuery(clients.Ctx, workitemtracking.GetQueryArgs{
			Project: &wiqProjectID,
			Query:   converter.String(wiqSharedQueryName),
			Depth:   converter.Int(1),
		}).
		Return(&workitemtracking.QueryHierarchyItem{
			Id:   &wiqSharedQueryID,
			Name: converter.String(wiqSharedQueryName),
			Children: &[]workitemtracking.QueryHierarchyItem{
				{
					Id:   &wiqFldrID,
					Name: converter.String(wiqFldrName),
				},
			},
		}, nil).
		Times(1)

	workitemtrackingClient.
		EXPECT().
		GetQuery(clients.Ctx, workitemtracking.GetQueryArgs{
			Project: &wiqProjectID,
			Query:   converter.String(wiqFldrID.String()),
			Depth:   converter.Int(1),
		}).
		Return(nil, fmt.Errorf(errMsg)).
		Times(1)

	d := getWorkItemQueryPermissionsResource(t, wiqProjectID, "/folder")
	token, err := createWorkItemQueryToken(d, clients)
	assert.Empty(t, token)
	assert.NotNil(t, err)
}

func TestWorkItemQueryPermissions_CreateWorkItemQueryToken_ChildDoesNotExist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workitemtrackingClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{
		WorkItemTrackingClient: workitemtrackingClient,
		Ctx:                    context.Background(),
	}

	workitemtrackingClient.
		EXPECT().
		GetQuery(clients.Ctx, workitemtracking.GetQueryArgs{
			Project: &wiqProjectID,
			Query:   converter.String(wiqSharedQueryName),
			Depth:   converter.Int(1),
		}).
		Return(&workitemtracking.QueryHierarchyItem{
			Id:   &wiqSharedQueryID,
			Name: converter.String(wiqSharedQueryName),
			Children: &[]workitemtracking.QueryHierarchyItem{
				{
					Id:   &wiqFldrID,
					Name: converter.String(wiqFldrName),
				},
			},
		}, nil).
		Times(1)

	workitemtrackingClient.
		EXPECT().
		GetQuery(clients.Ctx, workitemtracking.GetQueryArgs{
			Project: &wiqProjectID,
			Query:   converter.String(wiqFldrID.String()),
			Depth:   converter.Int(1),
		}).
		Return(&workitemtracking.QueryHierarchyItem{
			Id:   &wiqFldrID,
			Name: converter.String(wiqFldrName),
		}, nil).
		Times(1)

	d := getWorkItemQueryPermissionsResource(t, wiqProjectID, "/folder/child")
	token, err := createWorkItemQueryToken(d, clients)
	assert.NotNil(t, err)
	assert.Empty(t, token)
}

func TestWorkItemQueryPermissions_CreateWorkItemQueryToken_ValidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	workitemtrackingClient := azdosdkmocks.NewMockWorkitemtrackingClient(ctrl)
	clients := &client.AggregatedClient{
		WorkItemTrackingClient: workitemtrackingClient,
		Ctx:                    context.Background(),
	}

	workitemtrackingClient.
		EXPECT().
		GetQuery(clients.Ctx, workitemtracking.GetQueryArgs{
			Project: &wiqProjectID,
			Query:   converter.String(wiqSharedQueryName),
			Depth:   converter.Int(1),
		}).
		Return(&workitemtracking.QueryHierarchyItem{
			Id:   &wiqSharedQueryID,
			Name: converter.String(wiqSharedQueryName),
			Children: &[]workitemtracking.QueryHierarchyItem{
				{
					Id:   &wiqFldrID,
					Name: converter.String(wiqFldrName),
				},
			},
		}, nil).
		Times(1)

	workitemtrackingClient.
		EXPECT().
		GetQuery(clients.Ctx, workitemtracking.GetQueryArgs{
			Project: &wiqProjectID,
			Query:   converter.String(wiqFldrID.String()),
			Depth:   converter.Int(1),
		}).
		Return(&workitemtracking.QueryHierarchyItem{
			Id:   &wiqFldrID,
			Name: converter.String(wiqFldrName),
		}, nil).
		Times(1)

	d := getWorkItemQueryPermissionsResource(t, wiqProjectID, "/folder")
	token, err := createWorkItemQueryToken(d, clients)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)
	ref := fmt.Sprintf("$/%s/%s/%s", wiqProjectID, wiqSharedQueryID, wiqFldrID)
	assert.Equal(t, ref, token)
}

func getWorkItemQueryPermissionsResource(t *testing.T, projectID string, path string) *schema.ResourceData {
	d := schema.TestResourceDataRaw(t, ResourceWorkItemQueryPermissions().Schema, nil)
	if projectID != "" {
		d.Set("project_id", projectID)
	}
	if path != "" {
		d.Set("path", path)
	}
	return d
}
