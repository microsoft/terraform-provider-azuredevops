package identity

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/identity"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

func TestIdentityGroupsDataSource_DoesNotSwallowProjectDescriptorLookupError_Generic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projectID := uuid.New()
	projectIDstring := projectID.String()
	resourceData := createIdentityGroupsDataSource(t, projectID.String())

	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{IdentityClient: identityClient, Ctx: context.Background()}

	expectedArgs := identity.ListGroupsArgs{ScopeIds: &projectIDstring}
	identityClient.
		EXPECT().
		ListGroups(clients.Ctx, expectedArgs).
		Return(nil, errors.New("GetDescriptor() Failed"))

	err := dataSourceIdentityGroupsRead(resourceData, clients)
	require.Contains(t, err.Error(), "GetDescriptor() Failed")
}

// verifies that the translation for project_id to project_descriptor has proper error handling
func TestIdentityGroupsDataSource_DoesNotSwallowProjectDescriptorLookupError_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projectID := uuid.New()
	projectIDstring := projectID.String()
	resourceData := createIdentityGroupsDataSource(t, projectID.String())

	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{IdentityClient: identityClient, Ctx: context.Background()}

	expectedArgs := identity.ListGroupsArgs{ScopeIds: &projectIDstring}
	identityClient.
		EXPECT().
		ListGroups(clients.Ctx, expectedArgs).
		Return(nil, azuredevops.WrappedError{
			StatusCode: converter.Int(404),
		})

	err := dataSourceIdentityGroupsRead(resourceData, clients)
	require.Contains(t, err.Error(), "was not found")
}

// verifies that the group lookup functionality has proper error handling
func TestIdentityGroupsDataSource_DoesNotSwallowListGroupError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projectID := uuid.New()
	projectIDstring := projectID.String()
	resourceData := createIdentityGroupsDataSource(t, projectID.String())

	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{IdentityClient: identityClient, Ctx: context.Background()}

	expectedProjectDescriptorLookupArgs := identity.ListGroupsArgs{ScopeIds: &projectIDstring}
	projectDescriptor := converter.String("descriptor")
	projectDescriptorResponse := identity.Identity{Descriptor: projectDescriptor}
	identityClient.
		EXPECT().
		GetDescriptorById(clients.Ctx, expectedProjectDescriptorLookupArgs).
		Return(&projectDescriptorResponse, nil)

	expectedListGroupArgs := identity.ListGroupsArgs{ScopeIds: projectDescriptor}
	identityClient.
		EXPECT().
		ListGroups(clients.Ctx, expectedListGroupArgs).
		Return(nil, errors.New("ListGroups() Failed"))

	err := dataSourceIdentityGroupsRead(resourceData, clients)
	require.Contains(t, err.Error(), "ListGroups() Failed")
}

func createIdentityGroupsDataSource(t *testing.T, projectID string) *schema.ResourceData {
	resourceData := schema.TestResourceDataRaw(t, DataIdentityGroups().Schema, nil)
	if projectID != "" {
		resourceData.Set("project_id", projectID)
	}
	return resourceData
}
