//go:build (all || core || data_sources || data_users) && (!exclude_data_sources || !exclude_data_users)
// +build all core data_sources data_users
// +build !exclude_data_sources !exclude_data_users

package identity

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/identity"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userName := "SomeUser"
	searchFilter := "General"

	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{
		IdentityClient: identityClient,
		Ctx:            context.Background(),
	}

	// Set up the mock expectations for ReadIdentities
	setUpMockReadIdentities(identityClient, clients.Ctx, userName, searchFilter, nil, errors.New("User not found"))

	// Set up the resource data with input values
	resourceData := schema.TestResourceDataRaw(t, DataIdentityUser().Schema, nil)
	resourceData.Set("name", userName)
	resourceData.Set("search_filter", searchFilter)

	// Execute the function and check for the expected error
	err := dataIdentitySourceUserRead(resourceData, clients)
	require.Contains(t, err.Error(), " Finding user with filter")
}

func TestErrorNotSwallowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userName := "SomeUser"
	searchFilter := "General"

	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{
		IdentityClient: identityClient,
		Ctx:            context.Background(),
	}

	// Set up the mock expectations for ReadIdentities
	setUpMockReadIdentities(identityClient, clients.Ctx, userName, searchFilter, nil, errors.New("Some other error"))

	// Set up the resource data with input values
	resourceData := schema.TestResourceDataRaw(t, DataIdentityUser().Schema, nil)
	resourceData.Set("name", userName)

	// Execute the function and check for the expected error
	err := dataIdentitySourceUserRead(resourceData, clients)
	require.Error(t, err)
	require.Contains(t, err.Error(), " Finding user with filter")
	require.Contains(t, err.Error(), "with filter "+searchFilter)
}

func setUpMockReadIdentities(identityClient *azdosdkmocks.MockIdentityClient, ctx context.Context, userName, searchFilter string, identities *[]identity.Identity, err error) {
	expectedArgs := identity.ReadIdentitiesArgs{
		FilterValue:  &userName,
		SearchFilter: &searchFilter,
	}
	identityClient.EXPECT().ReadIdentities(ctx, expectedArgs).Return(identities, err)
}
