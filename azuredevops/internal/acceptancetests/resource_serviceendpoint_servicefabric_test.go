//go:build (all || resource_serviceendpoint_service_fabric) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_service_fabric
// +build !exclude_serviceendpoints

package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// validates that an apply followed by another apply (i.e., resource update) will be reflected in AzDO and the
// underlying terraform state.
func TestAccServiceEndpointServiceFabric_CertificateCreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_servicefabric"
	tfSvcEpNode := resourceType + ".serviceendpoint"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServiceEndpointServiceFabricResource(projectName, serviceEndpointNameFirst, "Certificate"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "cluster_endpoint"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "certificate.0.server_certificate_lookup"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "certificate.0.server_certificate_thumbprint"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "certificate.0.client_certificate", "test"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "certificate.0.client_certificate_password", "test"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
				),
			},
			{
				Config: testutils.HclServiceEndpointServiceFabricResource(projectName, serviceEndpointNameSecond, "Certificate"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "cluster_endpoint"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "certificate.0.server_certificate_lookup"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "certificate.0.server_certificate_thumbprint"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "certificate.0.client_certificate", "test"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "certificate.0.client_certificate_password", "test"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
				),
			},
		},
	})
}

func TestAccServiceEndpointServiceFabric_UsernamePasswordCreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_servicefabric"
	tfSvcEpNode := resourceType + ".serviceendpoint"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServiceEndpointServiceFabricResource(projectName, serviceEndpointNameFirst, "UsernamePassword"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "cluster_endpoint"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azure_active_directory.0.server_certificate_lookup"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azure_active_directory.0.server_certificate_thumbprint"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "azure_active_directory.0.username", "test"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "azure_active_directory.0.password", "test"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
				),
			},
			{
				Config: testutils.HclServiceEndpointServiceFabricResource(projectName, serviceEndpointNameSecond, "UsernamePassword"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "cluster_endpoint"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azure_active_directory.0.server_certificate_lookup"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azure_active_directory.0.server_certificate_thumbprint"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "azure_active_directory.0.username", "test"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "azure_active_directory.0.password", "test"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
				),
			},
		},
	})
}

func TestAccServiceEndpointServiceFabric_NoneCreateAndUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	resourceType := "azuredevops_serviceendpoint_servicefabric"
	tfSvcEpNode := resourceType + ".serviceendpoint"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckServiceEndpointDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServiceEndpointServiceFabricResource(projectName, serviceEndpointNameFirst, "None"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "cluster_endpoint"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "none.0.unsecured", "false"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "none.0.cluster_spn", "test"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameFirst),
				),
			},
			{
				Config: testutils.HclServiceEndpointServiceFabricResource(projectName, serviceEndpointNameSecond, "None"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "cluster_endpoint"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "none.0.unsecured", "false"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "none.0.cluster_spn", "test"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
					testutils.CheckServiceEndpointExistsWithName(tfSvcEpNode, serviceEndpointNameSecond),
				),
			},
		},
	})
}
