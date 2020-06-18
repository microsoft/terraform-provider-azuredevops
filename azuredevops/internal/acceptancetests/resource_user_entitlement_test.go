// +build all resource_user_entitlement
// +build !exclude_resource_user_entitlement

package acceptancetests

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/memberentitlementmanagement"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils"
)

func TestAccUserEntitlement_Create(t *testing.T) {
	tfNode := "azuredevops_user_entitlement.user"
	principalName := os.Getenv("AZDO_TEST_AAD_USER_EMAIL")
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, &[]string{"AZDO_TEST_AAD_USER_EMAIL"}) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkUserEntitlementDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclUserEntitlementResource(principalName),
				Check: resource.ComposeTestCheckFunc(
					checkUserEntitlementExists(principalName),
					resource.TestCheckResourceAttrSet(tfNode, "descriptor"),
				),
			},
		},
	})
}

// Given the principalName of an AzDO userEntitlement, this will return a function that will check whether
// or not the userEntitlement (1) exists in the state and (2) exist in AzDO and (3) has the correct name
func checkUserEntitlementExists(expectedPrincipalName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources["azuredevops_user_entitlement.user"]
		if !ok {
			return fmt.Errorf("Did not find a UserEntitlement in the TF state")
		}

		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
		id, err := uuid.Parse(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing UserEntitlement ID, got %s: %v", resource.Primary.ID, err)
		}

		userEntitlement, err := clients.MemberEntitleManagementClient.GetUserEntitlement(clients.Ctx, memberentitlementmanagement.GetUserEntitlementArgs{
			UserId: &id,
		})

		if err != nil {
			return fmt.Errorf("UserEntitlement with ID=%s cannot be found!. Error=%v", id, err)
		}

		if !strings.EqualFold(strings.ToLower(*userEntitlement.User.PrincipalName), strings.ToLower(expectedPrincipalName)) {
			return fmt.Errorf("UserEntitlement with ID=%s has PrincipalName=%s, but expected Name=%s", resource.Primary.ID, *userEntitlement.User.PrincipalName, expectedPrincipalName)
		}

		return nil
	}
}

// verifies that all projects referenced in the state are destroyed. This will be invoked
// *after* terraform destroys the resource but *before* the state is wiped clean.
func checkUserEntitlementDestroyed(s *terraform.State) error {
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

	//verify that every users referenced in the state does not exist in AzDO
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_user_entitlement" {
			continue
		}

		id, err := uuid.Parse(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing UserEntitlement ID, got %s: %v", resource.Primary.ID, err)
		}

		userEntitlement, err := clients.MemberEntitleManagementClient.GetUserEntitlement(clients.Ctx, memberentitlementmanagement.GetUserEntitlementArgs{
			UserId: &id,
		})

		if err != nil {
			if utils.ResponseWasNotFound(err) {
				return nil
			}
			return fmt.Errorf("Bad: Get UserEntitlment :  %+v", err)
		}

		if userEntitlement != nil && userEntitlement.AccessLevel != nil && string(*userEntitlement.AccessLevel.Status) != "none" {
			return fmt.Errorf("Status should be none : %s with readUserEntitlement error %v", string(*userEntitlement.AccessLevel.Status), err)
		}
	}

	return nil
}

type matchAddUserEntitlementArgs struct {
	t *testing.T
	x memberentitlementmanagement.AddUserEntitlementArgs
}

func MatchAddUserEntitlementArgs(t *testing.T, x memberentitlementmanagement.AddUserEntitlementArgs) gomock.Matcher {
	return &matchAddUserEntitlementArgs{t, x}
}

func (m *matchAddUserEntitlementArgs) Matches(x interface{}) bool {
	args := x.(memberentitlementmanagement.AddUserEntitlementArgs)
	m.t.Logf("MatchAddUserEntitlementArgs:\nVALUE: account_license_type: [%s], licensing_source: [%s], origin: [%s], origin_id: [%s], principal_name: [%s]\n  REF: account_license_type: [%s], licensing_source: [%s], origin: [%s], origin_id: [%s], principal_name: [%s]\n",
		*args.UserEntitlement.AccessLevel.AccountLicenseType,
		*args.UserEntitlement.AccessLevel.LicensingSource,
		*args.UserEntitlement.User.Origin,
		*args.UserEntitlement.User.OriginId,
		*args.UserEntitlement.User.PrincipalName,
		*m.x.UserEntitlement.AccessLevel.AccountLicenseType,
		*m.x.UserEntitlement.AccessLevel.LicensingSource,
		*m.x.UserEntitlement.User.Origin,
		*m.x.UserEntitlement.User.OriginId,
		*m.x.UserEntitlement.User.PrincipalName)

	return *args.UserEntitlement.AccessLevel.AccountLicenseType == *m.x.UserEntitlement.AccessLevel.AccountLicenseType &&
		*args.UserEntitlement.User.Origin == *m.x.UserEntitlement.User.Origin &&
		*args.UserEntitlement.User.OriginId == *m.x.UserEntitlement.User.OriginId &&
		*args.UserEntitlement.User.PrincipalName == *m.x.UserEntitlement.User.PrincipalName
}

func (m *matchAddUserEntitlementArgs) String() string {
	return fmt.Sprintf("account_license_type: [%s], licensing_source: [%s], origin: [%s], origin_id: [%s], principal_name: [%s]",
		*m.x.UserEntitlement.AccessLevel.AccountLicenseType,
		*m.x.UserEntitlement.AccessLevel.LicensingSource,
		*m.x.UserEntitlement.User.Origin,
		*m.x.UserEntitlement.User.OriginId,
		*m.x.UserEntitlement.User.PrincipalName)
}
