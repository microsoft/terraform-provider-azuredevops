//go:build (all || core || resource_group_membership) && !exclude_resource_group_membership
// +build all core resource_group_membership
// +build !exclude_resource_group_membership

package acceptancetests

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v6/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

// Verifies that the following sequence of events occurs without error:
//
//	(1) TF apply creates resource
//	(2) TF state values are set
//	(3) Group membership exists and can be queried for
//	(4) TF destroy removes group memberships
//
// Note: This will be uncommented in https://github.com/microsoft/terraform-provider-azuredevops/issues/174
func TestAccGroupMembership_CreateAndRemove(t *testing.T) {
	t.Skip("Skipping test TestAccGroupMembership_CreateAndRemove due to service inconsistent")
	projectName := testutils.GenerateResourceName()
	userPrincipalName := os.Getenv("AZDO_TEST_AAD_USER_EMAIL")
	groupName := "Build Administrators"
	tfNode := "azuredevops_group_membership.membership"

	tfStanzaWithMembership := testutils.HclGroupMembershipResource(projectName, groupName, userPrincipalName)
	tfStanzaWithoutMembership := testutils.HclGroupMembershipDependencies(projectName, groupName, userPrincipalName)

	// This test differs from most other acceptance tests in the following ways:
	//	- The second step is the same as the first except it omits the group membership.
	//	  This lets us test that the membership is removed in isolation of the project being deleted
	//	- There is no CheckDestroy function because that is covered based on the above point
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: tfStanzaWithMembership,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrSet(tfNode, "group"),
					resource.TestCheckResourceAttr(tfNode, "members.#", "1"),
					checkGroupMembershipMatchesState(),
				),
			}, {
				// remove the group membership
				Config: tfStanzaWithoutMembership,
				Check:  checkGroupMembershipMatchesState(),
			},
		},
	})
}

// Verifies that the group membership in AzDO matches the group membership specified by the state
func checkGroupMembershipMatchesState() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		memberDescriptor := s.RootModule().Outputs["user_descriptor"].Value.(string)
		groupDescriptor := s.RootModule().Outputs["group_descriptor"].Value.(string)
		_, expectingMembership := s.RootModule().Resources["azuredevops_group_membership.membership"]

		// The sleep here is to take into account some propagation delay that can happen with Group Membership APIs.
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
		if !strings.EqualFold(strings.ToLower(actualMemberDescriptor), strings.ToLower(memberDescriptor)) {
			return fmt.Errorf("expected member with descriptor %s but member had descriptor %s", memberDescriptor, actualMemberDescriptor)
		}

		return nil
	}
}

// call AzDO API to query for group members
func getMembersOfGroup(groupDescriptor string) (*[]graph.GraphMembership, error) {
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
	return clients.GraphClient.ListMemberships(clients.Ctx, graph.ListMembershipsArgs{
		SubjectDescriptor: &groupDescriptor,
		Direction:         &graph.GraphTraversalDirectionValues.Down,
		Depth:             converter.Int(1),
	})
}
