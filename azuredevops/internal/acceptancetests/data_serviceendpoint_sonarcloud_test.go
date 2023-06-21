//go:build (all || data_sources || data_serviceendpoint_sonarcloud) && (!exclude_data_sources || !exclude_data_serviceendpoint_sonarcloud)
// +build all data_sources data_serviceendpoint_sonarcloud
// +build !exclude_data_sources !exclude_data_serviceendpoint_sonarcloud

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointSonarCloud_dataSource(t *testing.T) {
	name := testutils.GenerateResourceName()

	tfNode := "data.azuredevops_serviceendpoint_sonarcloud.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclServiceEndpointSonarCloudDataSource(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "service_endpoint_name", name),
					resource.TestCheckResourceAttrSet(tfNode, "service_endpoint_name"),
				),
			},
		},
	})
}

func hclServiceEndpointSonarCloudDataSource(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_sonarcloud" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%[1]s"
  token                 = "0000000000000000000000000000000000000000"
  description           = "Managed by Terraform"
}

data "azuredevops_serviceendpoint_sonarcloud" "test" {
  project_id          = azuredevops_project.test.id
  service_endpoint_name = azuredevops_serviceendpoint_sonarcloud.test.service_endpoint_name
}
`, name)
}
