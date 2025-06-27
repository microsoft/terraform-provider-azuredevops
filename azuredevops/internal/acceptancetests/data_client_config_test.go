//go:build all || core

package acceptancetests

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccClientConfig_LoadsCorrectProperties(t *testing.T) {
	tfNode := "data.azuredevops_client_config.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: `data "azuredevops_client_config" "test" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttrSet(tfNode, "status"),
					resource.TestCheckResourceAttrSet(tfNode, "tenant_id"),
					resource.TestCheckResourceAttrSet(tfNode, "owner_id"),
					resource.TestCheckResourceAttr(tfNode, "organization_url", os.Getenv("AZDO_ORG_SERVICE_URL")),
				),
			},
		},
	})
}
