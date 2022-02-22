//go:build (all || core || data_sources || resource_project || data_projects) && (!data_sources || !exclude_data_projects)
// +build all core data_sources resource_project data_projects
// +build !data_sources !exclude_data_projects

package core

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/core"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/testhelper"
	"github.com/stretchr/testify/require"
)

var prjListStateWellFormed = []core.TeamProjectReference{
	{
		Name:  converter.String("vsteam-0177"),
		Id:    testhelper.CreateUUID(),
		State: &core.ProjectStateValues.WellFormed,
		Url:   nil,
	},
	{
		Name:  converter.String("vsteam-0178"),
		Id:    testhelper.CreateUUID(),
		State: &core.ProjectStateValues.WellFormed,
		Url:   nil,
	},
	{
		Name:  converter.String("vsteam-0179"),
		Id:    testhelper.CreateUUID(),
		State: &core.ProjectStateValues.WellFormed,
		Url:   nil,
	},
}

var prjListStateWellFormed2 = []core.TeamProjectReference{
	{
		Name:  converter.String("vsteam-0277"),
		Id:    testhelper.CreateUUID(),
		State: &core.ProjectStateValues.WellFormed,
		Url:   nil,
	},
	{
		Name:  converter.String("vsteam-0278"),
		Id:    testhelper.CreateUUID(),
		State: &core.ProjectStateValues.WellFormed,
		Url:   nil,
	},
	{
		Name:  converter.String("vsteam-0279"),
		Id:    testhelper.CreateUUID(),
		State: &core.ProjectStateValues.WellFormed,
		Url:   nil,
	},
}

var duplicatePrjID *uuid.UUID = testhelper.CreateUUID()
var prjListDoubleID = []core.TeamProjectReference{
	{
		Name:  converter.String("vsteam-0177"),
		Id:    duplicatePrjID,
		State: &core.ProjectStateValues.WellFormed,
		Url:   nil,
	},
	{
		Name:  converter.String("vsteam-0178"),
		Id:    testhelper.CreateUUID(),
		State: &core.ProjectStateValues.WellFormed,
		Url:   nil,
	},
	{
		Name:  converter.String("vsteam-0179"),
		Id:    duplicatePrjID,
		State: &core.ProjectStateValues.WellFormed,
		Url:   nil,
	},
}

func TestDataSourceProjects_Read_TestFindProjectByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	clients := &client.AggregatedClient{
		CoreClient: coreClient,
		Ctx:        context.Background(),
	}

	expectedGetProjectsArgs := core.GetProjectsArgs{
		StateFilter: &core.ProjectStateValues.WellFormed,
	}

	coreClient.
		EXPECT().
		GetProjects(clients.Ctx, expectedGetProjectsArgs).
		Return(&core.GetProjectsResponseValue{
			Value:             prjListStateWellFormed,
			ContinuationToken: "",
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataProjects().Schema, nil)
	resourceData.Set("name", "vsteam-0178")
	resourceData.Set("state", "wellFormed")
	err := dataSourceProjectsRead(clients.Ctx, resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, "wellFormed", resourceData.Get("state").(string))
	require.Equal(t, "vsteam-0178", resourceData.Get("name").(string))
	projectSet := resourceData.Get("projects").(*schema.Set)
	require.NotNil(t, projectSet)
	require.Equal(t, 1, projectSet.Len())
	projectReference := projectSet.List()[0].(map[string]interface{})
	require.NotNil(t, projectReference)
	require.Equal(t, "vsteam-0178", projectReference["name"])
	require.Equal(t, prjListStateWellFormed[1].Id.String(), projectReference["project_id"])
	require.Equal(t, "wellFormed", projectReference["state"])
}

func TestDataSourceProjects_Read_TestEmptyProjectList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	clients := &client.AggregatedClient{
		CoreClient: coreClient,
		Ctx:        context.Background(),
	}

	expectedGetProjectsArgs := core.GetProjectsArgs{
		StateFilter: &core.ProjectStateValues.All,
	}

	coreClient.
		EXPECT().
		GetProjects(clients.Ctx, expectedGetProjectsArgs).
		Return(&core.GetProjectsResponseValue{
			Value:             []core.TeamProjectReference{},
			ContinuationToken: "",
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataProjects().Schema, nil)
	err := dataSourceProjectsRead(clients.Ctx, resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, "all", resourceData.Get("state").(string))
	require.Equal(t, "", resourceData.Get("name").(string))
	projectSet := resourceData.Get("projects").(*schema.Set)
	require.NotNil(t, projectSet)
	require.Equal(t, 0, projectSet.Len())
}

