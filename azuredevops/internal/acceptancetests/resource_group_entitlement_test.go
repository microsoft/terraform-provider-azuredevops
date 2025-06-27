//go:build (all || resource_group_entitlement) && !exclude_resource_group_entitlement

package acceptancetests

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/memberentitlementmanagement"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils"
)

func TestAccGroupEntitlement_Create(t *testing.T) {
	tfNode := "azuredevops_group_entitlement.test"
	displayName := "group-038c153d-c86e-443c-b6f6-3d97378025d0"
	resource.ParallelTest(t, resource.TestCase{
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkGroupEntitlementDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclGroupEntitlementResourceBasic(displayName),
				Check: resource.ComposeTestCheckFunc(
					checkGroupEntitlementExists(),
					resource.TestCheckResourceAttrSet(tfNode, "descriptor"),
				),
			},
		},
	})
}

func TestAccGroupEntitlement_AAD_Create(t *testing.T) {
	if os.Getenv("AZDO_TEST_AAD_GROUP_ID") == "" {
		t.Skip("Skip test dueto `AZDO_TEST_AAD_GROUP_ID` not set")
	}
	tfNode := "azuredevops_group_entitlement.test"
	originId := os.Getenv("AZDO_TEST_AAD_GROUP_ID")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, &[]string{"AZDO_TEST_AAD_GROUP_ID"}) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkGroupEntitlementDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclGroupEntitlementResourceAAD(originId),
				Check: resource.ComposeTestCheckFunc(
					checkGroupEntitlementExists(),
					resource.TestCheckResourceAttrSet(tfNode, "descriptor"),
					resource.TestCheckResourceAttr(tfNode, "origin_id", originId),
				),
			},
		},
	})
}

// Given the principalName of an AzDO groupEntitlement, this will return a function that will check whether
// or not the groupEntitlement (1) exists in the state and (2) exist in AzDO and (3) has the correct name
func checkGroupEntitlementExists() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resource, ok := s.RootModule().Resources["azuredevops_group_entitlement.test"]
		if !ok {
			return fmt.Errorf("Did not find a GroupEntitlement in the TF state")
		}

		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
		id, err := uuid.Parse(resource.Primary.ID)
		if err != nil {
			return fmt.Errorf("Parsing GroupEntitlement ID, got %s: %v", resource.Primary.ID, err)
		}

		groupEntitlement, err := clients.MemberEntitleManagementClient.GetGroupEntitlement(clients.Ctx, memberentitlementmanagement.GetGroupEntitlementArgs{
			GroupId: &id,
		})

		if err != nil {
			return fmt.Errorf("GroupEntitlement with ID=%s cannot be found!. Error=%v", id, err)
		}

		if groupEntitlement == nil || groupEntitlement.Id == nil {
			return fmt.Errorf("GroupEntitlement with ID=%s cannot be found.", id)
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
			return fmt.Errorf("Parsing GroupEntitlement ID, got %s: %v", resource.Primary.ID, err)
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

func hclGroupEntitlementResourceBasic(displayName string) string {
	return fmt.Sprintf(`
resource "azuredevops_group_entitlement" "test" {
  display_name         = "%s"
  account_license_type = "express"
}`, displayName)
}

func hclGroupEntitlementResourceAAD(originId string) string {
	return fmt.Sprintf(`
resource "azuredevops_group_entitlement" "test" {
  origin_id            = "%s"
  origin               = "aad"
  account_license_type = "express"
}`, originId)
}
