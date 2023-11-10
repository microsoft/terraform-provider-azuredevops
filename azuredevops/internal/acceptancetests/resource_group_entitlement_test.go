//go:build (all || resource_group_entitlement) && !exclude_resource_group_entitlement
// +build all resource_group_entitlement
// +build !exclude_resource_group_entitlement

package acceptancetests

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/memberentitlementmanagement"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
)

func TestAccGroupEntitlement_Create(t *testing.T) {
	tfNode := "azuredevops_group_entitlement.group"
	displayName := "group-038c153d-c86e-443c-b6f6-3d97378025d0"
	resource.ParallelTest(t, resource.TestCase{
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkGroupEntitlementDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclGroupEntitlementResource(displayName),
				Check: resource.ComposeTestCheckFunc(
					checkGroupEntitlementExists(displayName),
					resource.TestCheckResourceAttrSet(tfNode, "descriptor"),
				),
			},
		},
	})
}

func TestAccGroupEntitlement_AAD_Create(t *testing.T) {
	tfNode := "azuredevops_group_entitlement.group_aad"
	originId := os.Getenv("AZDO_TEST_AAD_GROUP_ID")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, &[]string{"AZDO_TEST_AAD_GROUP_ID"}) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkGroupEntitlementDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclGroupEntitlementResourceAAD(originId),
				Check: resource.ComposeTestCheckFunc(
					checkGroupEntitlementAADExists(originId),
					resource.TestCheckResourceAttrSet(tfNode, "descriptor"),
				),
			},
		},
	})
}

// Given the principalName of an AzDO groupEntitlement, this will return a function that will check whether
// or not the groupEntitlement (1) exists in the state and (2) exist in AzDO and (3) has the correct name
func checkGroupEntitlementExists(expectedDisplayName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources["azuredevops_group_entitlement.group"]
		if !ok {
			return fmt.Errorf("Did not find a GroupEntitlement in the TF state")
		}

		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
		id, err := uuid.Parse(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing GroupEntitlement ID, got %s: %v", resource.Primary.ID, err)
		}

		groupEntitlement, err := clients.MemberEntitleManagementClient.GetGroupEntitlement(clients.Ctx, memberentitlementmanagement.GetGroupEntitlementArgs{
			GroupId: &id,
		})

		if err != nil {
			return fmt.Errorf("GroupEntitlement with ID=%s cannot be found!. Error=%v", id, err)
		}

		if !strings.EqualFold(strings.ToLower(*groupEntitlement.Group.DisplayName), strings.ToLower(expectedDisplayName)) {
			return fmt.Errorf("GroupEntitlement with ID=%s and principalName=%s has displayName=%s, but expected displayName=%s", resource.Primary.ID, *groupEntitlement.Group.PrincipalName, *groupEntitlement.Group.DisplayName, expectedDisplayName)
		}

		return nil
	}
}

func checkGroupEntitlementAADExists(expectedOriginId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources["azuredevops_group_entitlement.group_aad"]
		if !ok {
			return fmt.Errorf("Did not find a GroupEntitlement in the TF state")
		}

		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
		id, err := uuid.Parse(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing GroupEntitlement ID, got %s: %v", resource.Primary.ID, err)
		}

		groupEntitlement, err := clients.MemberEntitleManagementClient.GetGroupEntitlement(clients.Ctx, memberentitlementmanagement.GetGroupEntitlementArgs{
			GroupId: &id,
		})

		if err != nil {
			return fmt.Errorf("GroupEntitlement with ID=%s cannot be found!. Error=%v", id, err)
		}

		if !strings.EqualFold(strings.ToLower(*groupEntitlement.Group.OriginId), strings.ToLower(expectedOriginId)) {
			return fmt.Errorf("GroupEntitlement with ID=%s has originId=%s, but expected originId=%s", resource.Primary.ID, *groupEntitlement.Group.OriginId, expectedOriginId)
		}

		return nil
	}
}

