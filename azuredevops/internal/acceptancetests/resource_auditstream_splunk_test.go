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
				Config: testutils.HclAuditStreamSplunk(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfNode, "url"),
					testutils.CheckAuditStreamExists(tfNode, streamType),
				),
			},
		},
	})
}
