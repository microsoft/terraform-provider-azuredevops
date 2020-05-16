// +build all core

package azuredevops

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/testhelper"
)

/**
 * Begin acceptance tests
 */

// Verifies that the client config data source loads the configured AzDO org
func TestAccClientConfig_LoadsCorrectProperties(t *testing.T) {
	tfNode := "data.azuredevops_client_config.c"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testhelper.TestAccPreCheck(t, nil) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `data "azuredevops_client_config" "c" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttr(tfNode, "organization_url", os.Getenv("AZDO_ORG_SERVICE_URL")),
				),
			},
		},
	})
}

func init() {
	InitProvider()
}
