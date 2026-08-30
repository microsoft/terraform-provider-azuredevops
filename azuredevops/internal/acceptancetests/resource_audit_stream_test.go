package acceptancetests

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/audit"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

func TestAccAuditStream_basic(t *testing.T) {
	streamName := testutils.GenerateResourceName()
	tfNode := "azuredevops_audit_stream.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkAuditStreamDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclAuditStreamBasic(streamName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "display_name", streamName),
					resource.TestCheckResourceAttr(tfNode, "consumer_type", "splunk"),
					resource.TestCheckResourceAttr(tfNode, "status", "enabled"),
					resource.TestCheckResourceAttr(tfNode, "consumer_inputs.#", "2"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"consumer_inputs", // Secrets are not returned by the API
				},
			},
		},
	})
}

func TestAccAuditStream_update(t *testing.T) {
	streamName := testutils.GenerateResourceName()
	tfNode := "azuredevops_audit_stream.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkAuditStreamDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclAuditStreamBasic(streamName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "display_name", streamName),
					resource.TestCheckResourceAttr(tfNode, "status", "enabled"),
				),
			},
			{
				Config: hclAuditStreamUpdate(streamName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "display_name", streamName+"_updated"),
					resource.TestCheckResourceAttr(tfNode, "status", "disabledByUser"),
				),
			},
		},
	})
}

func checkAuditStreamDestroyed(s *terraform.State) error {
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "azuredevops_audit_stream" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = clients.AuditClient.QueryStreamById(clients.Ctx, audit.QueryStreamByIdArgs{
			StreamId: &id,
		})

		if err == nil {
			return fmt.Errorf("Audit Stream with ID %d still exists", id)
		}
	}

	return nil
}

func hclAuditStreamBasic(name string) string {
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
}`, name)
}

func hclAuditStreamUpdate(name string) string {
	return fmt.Sprintf(`
resource "azuredevops_audit_stream" "test" {
  display_name  = "%s_updated"
  consumer_type = "splunk"
  status        = "disabledByUser"

  consumer_inputs {
    key   = "url"
    value = "https://splunk.example.com/updated"
  }

  consumer_inputs {
    key   = "token"
    value = "updated-token"
  }
}`, name)
}
