package azuredevops

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/core"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"
)

func init() {
	/* add code for test setup here */
}

var idList = []uuid.UUID{
	uuid.New(),
	uuid.New(),
	uuid.New(),
	uuid.New(),
	uuid.New(),
	uuid.New(),
}

var prjListEmpty = []core.TeamProjectReference{}

var prjListStateMixed = []core.TeamProjectReference{
	{
		Name:  converter.String("vsteam-0177"),
		Id:    &idList[0],
		State: &core.ProjectStateValues.WellFormed,
		Url:   nil,
	},
	{
		Name:  converter.String("vsteam-0178"),
		Id:    &idList[1],
		State: &core.ProjectStateValues.Deleted,
		Url:   nil,
	},
	{
		Name:  converter.String("vsteam-0179"),
		Id:    &idList[2],
		State: &core.ProjectStateValues.New,
		Url:   nil,
	},
}

var prjListStateWellFormed = []core.TeamProjectReference{
	{
		Name:  converter.String("vsteam-0177"),
		Id:    &idList[0],
		State: &core.ProjectStateValues.WellFormed,
		Url:   nil,
	},
	{
		Name:  converter.String("vsteam-0178"),
		Id:    &idList[1],
		State: &core.ProjectStateValues.WellFormed,
		Url:   nil,
	},
	{
		Name:  converter.String("vsteam-0179"),
		Id:    &idList[2],
		State: &core.ProjectStateValues.WellFormed,
		Url:   nil,
	},
}

var prjListStateWellFormed2 = []core.TeamProjectReference{
	{
		Name:  converter.String("vsteam-0277"),
		Id:    &idList[0+len(prjListStateWellFormed)],
		State: &core.ProjectStateValues.WellFormed,
		Url:   nil,
	},
	{
		Name:  converter.String("vsteam-0278"),
		Id:    &idList[1+len(prjListStateWellFormed)],
		State: &core.ProjectStateValues.WellFormed,
		Url:   nil,
	},
	{
		Name:  converter.String("vsteam-0279"),
		Id:    &idList[2+len(prjListStateWellFormed)],
		State: &core.ProjectStateValues.WellFormed,
		Url:   nil,
	},
}

var prjListDoubleID = []core.TeamProjectReference{
	{
		Name:  converter.String("vsteam-0177"),
		Id:    &idList[0],
		State: &core.ProjectStateValues.WellFormed,
		Url:   nil,
	},
	{
		Name:  converter.String("vsteam-0178"),
		Id:    &idList[1],
		State: &core.ProjectStateValues.WellFormed,
		Url:   nil,
	},
	{
		Name:  converter.String("vsteam-0179"),
		Id:    &idList[0],
		State: &core.ProjectStateValues.WellFormed,
		Url:   nil,
	},
}

/**
 * Begin unit tests
 */

func TestDataSourceProjects_Read_TestFindProjectByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	clients := &config.AggregatedClient{
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

	resourceData := schema.TestResourceDataRaw(t, dataProjects().Schema, nil)
	resourceData.Set("project_name", "vsteam-0178")
	resourceData.Set("state", "wellFormed")
	err := dataSourceProjectsRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, "wellFormed", resourceData.Get("state").(string))
	require.Equal(t, "vsteam-0178", resourceData.Get("project_name").(string))
	projectSet := resourceData.Get("projects").(*schema.Set)
	require.NotNil(t, projectSet)
	require.Equal(t, 1, projectSet.Len())
	projectReference := projectSet.List()[0].(map[string]interface{})
	require.NotNil(t, projectReference)
	require.Equal(t, "vsteam-0178", projectReference["name"])
	require.Equal(t, idList[1].String(), projectReference["project_id"])
	require.Equal(t, "wellFormed", projectReference["state"])
}

func TestDataSourceProjects_Read_TestEmptyProjectList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	clients := &config.AggregatedClient{
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
			Value:             prjListEmpty,
			ContinuationToken: "",
		}, nil).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, dataProjects().Schema, nil)
	err := dataSourceProjectsRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, "all", resourceData.Get("state").(string))
	require.Equal(t, "", resourceData.Get("project_name").(string))
	projectSet := resourceData.Get("projects").(*schema.Set)
	require.NotNil(t, projectSet)
	require.Equal(t, 0, projectSet.Len())
}

func TestDataSourceProjects_Read_TestFindAllProjects(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	clients := &config.AggregatedClient{
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

	resourceData := schema.TestResourceDataRaw(t, dataProjects().Schema, nil)
	err := dataSourceProjectsRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, "all", resourceData.Get("state").(string))
	require.Equal(t, "", resourceData.Get("project_name").(string))
	projectSet := resourceData.Get("projects").(*schema.Set)
	require.NotNil(t, projectSet)
	require.Equal(t, 3, projectSet.Len())
}

func TestDataSourceProjects_Read_TestDuplicateProjectId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	clients := &config.AggregatedClient{
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

	resourceData := schema.TestResourceDataRaw(t, dataProjects().Schema, nil)
	err := dataSourceProjectsRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, "all", resourceData.Get("state").(string))
	require.Equal(t, "", resourceData.Get("project_name").(string))
	projectSet := resourceData.Get("projects").(*schema.Set)
	require.NotNil(t, projectSet)
	require.Equal(t, 2, projectSet.Len())
}

func TestDataSourceProjects_Read_TestFindProjectsWithState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	clients := &config.AggregatedClient{
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

	resourceData := schema.TestResourceDataRaw(t, dataProjects().Schema, nil)
	resourceData.Set("state", "wellFormed")
	err := dataSourceProjectsRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, "wellFormed", resourceData.Get("state").(string))
	require.Equal(t, "", resourceData.Get("project_name").(string))
	projectSet := resourceData.Get("projects").(*schema.Set)
	require.NotNil(t, projectSet)
	require.Equal(t, 3, projectSet.Len())
}

func TestDataSourceProjects_Read_TestHandleError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	clients := &config.AggregatedClient{
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

	resourceData := schema.TestResourceDataRaw(t, dataProjects().Schema, nil)
	err := dataSourceProjectsRead(resourceData, clients)
	require.Error(t, err)
	require.Contains(t, err.Error(), "GetProjects() Failed")
}

func TestDataSourceProjects_Read_TestContinuationToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	clients := &config.AggregatedClient{
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

	resourceData := schema.TestResourceDataRaw(t, dataProjects().Schema, nil)
	err := dataSourceProjectsRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, "all", resourceData.Get("state").(string))
	require.Equal(t, "", resourceData.Get("project_name").(string))
	projectSet := resourceData.Get("projects").(*schema.Set)
	require.NotNil(t, projectSet)
	require.Equal(t, 6, projectSet.Len())
}
