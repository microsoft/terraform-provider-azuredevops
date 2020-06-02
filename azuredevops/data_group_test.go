// +build all core data_sources data_group
// +build !exclude_data_sources !exclude_data_group

package azuredevops

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/stretchr/testify/require"
)

// A helper type that is used in some of these tests to make initializing
// graph entities easier
type groupMeta struct {
	name       string
	descriptor string
}

// verifies that the translation for project_id to project_descriptor has proper error handling
func TestGroupDataSource_DoesNotSwallowProjectDescriptorLookupError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projectID := uuid.New()
	resourceData := createResourceData(t, projectID.String(), "group-name")

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &config.AggregatedClient{GraphClient: graphClient, Ctx: context.Background()}

	expectedArgs := graph.GetDescriptorArgs{StorageKey: &projectID}
	graphClient.
		EXPECT().
		GetDescriptor(clients.Ctx, expectedArgs).
		Return(nil, errors.New("GetDescriptor() Failed"))

	err := dataSourceGroupRead(resourceData, clients)
	require.Contains(t, err.Error(), "GetDescriptor() Failed")
}

// verifies that the group lookup functionality has proper error handling
func TestGroupDataSource_DoesNotSwallowListGroupError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projectID := uuid.New()
	resourceData := createResourceData(t, projectID.String(), "group-name")

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &config.AggregatedClient{GraphClient: graphClient, Ctx: context.Background()}

	expectedProjectDescriptorLookupArgs := graph.GetDescriptorArgs{StorageKey: &projectID}
	projectDescriptor := converter.String("descriptor")
	projectDescriptorResponse := graph.GraphDescriptorResult{Value: projectDescriptor}
	graphClient.
		EXPECT().
		GetDescriptor(clients.Ctx, expectedProjectDescriptorLookupArgs).
		Return(&projectDescriptorResponse, nil)

	expectedListGroupArgs := graph.ListGroupsArgs{ScopeDescriptor: projectDescriptor}
	graphClient.
		EXPECT().
		ListGroups(clients.Ctx, expectedListGroupArgs).
		Return(nil, errors.New("ListGroups() Failed"))

	err := dataSourceGroupRead(resourceData, clients)
	require.Contains(t, err.Error(), "ListGroups() Failed")
}

// verifies that the group lookup functionality will make multiple API calls using the continuation token
// returned from the `ListGroups` api, until the API no longer returns a token
func TestGroupDataSource_HandlesContinuationToken_And_SelectsCorrectGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	projectID := uuid.New()
	resourceData := createResourceData(t, projectID.String(), "name1")

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &config.AggregatedClient{GraphClient: graphClient, Ctx: context.Background()}

	expectedProjectDescriptorLookupArgs := graph.GetDescriptorArgs{StorageKey: &projectID}
	projectDescriptor := converter.String("descriptor")
	projectDescriptorResponse := graph.GraphDescriptorResult{Value: projectDescriptor}
	graphClient.
		EXPECT().
		GetDescriptor(clients.Ctx, expectedProjectDescriptorLookupArgs).
		Return(&projectDescriptorResponse, nil)

	firstListGroupCallArgs := graph.ListGroupsArgs{ScopeDescriptor: projectDescriptor}
	continuationToken := "continuation-token"
	firstListGroupCallResponse := createPaginatedResponse(continuationToken, groupMeta{name: "name1", descriptor: "descriptor1"})
	firstCall := graphClient.
		EXPECT().
		ListGroups(clients.Ctx, firstListGroupCallArgs).
		Return(firstListGroupCallResponse, nil)

	secondListGroupCallArgs := graph.ListGroupsArgs{ScopeDescriptor: projectDescriptor, ContinuationToken: &continuationToken}
	secondListGroupCallResponse := createPaginatedResponse("", groupMeta{name: "name2", descriptor: "descriptor2"})
	secondCall := graphClient.
		EXPECT().
		ListGroups(clients.Ctx, secondListGroupCallArgs).
		Return(secondListGroupCallResponse, nil)

	gomock.InOrder(firstCall, secondCall)

	err := dataSourceGroupRead(resourceData, clients)
	require.Nil(t, err)
	require.Equal(t, "descriptor1", resourceData.Id())
}

func createPaginatedResponse(continuationToken string, groups ...groupMeta) *graph.PagedGraphGroups {
	continuationTokenList := []string{continuationToken}
	return &graph.PagedGraphGroups{
		ContinuationToken: &continuationTokenList,
		GraphGroups:       createGroupsWithDescriptors(groups...),
	}
}

func createGroupsWithDescriptors(groups ...groupMeta) *[]graph.GraphGroup {
	var graphs []graph.GraphGroup
	for _, group := range groups {
		graphs = append(graphs, graph.GraphGroup{Descriptor: &group.descriptor, DisplayName: &group.name})
	}

	return &graphs
}

func createResourceData(t *testing.T, projectID string, groupName string) *schema.ResourceData {
	resourceData := schema.TestResourceDataRaw(t, dataGroup().Schema, nil)
	resourceData.Set("project_id", projectID)
	resourceData.Set("name", groupName)
	return resourceData
}
