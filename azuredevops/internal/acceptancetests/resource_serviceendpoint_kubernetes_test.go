//go:build (all || resource_serviceendpoint_kubernetes) && !exclude_serviceendpoints

package acceptancetests

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
)

func TestAccServiceEndpointKubernetes_azureSubscription(t *testing.T) {
	if os.Getenv("KUBE_ARM_SUBSCRIPTION_ID") == "" || os.Getenv("KUBE_ARM_SUBSCRIPTION_NAME") == "" ||
		os.Getenv("KUBE_ARM_TENANT_ID") == "" || os.Getenv("KUBE_ARM_RESOURCE_GROUP") == "" ||
		os.Getenv("KUBE_ARM_KUBE_NAMESPACE") == "" || os.Getenv("KUBE_ARM_KUBE_CLUSTER_NAME") == "" ||
		os.Getenv("KUBE_ARM_KUBE_API_URL") == "" {
		t.Skip("Skipping tests due to missing Kubernetes Resource")
	}

	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	tfSvcEpNode := "azuredevops_serviceendpoint_kubernetes.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testutils.PreCheck(t, &[]string{
				"KUBE_ARM_SUBSCRIPTION_ID",
				"KUBE_ARM_SUBSCRIPTION_NAME",
				"KUBE_ARM_TENANT_ID",
				"KUBE_ARM_RESOURCE_GROUP",
				"KUBE_ARM_KUBE_API_URL",
				"KUBE_ARM_KUBE_NAMESPACE",
				"KUBE_ARM_KUBE_CLUSTER_NAME",
			})
		},
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkSvcEndpointKubernetesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclServiceEndpointKubernetesAzureSubscriptionResource(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					checkSvcEndpointKubernetesExists(serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azure_subscription.#"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "authorization_type", "AzureSubscription"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
				),
			}, {
				Config: hclServiceEndpointKubernetesAzureSubscriptionResource(projectName, serviceEndpointNameSecond),
				Check: resource.ComposeTestCheckFunc(
					checkSvcEndpointKubernetesExists(serviceEndpointNameSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "azure_subscription.#"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "authorization_type", "AzureSubscription"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
				),
			},
		},
	})
}

func TestAccServiceEndpointKubernetes_serviceAccount(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	tfSvcEpNode := "azuredevops_serviceendpoint_kubernetes.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkSvcEndpointKubernetesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclServiceEndpointKubernetesServiceAccount(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					checkSvcEndpointKubernetesExists(serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "service_account.#"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "service_account.0.accept_untrusted_certs"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "authorization_type", "ServiceAccount"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
				),
			}, {
				Config: hclServiceEndpointKubernetesServiceAccount(projectName, serviceEndpointNameSecond),
				Check: resource.ComposeTestCheckFunc(
					checkSvcEndpointKubernetesExists(serviceEndpointNameSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "service_account.#"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "authorization_type", "ServiceAccount"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
				),
			},
		},
	})
}

func TestAccServiceEndpointKubernetes_kubeConfig(t *testing.T) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	tfSvcEpNode := "azuredevops_serviceendpoint_kubernetes.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testutils.PreCheck(t, nil) },
		ProviderFactories: testutils.GetProviderFactories(),
		CheckDestroy:      checkSvcEndpointKubernetesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: hclServiceEndpointKubernetesKubeConfig(projectName, serviceEndpointNameFirst),
				Check: resource.ComposeTestCheckFunc(
					checkSvcEndpointKubernetesExists(serviceEndpointNameFirst),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "kubeconfig.#"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "authorization_type", "Kubeconfig"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
				),
			}, {
				Config: hclServiceEndpointKubernetesKubeConfig(projectName, serviceEndpointNameSecond),
				Check: resource.ComposeTestCheckFunc(
					checkSvcEndpointKubernetesExists(serviceEndpointNameSecond),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "project_id"),
					resource.TestCheckResourceAttrSet(tfSvcEpNode, "kubeconfig.#"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "authorization_type", "Kubeconfig"),
					resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
				),
			},
		},
	})
}

