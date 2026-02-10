package acceptancetests

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/taskagent"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

// TestAccDeploymentGroup_basic verifies that a deployment group can be created and imported
func TestAccDeploymentGroup_basic(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	deploymentGroupName := testutils.GenerateResourceName()
	tfNode := "azuredevops_deployment_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkDeploymentGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclDeploymentGroupBasic(projectName, deploymentGroupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", deploymentGroupName),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccDeploymentGroup_update verifies that a deployment group can be updated
func TestAccDeploymentGroup_update(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	deploymentGroupNameFirst := testutils.GenerateResourceName()
	deploymentGroupNameSecond := testutils.GenerateResourceName()
	tfNode := "azuredevops_deployment_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkDeploymentGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclDeploymentGroupBasic(projectName, deploymentGroupNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", deploymentGroupNameFirst),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: hclDeploymentGroupResource(projectName, deploymentGroupNameSecond, "Updated description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", deploymentGroupNameSecond),
					resource.TestCheckResourceAttr(tfNode, "description", "Updated description"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccDeploymentGroup_withPoolId verifies that a deployment group can be created with a pool_id
func TestAccDeploymentGroup_withPoolId(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	deploymentGroupName := testutils.GenerateResourceName()
	tfNode := "azuredevops_deployment_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkDeploymentGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclDeploymentGroupWithPoolId(projectName, deploymentGroupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", deploymentGroupName),
					resource.TestCheckResourceAttrSet(tfNode, "pool_id"),
				),
			},
			{
				ResourceName:      tfNode,
				ImportStateIdFunc: testutils.ComputeProjectQualifiedResourceImportID(tfNode),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func checkDeploymentGroupDestroyed(s *terraform.State) error {
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

	for _, res := range s.RootModule().Resources {
		if res.Type != "azuredevops_deployment_group" {
			continue
		}

		deploymentGroupId, err := strconv.Atoi(res.Primary.ID)
		if err != nil {
			return fmt.Errorf("Deployment group ID=%s cannot be parsed. Error=%v", res.Primary.ID, err)
		}
		projectID := res.Primary.Attributes["project_id"]

		if _, err := clients.TaskAgentClient.GetDeploymentGroup(clients.Ctx, taskagent.GetDeploymentGroupArgs{
			Project:           &projectID,
			DeploymentGroupId: &deploymentGroupId,
		}); err == nil {
			return fmt.Errorf("Deployment group ID %d should not exist", deploymentGroupId)
		}
	}

	return nil
}

func hclDeploymentGroupBasic(projectName string, deploymentGroupName string) string {
	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf(`
%s

resource "azuredevops_deployment_group" "test" {
  project_id = azuredevops_project.project.id
  name       = "%s"
}
`, projectResource, deploymentGroupName)
}

func hclDeploymentGroupResource(projectName string, deploymentGroupName string, description string) string {
	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf(`
%s

resource "azuredevops_deployment_group" "test" {
  project_id  = azuredevops_project.project.id
  name        = "%s"
  description = "%s"
}
`, projectResource, deploymentGroupName, description)
}

func hclDeploymentGroupWithPoolId(projectName string, deploymentGroupName string) string {
	projectResource := testutils.HclProjectResource(projectName)
	return fmt.Sprintf(`
%s

resource "azuredevops_deployment_group" "pool_source" {
  project_id = azuredevops_project.project.id
  name       = "%s-source"
}

resource "azuredevops_deployment_group" "test" {
  project_id = azuredevops_project.project.id
  name       = "%s"
  pool_id    = azuredevops_deployment_group.pool_source.pool_id
}
`, projectResource, deploymentGroupName, deploymentGroupName)
}
