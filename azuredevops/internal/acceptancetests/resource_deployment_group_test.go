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

// TestAccDeploymentGroup_basic verifies that a deployment group can be created and updated
func TestAccDeploymentGroup_basic(t *testing.T) {
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
				Config: hclDeploymentGroupResource(projectName, deploymentGroupNameFirst, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", deploymentGroupNameFirst),
					resource.TestCheckResourceAttr(tfNode, "description", ""),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					checkDeploymentGroupExists(deploymentGroupNameFirst),
				),
			},
			{
				Config: hclDeploymentGroupResource(projectName, deploymentGroupNameSecond, "Updated description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", deploymentGroupNameSecond),
					resource.TestCheckResourceAttr(tfNode, "description", "Updated description"),
					checkDeploymentGroupExists(deploymentGroupNameSecond),
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

// TestAccDeploymentGroup_withDescription verifies that a deployment group with description can be created
func TestAccDeploymentGroup_withDescription(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	deploymentGroupName := testutils.GenerateResourceName()
	tfNode := "azuredevops_deployment_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkDeploymentGroupDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclDeploymentGroupResource(projectName, deploymentGroupName, "Test description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", deploymentGroupName),
					resource.TestCheckResourceAttr(tfNode, "description", "Test description"),
					checkDeploymentGroupExists(deploymentGroupName),
				),
			},
		},
	})
}

func checkDeploymentGroupExists(expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		res, ok := s.RootModule().Resources["azuredevops_deployment_group.test"]
		if !ok {
			return fmt.Errorf("Did not find a deployment group in the TF state")
		}

		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
		deploymentGroupId, err := strconv.Atoi(res.Primary.ID)
		if err != nil {
			return fmt.Errorf("Parse ID error, ID: %v. Error= %v", res.Primary.ID, err)
		}
		projectID := res.Primary.Attributes["project_id"]

		deploymentGroup, err := clients.TaskAgentClient.GetDeploymentGroup(clients.Ctx, taskagent.GetDeploymentGroupArgs{
			Project:           &projectID,
			DeploymentGroupId: &deploymentGroupId,
		})
		if err != nil {
			return fmt.Errorf("Deployment group with ID=%d cannot be found. Error=%v", deploymentGroupId, err)
		}

		if *deploymentGroup.Name != expectedName {
			return fmt.Errorf("Deployment group with ID=%d has Name=%s, but expected Name=%s", deploymentGroupId, *deploymentGroup.Name, expectedName)
		}

		return nil
	}
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
