package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccDataSecurityNamespaceToken_collection(t *testing.T) {
	tfNode := "data.azuredevops_security_namespace_token.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSecurityNamespaceToken_collection(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "token", "NAMESPACE:"),
				),
			},
		},
	})
}

func hclDataSecurityNamespaceToken_collection() string {
	return `
data "azuredevops_security_namespace_token" "test" {
  namespace_name = "Collection"
}
`
}