// verifies that all projects referenced in the state are destroyed. This will be invoked
// *after* terraform destroys the resource but *before* the state is wiped clean.
func checkGroupEntitlementDestroyed(s *terraform.State) error {
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

	//verify that every users referenced in the state does not exist in AzDO
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_group_entitlement" {
			continue
		}

		id, err := uuid.Parse(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing GroupEntitlement ID, got %s: %v", resource.Primary.ID, err)
		}

		groupEntitlement, err := clients.MemberEntitleManagementClient.GetGroupEntitlement(clients.Ctx, memberentitlementmanagement.GetGroupEntitlementArgs{
			GroupId: &id,
		})

		if err != nil {
			if utils.ResponseWasNotFound(err) {
				return nil
			}
			return fmt.Errorf("Bad: Get GroupEntitlment :  %+v", err)
		}

		if groupEntitlement != nil && groupEntitlement.LicenseRule != nil && string(*groupEntitlement.LicenseRule.Status) != "none" {
			return fmt.Errorf("Status should be none : %s with readGroupEntitlement error %v", string(*groupEntitlement.LicenseRule.Status), err)
		}
	}

	return nil
}

type matchAddGroupEntitlementArgs struct {
	t *testing.T
	x memberentitlementmanagement.AddGroupEntitlementArgs
}

func MatchAddGroupEntitlementArgs(t *testing.T, x memberentitlementmanagement.AddGroupEntitlementArgs) gomock.Matcher {
	return &matchAddGroupEntitlementArgs{t, x}
}

func (m *matchAddGroupEntitlementArgs) Matches(x interface{}) bool {
	args := x.(memberentitlementmanagement.AddGroupEntitlementArgs)
	m.t.Logf("MatchAddGroupEntitlementArgs:\nVALUE: account_license_type: [%s], licensing_source: [%s], origin: [%s], origin_id: [%s], display_name: [%s], principal_name: [%s]\n  REF: account_license_type: [%s], licensing_source: [%s], origin: [%s], origin_id: [%s], display_name: [%s], principal_name: [%s]\n",
		*args.GroupEntitlement.LicenseRule.AccountLicenseType,
		*args.GroupEntitlement.LicenseRule.LicensingSource,
		*args.GroupEntitlement.Group.Origin,
		*args.GroupEntitlement.Group.OriginId,
		*args.GroupEntitlement.Group.DisplayName,
		*args.GroupEntitlement.Group.PrincipalName,
		*m.x.GroupEntitlement.LicenseRule.AccountLicenseType,
		*m.x.GroupEntitlement.LicenseRule.LicensingSource,
		*m.x.GroupEntitlement.Group.Origin,
		*m.x.GroupEntitlement.Group.OriginId,
		*m.x.GroupEntitlement.Group.DisplayName,
		*m.x.GroupEntitlement.Group.PrincipalName)

	return *args.GroupEntitlement.LicenseRule.AccountLicenseType == *m.x.GroupEntitlement.LicenseRule.AccountLicenseType &&
		*args.GroupEntitlement.Group.Origin == *m.x.GroupEntitlement.Group.Origin &&
		*args.GroupEntitlement.Group.OriginId == *m.x.GroupEntitlement.Group.OriginId &&
		*args.GroupEntitlement.Group.PrincipalName == *m.x.GroupEntitlement.Group.PrincipalName
}

func (m *matchAddGroupEntitlementArgs) String() string {
	return fmt.Sprintf("account_license_type: [%s], licensing_source: [%s], origin: [%s], origin_id: [%s], display_name: [%s], principal_name: [%s]",
		*m.x.GroupEntitlement.LicenseRule.AccountLicenseType,
		*m.x.GroupEntitlement.LicenseRule.LicensingSource,
		*m.x.GroupEntitlement.Group.Origin,
		*m.x.GroupEntitlement.Group.OriginId,
		*m.x.GroupEntitlement.Group.DisplayName,
		*m.x.GroupEntitlement.Group.PrincipalName)
}
