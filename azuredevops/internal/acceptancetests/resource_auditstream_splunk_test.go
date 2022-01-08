//go:build (all || resource_auditstream_splunk) && !exclude_auditstreams
// +build all resource_auditstream_splunk
// +build !exclude_auditstreams

package acceptancetests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccAuditStreamSplunk_CreateAndUpdate(t *testing.T) {
	t.Skip("Skipping test TestAccAuditStreamSplunk_CreateAndUpdate: Splunk not provisioned on test infrastructure")
	streamType := "Splunk"

	resourceType := "azuredevops_auditstream_splunk"
	tfNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckAuditStreamDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclAuditStreamSplunk(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "url"),
					testutils.CheckAuditStreamExists(tfNode, streamType),
				),
			},
		},
	})
}

func TestAccAuditStreamSplunk_CreateDisabled(t *testing.T) {
	t.Skip("Skipping test TestAccAuditStreamSplunk_CreateDisabled: Splunk not provisioned on test infrastructure")
	streamType := "Splunk"

	resourceType := "azuredevops_auditstream_splunk"
	tfNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckAuditStreamDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: testutils.HclAuditStreamAzureEventGrid(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "url"),
					resource.TestCheckResourceAttrSet(tfNode, "enabled"),
					resource.TestCheckResourceAttr(tfNode, "enabled", "false"),
					testutils.CheckAuditStreamExists(tfNode, streamType),
				),
			},
		},
	})
}
