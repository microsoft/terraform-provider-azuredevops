//go:build (all || data_sources || data_storage_key) && (!exclude_data_sources || !exclude_data_storage_key)

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccStorageKeyDatasource(t *testing.T) {
	name := testutils.GenerateResourceName() + "@contoso.com"
	tfNode := "data.azuredevops_storage_key.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testutils.PreCheck(t, nil) },
		ProviderFactories:         testutils.GetProviderFactories(),
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: hclStorageKeyDataSource(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "storage_key"),
				),
			},
		},
	})
}

func hclStorageKeyDataSource(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_user_entitlement" "test" {
  principal_name       = "%s"
  account_license_type = "express"
}

data "azuredevops_storage_key" "test" {
  descriptor = azuredevops_user_entitlement.test.descriptor
}`, name)
}
