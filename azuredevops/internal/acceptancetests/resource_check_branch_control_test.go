//go:build (all || resource_check_branch_control) && !exclude_approvalsandchecks
// +build all resource_check_branch_control
// +build !exclude_approvalsandchecks

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccCheckBranchControl_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	checkName := testutils.GenerateResourceName()
	branches := "refs/heads/main"

	resourceType := "azuredevops_check_branch_control"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckPipelineCheckDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclBranchControlCheckResourceBasic(projectName, checkName, branches),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckPipelineCheckExistsWithName(tfCheckNode, checkName),
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "allowed_branches", branches),
					resource.TestCheckResourceAttr(tfCheckNode, "display_name", checkName),
					resource.TestCheckResourceAttr(tfCheckNode, "timeout", "50000"),
				),
			},
		},
	})
}

func TestAccCheckBranchControl_complete(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	checkName := testutils.GenerateResourceName()
	branches := "refs/heads/main"

	resourceType := "azuredevops_check_branch_control"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckPipelineCheckDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclBranchControlCheckResourceComplete(projectName, checkName, branches),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckPipelineCheckExistsWithName(tfCheckNode, checkName),
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "allowed_branches", branches),
					resource.TestCheckResourceAttr(tfCheckNode, "display_name", checkName),
					resource.TestCheckResourceAttr(tfCheckNode, "timeout", "50000"),
					resource.TestCheckResourceAttr(tfCheckNode, "verify_branch_protection", "true"),
					resource.TestCheckResourceAttr(tfCheckNode, "ignore_unknown_protection_status", "false"),
				),
			},
		},
	})
}

func TestAccCheckBranchControl_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	checkNameFirst := testutils.GenerateResourceName()
	branchesFirst := "refs/heads/main"

	checkNameSecond := testutils.GenerateResourceName()
	branchesSecond := "refs/heads/master"

	resourceType := "azuredevops_check_branch_control"
	tfCheckNode := resourceType + ".test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckPipelineCheckDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclBranchControlCheckResourceBasic(projectName, checkNameFirst, branchesFirst),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckPipelineCheckExistsWithName(tfCheckNode, checkNameFirst),
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "allowed_branches", branchesFirst),
					resource.TestCheckResourceAttr(tfCheckNode, "display_name", checkNameFirst),
				),
			},
			{
				Config: hclBranchControlCheckResourceUpdate(projectName, checkNameSecond, branchesSecond),
				Check: resource.ComposeTestCheckFunc(
					testutils.CheckPipelineCheckExistsWithName(tfCheckNode, checkNameSecond),
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "allowed_branches", branchesSecond),
					resource.TestCheckResourceAttr(tfCheckNode, "display_name", checkNameSecond),
				),
			},
		},
	})
}

func hclBranchControlCheckResourceBasic(projectName string, checkName string, branches string) string {
	checkResource := fmt.Sprintf(`
resource "azuredevops_check_branch_control" "test" {
  project_id           = azuredevops_project.project.id
  display_name         = "%s"
  target_resource_id   = azuredevops_serviceendpoint_generic.test.id
  allowed_branches     = "%s"
  target_resource_type = "endpoint"
}`, checkName, branches)

	genericcheckResource := testutils.HclServiceEndpointGenericResource(projectName, "serviceendpoint", "https://test/", "test", "test")
	return fmt.Sprintf("%s\n%s", genericcheckResource, checkResource)
}

func hclBranchControlCheckResourceComplete(projectName string, checkName string, branches string) string {
	checkResource := fmt.Sprintf(`
resource "azuredevops_check_branch_control" "test" {
  project_id                       = azuredevops_project.project.id
  display_name                     = "%s"
  target_resource_id               = azuredevops_serviceendpoint_generic.test.id
  allowed_branches                 = "%s"
  verify_branch_protection         = true
  ignore_unknown_protection_status = false
  target_resource_type             = "endpoint"
}`, checkName, branches)

	genericcheckResource := testutils.HclServiceEndpointGenericResource(projectName, "serviceendpoint", "https://test/", "test", "test")
	return fmt.Sprintf("%s\n%s", genericcheckResource, checkResource)
}

func hclBranchControlCheckResourceUpdate(projectName string, checkName string, branches string) string {
	checkResource := fmt.Sprintf(`
resource "azuredevops_check_branch_control" "test" {
  project_id                       = azuredevops_project.project.id
  display_name                     = "%s"
  target_resource_id               = azuredevops_serviceendpoint_generic.test.id
  target_resource_type             = "endpoint"
  allowed_branches                 = "%s"
  verify_branch_protection         = true
  ignore_unknown_protection_status = false
  timeout                          = 50000
}`, checkName, branches)

	genericcheckResource := testutils.HclServiceEndpointGenericResource(projectName, "serviceendpoint", "https://test/", "test", "test")
	return fmt.Sprintf("%s\n%s", genericcheckResource, checkResource)
}
