//go:build (all || resource_auditstream_azuremonitorlogs) && !exclude_auditstreams
// +build all resource_auditstream_azuremonitorlogs
// +build !exclude_auditstreams

package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccAuditStreamAzureMonitorLogs_CreateAndUpdate(t *testing.T) {
	t.Skip("Skipping test TestAccAuditStreamAzureMonitorLogs_CreateAndUpdate: Azure Monitor not provisioned on test infrastructure")
	streamType := "AzureMonitorLogs"

	resourceType := "azuredevops_auditstream_azuremonitorlogs"
	tfNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckAuditStreamDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclAuditStreamAzureMonitorLogs(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "workspace_id"),
					resource.TestCheckResourceAttrSet(tfNode, "enabled"),
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttr(tfNode, "enabled", "true"),
					testutils.CheckAuditStreamExists(tfNode, streamType),
				),
			},
		},
	})
}

func TestAccAuditStreamAzureMonitorLogs_CreateDisabled(t *testing.T) {
	t.Skip("Skipping test TestAccAuditStreamAzureMonitorLogs_CreateDisabled: Azure Monitor not provisioned on test infrastructure")
	streamType := "AzureMonitorLogs"

	resourceType := "azuredevops_auditstream_azuremonitorlogs"
	tfNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckAuditStreamDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclAuditStreamAzureEventGrid(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "workspace_id"),
					resource.TestCheckResourceAttrSet(tfNode, "enabled"),
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttr(tfNode, "enabled", "false"),
					testutils.CheckAuditStreamExists(tfNode, streamType),
				),
			},
		},
	})
}
