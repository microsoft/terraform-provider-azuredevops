package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccAzureDevOpsSecurityroles_DataSource_Securityrole_Definitions(t *testing.T) {
	securityroleDefinitionsData := testutils.HclSecurityroleDefinitionsDataSource()

	tfNode := "data.azuredevops_securityrole_definitions.definitions-list"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: securityroleDefinitionsData,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "definitions.#"),
				),
			},
		},
	})
}
