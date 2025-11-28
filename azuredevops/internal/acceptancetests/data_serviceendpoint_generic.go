package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointGeneric_dataSource_with_serviceEndpointID(t *testing.T) {
	serviceEndpointGenericName := testutils.GenerateResourceName()
	projectName := testutils.GenerateResourceName()
	serverUrl := testutils.GenerateResourceName()
	username := testutils.GenerateResourceName()
	password := testutils.GenerateResourceName()
	createServiceEndpointGenericWithServiceEndpointIDData := fmt.Sprintf("%s\n%s",
		testutils.HclServiceEndpointGenericResource(projectName, serviceEndpointGenericName, serverUrl, username, password),
		testutils.HclServiceEndpointGenericDataSourceWithServiceEndpointID(),
	)

	tfNode := "data.azuredevops_serviceendpoint_github.serviceendpoint"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: createServiceEndpointGenericWithServiceEndpointIDData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "service_endpoint_name", serviceEndpointGenericName),
					resource.TestCheckResourceAttrSet(tfNode, "service_endpoint_id"),
				),
			},
		},
	})
}

func TestAccServiceEndpointGeneric_dataSource_with_serviceEndpointName_DataSource(t *testing.T) {
	serviceEndpointGenericName := testutils.GenerateResourceName()
	projectName := testutils.GenerateResourceName()
	serverUrl := testutils.GenerateResourceName()
	username := testutils.GenerateResourceName()
	password := testutils.GenerateResourceName()
	createServiceEndpointGenericWithServiceEndpointNameData := fmt.Sprintf("%s\n%s",
		testutils.HclServiceEndpointGenericResource(projectName, serviceEndpointGenericName, serverUrl, username, password),
		testutils.HclServiceEndpointGenericDataSourceWithServiceEndpointID(),
	)

	tfNode := "data.azuredevops_serviceendpoint_github.serviceendpoint"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: createServiceEndpointGenericWithServiceEndpointNameData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "service_endpoint_name", serviceEndpointGenericName),
					resource.TestCheckResourceAttrSet(tfNode, "service_endpoint_id"),
				),
			},
		},
	})
}
