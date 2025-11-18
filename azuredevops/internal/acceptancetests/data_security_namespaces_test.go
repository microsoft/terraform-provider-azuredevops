package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccDataSecurityNamespaces_basic(t *testing.T) {
	tfNode := "data.azuredevops_security_namespaces.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSecurityNamespaces_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "namespaces.#"),
				),
			},
		},
	})
}

func hclDataSecurityNamespaces_basic() string {
	return `
data "azuredevops_security_namespaces" "test" {}
`
}
