//go:build (all || core || resource_group_membership) && !exclude_resource_group_membership

package acceptancetests

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/graph"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

func TestAccGroupMembership_createAndRemove(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	userPrincipalName := testutils.GenerateResourceName() + "@msaztest.com"
	tfNode := "azuredevops_group_membership.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclMemberShipBasic(projectName, userPrincipalName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrSet(tfNode, "group"),
					resource.TestCheckResourceAttr(tfNode, "members.#", "1"),
				),
			}, {
				Config: hclMemberShipRemove(projectName, userPrincipalName),
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
		_, expectingMembership := s.RootModule().Resources["azuredevops_group_membership.test"]

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

func hclMemberShipTemplate(projectName, userPrincipalName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}
data "azuredevops_group" "test" {
  project_id = azuredevops_project.test.id
  name       = "Build Administrators"
}
resource "azuredevops_user_entitlement" "test" {
  principal_name       = "%s"
  account_license_type = "express"
}`, projectName, userPrincipalName)
}

func hclMemberShipBasic(projectName, userPrincipalName string) string {
	return fmt.Sprintf(`
%s
output "group_descriptor" {
  value = data.azuredevops_group.test.descriptor
}
output "user_descriptor" {
  value = azuredevops_user_entitlement.test.descriptor
}
resource "azuredevops_group_membership" "test" {
  group   = data.azuredevops_group.test.descriptor
  members = [azuredevops_user_entitlement.test.descriptor]
}`, hclMemberShipTemplate(projectName, userPrincipalName))
}

func hclMemberShipRemove(projectName, userPrincipalName string) string {
	return fmt.Sprintf(`
%s
output "group_descriptor" {
  value = data.azuredevops_group.test.descriptor
}
output "user_descriptor" {
  value = azuredevops_user_entitlement.test.descriptor
}`, hclMemberShipTemplate(projectName, userPrincipalName))
}
