//go:build (all || data_sources || data_serviceendpoint_azurerm) && (!exclude_data_sources || !exclude_data_serviceendpoint_azurerm)
// +build all data_sources data_serviceendpoint_azurerm
// +build !exclude_data_sources !exclude_data_serviceendpoint_azurerm

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointAzureRM_with_serviceEndpointID_DataSource(t *testing.T) {
	serviceEndpointAzureRMName := testutils.GenerateResourceName()
	serviceprincipalid := uuid.New().String()
	serviceprincipalkey := uuid.New().String()
	projectName := testutils.GenerateResourceName()
	createServiceEndpointAzureRMWithServiceEndpointIDData := fmt.Sprintf("%s\n%s",
		testutils.HclServiceEndpointAzureRMResource(projectName, serviceEndpointAzureRMName, serviceprincipalid, serviceprincipalkey),
		testutils.HclServiceEndpointAzureRMDataSourceWithServiceEndpointID(),
	)

	tfNode := "data.azuredevops_serviceendpoint_azurerm.serviceendpointrm"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: createServiceEndpointAzureRMWithServiceEndpointIDData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "service_endpoint_name", serviceEndpointAzureRMName),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
		},
	})
}

func TestAccServiceEndpointAzureRM_with_serviceEndpointName_DataSource(t *testing.T) {
	serviceEndpointAzureRMName := testutils.GenerateResourceName()
	projectName := testutils.GenerateResourceName()
	serviceprincipalid := uuid.New().String()
	serviceprincipalkey := uuid.New().String()
	createServiceEndpointAzureRMWithServiceEndpointNameData := fmt.Sprintf("%s\n%s",
		testutils.HclServiceEndpointAzureRMResource(projectName, serviceEndpointAzureRMName, serviceprincipalid, serviceprincipalkey),
		testutils.HclServiceEndpointAzureRMDataSourceWithServiceEndpointName(serviceEndpointAzureRMName),
	)

	tfNode := "data.azuredevops_serviceendpoint_azurerm.serviceendpointrm"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: createServiceEndpointAzureRMWithServiceEndpointNameData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "service_endpoint_name", serviceEndpointAzureRMName),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
		},
	})
}
