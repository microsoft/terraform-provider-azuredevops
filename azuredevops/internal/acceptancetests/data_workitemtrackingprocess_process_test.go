package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccWorkitemtrackingprocessProcess_DataSource_Get(t *testing.T) {
	tfNode := "data.azuredevops_workitemtrackingprocess_process.agile"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: hclDataSourceAgileSystemProcess(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "id", agileSystemProcessTypeId),
				),
			},
		},
	})
}

func hclDataSourceAgileSystemProcess() string {
	return fmt.Sprintf(`
data "azuredevops_workitemtrackingprocess_process" "agile" {
  id = "%s"
}
`, agileSystemProcessTypeId)
}
