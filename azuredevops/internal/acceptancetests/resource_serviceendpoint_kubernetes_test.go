// +build all resource_serviceendpoint_kubernetes
// +build !exclude_serviceendpoints

package acceptancetests

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/acceptancetests/testutils"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
)

const terraformServiceEndpointNode = "azuredevops_serviceendpoint_kubernetes.serviceendpoint"

// validates that an apply followed by another apply (i.e., resource update) will be reflected in AzDO and the
// underlying terraform state.
func TestAccServiceEndpointKubernetesForAzureSubscriptionCreateAndUpdate(t *testing.T) {
	authorizationType := "AzureSubscription"

	var attrTestCheckFuncList []resource.TestCheckFunc
	attrTestCheckFuncList = append(
		attrTestCheckFuncList,
		resource.TestCheckResourceAttrSet(terraformServiceEndpointNode, "azure_subscription.#"),
	)
	runSvcEndpointAcceptanceTest(t, attrTestCheckFuncList, authorizationType)
}

// validates that an apply followed by another apply (i.e., resource update) will be reflected in AzDO and the
// underlying terraform state.
func TestAccServiceEndpointKubernetesForServiceAccountCreateAndUpdate(t *testing.T) {
	authorizationType := "ServiceAccount"

	var attrTestCheckFuncList []resource.TestCheckFunc
	attrTestCheckFuncList = append(
		attrTestCheckFuncList,
		resource.TestCheckResourceAttrSet(terraformServiceEndpointNode, "service_account.#"),
	)

	runSvcEndpointAcceptanceTest(t, attrTestCheckFuncList, authorizationType)
}

// validates that an apply followed by another apply (i.e., resource update) will be reflected in AzDO and the
// underlying terraform state.
func TestAccServiceEndpointKubernetesForKubeconfigCreateAndUpdate(t *testing.T) {
	authorizationType := "Kubeconfig"

	var attrTestCheckFuncList []resource.TestCheckFunc
	attrTestCheckFuncList = append(
		attrTestCheckFuncList,
		resource.TestCheckResourceAttrSet(terraformServiceEndpointNode, "kubeconfig.#"),
	)
	runSvcEndpointAcceptanceTest(t, attrTestCheckFuncList, authorizationType)
}

func runSvcEndpointAcceptanceTest(t *testing.T, attrTestCheckFuncList []resource.TestCheckFunc, authorizationType string) {
	projectName := testutils.GenerateResourceName()
	serviceEndpointNameFirst := testutils.GenerateResourceName()
	serviceEndpointNameSecond := testutils.GenerateResourceName()

	tfSvcEpNode := terraformServiceEndpointNode

	attrTestCheckFuncList = append(
		attrTestCheckFuncList,
		resource.TestCheckResourceAttrSet(terraformServiceEndpointNode, "project_id"),
		resource.TestCheckResourceAttr(terraformServiceEndpointNode, "authorization_type", authorizationType),
	)
	attrTestCheckFuncListNameFirst := append(
		attrTestCheckFuncList,
		resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
		checkSvcEndpointKubernetesExists(serviceEndpointNameFirst),
	)

	attrTestCheckFuncListNameSecond := append(
		attrTestCheckFuncList,
		resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
		checkSvcEndpointKubernetesExists(serviceEndpointNameSecond),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testutils.PreCheck(t, nil) },
		Providers:    testutils.GetProviders(),
		CheckDestroy: checkSvcEndpointKubernetesDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testutils.HclServiceEndpointKubernetesResource(projectName, serviceEndpointNameFirst, authorizationType),
				Check:  resource.ComposeTestCheckFunc(attrTestCheckFuncListNameFirst...),
			}, {
				Config: testutils.HclServiceEndpointKubernetesResource(projectName, serviceEndpointNameSecond, authorizationType),
				Check:  resource.ComposeTestCheckFunc(attrTestCheckFuncListNameSecond...),
			},
		},
	})
}

// Given the name of an AzDO service endpoint, this will return a function that will check whether
// or not the resource (1) exists in the state and (2) exist in AzDO and (3) has the correct name
func checkSvcEndpointKubernetesExists(expectedName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		serviceEndpointDef, ok := s.RootModule().Resources[terraformServiceEndpointNode]
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

// verifies that all service endpoints referenced in the state are destroyed. This will be invoked
// *after* terrafform destroys the resource but *before* the state is wiped clean.
func checkSvcEndpointKubernetesDestroyed(s *terraform.State) error {
	for _, resource := range s.RootModule().Resources {
		if resource.Type != "azuredevops_serviceendpoint_kubernetes" {
			continue
		}

		// indicates the service endpoint still exists - this should fail the test
		if _, err := getServiceEndpointKubernetesFromResource(resource); err == nil {
			return fmt.Errorf("Unexpectedly found a service endpoint that should be deleted")
		}
	}

	return nil
}

// given a resource from the state, return a service endpoint (and error)
func getServiceEndpointKubernetesFromResource(resource *terraform.ResourceState) (*serviceendpoint.ServiceEndpoint, error) {
	serviceEndpointDefID, err := uuid.Parse(resource.Primary.ID)
	if err != nil {
		return nil, err
	}

	projectID := resource.Primary.Attributes["project_id"]
	clients := testutils.GetProvider().Meta().(*config.AggregatedClient)
	return clients.ServiceEndpointClient.GetServiceEndpointDetails(clients.Ctx, serviceendpoint.GetServiceEndpointDetailsArgs{
		Project:    &projectID,
		EndpointId: &serviceEndpointDefID,
	})
}

func configureServiceAccount(d *schema.ResourceData) {
	d.Set("service_account", &[]map[string]interface{}{
		{
			"token":   "kubernetes_TEST_api_token",
			"ca_cert": "kubernetes_TEST_ca_cert",
		},
	})
}

func configureKubeconfig(d *schema.ResourceData) {
	d.Set("kubeconfig", &[]map[string]interface{}{
		{
			"kube_config": `<<EOT
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
							EOT`,
			"accept_untrusted_certs": true,
			"cluster_context":        "dev-frontend",
		},
	})
}
