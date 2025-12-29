package acceptancetests

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[tfNode]
						if !ok {
							return fmt.Errorf("Not found: %s", tfNode)
						}

						namespaces := rs.Primary.Attributes["namespaces.#"]
						count, err := strconv.Atoi(namespaces)
						if count == 0 || err != nil {
							return fmt.Errorf("Security namespaces list is empty")
						}
						return nil
					},
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