func checkSvcEndpointKubernetesExists(expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		serviceEndpointDef, ok := s.RootModule().Resources["azuredevops_serviceendpoint_kubernetes.test"]
		if !ok {
			return fmt.Errorf("Did not find a service endpoint in the TF state")
		}

		serviceEndpoint, err := getServiceEndpointKubernetesFromResource(serviceEndpointDef)
		if err != nil {
			return err
		}

		if *serviceEndpoint.Name != expectedName {
			return fmt.Errorf("Service Endpoint has Name=%s, but expected Name=%s", *serviceEndpoint.Name, expectedName)
		}

		return nil
	}
}

func checkSvcEndpointKubernetesDestroyed(s *terraform.State) error {
	for _, res := range s.RootModule().Resources {
		if res.Type != "azuredevops_serviceendpoint_kubernetes" {
			continue
		}

		if _, err := getServiceEndpointKubernetesFromResource(res); err == nil {
			return fmt.Errorf("Unexpectedly found a service endpoint that should be deleted")
		}
	}

	return nil
}

func getServiceEndpointKubernetesFromResource(resource *terraform.ResourceState) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpointDefID, err := uuid.Parse(resource.Primary.ID)
	if err != nil {
		return nil, err
	}

	projectID := resource.Primary.Attributes["project_id"]
	clients := testutils.GetProvider().Meta().(*client.AggregatedClient)
	return clients.ServiceEndpointClient.GetServiceEndpointDetails(clients.Ctx, serviceendpoint.GetServiceEndpointDetailsArgs{
		Project:    &projectID,
		EndpointId: &serviceEndpointDefID,
	})
}

func hclServiceEndpointKubernetesAzureSubscriptionResource(projectName, serviceEndpointName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_kubernetes" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  authorization_type    = "AzureSubscription"
  apiserver_url         = "%s"
  azure_subscription {
    subscription_id   = "%s"
    subscription_name = "%s"
    tenant_id         = "%s"
    resourcegroup_id  = "%s"
    namespace         = "%s"
    cluster_name      = "%s"
  }
}
`, projectName, serviceEndpointName,
		os.Getenv("KUBE_ARM_KUBE_API_URL"),
		os.Getenv("KUBE_ARM_SUBSCRIPTION_ID"),
		os.Getenv("KUBE_ARM_SUBSCRIPTION_NAME"),
		os.Getenv("KUBE_ARM_TENANT_ID"),
		os.Getenv("KUBE_ARM_RESOURCE_GROUP"),
		os.Getenv("KUBE_ARM_KUBE_NAMESPACE"),
		os.Getenv("KUBE_ARM_KUBE_CLUSTER_NAME"))
}

func hclServiceEndpointKubernetesServiceAccount(projectName, serviceEndpointName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_kubernetes" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  apiserver_url         = "https://sample-kubernetes-cluster.hcp.westeurope.azmk8s.io"
  authorization_type    = "ServiceAccount"
  service_account {
    accept_untrusted_certs = false
    token                  = "kubernetes_TEST_api_token"
    ca_cert                = "kubernetes_TEST_ca_cert"
  }
}
`, projectName, serviceEndpointName)
}

func hclServiceEndpointKubernetesKubeConfig(projectName, serviceEndpointName string) string {
	return fmt.Sprintf(`
resource "azuredevops_project" "test" {
  name = "%s"
}

resource "azuredevops_serviceendpoint_kubernetes" "test" {
  project_id            = azuredevops_project.test.id
  service_endpoint_name = "%s"
  apiserver_url         = "https://sample-kubernetes-cluster.hcp.westeurope.azmk8s.io"
  authorization_type    = "Kubeconfig"
  kubeconfig {
    accept_untrusted_certs = true
    cluster_context        = "dev-frontend"
    kube_config            = <<EOT
								apiVersion: v1
								clusters:
								- cluster:
									certificate-authority: fake-ca-file
									server: https://1.2.3.4
								name: development
								contexts:
								- context:
									cluster: development
									namespace: frontend
									user: developer
								name: dev-frontend
								current-context: dev-frontend
								kind: Config
								preferences: {}
								users:
								- name: developer
								user:
									client-certificate: fake-cert-file
									client-key: fake-key-file
								EOT
  }
}
`, projectName, serviceEndpointName)
}
