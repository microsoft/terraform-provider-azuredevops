package azuredevops

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

/**
 * Begin unit tests
 */

func TestGroupMembership_ComputeMembershipDiff_ResolvesDiffProperly(t *testing.T) {
	// If you are curious about the use of map here, have a read through this article:
	//	https://stackoverflow.com/questions/34018908/golang-why-dont-we-have-a-set-datastructure
	type membershipsTestMeta struct {
		old      map[string]bool
		new      map[string]bool
		toAdd    map[string]bool
		toRemove map[string]bool
	}
	// table of tests for computing membership diffs
	tests := []membershipsTestMeta{{
		// add single member
		old:      toStringSet("A"),
		new:      toStringSet("A", "B"),
		toAdd:    toStringSet("B"),
		toRemove: toStringSet(),
	}, {
		// remove single member
		old:      toStringSet("A", "B"),
		new:      toStringSet("A"),
		toAdd:    toStringSet(),
		toRemove: toStringSet("B"),
	}, {
		// add and remove members
		old:      toStringSet("A", "B", "C", "D"),
		new:      toStringSet("A", "B", "E", "F"),
		toAdd:    toStringSet("E", "F"),
		toRemove: toStringSet("C", "D"),
	}, {
		// no change to members
		old:      toStringSet("A"),
		new:      toStringSet("A"),
		toAdd:    toStringSet(),
		toRemove: toStringSet(),
	}}

	for _, test := range tests {
		toAdd, toRemove := computeMembershipDiff("", test.old, test.new)
		require.Equal(t, len(test.toAdd), len(*toAdd))
		require.Equal(t, len(test.toRemove), len(*toRemove))

		for _, membership := range *toAdd {
			if _, exists := test.toAdd[*membership.MemberDescriptor]; !exists {
				require.Fail(t, fmt.Sprintf("%s was unexpectedly not in the list of membershps to add!", *membership.MemberDescriptor))
			}
		}

		for _, membership := range *toRemove {
			if _, exists := test.toRemove[*membership.MemberDescriptor]; !exists {
				require.Fail(t, fmt.Sprintf("%s was unexpectedly not in the list of membershps to remove!", *membership.MemberDescriptor))
			}
		}
	}
}

func getGroupMembershipResourceData(t *testing.T, group string, members ...string) *schema.ResourceData {
	d := schema.TestResourceDataRaw(t, resourceGroupMembership().Schema, nil)
	d.Set("group", group)
	d.Set("members", members)
	return d
}

func TestGroupMembership_Create_DoesNotSwallowErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &aggregatedClient{GraphClient: graphClient, ctx: context.Background()}

	expectedArgs := graph.AddMembershipArgs{
		ContainerDescriptor: converter.String("TEST_GROUP"),
		SubjectDescriptor:   converter.String("TEST_MEMBER_1"),
	}
	graphClient.
		EXPECT().
		AddMembership(clients.ctx, expectedArgs).
		Return(nil, errors.New("AddMembership() Failed"))

	resourceData := getGroupMembershipResourceData(t, "TEST_GROUP", "TEST_MEMBER_1")
	err := resourceGroupMembershipCreate(resourceData, clients)
	require.Contains(t, err.Error(), "AddMembership() Failed")
}

func TestGroupMembership_Destroy_DoesNotSwallowErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &aggregatedClient{GraphClient: graphClient, ctx: context.Background()}

	expectedArgs := graph.RemoveMembershipArgs{
		ContainerDescriptor: converter.String("TEST_GROUP"),
		SubjectDescriptor:   converter.String("TEST_MEMBER_1"),
	}
	graphClient.
		EXPECT().
		RemoveMembership(clients.ctx, expectedArgs).
		Return(errors.New("RemoveMembership() Failed"))

	resourceData := getGroupMembershipResourceData(t, "TEST_GROUP", "TEST_MEMBER_1")
	err := resourceGroupMembershipDelete(resourceData, clients)
	require.Contains(t, err.Error(), "RemoveMembership() Failed")
}

func TestGroupMembership_Read_DoesNotSwallowErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	graphClient := azdosdkmocks.NewMockGraphClient(ctrl)
	clients := &aggregatedClient{GraphClient: graphClient, ctx: context.Background()}

	expectedArgs := graph.ListMembershipsArgs{
		SubjectDescriptor: converter.String("TEST_GROUP"),
		Direction:         &graph.GraphTraversalDirectionValues.Down,
		Depth:             converter.Int(1),
	}
	graphClient.
		EXPECT().
		ListMemberships(clients.ctx, expectedArgs).
		Return(nil, errors.New("ListMemberships() Failed"))

	resourceData := getGroupMembershipResourceData(t, "TEST_GROUP", "TEST_MEMBER_1")
	err := resourceGroupMembershipRead(resourceData, clients)
	require.Contains(t, err.Error(), "ListMemberships() Failed")
}

/**
 * Begin acceptance tests
 */

