//go:build (all || core || data_sources || data_teams) && (!exclude_data_sources || !exclude_data_teams)
// +build all core data_sources data_teams
// +build !exclude_data_sources !exclude_data_teams

package core

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/security"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	securityhelper "github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/service/permissions/utils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

func TestDataTeams_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)

	clients := &client.AggregatedClient{
		CoreClient: coreClient,
		Ctx:        context.Background(),
	}

	testProjectID := uuid.New()

	coreClient.
		EXPECT().
		GetTeams(clients.Ctx, core.GetTeamsArgs{
			ProjectId:      converter.String(testProjectID.String()),
			Mine:           converter.Bool(false),
			Top:            converter.Int(100),
			ExpandIdentity: converter.Bool(false),
		}).
		Return(nil, fmt.Errorf("@@GetTeams@@failed@@")).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataTeams().Schema, nil)
	resourceData.Set("project_id", testProjectID.String())
	resourceData.Set("project_id", testProjectID.String())
	err := dataTeamsRead(resourceData, clients)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "@@GetTeams@@failed@@")
}

func TestDataTeams_Read_DoesNotSwallowErrorAllProjects(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)

	clients := &client.AggregatedClient{
		CoreClient: coreClient,
		Ctx:        context.Background(),
	}

	coreClient.EXPECT().
		GetProjects(clients.Ctx, core.GetProjectsArgs{
			StateFilter: &core.ProjectStateValues.All,
		}).
		Return(nil, fmt.Errorf("@@GetProjects@@failed@@")).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataTeams().Schema, nil)
	err := dataTeamsRead(resourceData, clients)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), "@@GetProjects@@failed@@")
}

func TestDataTeams_Read_EnsureAllByProject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)

	clients := &client.AggregatedClient{
		CoreClient:     coreClient,
		IdentityClient: identityClient,
		SecurityClient: securityClient,
		Ctx:            context.Background(),
	}

	testProjectID := uuid.New()
	teamList := []struct {
		name        string
		id          uuid.UUID
		description string
	}{
		{
			name:        "@@TEST TEAM@@",
			id:          uuid.New(),
			description: "@@TEST TEAM@@DESCRIPTION@@",
		},
		{
			name:        "@@TEST TEAM@@2",
			id:          uuid.New(),
			description: "@@TEST TEAM@@DESCRIPTION@@2",
		},
	}

	coreClient.
		EXPECT().
		GetTeams(clients.Ctx, core.GetTeamsArgs{
			ProjectId:      converter.String(testProjectID.String()),
			Mine:           converter.Bool(false),
			Top:            converter.Int(100),
			ExpandIdentity: converter.Bool(false),
		}).
		Return(&[]core.WebApiTeam{
			{
				Id:          &teamList[0].id,
				Name:        &teamList[0].name,
				Description: &teamList[0].description,
				ProjectId:   &testProjectID,
			},
			{
				Id:          &teamList[1].id,
				Name:        &teamList[1].name,
				Description: &teamList[1].description,
				ProjectId:   &testProjectID,
			},
		}, nil).
		Times(1)

	identityClient.
		EXPECT().
		ReadMembers(clients.Ctx, gomock.Any()).
		Return(nil, nil).
		Times(2)

	nsID := uuid.UUID(securityhelper.SecurityNamespaceIDValues.Identity)
	securityClient.
		EXPECT().
		QuerySecurityNamespaces(clients.Ctx, security.QuerySecurityNamespacesArgs{
			SecurityNamespaceId: &nsID,
		}).
		Return(&[]security.SecurityNamespaceDescription{
			{
				Actions: &[]security.ActionDefinition{
					{
						Bit:  converter.Int(8),
						Name: converter.String("ManageMembership"),
					},
				},
			},
		}, nil).
		Times(2)

	securityClient.
		EXPECT().
		QueryAccessControlLists(clients.Ctx, gomock.Any()).
		Return(&[]security.AccessControlList{}, nil).
		Times(2)

	resourceData := schema.TestResourceDataRaw(t, DataTeams().Schema, nil)
	resourceData.Set("project_id", testProjectID.String())
	err := dataTeamsRead(resourceData, clients)

	require.Nil(t, err)
	require.Equal(t, testProjectID.String(), resourceData.Get("project_id"))

	data, ok := resourceData.GetOk("teams")
	require.True(t, ok)
	teamSet := data.([]interface{})
	require.NotNil(t, teamSet)
	require.Equal(t, len(teamList), len(teamSet))

	teamMap := make(map[string]map[string]interface{}, len(teamSet))
	for _, e := range teamSet {
		team := e.(map[string]interface{})
		teamMap[team["id"].(string)] = team
	}
	for _, item := range teamList {
		team, ok := teamMap[item.id.String()]
		require.True(t, ok)
		require.Equal(t, item.name, team["name"])
		require.Equal(t, testProjectID.String(), team["project_id"])
		require.Equal(t, item.description, team["description"])
	}
}

