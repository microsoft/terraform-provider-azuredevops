//go:build (all || resource_check_branch_control) && !exclude_approvalsandchecks

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
)

func TestAccCheckEndpoint(t *testing.T) {
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
				Config: hclBranchControlCheckResourceBasicEndpoint(projectName, checkName, branches),
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

func hclBranchControlCheckResourceBasicEndpoint(projectName string, checkName string, branches string) string {
	serviceEndpointResource := fmt.Sprintf(`
resource "azuredevops_check_branch_control" "test" {
  project_id           = azuredevops_project.project.id
  display_name         = "%s"
  target_resource_id   = azuredevops_serviceendpoint_generic.test.id
  allowed_branches     = "%s"
  target_resource_type = "endpoint"
}`, checkName, branches)

	genericServiceEndpointResource := testutils.HclServiceEndpointGenericResource(projectName, "serviceendpoint", "https://test/", "test", "test")
	return fmt.Sprintf("%s\n%s", genericServiceEndpointResource, serviceEndpointResource)
}

func TestAccCheckEnvironment(t *testing.T) {
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
				Config: hclBranchControlCheckResourceBasicEnvironment(projectName, checkName, branches),
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

func hclBranchControlCheckResourceBasicEnvironment(projectName string, checkName string, branches string) string {
	checkResource := fmt.Sprintf(`
resource "azuredevops_check_branch_control" "test" {
  project_id           = azuredevops_project.project.id
  display_name         = "%s"
  target_resource_id   = azuredevops_environment.environment.id
  target_resource_type = "environment"
  allowed_branches     = "%s"
}`, checkName, branches)

	environmentResource := testutils.HclEnvironmentResource(projectName, "environment_test")
	return fmt.Sprintf("%s\n%s", environmentResource, checkResource)
}

func TestAccCheckQueue(t *testing.T) {
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
				Config: hclBranchControlCheckResourceBasicQueue(projectName, checkName, branches),
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

func hclBranchControlCheckResourceBasicQueue(projectName string, checkName string, branches string) string {
	checkResource := fmt.Sprintf(`
resource "azuredevops_check_branch_control" "test" {
  project_id           = azuredevops_project.p.id
  display_name         = "%s"
  target_resource_id   = azuredevops_agent_queue.q.id
  target_resource_type = "queue"
  allowed_branches     = "%s"
}`, checkName, branches)

	queueResource := testutils.HclAgentQueueResource(projectName, "test_queue")
	return fmt.Sprintf("%s\n%s", queueResource, checkResource)
}

func TestAccCheckRepo(t *testing.T) {
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
				Config: hclBranchControlCheckResourceBasicRepo(projectName, checkName, branches),
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

func hclBranchControlCheckResourceBasicRepo(projectName string, checkName string, branches string) string {
	checkResource := fmt.Sprintf(`
resource "azuredevops_check_branch_control" "test" {
  project_id           = azuredevops_project.project.id
  display_name         = "%s"
  target_resource_id   = "${azuredevops_project.project.id}.${azuredevops_git_repository.repository.id}"
  target_resource_type = "repository"
  allowed_branches     = "%s"
}`, checkName, branches)

	projectAndRepoResource := testutils.HclGitRepoResource(projectName, testutils.GenerateResourceName(), "Clean")
	return fmt.Sprintf("%s\n%s", projectAndRepoResource, checkResource)
}

func TestAccCheckVariableGroup(t *testing.T) {
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
				Config: hclBranchControlCheckResourceBasicVariableGroup(projectName, checkName, branches),
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

func hclBranchControlCheckResourceBasicVariableGroup(projectName string, checkName string, branches string) string {
	checkResource := fmt.Sprintf(`
resource "azuredevops_check_branch_control" "test" {
  project_id           = azuredevops_project.project.id
  display_name         = "%s"
  target_resource_id   = azuredevops_variable_group.vg.id
  target_resource_type = "variablegroup"
  allowed_branches     = "%s"
}`, checkName, branches)

	variableGroupResource := testutils.HclVariableGroupResource(testutils.GenerateResourceName(), true)
	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf("%s\n%s\n%s", variableGroupResource, projectResource, checkResource)
}
