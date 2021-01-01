// +build all resource_serviceendpoint_service_fabric
// +build !exclude_serviceendpoints

package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// validates that an apply followed by another apply (i.e., resource update) will be reflected in AzDO and the
// underlying terraform state.
func TestAccServiceEndpointServiceFabric_CreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_servicefabric"
	tfSvcEpNode := resourceType + ".serviceendpoint"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServiceEndpointServiceFabricResource(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "cluster_endpoint"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "certificate.0.server_certificate_lookup"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "certificate.0.server_certificate_thumbprint"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "certificate.0.client_certificate_hash"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "certificate.0.client_certificate_password_hash"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "certificate.0.client_certificate", ""),
					resource.TestCheckResourceAttr(tfSvcEpNode, "certificate.0.client_certificate_password", ""),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
				),
			},
			{
				Config: testutils.HclServiceEndpointServiceFabricResource(projectName, serviceEndpointNameSecond),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "cluster_endpoint"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "certificate.0.server_certificate_lookup"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "certificate.0.server_certificate_thumbprint"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "certificate.0.client_certificate_hash"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "certificate.0.client_certificate_password_hash"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "certificate.0.client_certificate", ""),
					resource.TestCheckResourceAttr(tfSvcEpNode, "certificate.0.client_certificate_password", ""),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
				),
			},
		},
	})
}
