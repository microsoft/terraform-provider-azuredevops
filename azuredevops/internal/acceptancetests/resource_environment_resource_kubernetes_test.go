//go:build (all || resource_environment_resource_kubernetes) && !exclude_resource_environment_resource_kubernetes
// +build all resource_environment_resource_kubernetes
// +build !exclude_resource_environment_resource_kubernetes

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

func TestAccEnvironmentKubernetes_createUpdate(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	environmentName := testutils.GenerateResourceName()
	serviceEndpointName := testutils.GenerateResourceName()
	resourceNameFirst := testutils.GenerateResourceName()
	resourceNameSecond := testutils.GenerateResourceName()
	tfNode := "azuredevops_environment_resource_kubernetes.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkEnvironmentKubernetesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclEnvironmentKubernetes(projectName, environmentName, serviceEndpointName, resourceNameFirst),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", resourceNameFirst),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "environment_id"),
					resource.TestCheckResourceAttrSet(tfNode, "service_endpoint_id"),
					checkEnvironmentKubernetesExists(tfNode, resourceNameFirst),
				),
			},
			{
				Config: hclEnvironmentKubernetes(projectName, environmentName, serviceEndpointName, resourceNameSecond),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(tfNode, "name", resourceNameSecond),
					resource.TestCheckResourceAttrSet(tfNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfNode, "environment_id"),
					resource.TestCheckResourceAttrSet(tfNode, "service_endpoint_id"),
					checkEnvironmentKubernetesExists(tfNode, resourceNameSecond),
				),
			},
		},
	})
}

// Given the name of a resource, this will return a function that will check whether
// the resource (1) exists in the state and (2) exist in AzDO and (3) has the correct name
func checkEnvironmentKubernetesExists(tfNode string, expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		res, ok := s.RootModule().Resources[tfNode]
		if !ok {
			return fmt.Errorf(" Did not find an resource in the TF state")
		}

		clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
		id, err := strconv.Atoi(res.Primary.ID)
		if err != nil {
			return fmt.Errorf("Parse resource id, ID:  %v !. Error= %v", res.Primary.ID, err)
		}
		projectId := res.Primary.Attributes["project_id"]
		environmentIdStr := res.Primary.Attributes["environment_id"]
		environmentId, err := strconv.Atoi(environmentIdStr)
		if err != nil {
			return fmt.Errorf("Parse environment_id error, ID:  %v !. Error= %v", environmentIdStr, err)
		}

		readResource, err := readEnvironmentKubernetes(clients, projectId, environmentId, id)
		if err != nil {
			return fmt.Errorf(" Resource with ID=%d cannot be found!. Error=%v", id, err)
		}

		if *readResource.Name != expectedName {
			return fmt.Errorf(" Resource with ID=%d has Name=%s, but expected Name=%s", id, *readResource.Name, expectedName)
		}
		return nil
	}
}

// verifies that environment referenced in the state is destroyed. This will be invoked
// *after* terraform destroys the resource but *before* the state is wiped clean.
func checkEnvironmentKubernetesDestroyed(s *terraform.State) error {
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)

	// verify that every environment referenced in the state does not exist in AzDO
	for _, res := range s.RootModule().Resources {
		if res.Type != "azuredevops_environment_kubernetes" {
			continue
		}

		id, err := strconv.Atoi(res.Primary.ID)
		if err != nil {
			return fmt.Errorf("Parse resource id, ID:  %v !. Error= %v", res.Primary.ID, err)
		}
		projectId := res.Primary.Attributes["project_id"]
		environmentIdStr := res.Primary.Attributes["environment_id"]
		environmentId, err := strconv.Atoi(environmentIdStr)
		if err != nil {
			return fmt.Errorf("Parse environment_id error, ID:  %v !. Error= %v", environmentIdStr, err)
		}

		// indicates the environment still exists - this should fail the test
		if _, err := readEnvironmentKubernetes(clients, projectId, environmentId, id); err == nil {
			return fmt.Errorf(" Resource ID %d should not exist", id)
		}
	}
	return nil
}

// Lookup an Environment using the ID and the project ID.
func readEnvironmentKubernetes(clients *client.AggregatedClient, projectId string, environmentId int, resourceId int) (*taskagent.KubernetesResource, error) {
	return clients.TaskAgentClient.GetKubernetesResource(clients.Ctx,
		taskagent.GetKubernetesResourceArgs{
			Project:       &projectId,
			EnvironmentId: &environmentId,
			ResourceId:    &resourceId,
		},
	)
}

func hclEnvironmentKubernetes(projectName, environmentName, serviceEndpointName, k8sName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_environment" "test" {
  project_id = azuredevops_project.test.id
  name       = "%s"
}

resource "azuredevops_serviceendpoint_kubernetes" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  apiserver_url         = "https://test-dns-r9lconkh.hcp.eastus.azmk8s.io:443"
  authorization_type    = "ServiceAccount"
  service_account {
    token   = "test"
    ca_cert = "test"
  }
}

resource "azuredevops_environment_resource_kubernetes" "test" {
  project_id          = azuredevops_project.test.id
  environment_id      = azuredevops_environment.test.id
  service_endpoint_id = azuredevops_serviceendpoint_kubernetes.test.id
  name                = "%s"
  namespace           = "test"
  cluster_name        = "k8scluster"
  tags                = ["tag1", "tag2"]
}
`, projectName, environmentName, serviceEndpointName, k8sName)
}
