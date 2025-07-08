//go:build (all || core || data_sources || data_service_principal) && (!exclude_data_sources || !exclude_data_service_principal)

package acceptancetests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

// Validates that a configuration containing a project group lookup is able to read the resource correctly.
// Because this is a data source, there are no resources to inspect in AzDO
func TestAccServicePrincipalDataSource_Read_HappyPath(t *testing.T) {
	if os.Getenv("AZDO_TEST_AAD_SERVICE_PRINCIPAL_OBJECT_ID") == "" {
		t.Skip("Skip test due to `AZDO_TEST_AAD_SERVICE_PRINCIPAL_OBJECT_ID` not set")
	}
	servicePrincipalObjectId := os.Getenv("AZDO_TEST_AAD_SERVICE_PRINCIPAL_OBJECT_ID")

	tfBuildDefNode := "data.azuredevops_service_principal.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclServicePrincipalDataBasic(servicePrincipalObjectId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "display_name"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "id"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "origin"),
					resource.TestCheckResourceAttrSet(tfBuildDefNode, "origin_id"),
				),
			},
		},
	})
}

func hclServicePrincipalDataBasic(servicePrincipalObjectId string) string {
	return fmt.Sprintf(`
%s
data "azuredevops_service_principal" "test" {
  display_name = azuredevops_service_principal_entitlement.test.display_name
}`, testutils.HclServicePrincipleEntitlementResource(servicePrincipalObjectId))
}
