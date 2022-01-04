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
				Config: testutils.HclAuditStreamAzureMonitorLogs(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "workspace_id"),
					testutils.CheckAuditStreamExists(tfNode, streamType),
				),
			},
		},
	})
}
