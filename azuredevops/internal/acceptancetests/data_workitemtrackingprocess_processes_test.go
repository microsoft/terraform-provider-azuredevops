package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccWorkitemtrackingprocessProcesses_DataSource_AllProcesses(t *testing.T) {
	tfNode := "data.azuredevops_workitemtrackingprocess_processes.all"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSourceAllProcesses(),
				Check:  testutils.TestCheckAttrGreaterThan(tfNode, "processes.#", 0),
			},
		},
	})
}

func hclDataSourceAllProcesses() string {
	return `
data "azuredevops_workitemtrackingprocess_processes" "all" {
}
`
}
