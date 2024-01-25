//go:build (all || core || data_sources || data_group) && (!exclude_data_sources || !exclude_data_group)
// +build all core data_sources data_group
// +build !exclude_data_sources !exclude_data_group

package identity

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/identity"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

// A helper type that is used in some of these tests to make initializing
// identity entities easier
type groupMeta struct {
	name       string
	descriptor string
	domain     string
	origin     string
	originId   string
}

// verifies that the translation for project_id to project_descriptor has proper error handling
func TestIdentityGroupDataSource_DoesNotSwallowProjectDescriptorLookupError_Generic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projectID := uuid.New()
	projectIDstring := projectID.String()
	resourceData := createResourceData(t, projectID.String(), "group-name")

	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{IdentityClient: identityClient, Ctx: context.Background()}

	expectedArgs := identity.ListGroupsArgs{ScopeIds: &projectIDstring}
	identityClient.
		EXPECT().
		ListGroups(clients.Ctx, expectedArgs).
		Return(nil, errors.New("ListGroups() Failed"))

	err := dataSourceIdentityGroupRead(resourceData, clients)
	require.Contains(t, err.Error(), "ListGroups() Failed")
}

// verifies that the translation for project_id to project_descriptor has proper error handling
func TestIdentityGroupDataSource_DoesNotSwallowProjectDescriptorLookupError_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projectID := uuid.New()
	projectIDstring := projectID.String()
	resourceData := createResourceData(t, projectID.String(), "group-name")

	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{IdentityClient: identityClient, Ctx: context.Background()}

	expectedArgs := identity.ListGroupsArgs{ScopeIds: &projectIDstring}
	identityClient.
		EXPECT().
		ListGroups(clients.Ctx, expectedArgs).
		Return(nil, azuredevops.WrappedError{
			StatusCode: converter.Int(404),
		})

	err := dataSourceIdentityGroupRead(resourceData, clients)
	require.Contains(t, err.Error(), "was not found")
}

// verifies that the group lookup functionality has proper error handling
func TestIdentityGroupDataSource_DoesNotSwallowListGroupError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projectID := uuid.New()
	projectIDstring := projectID.String()
	resourceData := createResourceData(t, projectID.String(), "group-name")

	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{IdentityClient: identityClient, Ctx: context.Background()}

	expectedProjectDescriptorLookupArgs := identity.ListGroupsArgs{ScopeIds: &projectIDstring}
	projectDescriptor := converter.String("descriptor")
	projectDescriptorResponse := identity.Identity{Descriptor: projectDescriptor}
	identityClient.
		EXPECT().
		ListGroups(clients.Ctx, expectedProjectDescriptorLookupArgs).
		Return(&projectDescriptorResponse, nil)

	expectedListGroupArgs := identity.ListGroupsArgs{ScopeIds: projectDescriptor}
	identityClient.
		EXPECT().
		ListGroups(clients.Ctx, expectedListGroupArgs).
		Return(nil, errors.New("ListGroups() Failed"))

	err := dataSourceIdentityGroupRead(resourceData, clients)
	require.Contains(t, err.Error(), "ListGroups() Failed")
}

func createGroupsWithDescriptors(groups ...groupMeta) *[]identity.Identity {
	var identitys []identity.Identity
	for _, group := range groups {
		identitys = append(identitys, identity.Identity{
			Descriptor:          converter.String(group.descriptor),
			ProviderDisplayName: converter.String(group.name),
		})
	}

	return &identitys
}

func createResourceData(t *testing.T, projectID string, groupName string) *schema.ResourceData {
	resourceData := schema.TestResourceDataRaw(t, DataIdentityGroup().Schema, nil)
	resourceData.Set("name", groupName)
	if projectID != "" {
		resourceData.Set("project_id", projectID)
	}
	return resourceData
}
