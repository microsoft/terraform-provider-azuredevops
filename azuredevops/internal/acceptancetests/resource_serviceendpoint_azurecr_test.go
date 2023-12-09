//go:build (all || resource_serviceendpoint_azurecr) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_azurecr
// +build !exclude_serviceendpoints

package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// validates that an apply followed by another apply (i.e., resource update) will be reflected in AzDO and the
// underlying terraform state.
func TestAccServiceEndpointAzureCR_CreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_azurecr"
	tfSvcEpNode := resourceType + ".serviceendpoint"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, &[]string{"TEST_ARM_SUBSCRIPTION_ID", "TEST_ARM_SUBSCRIPTION_NAME", "TEST_ARM_TENANT_ID"})
		},
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServiceEndpointAzureCRResource(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_spn_tenantid"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
				),
			},
			{
				Config: testutils.HclServiceEndpointAzureCRResource(projectName, serviceEndpointNameSecond),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_spn_tenantid"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_subscription_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azurecr_subscription_name"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
				),
			},
		},
	})
}
