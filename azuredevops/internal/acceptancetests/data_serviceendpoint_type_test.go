package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccDataServiceEndpointType_basic(t *testing.T) {
	tfNode := "data.azuredevops_serviceendpoint_type.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataServiceEndpointType_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", "generic"),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrSet(tfNode, "display_name"),
					resource.TestCheckResourceAttrSet(tfNode, "authentication_schemes.#"),
					resource.TestCheckResourceAttrSet(tfNode, "parameters.%"),
				),
			},
		},
	})
}

func TestAccDataServiceEndpointType_withAuthScheme(t *testing.T) {
	tfNode := "data.azuredevops_serviceendpoint_type.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataServiceEndpointType_withAuthScheme(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", "generic"),
					resource.TestCheckResourceAttr(tfNode, "authorization_scheme", "UsernamePassword"),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrSet(tfNode, "display_name"),
					resource.TestCheckResourceAttrSet(tfNode, "authentication_schemes.#"),
					resource.TestCheckResourceAttrSet(tfNode, "parameters.%"),
					resource.TestCheckResourceAttrSet(tfNode, "authorization_parameters.%"),
				),
			},
		},
	})
}

func hclDataServiceEndpointType_basic() string {
	return `
data "azuredevops_serviceendpoint_type" "test" {
  name = "generic"
}
`
}

func hclDataServiceEndpointType_withAuthScheme() string {
	return `
data "azuredevops_serviceendpoint_type" "test" {
  name                 = "generic"
  authorization_scheme = "UsernamePassword"
}
`
}
