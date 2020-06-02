// +build all resource_serviceendpoint_kubernetes
// +build !exclude_serviceendpoints

package azuredevops

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/config"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/converter"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/utils/testhelper"
	"github.com/stretchr/testify/require"

	"github.com/google/uuid"

	"github.com/microsoft/azure-devops-go-api/azuredevops/serviceendpoint"
)

const kubernetesApiserverUrl = "https://kubernetes.apiserver.com/"
const terraformServiceEndpointNode = "azuredevops_serviceendpoint_kubernetes.serviceendpoint"

const errMsgCreateServiceEndpoint = "CreateServiceEndpoint() Failed"
const errMsgUpdateServiceEndpoint = "UpdateServiceEndpoint() Failed"
const errMsgGetServiceEndpoint = "GetServiceEndpoint() Failed"
const errMsgDeleteServiceEndpoint = "DeleteServiceEndpoint() Failed"

var kubernetesTestServiceEndpointID = uuid.New()
var kubernetesRandomServiceEndpointProjectID = uuid.New().String()
var kubernetesTestServiceEndpointProjectID = &kubernetesRandomServiceEndpointProjectID

var kubernetesTestServiceEndpoint = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{},
	Id:            &kubernetesTestServiceEndpointID,
	Name:          converter.String("UNIT_TEST_CONN_NAME"),
	Owner:         converter.String("library"), // Supported values are "library", "agentcloud"
	Type:          converter.String("kubernetes"),
	Url:           converter.String(kubernetesApiserverUrl),
	Description:   converter.String("description"),
}

func createkubernetesTestServiceEndpointForAzureSubscription() *serviceendpoint.ServiceEndpoint {
	serviceEndpoint := kubernetesTestServiceEndpoint
	serviceEndpoint.Authorization.Scheme = converter.String("Kubernetes")
	serviceEndpoint.Authorization.Parameters = &map[string]string{
		"azureEnvironment": "AzureCloud",
		"azureTenantId":    "kubernetes_TEST_tenant_id",
	}
	serviceEndpoint.Data = &map[string]string{
		"authorizationType":     "AzureSubscription",
		"azureSubscriptionId":   "kubernetes_TEST_subscription_id",
		"azureSubscriptionName": "kubernetes_TEST_subscription_name",
		"clusterId":             "/subscriptions/kubernetes_TEST_subscription_id/resourcegroups/kubernetes_TEST_resource_group_id/providers/Microsoft.ContainerService/managedClusters/kubernetes_TEST_cluster_name",
		"namespace":             "default",
	}

	return &serviceEndpoint
}

