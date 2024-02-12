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

type groupMeta struct {
	name       string
	descriptor string
	domain     string
	origin     string
	originId   string
}

func TestIdentityGroupDataSource_ProjectDescriptorLookupError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projectID := uuid.NewString()
	resourceData := createIdentityGroupDataSource(t, projectID, "group-name")

	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{IdentityClient: identityClient, Ctx: context.Background()}

	expectedArgs := identity.ListGroupsArgs{ScopeIds: &projectID}
	identityClient.
		EXPECT().
		ListGroups(clients.Ctx, expectedArgs).
		Return(nil, errors.New("ListGroups() Failed"))

	err := dataSourceIdentityGroupRead(resourceData, clients)
	require.Contains(t, err.Error(), "ListGroups() Failed")
}

func TestIdentityGroupDataSource_ProjectDescriptorLookupErrorNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projectID := uuid.NewString()
	resourceData := createIdentityGroupDataSource(t, projectID, "group-name")

	identityClient := azdosdkmocks.NewMockIdentityClient(ctrl)
	clients := &client.AggregatedClient{IdentityClient: identityClient, Ctx: context.Background()}

	expectedArgs := identity.ListGroupsArgs{ScopeIds: &projectID}
	identityClient.
		EXPECT().
		ListGroups(clients.Ctx, expectedArgs).
		Return(nil, azuredevops.WrappedError{
			StatusCode: converter.Int(404),
		})

	err := dataSourceIdentityGroupRead(resourceData, clients)
	require.Contains(t, err.Error(), "Error finding groups")
}

func createIdentityGroupDataSource(t *testing.T, projectID string, groupName string) *schema.ResourceData {
	resourceData := schema.TestResourceDataRaw(t, DataIdentityGroup().Schema, nil)
	resourceData.Set("name", groupName)
	if projectID != "" {
		resourceData.Set("project_id", projectID)
	}
	return resourceData
}