func TestDataTeams_Read_EnsureAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	coreClient := azdosdkmocks.NewMockCoreClient(ctrl)
	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	securityClient := azdosdkmocks.NewMockSecurityClient(ctrl)

	clients := &client.AggregatedClient{
		CoreClient:     coreClient,
		IdentityClient: identityClient,
		SecurityClient: securityClient,
		Ctx:            context.Background(),
	}

	testProjectID := uuid.New()

	coreClient.EXPECT().
		GetProjects(clients.Ctx, core.GetProjectsArgs{
			StateFilter: &core.ProjectStateValues.All,
		}).
		Return(&core.GetProjectsResponseValue{
			Value: []core.TeamProjectReference{
				{
					Id: &testProjectID,
				},
			},
			ContinuationToken: "",
		}, nil).
		Times(1)

	teamList := []struct {
		name        string
		id          uuid.UUID
		description string
	}{
		{
			name:        "@@TEST TEAM@@",
			id:          uuid.New(),
			description: "@@TEST TEAM@@DESCRIPTION@@",
		},
		{
			name:        "@@TEST TEAM@@2",
			id:          uuid.New(),
			description: "@@TEST TEAM@@DESCRIPTION@@2",
		},
	}

	coreClient.
		EXPECT().
		GetTeams(clients.Ctx, core.GetTeamsArgs{
			ProjectId:      converter.String(testProjectID.String()),
			Mine:           converter.Bool(false),
			ExpandIdentity: converter.Bool(false),
			Top:            converter.Int(100),
		}).
		Return(&[]core.WebApiTeam{
			{
				Id:          &teamList[0].id,
				Name:        &teamList[0].name,
				Description: &teamList[0].description,
				ProjectId:   &testProjectID,
			},
			{
				Id:          &teamList[1].id,
				Name:        &teamList[1].name,
				Description: &teamList[1].description,
				ProjectId:   &testProjectID,
			},
		}, nil).
		Times(1)

	identityClient.
		EXPECT().
		ReadMembers(clients.Ctx, gomock.Any()).
		Return(nil, nil).
		Times(2)

	nsID := uuid.UUID(securityhelper.SecurityNamespaceIDValues.Identity)
	securityClient.
		EXPECT().
		QuerySecurityNamespaces(clients.Ctx, security.QuerySecurityNamespacesArgs{
			SecurityNamespaceId: &nsID,
		}).
		Return(&[]security.SecurityNamespaceDescription{
			{
				Actions: &[]security.ActionDefinition{
					{
						Bit:  converter.Int(8),
						Name: converter.String("ManageMembership"),
					},
				},
			},
		}, nil).
		Times(2)

	securityClient.
		EXPECT().
		QueryAccessControlLists(clients.Ctx, gomock.Any()).
		Return(&[]security.AccessControlList{}, nil).
		Times(2)

	resourceData := schema.TestResourceDataRaw(t, DataTeams().Schema, nil)
	err := dataTeamsRead(resourceData, clients)

	require.Nil(t, err)
	require.Zero(t, resourceData.Get("project_id"))

	data, ok := resourceData.GetOk("teams")
	require.True(t, ok)
	teamSet := data.([]interface{})
	require.NotNil(t, teamSet)
	require.Equal(t, len(teamList), len(teamSet))

	teamMap := make(map[string]map[string]interface{}, len(teamSet))
	for _, e := range teamSet {
		team := e.(map[string]interface{})
		teamMap[team["id"].(string)] = team
	}
	for _, item := range teamList {
		team, ok := teamMap[item.id.String()]
		require.True(t, ok)
		require.Equal(t, item.name, team["name"])
		require.Equal(t, item.description, team["description"])
		require.Equal(t, testProjectID.String(), team["project_id"])
	}
}
