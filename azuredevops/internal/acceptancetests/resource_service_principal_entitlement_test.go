//go:build (all || resource_user_entitlement) && !exclude_resource_user_entitlement
// +build all resource_user_entitlement
// +build !exclude_resource_user_entitlement

package acceptancetests

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/memberentitlementmanagement"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
)

func TestAccServicePrincipalEntitlement_Create(t *testing.T) {
	if os.Getenv("AZDO_TEST_AAD_SERVICE_PRINCIPAL_ID") == "" {
		t.Skip("Skip test due to `AZDO_TEST_AAD_SERVICE_PRINCIPAL_ID` not set")
	}
	tfNode := "azuredevops_service_principal_entitlement.service_principal"
	ServicePrincipalId := os.Getenv("AZDO_TEST_AAD_SERVICE_PRINCIPAL_ID")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, &[]string{"AZDO_TEST_AAD_SERVICE_PRINCIPAL_ID"}) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkServicePrincipalEntitlementDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclServicePrincipalEntitlementResource(ServicePrincipalId),
				Check: resource.ComposeTestCheckFunc(
					checkServicePrincipalEntitlementExists(ServicePrincipalId),
					resource.TestCheckResourceAttrSet(tfNode, "descriptor"),
				),
			},
		},
	})
}

// Given the principalName of an AzDO userEntitlement, this will return a function that will check whether
// or not the userEntitlement (1) exists in the state and (2) exist in AzDO and (3) has the correct name
func checkServicePrincipalEntitlementExists(expectedServicePrincipalId string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources["azuredevops_user_entitlement.user"]
		if !ok {
			return fmt.Errorf(" Did not find a UserEntitlement in the TF state")
		}

		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
		id, err := uuid.Parse(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf(" Parsing ServicePrincipalEntitlement ID, got %s: %v", resource.Primary.ID, err)
		}

		servicePrincipalEntitlement, err := clients.MemberEntitleManagementClient.GetServicePrincipalEntitlement(clients.Ctx, memberentitlementmanagement.GetServicePrincipalEntitlementArgs{
			ServicePrincipalId: &id,
		})

		if err != nil {
			return fmt.Errorf(" ServicePrincipalEntitlement with ID=%s cannot be found!. Error=%v", id, err)
		}

		if !strings.EqualFold(strings.ToLower(*servicePrincipalEntitlement.ServicePrincipal.OriginId), strings.ToLower(expectedServicePrincipalId)) {
			return fmt.Errorf("ServicePrincipalEntitlement with ID=%s has OriginId=%s, but expected OriginId=%s", resource.Primary.ID, *servicePrincipalEntitlement.ServicePrincipal.OriginId, expectedServicePrincipalId)
		}

		return nil
	}
}

// verifies that all projects referenced in the state are destroyed. This will be invoked
// *after* terraform destroys the resource but *before* the state is wiped clean.
func checkServicePrincipalEntitlementDestroyed(s *terraform.State) error {
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

	//verify that every users referenced in the state does not exist in AzDO
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_service_principal_entitlement" {
			continue
		}

		id, err := uuid.Parse(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf(" Parsing ServicePrincipalEntitlement ID, got %s: %v", resource.Primary.ID, err)
		}

		servicePrincipalEntitlement, err := clients.MemberEntitleManagementClient.GetServicePrincipalEntitlement(clients.Ctx, memberentitlementmanagement.GetServicePrincipalEntitlementArgs{
			ServicePrincipalId: &id,
		})

		if err != nil {
			if utils.ResponseWasNotFound(err) {
				return nil
			}
			return fmt.Errorf("Bad: Get ServicePrincipalEntitlment :  %+v", err)
		}

		if servicePrincipalEntitlement != nil && servicePrincipalEntitlement.AccessLevel != nil && string(*servicePrincipalEntitlement.AccessLevel.Status) != "none" {
			return fmt.Errorf(" Status should be none : %s with readUserEntitlement error %v", string(*servicePrincipalEntitlement.AccessLevel.Status), err)
		}
	}

	return nil
}

func hclServicePrincipalEntitlementResource(servicePrincipalId string) string {
	return fmt.Sprintf(`
resource "azuredevops_service_principal_entitlement" "service_principal" {
  origin_id       = "%s"
  account_license_type = "express"
}`, servicePrincipalId)
}
