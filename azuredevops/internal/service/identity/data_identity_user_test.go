//go:build (all || core || data_sources || data_users) && (!exclude_data_sources || !exclude_data_users)
// +build all core data_sources data_users
// +build !exclude_data_sources !exclude_data_users

package identity

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"errors"
	"fmt"
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

	searchFilter := "General"
	userName := "NonExistentUser"

	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{IdentityClient: identityClient, Ctx: context.Background()}
	// Set up the mock expectations for ReadIdentities
	expectedArgs := identity.ReadIdentitiesArgs{FilterValue: &userName, SearchFilter: &searchFilter}
	identityClient.
		EXPECT().
		ReadIdentities(clients.Ctx, expectedArgs).
		Return(nil, errors.New("User not found"))

	// Set up the resource data with input values
	resourceData := createUserResourceData(t, "group-name")
	// Execute the function and check for the expected error
	err := dataIdentitySourceUserRead(resourceData, clients)
	require.Error(t, err)
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

	// Set up the mock expectations for ReadIdentities
	expectedArgs := []identity.ReadIdentitiesArgs{{FilterValue: &userName, SearchFilter: &searchFilter}}
	// Set up the mock expectations for ReadIdentities
	identityClient.
		EXPECT().
		ReadIdentities(clients.Ctx, expectedArgs).
		Return(nil, errors.New("Some other error"))

	// Set up the resource data with input values
	resourceData := createUserResourceData(t, "group-name")
	// Execute the function and check for the expected error
	err := dataIdentitySourceUserRead(resourceData, clients)
	require.Error(t, err)
	require.Contains(t, err.Error(), "Error finding user")
	require.Contains(t, err.Error(), "with filter "+searchFilter)
}

// Helper function to simulate the behavior of Read method
func testReadFunction(d *schema.ResourceData, m interface{}) error {
	// Convert interface{} to *client.AggregatedClient
	clients := m.(*client.AggregatedClient)

	// Call the actual dataIdentitySourceUserRead function
	//return dataIdentitySourceUserRead(d, clients)
	if err := dataIdentitySourceUserRead(d, clients); err != nil {
		return nil
	}
	id, idExists := d.GetOk("descriptor")
	if !idExists {
		return fmt.Errorf("id field not set in ResourceData")
	}
	d.SetId(id.(string))
	return nil
}

func createUserResourceData(t *testing.T, groupName string) *schema.ResourceData {
	resourceData := schema.TestResourceDataRaw(t, DataIdentityUser().Schema, nil)
	resourceData.Set("name", groupName)
	return resourceData
}
