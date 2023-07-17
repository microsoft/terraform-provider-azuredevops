//go:build (all || resource_check_exclusive_lock) && !exclude_approvalsandchecks
// +build all resource_check_exclusive_lock
// +build !exclude_approvalsandchecks

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccCheckExclusiveLock_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	timeout := 43200

	resourceType := "azuredevops_check_exclusive_lock"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckPipelineCheckDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclCheckExclusiveLockResourceBasic(projectName, timeout),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfCheckNode, "target_resource_id"),
					resource.TestCheckResourceAttrSet(tfCheckNode, "target_resource_type"),
					resource.TestCheckResourceAttr(tfCheckNode, "timeout", fmt.Sprintf("%d", timeout)),
				),
			},
		},
	})
}

func hclCheckExclusiveLockResourceBasic(projectName string, timeout int) string {
	checkResource := fmt.Sprintf(`
resource "azuredevops_check_exclusive_lock" "test" {
  project_id           = azuredevops_project.project.id
  target_resource_id   = azuredevops_serviceendpoint_generic.test.id
  target_resource_type = "endpoint"
  timeout              = %d
}`, timeout)

	genericServiceEndpointResource := testutils.HclServiceEndpointGenericResource(projectName, "serviceendpoint", "https://test/", "test", "test")
	return fmt.Sprintf("%s\n%s", genericServiceEndpointResource, checkResource)
}