func createkubernetesTestServiceEndpointForKubeconfig() *serviceendpoint.ServiceEndpoint {
	serviceEndpoint := kubernetesTestServiceEndpoint
	serviceEndpoint.Authorization.Scheme = converter.String("Kubernetes")
	serviceEndpoint.Authorization.Parameters = &map[string]string{
		"kubeconfig": `<<EOT
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
		"clusterContext": "dev-frontend",
	}
	serviceEndpoint.Data = &map[string]string{
		"authorizationType":    "Kubeconfig",
		"acceptUntrustedCerts": "true",
	}

	return &serviceEndpoint
}

func createkubernetesTestServiceEndpointForServiceAccount() *serviceendpoint.ServiceEndpoint {
	serviceEndpoint := kubernetesTestServiceEndpoint
	serviceEndpoint.Authorization.Scheme = converter.String("Token")
	serviceEndpoint.Authorization.Parameters = &map[string]string{
		"apiToken":                  "kubernetes_TEST_api_token",
		"serviceAccountCertificate": "kubernetes_TEST_ca_cert",
	}
	serviceEndpoint.Data = &map[string]string{
		"authorizationType": "ServiceAccount",
	}

	return &serviceEndpoint
}

/**
 * Begin unit tests
 */

// verifies that the flatten/expand round trip yields the same service endpoint for autorization type "AzureSubscription"
func TestAzureDevOpsServiceEndpointKubernetesForAzureSubscriptionExpandFlattenRoundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceServiceEndpointKubernetes().Schema, nil)
	kubernetesTestServiceEndpointForAzureSubscription := createkubernetesTestServiceEndpointForAzureSubscription()
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForAzureSubscription, kubernetesTestServiceEndpointProjectID)

	serviceEndpointAfterRoundTrip, projectID, err := expandServiceEndpointKubernetes(resourceData)

	require.Nil(t, err)
	require.Equal(t, *kubernetesTestServiceEndpointForAzureSubscription, *serviceEndpointAfterRoundTrip)
	require.Equal(t, kubernetesTestServiceEndpointProjectID, projectID)
}

// verifies that if an error is produced on create, the error is not swallowed
func TestAzureDevOpsServiceEndpointKubernetesForAzureSubscriptionCreateDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	kubernetesTestServiceEndpointForAzureSubscription := createkubernetesTestServiceEndpointForAzureSubscription()
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForAzureSubscription, kubernetesTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: kubernetesTestServiceEndpointForAzureSubscription, Project: kubernetesTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New(errMsgCreateServiceEndpoint)).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), errMsgCreateServiceEndpoint)
}

// verifies that if an error is produced on a read, it is not swallowed
func TestAzureDevOpsServiceEndpointKubernetesForAzureSubscriptionReadDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	kubernetesTestServiceEndpointForAzureSubscription := createkubernetesTestServiceEndpointForAzureSubscription()
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForAzureSubscription, kubernetesTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{EndpointId: kubernetesTestServiceEndpointForAzureSubscription.Id, Project: kubernetesTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, expectedArgs).
		Return(nil, errors.New(errMsgGetServiceEndpoint)).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), errMsgGetServiceEndpoint)
}

// verifies that if an error is produced on a delete, it is not swallowed
func TestAzureDevOpsServiceEndpointKubernetesForAzureSubscriptionDeleteDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	kubernetesTestServiceEndpointForAzureSubscription := createkubernetesTestServiceEndpointForAzureSubscription()
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForAzureSubscription, kubernetesTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{EndpointId: kubernetesTestServiceEndpointForAzureSubscription.Id, Project: kubernetesTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		DeleteServiceEndpoint(clients.Ctx, expectedArgs).
		Return(errors.New(errMsgDeleteServiceEndpoint)).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Contains(t, err.Error(), errMsgDeleteServiceEndpoint)
}

// verifies that if an error is produced on an update, it is not swallowed
func TestAzureDevOpsServiceEndpointKubernetesForAzureSubscriptionUpdateDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	kubernetesTestServiceEndpointForAzureSubscription := createkubernetesTestServiceEndpointForAzureSubscription()
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForAzureSubscription, kubernetesTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   kubernetesTestServiceEndpointForAzureSubscription,
		EndpointId: kubernetesTestServiceEndpointForAzureSubscription.Id,
		Project:    kubernetesTestServiceEndpointProjectID,
	}

	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New(errMsgUpdateServiceEndpoint)).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), errMsgUpdateServiceEndpoint)
}

// verifies that the flatten/expand round trip yields the same service endpoint for autorization type "Kubeconfig"
func TestAzureDevOpsServiceEndpointKubernetesForKubeconfigExpandFlattenRoundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceServiceEndpointKubernetes().Schema, nil)
	configureKubeconfig(resourceData)
	kubernetesTestServiceEndpointForKubeconfig := createkubernetesTestServiceEndpointForKubeconfig()
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForKubeconfig, kubernetesTestServiceEndpointProjectID)

	serviceEndpointAfterRoundTrip, projectID, err := expandServiceEndpointKubernetes(resourceData)
	require.Nil(t, err)
	require.Equal(t, *kubernetesTestServiceEndpointForKubeconfig, *serviceEndpointAfterRoundTrip)
	require.Equal(t, kubernetesTestServiceEndpointProjectID, projectID)
}

// verifies that if an error is produced on a read, it is not swallowed
func TestAzureDevOpsServiceEndpointKubernetesForKubeconfigCreateDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	configureKubeconfig(resourceData)
	kubernetesTestServiceEndpointForKubeconfig := createkubernetesTestServiceEndpointForKubeconfig()
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForKubeconfig, kubernetesTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: kubernetesTestServiceEndpointForKubeconfig, Project: kubernetesTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New(errMsgCreateServiceEndpoint)).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), errMsgCreateServiceEndpoint)
}

// verifies that if an error is produced on a read, it is not swallowed
func TestAzureDevOpsServiceEndpointKubernetesForKubeconfigReadDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	configureKubeconfig(resourceData)
	kubernetesTestServiceEndpointForKubeconfig := createkubernetesTestServiceEndpointForKubeconfig()
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForKubeconfig, kubernetesTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{EndpointId: kubernetesTestServiceEndpointForKubeconfig.Id, Project: kubernetesTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, expectedArgs).
		Return(nil, errors.New(errMsgGetServiceEndpoint)).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), errMsgGetServiceEndpoint)
}

// verifies that if an error is produced on a delete, it is not swallowed
func TestAzureDevOpsServiceEndpointKubernetesForKubeconfigDeleteDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	configureKubeconfig(resourceData)
	kubernetesTestServiceEndpointForKubeconfig := createkubernetesTestServiceEndpointForKubeconfig()
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForKubeconfig, kubernetesTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{EndpointId: kubernetesTestServiceEndpointForKubeconfig.Id, Project: kubernetesTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		DeleteServiceEndpoint(clients.Ctx, expectedArgs).
		Return(errors.New(errMsgDeleteServiceEndpoint)).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Contains(t, err.Error(), errMsgDeleteServiceEndpoint)
}

// verifies that if an error is produced on an update, it is not swallowed
func TestAzureDevOpsServiceEndpointKubernetesForKubeconfigUpdateDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	configureKubeconfig(resourceData)
	kubernetesTestServiceEndpointForKubeconfig := createkubernetesTestServiceEndpointForKubeconfig()
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForKubeconfig, kubernetesTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   kubernetesTestServiceEndpointForKubeconfig,
		EndpointId: kubernetesTestServiceEndpointForKubeconfig.Id,
		Project:    kubernetesTestServiceEndpointProjectID,
	}

	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New(errMsgUpdateServiceEndpoint)).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), errMsgUpdateServiceEndpoint)
}

// verifies that the flatten/expand round trip yields the same service endpoint for autorization type "ServiceAccount"
func TestAzureDevOpsServiceEndpointKubernetesForServiceAccountExpandFlattenRoundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, resourceServiceEndpointKubernetes().Schema, nil)
	configureServiceAccount(resourceData)
	kubernetesTestServiceEndpointForServiceAccount := createkubernetesTestServiceEndpointForServiceAccount()
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForServiceAccount, kubernetesTestServiceEndpointProjectID)

	serviceEndpointAfterRoundTrip, projectID, err := expandServiceEndpointKubernetes(resourceData)

	require.Nil(t, err)
	require.Equal(t, *kubernetesTestServiceEndpointForServiceAccount, *serviceEndpointAfterRoundTrip)
	require.Equal(t, kubernetesTestServiceEndpointProjectID, projectID)
}

// verifies that if an error is produced on a read, it is not swallowed
func TestAzureDevOpsServiceEndpointKubernetesForServiceAccountCreateDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	configureServiceAccount(resourceData)
	kubernetesTestServiceEndpointForServiceAccount := createkubernetesTestServiceEndpointForServiceAccount()
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForServiceAccount, kubernetesTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: kubernetesTestServiceEndpointForServiceAccount, Project: kubernetesTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New(errMsgCreateServiceEndpoint)).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), errMsgCreateServiceEndpoint)
}

// verifies that if an error is produced on a read, it is not swallowed
func TestAzureDevOpsServiceEndpointKubernetesForServiceAccountReadDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	configureServiceAccount(resourceData)
	kubernetesTestServiceEndpointForServiceAccount := createkubernetesTestServiceEndpointForServiceAccount()
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForServiceAccount, kubernetesTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{EndpointId: kubernetesTestServiceEndpointForServiceAccount.Id, Project: kubernetesTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, expectedArgs).
		Return(nil, errors.New(errMsgGetServiceEndpoint)).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), errMsgGetServiceEndpoint)
}

// verifies that if an error is produced on a delete, it is not swallowed
func TestAzureDevOpsServiceEndpointKubernetesForServiceAccountDeleteDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	configureServiceAccount(resourceData)
	kubernetesTestServiceEndpointForServiceAccount := createkubernetesTestServiceEndpointForServiceAccount()
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForServiceAccount, kubernetesTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{EndpointId: kubernetesTestServiceEndpointForServiceAccount.Id, Project: kubernetesTestServiceEndpointProjectID}
	buildClient.
		EXPECT().
		DeleteServiceEndpoint(clients.Ctx, expectedArgs).
		Return(errors.New(errMsgDeleteServiceEndpoint)).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Contains(t, err.Error(), errMsgDeleteServiceEndpoint)
}

// verifies that if an error is produced on an update, it is not swallowed
func TestAzureDevOpsServiceEndpointKubernetesForServiceAccountUpdateDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := resourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	configureServiceAccount(resourceData)
	kubernetesTestServiceEndpointForServiceAccount := createkubernetesTestServiceEndpointForServiceAccount()
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForServiceAccount, kubernetesTestServiceEndpointProjectID)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &config.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   kubernetesTestServiceEndpointForServiceAccount,
		EndpointId: kubernetesTestServiceEndpointForServiceAccount.Id,
		Project:    kubernetesTestServiceEndpointProjectID,
	}

	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New(errMsgUpdateServiceEndpoint)).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), errMsgUpdateServiceEndpoint)
}

/**
 * Begin acceptance tests
 */

// validates that an apply followed by another apply (i.e., resource update) will be reflected in AzDO and the
// underlying terraform state.
func TestAccAzureDevOpsServiceEndpointKubernetesForAzureSubscriptionCreateAndUpdate(t *testing.T) {
	authorizationType := "AzureSubscription"

	var attrTestCheckFuncList []resource.TestCheckFunc
	attrTestCheckFuncList = append(
		attrTestCheckFuncList,
		resource.TestCheckResourceAttrSet(terraformServiceEndpointNode, "azure_subscription.#"),
	)
	testAccAzureDevOpsServiceEndpoint(t, attrTestCheckFuncList, authorizationType)
}

// validates that an apply followed by another apply (i.e., resource update) will be reflected in AzDO and the
// underlying terraform state.
func TestAccAzureDevOpsServiceEndpointKubernetesForServiceAccountCreateAndUpdate(t *testing.T) {
	authorizationType := "ServiceAccount"

	var attrTestCheckFuncList []resource.TestCheckFunc
	attrTestCheckFuncList = append(
		attrTestCheckFuncList,
		resource.TestCheckResourceAttrSet(terraformServiceEndpointNode, "service_account.#"),
	)

	testAccAzureDevOpsServiceEndpoint(t, attrTestCheckFuncList, authorizationType)
}

// validates that an apply followed by another apply (i.e., resource update) will be reflected in AzDO and the
// underlying terraform state.
func TestAccAzureDevOpsServiceEndpointKubernetesForKubeconfigCreateAndUpdate(t *testing.T) {
	authorizationType := "Kubeconfig"

	var attrTestCheckFuncList []resource.TestCheckFunc
	attrTestCheckFuncList = append(
		attrTestCheckFuncList,
		resource.TestCheckResourceAttrSet(terraformServiceEndpointNode, "kubeconfig.#"),
	)
	testAccAzureDevOpsServiceEndpoint(t, attrTestCheckFuncList, authorizationType)
}

func testAccAzureDevOpsServiceEndpoint(t *testing.T, attrTestCheckFuncList []resource.TestCheckFunc, authorizationType string) {
	projectName := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	serviceEndpointNameFirst := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	serviceEndpointNameSecond := testhelper.TestAccResourcePrefix + acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	tfSvcEpNode := terraformServiceEndpointNode

	attrTestCheckFuncList = append(
		attrTestCheckFuncList,
		resource.TestCheckResourceAttrSet(terraformServiceEndpointNode, "project_id"),
		resource.TestCheckResourceAttr(terraformServiceEndpointNode, "authorization_type", authorizationType),
	)
	attrTestCheckFuncListNameFirst := append(
		attrTestCheckFuncList,
		resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameFirst),
		testAccCheckServiceEndpointKubernetesResourceExists(serviceEndpointNameFirst),
	)

	attrTestCheckFuncListNameSecond := append(
		attrTestCheckFuncList,
		resource.TestCheckResourceAttr(tfSvcEpNode, "service_endpoint_name", serviceEndpointNameSecond),
		testAccCheckServiceEndpointKubernetesResourceExists(serviceEndpointNameSecond),
	)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testhelper.TestAccPreCheck(t, nil) },
		Providers:    testAccProviders,
		CheckDestroy: testAccServiceEndpointKubernetesCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testhelper.TestAccServiceEndpointKubernetesResource(projectName, serviceEndpointNameFirst, authorizationType),
				Check:  resource.ComposeTestCheckFunc(attrTestCheckFuncListNameFirst...),
			}, {
				Config: testhelper.TestAccServiceEndpointKubernetesResource(projectName, serviceEndpointNameSecond, authorizationType),
				Check:  resource.ComposeTestCheckFunc(attrTestCheckFuncListNameSecond...),
			},
		},
	})
}

// Given the name of an AzDO service endpoint, this will return a function that will check whether
// or not the resource (1) exists in the state and (2) exist in AzDO and (3) has the correct name
func testAccCheckServiceEndpointKubernetesResourceExists(expectedName string) resource.TestCheckFunc {
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
func testAccServiceEndpointKubernetesCheckDestroy(s *terraform.State) error {
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
	clients := testAccProvider.Meta().(*config.AggregatedClient)
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

func init() {
	InitProvider()
}
