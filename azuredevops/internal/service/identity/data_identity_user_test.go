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

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/identity"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/stretchr/testify/require"
)

// Test when user is not found
func TestDataSourceIdentityUser_UserNotFound(t *testing.T) {
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
	expectedArgs := identity.ReadIdentitiesArgs{FilterValue: &userName, SearchFilter: &searchFilter}
	identityClient.
		EXPECT().
		ReadIdentities(clients.Ctx, expectedArgs).
		Return(nil, errors.New("User not found"))
	// Set up the resource data with input values
	t.Log("after executing")
	resourceData := schema.TestResourceDataRaw(t, DataIdentityUser().Schema, nil)
	resourceData.Set("name", "SomeUser")
	resourceData.Set("search_filter", "General")
	t.Log("after executing")
	// Execute the function and check for the expected error
	err := dataIdentitySourceUserRead(resourceData, clients)
	require.Contains(t, err.Error(), "Could not find user with name "+userName)
}

// Test to validate that the error is not swallowed
func TestDataSourceIdentityUser_ErrorNotSwallowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userName := "SomeUser"
	searchFilter := "General"

	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{
		IdentityClient: identityClient,
		Ctx:            context.Background(),
	}
	// Set up the resource data with input values
	resourceData := schema.TestResourceDataRaw(t, DataIdentityUser().Schema, nil)
	resourceData.Set("name", userName)

	// Set up the mock expectations for ReadIdentities
	expectedArgs := identity.ReadIdentitiesArgs{FilterValue: &userName, SearchFilter: &searchFilter}
	// Set up the mock expectations for ReadIdentities
	identityClient.
		EXPECT().
		ReadIdentities(clients.Ctx, expectedArgs).
		Return(nil, errors.New("Some other error"))

	// Execute the function and check for the expected error
	err := dataIdentitySourceUserRead(resourceData, clients)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Error finding user")
	require.Contains(t, err.Error(), "with filter "+searchFilter)
}

//
