//go:build (all || resource_check_branch_control) && !exclude_approvalsandchecks

package acceptancetests

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccCheckApproval_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()

	resourceType := "azuredevops_check_approval"
	tfCheckNode := resourceType + ".test"
	principalName := os.Getenv("AZDO_TEST_AAD_USER_EMAIL")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, &[]string{"AZDO_TEST_AAD_USER_EMAIL"}) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckPipelineCheckDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclCheckApprovalResourceBasic(projectName, principalName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "requester_can_approve", "false"),
					resource.TestCheckResourceAttr(tfCheckNode, "timeout", "43200"),
					resource.TestCheckResourceAttr(tfCheckNode, "approvers.#", "1"),
				),
			},
		},
	})
}

func TestAccCheckApproval_complete(t *testing.T) {
	projectName := testutils.GenerateResourceName()

	resourceType := "azuredevops_check_approval"
	tfCheckNode := resourceType + ".test"
	principalName := os.Getenv("AZDO_TEST_AAD_USER_EMAIL")
	azdoGroupName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckPipelineCheckDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclCheckApprovalResourceComplete(projectName, principalName, azdoGroupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "requester_can_approve", "true"),
					resource.TestCheckResourceAttr(tfCheckNode, "timeout", "40000"),
					resource.TestCheckResourceAttr(tfCheckNode, "approvers.#", "2"),
				),
			},
		},
	})
}

func TestAccCheckApproval_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()

	resourceType := "azuredevops_check_approval"
	tfCheckNode := resourceType + ".test"
	principalName := os.Getenv("AZDO_TEST_AAD_USER_EMAIL")
	azdoGroupName := testutils.GenerateResourceName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, &[]string{"AZDO_TEST_AAD_USER_EMAIL"}) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckPipelineCheckDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclCheckApprovalResourceBasic(projectName, principalName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "approvers.#", "1"),
				),
			},
			{
				Config: hclCheckApprovalResourceComplete(projectName, principalName, azdoGroupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "approvers.#", "2"),
					resource.TestCheckResourceAttr(tfCheckNode, "version", "2"),
				),
			},
		},
	})
}

func hclCheckApprovalResourceBasic(projectName string, principalName string) string {
	checkResource := fmt.Sprintf(`
data "azuredevops_users" "test" {
  principal_name = "%s"
}

resource "azuredevops_check_approval" "test" {
  project_id           = azuredevops_project.project.id
  target_resource_id   = azuredevops_serviceendpoint_generic.test.id
  target_resource_type = "endpoint"

  requester_can_approve = false
  approvers = [
    one(data.azuredevops_users.test.users).id,
  ]
}
`, principalName)

	genericcheckResource := testutils.HclServiceEndpointGenericResource(projectName, "serviceendpoint", "https://test/", "test", "test")
	return fmt.Sprintf("%s\n%s", genericcheckResource, checkResource)
}

func hclCheckApprovalResourceComplete(projectName string, principalName string, azdoGroupName string) string {
	checkResource := fmt.Sprintf(`
data "azuredevops_users" "test" {
  principal_name = "%s"
}

resource "azuredevops_group" "test" {
  display_name = "%s"
}

resource "azuredevops_check_approval" "test" {
  project_id           = azuredevops_project.project.id
  target_resource_id   = azuredevops_serviceendpoint_generic.test.id
  target_resource_type = "endpoint"

  requester_can_approve = true
  approvers = [
    one(data.azuredevops_users.test.users).id,
    azuredevops_group.test.origin_id,
  ]

  timeout = 40000
}
`, principalName, azdoGroupName)

	genericcheckResource := testutils.HclServiceEndpointGenericResource(projectName, "serviceendpoint", "https://test/", "test", "test")
	return fmt.Sprintf("%s\n%s", genericcheckResource, checkResource)
}
