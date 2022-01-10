//go:build (all || resource_auditstream_azureeventgrid) && !exclude_auditstreams
// +build all resource_auditstream_azureeventgrid
// +build !exclude_auditstreams

package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccAuditStreamAzureEventGrid_CreateAndUpdate(t *testing.T) {
	t.Skip("Skipping test TestAccAuditStreamAzureEventGrid_CreateAndUpdate: event grid not provisioned on test infrastructure")
	streamType := "AzureEventGrid"

	resourceType := "azuredevops_auditstream_azureeventgrid"
	tfNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckAuditStreamDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclAuditStreamAzureEventGrid(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "topic_url"),
					resource.TestCheckResourceAttrSet(tfNode, "enabled"),
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttr(tfNode, "enabled", "true"),
					testutils.CheckAuditStreamExists(tfNode, streamType),
					testutils.CheckAuditStreamStatus(tfNode, true),
				),
			},
			{
				Config: testutils.HclAuditStreamAzureEventGrid(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "topic_url"),
					resource.TestCheckResourceAttrSet(tfNode, "enabled"),
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttr(tfNode, "enabled", "false"),
					testutils.CheckAuditStreamExists(tfNode, streamType),
					testutils.CheckAuditStreamStatus(tfNode, false),
				),
			},
		},
	})
}

func TestAccAuditStreamAzureEventGrid_CreateDisabled(t *testing.T) {
	t.Skip("Skipping test TestAccAuditStreamAzureEventGrid_CreateDisabled: event grid not provisioned on test infrastructure")
	streamType := "AzureEventGrid"

	resourceType := "azuredevops_auditstream_azureeventgrid"
	tfNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckAuditStreamDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclAuditStreamAzureEventGrid(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "topic_url"),
					resource.TestCheckResourceAttrSet(tfNode, "enabled"),
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttr(tfNode, "enabled", "false"),
					testutils.CheckAuditStreamExists(tfNode, streamType),
					testutils.CheckAuditStreamStatus(tfNode, false),
				),
			},
			{
				Config: testutils.HclAuditStreamAzureEventGrid(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "topic_url"),
					resource.TestCheckResourceAttrSet(tfNode, "enabled"),
					resource.TestCheckResourceAttrSet(tfNode, "name"),
					resource.TestCheckResourceAttr(tfNode, "enabled", "true"),
					testutils.CheckAuditStreamExists(tfNode, streamType),
					testutils.CheckAuditStreamStatus(tfNode, true),
				),
			},
		},
	})
}
