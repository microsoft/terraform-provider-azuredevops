package acceptancetests

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccDataServiceEndpointTypes_basic(t *testing.T) {
	tfNode := "data.azuredevops_serviceendpoint_types.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataServiceEndpointTypes_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "types.#"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[tfNode]
						if !ok {
							return fmt.Errorf("Not found: %s", tfNode)
						}

						types := rs.Primary.Attributes["types.#"]
						count, err := strconv.Atoi(types)
						if count == 0 || err != nil {
							return fmt.Errorf("Service endpoint types list is empty")
						}
						return nil
					},
				),
			},
		},
	})
}

func hclDataServiceEndpointTypes_basic() string {
	return `
data "azuredevops_serviceendpoint_types" "test" {}
`
}
