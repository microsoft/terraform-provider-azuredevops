//go:build (all || data_sources || data_serviceendpoint_azurecr) && (!exclude_data_sources || !exclude_data_serviceendpoint_azurecr)
// +build all data_sources data_serviceendpoint_azurecr
// +build !exclude_data_sources !exclude_data_serviceendpoint_azurecr

package acceptancetests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccServiceEndpointAzureCR_dataSource(t *testing.T) {
	name := testutils.GenerateResourceName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, &[]string{"TEST_ARM_SUBSCRIPTION_ID", "TEST_ARM_SUBSCRIPTION_NAME", "TEST_ARM_TENANT_ID", "TEST_ARM_RESOURCE_GROUP", "TEST_ARM_ACR_NAME"})
		},
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclServiceEndpointAzureCRDataSource(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.azuredevops_serviceendpoint_azurecr.test", "service_endpoint_name", name),
				),
			},
		},
	})
}

func hclServiceEndpointAzureCRDataSource(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name               = "%[1]s"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_azurecr" "test" {
  project_id                = azuredevops_project.test.id
  service_endpoint_name     = "%[1]s"
  azurecr_subscription_id   = "%[2]s"
  azurecr_subscription_name = "%[3]s"
  azurecr_spn_tenantid      = "%[4]s"
  resource_group            = "%[5]s"
  azurecr_name              = "%[6]s"
}

data "azuredevops_serviceendpoint_azurecr" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = azuredevops_serviceendpoint_azurecr.test.service_endpoint_name
}
`, name, os.Getenv("TEST_ARM_SUBSCRIPTION_ID"),
		os.Getenv("TEST_ARM_SUBSCRIPTION_NAME"), os.Getenv("TEST_ARM_TENANT_ID"),
		os.Getenv("TEST_ARM_RESOURCE_GROUP"), os.Getenv("TEST_ARM_ACR_NAME"))
}