// Verifies that the following sequence of events occurrs without error:
//	(1) TF apply creates resource
//	(2) TF state values are set
//	(3) Group membership exists and can be queried for
// 	(4) TF destroy removes group memberships
//
// Note: This will be uncommented in https://github.com/microsoft/terraform-provider-azuredevops/issues/174
//
// func TestAccGroupMembership_CreateAndRemove(t *testing.T) {
// 	projectName := testAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
// 	userPrincipalName := os.Getenv("AZDO_TEST_AAD_USER_EMAIL")
// 	groupName := "Build Administrators"
// 	tfNode := "azuredevops_group_membership.membership"

// 	tfStanzaWithMembership := testAccGroupMembershipResource(projectName, groupName, userPrincipalName)
// 	tfStanzaWithoutMembership := testAccGroupMembershipDependencies(projectName, groupName, userPrincipalName)

// 	// This test differs from most other acceptance tests in the following ways:
// 	//	- The second step is the same as the first except it omits the group membershp.
// 	//	  This lets us test that the membership is removed in isolation of the project being deleted
// 	//	- There is no CheckDestroy function because that is covered based on the above point
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:  func() { testAccPreCheck(t) },
// 		Providers: testAccProviders,
// 		Steps: []resource.TestStep{
// 			{
// 				// add the group membership
// 				Config: tfStanzaWithMembership,
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttrSet(tfNode, "id"),
// 					resource.TestCheckResourceAttrSet(tfNode, "group"),
// 					// this attribute specifies the number of members in the resource state. the
// 					// syntax is how terraform maps complex types into a flattened map.
// 					resource.TestCheckResourceAttr(tfNode, "members.#", "1"),
// 					testAccVerifyGroupMembershipMatchesState(),
// 				),
// 			}, {
// 				// remove the group membership
// 				Config: tfStanzaWithoutMembership,
// 				Check:  testAccVerifyGroupMembershipMatchesState(),
// 			},
// 		},
// 	})
// }

// Verifies that the group membership in AzDO matches the group membership specified by the state
func testAccVerifyGroupMembershipMatchesState() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		memberDescriptor := s.RootModule().Outputs["user_descriptor"].Value.(string)
		groupDescriptor := s.RootModule().Outputs["group_descriptor"].Value.(string)
		_, expectingMembership := s.RootModule().Resources["azuredevops_group_membership.membership"]

		// The sleep here is to take into account some propegation delay that can happen with Group Membership APIs.
		// If we want to go inspect the behavior of the service after a Terraform Apply, we'll need to wait a little bit
		// before making the API call.
		//
		// Note: some thought was put behind keeping the time.sleep here vs in the provider implementation. After consideration,
		// I decided to keep it here. Moving to the provider would (1) provide no functional benefit to the end user, (2) increase
		// complexity and (3) be inconsistent with the UI and CLI behavior for the same operation.
		time.Sleep(5 * time.Second)
		memberships, err := getMembersOfGroup(groupDescriptor)
		if err != nil {
			return err
		}

		if !expectingMembership && len(*memberships) == 0 {
			return nil
		}

		if !expectingMembership && len(*memberships) > 0 {
			return fmt.Errorf("unexpectedly found group members: %+v", memberships)
		}

		if expectingMembership && len(*memberships) == 0 {
			return fmt.Errorf("unexpectedly did not find memberships")
		}

		actualMemberDescriptor := *(*memberships)[0].MemberDescriptor
		if strings.ToLower(actualMemberDescriptor) != strings.ToLower(memberDescriptor) {
			return fmt.Errorf("expected member with descriptor %s but member had descriptor %s", memberDescriptor, actualMemberDescriptor)
		}

		return nil
	}

}

// call AzDO API to query for group members
func getMembersOfGroup(groupDescriptor string) (*[]graph.GraphMembership, error) {
	clients := testAccProvider.Meta().(*aggregatedClient)
	return clients.GraphClient.ListMemberships(clients.ctx, graph.ListMembershipsArgs{
		SubjectDescriptor: &groupDescriptor,
		Direction:         &graph.GraphTraversalDirectionValues.Down,
		Depth:             converter.Int(1),
	})
}

// full terraform stanza to standup a group membership
func testAccGroupMembershipResource(projectName, groupName, userPrincipalName string) string {
	membershipDependenciesStanza := testAccGroupMembershipDependencies(projectName, groupName, userPrincipalName)
	membershipStanza := `
resource "azuredevops_group_membership" "membership" {
	group = data.azuredevops_group.group.descriptor
	members = [azuredevops_user_entitlement.user.descriptor]
}`

	return membershipDependenciesStanza + "\n" + membershipStanza
}

// all the dependencies needed to configure a group membership
func testAccGroupMembershipDependencies(projectName, groupName, userPrincipalName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "project" {
	project_name = "%s"
}
data "azuredevops_group" "group" {
	project_id = azuredevops_project.project.id
	name       = "%s"
}
resource "azuredevops_user_entitlement" "user" {
	principal_name       = "%s"
	account_license_type = "express"
}

output "group_descriptor" {
	value = data.azuredevops_group.group.descriptor
}
output "user_descriptor" {
	value = azuredevops_user_entitlement.user.descriptor
}
`, projectName, groupName, userPrincipalName)
}
