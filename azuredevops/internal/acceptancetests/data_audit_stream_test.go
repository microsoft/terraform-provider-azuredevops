package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccAuditStream_DataSource(t *testing.T) {
	streamName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_audit_stream.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclAuditStreamDataSource(streamName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "display_name", streamName),
					resource.TestCheckResourceAttr(tfNode, "consumer_type", "splunk"),
					resource.TestCheckResourceAttrSet(tfNode, "id"),
				),
			},
		},
	})
}

func TestAccAuditStreams_DataSource(t *testing.T) {
	streamName := testutils.GenerateResourceName()
	tfNode := "data.azuredevops_audit_streams.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testutils.PreCheck(t, nil) },
		Providers: testutils.GetProviders(),
		Steps: []resource.TestStep{
			{
				Config: hclAuditStreamDataSource(streamName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "id"),
					resource.TestCheckResourceAttrSet(tfNode, "streams.#"),
				),
			},
		},
	})
}

func hclAuditStreamDataSource(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_audit_stream" "test" {
  display_name  = "%s"
  consumer_type = "splunk"
  status        = "enabled"

  consumer_inputs {
    key   = "url"
    value = "https://splunk.example.com"
  }

  consumer_inputs {
    key   = "token"
    value = "dummy-token"
  }
}

data "azuredevops_audit_stream" "test" {
  display_name = azuredevops_audit_stream.test.display_name
}

data "azuredevops_audit_streams" "test" {
  depends_on = [azuredevops_audit_stream.test]
}
`, name)
}
