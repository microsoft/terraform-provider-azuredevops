//go:build (all || core || data_sources || data_team) && (!exclude_data_sources || !exclude_data_team)
// +build all core data_sources data_team
// +build !exclude_data_sources !exclude_data_team

package core

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/core"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

func TestDataTeam_Read_DoesNotSwallowError(t *testing.T) {
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
	testTeamName := "@@TEST TEAM@@"

	coreClient.
		EXPECT().
		GetTeam(clients.Ctx, core.GetTeamArgs{
			ProjectId: converter.String(testProjectID.String()),
			TeamId:    converter.String(testTeamName),
		}).
		Return(&core.WebApiTeam{}, errors.New("@@GetTeams@@failed@@")).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataTeam().Schema, nil)
	resourceData.Set("project_id", testProjectID.String())
	resourceData.Set("name", testTeamName)
	err := dataTeamRead(clients.Ctx, resourceData, clients)

	require.NotNil(t, err)
	require.Contains(t, err[0].Summary, "@@GetTeams@@failed@@")

	require.Equal(t, testProjectID.String(), resourceData.Get("project_id"))
	require.Equal(t, testTeamName, resourceData.Get("name"))
}

func TestDataTeam_Read_FailOnNotFound(t *testing.T) {
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
	testTeamName := "@@TEST TEAM@@"

	coreClient.
		EXPECT().
		GetTeam(clients.Ctx, core.GetTeamArgs{
			ProjectId: converter.String(testProjectID.String()),
			TeamId:    &testTeamName,
		}).
		Return(nil, fmt.Errorf("Unable to find Team with name")).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, DataTeam().Schema, nil)
	resourceData.Set("project_id", testProjectID.String())
	resourceData.Set("name", testTeamName)
	err := dataTeamRead(clients.Ctx, resourceData, clients)

	require.NotNil(t, err)
	require.Contains(t, err[0].Summary, "Unable to find Team with name")
	require.Equal(t, testProjectID.String(), resourceData.Get("project_id"))
	require.Equal(t, testTeamName, resourceData.Get("name"))
	require.Zero(t, resourceData.Get("description"))
}
