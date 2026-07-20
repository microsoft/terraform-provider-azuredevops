package acceptancetests

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
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
			{
				ResourceName:      tfCheckNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfCheckNode),
				ImportState:       true,
				ImportStateVerify: true,
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
			{
				ResourceName:      tfCheckNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfCheckNode),
				ImportState:       true,
				ImportStateVerify: true,
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

func TestAccCheckApproval_targetResourceDeletedOutOfBand(t *testing.T) {
	if os.Getenv("AZDO_TEST_AAD_USER_EMAIL") == "" {
		t.Skip("Skip test due to `AZDO_TEST_AAD_USER_EMAIL` not set")
	}

	projectName := testutils.GenerateResourceName()

	resourceType := "azuredevops_check_approval"
	tfCheckNode := resourceType + ".test"
	tfEnvNode := "azuredevops_environment.environment"
	principalName := os.Getenv("AZDO_TEST_AAD_USER_EMAIL")

	var projectID, environmentID string

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, &[]string{"AZDO_TEST_AAD_USER_EMAIL"}) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: testutils.CheckPipelineCheckDestroyed(resourceType),
		Steps: []resource.TestStep{
			{
				Config: hclCheckApprovalResourceEnvironment(projectName, principalName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "target_resource_type", "environment"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[tfEnvNode]
						if !ok {
							return fmt.Errorf("environment %q not found in state", tfEnvNode)
						}
						environmentID = rs.Primary.ID
						projectID = rs.Primary.Attributes["project_id"]
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
					envID, err := strconv.Atoi(environmentID)
					if err != nil {
						t.Fatalf("parsing environment ID %q: %v", environmentID, err)
					}
					if err := clients.TaskAgentClient.DeleteEnvironment(clients.Ctx, taskagent.DeleteEnvironmentArgs{
						Project:       &projectID,
						EnvironmentId: &envID,
					}); err != nil {
						t.Fatalf("deleting environment %d out of band: %v", envID, err)
					}
				},
				Config: hclCheckApprovalResourceEnvironment(projectName, principalName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(tfCheckNode, "project_id"),
					resource.TestCheckResourceAttr(tfCheckNode, "target_resource_type", "environment"),
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

func hclCheckApprovalResourceEnvironment(projectName string, principalName string) string {
	checkResource := fmt.Sprintf(`
data "azuredevops_users" "test" {
  principal_name = "%s"
}

resource "azuredevops_check_approval" "test" {
  project_id           = azuredevops_project.project.id
  target_resource_id   = azuredevops_environment.environment.id
  target_resource_type = "environment"

  requester_can_approve = false
  approvers = [
    one(data.azuredevops_users.test.users).id,
  ]
}
`, principalName)

	environmentResource := testutils.HclEnvironmentResource(projectName, "environment_test")
	return fmt.Sprintf("%s\n%s", environmentResource, checkResource)
}
