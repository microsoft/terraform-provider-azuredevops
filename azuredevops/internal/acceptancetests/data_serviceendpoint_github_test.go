//go:build (all || data_sources || data_serviceendpoint_github) && (!exclude_data_sources || !exclude_data_serviceendpoint_github)
// +build all data_sources data_serviceendpoint_github
// +build !exclude_data_sources !exclude_data_serviceendpoint_github

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointGitHub_with_serviceEndpointID_DataSource(t *testing.T) {
	serviceEndpointGitHubName := testutils.GenerateResourceName()
	projectName := testutils.GenerateResourceName()
	createServiceEndpointGitHubWithServiceEndpointIDData := fmt.Sprintf("%s\n%s",
		testutils.HclServiceEndpointGitHubResource(projectName, serviceEndpointGitHubName),
		testutils.HclServiceEndpointGitHubDataSourceWithServiceEndpointID(),
	)

	tfNode := "data.azuredevops_serviceendpoint_github.serviceendpoint"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: createServiceEndpointGitHubWithServiceEndpointIDData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "service_endpoint_name", serviceEndpointGitHubName),
					resource.TestCheckResourceAttrSet(tfNode, "service_endpoint_id"),
				),
			},
		},
	})
}

func TestAccServiceEndpointGitHub_with_serviceEndpointName_DataSource(t *testing.T) {
	serviceEndpointGitHubName := testutils.GenerateResourceName()
	projectName := testutils.GenerateResourceName()
	createServiceEndpointGitHubWithServiceEndpointNameData := fmt.Sprintf("%s\n%s",
		testutils.HclServiceEndpointGitHubResource(projectName, serviceEndpointGitHubName),
		testutils.HclServiceEndpointGitHubDataSourceWithServiceEndpointName(serviceEndpointGitHubName),
	)

	tfNode := "data.azuredevops_serviceendpoint_github.serviceendpoint"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: createServiceEndpointGitHubWithServiceEndpointNameData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "service_endpoint_name", serviceEndpointGitHubName),
					resource.TestCheckResourceAttrSet(tfNode, "service_endpoint_id"),
				),
			},
		},
	})
}
