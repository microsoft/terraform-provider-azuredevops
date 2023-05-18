//go:build (all || core || resource_team_administrators) && !exclude_resource_team_administrators
// +build all core resource_team_administrators
// +build !exclude_resource_team_administrators

package core

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

func TestTeamAdministrators_Create_DontSwallowError(t *testing.T) {
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
	testTeamID := uuid.New()
	errMsg := "@@GetTeam@@failed@@"

	coreClient.
		EXPECT().
		GetTeam(clients.Ctx, gomock.Any()).
		Return(nil, fmt.Errorf(errMsg)).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceTeamAdministrators().Schema, nil)
	resourceData.Set("project_id", testProjectID.String())
	resourceData.Set("team_id", testTeamID.String())
	err := resourceTeamAdministratorsCreate(resourceData, clients)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), errMsg)
}

func TestTeamAdministrators_Read_DontSwallowError(t *testing.T) {
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
	testTeamID := uuid.New()
	errMsg := "@@GetTeam@@failed@@"

	coreClient.
		EXPECT().
		GetTeam(clients.Ctx, gomock.Any()).
		Return(nil, fmt.Errorf(errMsg)).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceTeamAdministrators().Schema, nil)
	resourceData.Set("project_id", testProjectID.String())
	resourceData.Set("team_id", testTeamID.String())
	err := resourceTeamAdministratorsRead(resourceData, clients)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), errMsg)
}

func TestTeamAdministrators_Read_HandleMissingTeamCorrectly(t *testing.T) {
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
	testTeamID := uuid.New()

	coreClient.
		EXPECT().
		GetTeam(clients.Ctx, gomock.Any()).
		Return(nil, azuredevops.WrappedError{
			StatusCode: converter.Int(http.StatusNotFound),
		}).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceTeamAdministrators().Schema, nil)
	resourceData.Set("project_id", testProjectID.String())
	resourceData.Set("team_id", testTeamID.String())
	err := resourceTeamAdministratorsRead(resourceData, clients)

	require.Nil(t, err)
}

func TestTeamAdministrators_Delete_DontSwallowError(t *testing.T) {
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
	testTeamID := uuid.New()
	errMsg := "@@GetTeam@@failed@@"

	coreClient.
		EXPECT().
		GetTeam(clients.Ctx, gomock.Any()).
		Return(nil, fmt.Errorf(errMsg)).
		Times(1)

	resourceData := schema.TestResourceDataRaw(t, ResourceTeamAdministrators().Schema, nil)
	resourceData.Set("project_id", testProjectID.String())
	resourceData.Set("team_id", testTeamID.String())
	err := resourceTeamAdministratorsDelete(resourceData, clients)

	require.NotNil(t, err)
	require.Contains(t, err.Error(), errMsg)
}
