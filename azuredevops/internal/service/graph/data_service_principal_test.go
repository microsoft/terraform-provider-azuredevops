//go:build (all || core || data_sources || data_service_principal) && (!exclude_data_sources || !exclude_service_principal)
// +build all core data_sources data_service_principal
// +build !exclude_data_sources !exclude_service_principal

package graph

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/identity"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

func TestServicePrincipalNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	searchFilter := "General"
	displayName := "sp-test"

	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{
		IdentityClient: identityClient,
		Ctx:            context.Background(),
	}

	// Set up the mock expectations for ReadIdentities
	setUpMockReadIdentities(identityClient, clients.Ctx, searchFilter, displayName, nil, errors.New("Service principal not found"))

	// Set up the resource data with input values
	resourceData := schema.TestResourceDataRaw(t, DataServicePrincipal().Schema, nil)
	resourceData.Set("display_name", displayName)

	// Execute the function and check for the expected error
	err := dataSourceServicePrincipalRead(resourceData, clients)
	require.Contains(t, err.Error(), " Finding service principal with filter")
}

// verifies that the translation for display_name to descriptor has proper error handling
func TestServicePrincipalDataSource_DoesNotSwallowServicePrincipalDescriptorLookupError_Generic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	displayName := "sp-test"
	resourceData := createServicePrincipalDataSource(t, displayName)

	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{IdentityClient: identityClient, Ctx: context.Background()}

	expectedArgs := identity.ReadIdentitiesArgs{
		SearchFilter: converter.String("General"),
		FilterValue:  &displayName,
	}
	identityClient.
		EXPECT().
		ReadIdentities(clients.Ctx, expectedArgs).
		Return(nil, errors.New("ReadIdentities() Failed"))

	err := dataSourceServicePrincipalRead(resourceData, clients)
	require.Contains(t, err.Error(), "ReadIdentities() Failed")
}

// verifies that the translation for origin_id to descriptor has proper error handling
func TestServicePrincipalDataSource_DoesNotSwallowServicePrincipalLookupError_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	servicePrincipalDescriptor := uuid.New().String()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &client.AggregatedClient{GraphClient: graphClient, Ctx: context.Background()}

	expectedArgs := graph.GetServicePrincipalArgs{ServicePrincipalDescriptor: &servicePrincipalDescriptor}
	graphClient.
		EXPECT().
		GetServicePrincipal(clients.Ctx, expectedArgs).
		Return(nil, azuredevops.WrappedError{
			StatusCode: converter.Int(404),
		})

	servicePrincipal, err := getServicePrincipal(clients, &servicePrincipalDescriptor)
	require.Contains(t, err.Error(), "404")
	require.Same(t, servicePrincipal, (*graph.GraphServicePrincipal)(nil))
}

func createServicePrincipalDataSource(t *testing.T, displayName string) *schema.ResourceData {
	resourceData := schema.TestResourceDataRaw(t, DataServicePrincipal().Schema, nil)
	if displayName != "" {
		resourceData.Set("display_name", displayName)
	}
	return resourceData
}

func setUpMockReadIdentities(identityClient *azdosdkmocks.MockIdentityClient, ctx context.Context, searchFilter string, filterValue string, identities *[]identity.Identity, err error) {
	expectedArgs := identity.ReadIdentitiesArgs{
		SearchFilter: &searchFilter,
		FilterValue:  &filterValue,
	}
	identityClient.EXPECT().ReadIdentities(ctx, expectedArgs).Return(identities, err)
}
