//go:build (all || resource_serviceendpoint_generic) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_generic
// +build !exclude_serviceendpoints

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccBranchControlCheck_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	checkName := testutils.GenerateResourceName()
	branches := "refs/heads/main"

	resourceType := "azuredevops_serviceendpoint_check_branch_control"
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
				),
			},
		},
	})
}

func TestAccBranchControlCheck_complete(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	checkName := testutils.GenerateResourceName()
	branches := "refs/heads/main"

	resourceType := "azuredevops_serviceendpoint_check_branch_control"
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
					resource.TestCheckResourceAttr(tfCheckNode, "verify_branch_protection", "true"),
					resource.TestCheckResourceAttr(tfCheckNode, "ignore_unknown_protection_status", "false"),
				),
			},
		},
	})
}

func TestAccBranchControlCheck_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	checkNameFirst := testutils.GenerateResourceName()
	branchesFirst := "refs/heads/main"

	checkNameSecond := testutils.GenerateResourceName()
	branchesSecond := "refs/heads/master"

	resourceType := "azuredevops_serviceendpoint_check_branch_control"
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
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_check_branch_control" "test" {
	project_id       = azuredevops_project.project.id
	display_name     = "%s"
	endpoint_id      = azuredevops_serviceendpoint_generic.test.id
	allowed_branches = "%s"
}`, checkName, branches)

	genericServiceEndpointResource := testutils.HclServiceEndpointGenericResource(projectName, "serviceendpoint", "https://test/", "test", "test")
	return fmt.Sprintf("%s\n%s", genericServiceEndpointResource, serviceEndpointResource)
}

func hclBranchControlCheckResourceComplete(projectName string, checkName string, branches string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_check_branch_control" "test" {
	project_id                       = azuredevops_project.project.id
	display_name                     = "%s"
	endpoint_id                      = azuredevops_serviceendpoint_generic.test.id
	allowed_branches                 = "%s"
	verify_branch_protection         = true
	ignore_unknown_protection_status = false
}`, checkName, branches)

	genericServiceEndpointResource := testutils.HclServiceEndpointGenericResource(projectName, "serviceendpoint", "https://test/", "test", "test")
	return fmt.Sprintf("%s\n%s", genericServiceEndpointResource, serviceEndpointResource)
}

func hclBranchControlCheckResourceUpdate(projectName string, checkName string, branches string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_serviceendpoint_check_branch_control" "test" {
	project_id                       = azuredevops_project.project.id
	display_name                     = "%s"
	endpoint_id                      = azuredevops_serviceendpoint_generic.test.id
	allowed_branches                 = "%s"
	verify_branch_protection         = true
	ignore_unknown_protection_status = false
}`, checkName, branches)

	genericServiceEndpointResource := testutils.HclServiceEndpointGenericResource(projectName, "serviceendpoint", "https://test/", "test", "test")
	return fmt.Sprintf("%s\n%s", genericServiceEndpointResource, serviceEndpointResource)
}