func TestDataSourceProjects_Read_TestFindAllProjects(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	clients := &client.AggregatedClient{
		CoreClient: coreClient,
		Ctx:        context.Background(),
	}

	expectedGetProjectsArgs := core.GetProjectsArgs{
		StateFilter: &core.ProjectStateValues.All,
	}

	coreClient.
		EXPECT().
		GetProjects(clients.Ctx, expectedGetProjectsArgs).
		Return(&core.GetProjectsResponseValue{
			Value:             prjListStateWellFormed,
			ContinuationToken: "",
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataProjects().Schema, nil)
	err := dataSourceProjectsRead(clients.Ctx, resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, "all", resourceData.Get("state").(string))
	require.Equal(t, "", resourceData.Get("name").(string))
	projectSet := resourceData.Get("projects").(*schema.Set)
	require.NotNil(t, projectSet)
	require.Equal(t, 3, projectSet.Len())
}

func TestDataSourceProjects_Read_TestDuplicateProjectId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	clients := &client.AggregatedClient{
		CoreClient: coreClient,
		Ctx:        context.Background(),
	}

	expectedGetProjectsArgs := core.GetProjectsArgs{
		StateFilter: &core.ProjectStateValues.All,
	}

	coreClient.
		EXPECT().
		GetProjects(clients.Ctx, expectedGetProjectsArgs).
		Return(&core.GetProjectsResponseValue{
			Value:             prjListDoubleID,
			ContinuationToken: "",
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataProjects().Schema, nil)
	err := dataSourceProjectsRead(clients.Ctx, resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, "all", resourceData.Get("state").(string))
	require.Equal(t, "", resourceData.Get("name").(string))
	projectSet := resourceData.Get("projects").(*schema.Set)
	require.NotNil(t, projectSet)
	require.Equal(t, 2, projectSet.Len())
}

func TestDataSourceProjects_Read_TestFindProjectsWithState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	clients := &client.AggregatedClient{
		CoreClient: coreClient,
		Ctx:        context.Background(),
	}

	expectedGetProjectsArgs := core.GetProjectsArgs{
		StateFilter: &core.ProjectStateValues.WellFormed,
	}

	coreClient.
		EXPECT().
		GetProjects(clients.Ctx, expectedGetProjectsArgs).
		Return(&core.GetProjectsResponseValue{
			Value:             prjListStateWellFormed,
			ContinuationToken: "",
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataProjects().Schema, nil)
	resourceData.Set("state", "wellFormed")
	err := dataSourceProjectsRead(clients.Ctx, resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, "wellFormed", resourceData.Get("state").(string))
	require.Equal(t, "", resourceData.Get("name").(string))
	projectSet := resourceData.Get("projects").(*schema.Set)
	require.NotNil(t, projectSet)
	require.Equal(t, 3, projectSet.Len())
}

func TestDataSourceProjects_Read_TestHandleError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	clients := &client.AggregatedClient{
		CoreClient: coreClient,
		Ctx:        context.Background(),
	}

	expectedGetProjectsArgs := core.GetProjectsArgs{
		StateFilter: &core.ProjectStateValues.All,
	}

	coreClient.
		EXPECT().
		GetProjects(clients.Ctx, expectedGetProjectsArgs).
		Return(nil, errors.New("GetProjects() Failed")).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataProjects().Schema, nil)
	err := dataSourceProjectsRead(clients.Ctx, resourceData, clients)
	require.Equal(t, err.HasError(), true)
	require.Contains(t, err[0].Summary, "GetProjects() Failed")
}

func TestDataSourceProjects_Read_TestContinuationToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	clients := &client.AggregatedClient{
		CoreClient: coreClient,
		Ctx:        context.Background(),
	}

	var calls []*gomock.Call
	calls = append(calls, coreClient.
		EXPECT().
		GetProjects(clients.Ctx, core.GetProjectsArgs{
			StateFilter: &core.ProjectStateValues.All,
		}).
		Return(&core.GetProjectsResponseValue{
			Value:             prjListStateWellFormed,
			ContinuationToken: "2",
		}, nil).
		Times(1))

	calls = append(calls, coreClient.
		EXPECT().
		GetProjects(clients.Ctx, core.GetProjectsArgs{
			StateFilter:       &core.ProjectStateValues.All,
			ContinuationToken: converter.String("2"),
		}).
		Return(&core.GetProjectsResponseValue{
			Value:             prjListStateWellFormed2,
			ContinuationToken: "",
		}, nil).
		Times(1))

	gomock.InOrder(calls...)

	resourceData := schema.TestResourceDataRaw(t, DataProjects().Schema, nil)
	err := dataSourceProjectsRead(clients.Ctx, resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, "all", resourceData.Get("state").(string))
	require.Equal(t, "", resourceData.Get("name").(string))
	projectSet := resourceData.Get("projects").(*schema.Set)
	require.NotNil(t, projectSet)
	require.Equal(t, 6, projectSet.Len())
}
