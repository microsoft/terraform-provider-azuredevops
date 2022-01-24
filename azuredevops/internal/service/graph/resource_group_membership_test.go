//go:build (all || core || resource_group_membership) && !exclude_resource_group_membership
// +build all core resource_group_membership
// +build !exclude_resource_group_membership

package graph

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

func TestGroupMembership_Create_DoesNotSwallowErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &client.AggregatedClient{GraphClient: graphClient, Ctx: context.Background()}

	expectedArgs := graph.AddMembershipArgs{
		ContainerDescriptor: converter.String("TEST_GROUP"),
		SubjectDescriptor:   converter.String("TEST_MEMBER_1"),
	}
	graphClient.
		EXPECT().
		AddMembership(clients.Ctx, expectedArgs).
		Return(nil, errors.New("AddMembership() Failed"))

	resourceData := getGroupMembershipResourceData(t, "TEST_GROUP", "TEST_MEMBER_1")
	err := resourceGroupMembershipCreate(resourceData, clients)
	require.Contains(t, err.Error(), "AddMembership() Failed")
}

func TestGroupMembership_Destroy_DoesNotSwallowErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &client.AggregatedClient{GraphClient: graphClient, Ctx: context.Background()}

	expectedArgs := graph.RemoveMembershipArgs{
		ContainerDescriptor: converter.String("TEST_GROUP"),
		SubjectDescriptor:   converter.String("TEST_MEMBER_1"),
	}
	graphClient.
		EXPECT().
		RemoveMembership(clients.Ctx, expectedArgs).
		Return(errors.New("RemoveMembership() Failed"))

	resourceData := getGroupMembershipResourceData(t, "TEST_GROUP", "TEST_MEMBER_1")
	err := resourceGroupMembershipDelete(resourceData, clients)
	require.Contains(t, err.Error(), "RemoveMembership() Failed")
}

func TestGroupMembership_Read_DoesNotSwallowErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &client.AggregatedClient{GraphClient: graphClient, Ctx: context.Background()}

	expectedArgs := graph.ListMembershipsArgs{
		SubjectDescriptor: converter.String("TEST_GROUP"),
		Direction:         &graph.GraphTraversalDirectionValues.Down,
		Depth:             converter.Int(1),
	}
	graphClient.
		EXPECT().
		ListMemberships(clients.Ctx, expectedArgs).
		Return(nil, errors.New("ListMemberships() Failed"))

	resourceData := getGroupMembershipResourceData(t, "TEST_GROUP", "TEST_MEMBER_1")
	err := resourceGroupMembershipRead(resourceData, clients)
	require.Contains(t, err.Error(), "ListMemberships() Failed")
}

func getGroupMembershipResourceData(t *testing.T, group string, members ...string) *schema.ResourceData {
	d := schema.TestResourceDataRaw(t, ResourceGroupMembership().Schema, nil)
	d.Set("group", group)
	d.Set("members", members)
	return d
}
