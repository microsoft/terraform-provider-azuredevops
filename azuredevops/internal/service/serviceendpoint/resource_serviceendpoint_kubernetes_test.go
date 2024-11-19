//go:build (all || resource_serviceendpoint_kubernetes) && !exclude_serviceendpoints
// +build all resource_serviceendpoint_kubernetes
// +build !exclude_serviceendpoints

package serviceendpoint

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/v7/serviceendpoint"
	"github.com/microsoft/terraform-provider-azuredevops/azdosdkmocks"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/microsoft/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
	"github.com/stretchr/testify/require"
)

const errMsgCreateServiceEndpoint = "CreateServiceEndpoint() Failed"
const errMsgUpdateServiceEndpoint = "UpdateServiceEndpoint() Failed"
const errMsgGetServiceEndpoint = "GetServiceEndpoint() Failed"
const errMsgDeleteServiceEndpoint = "DeleteServiceEndpoint() Failed"

var kubernetesTestServiceEndpointID = uuid.New()
var kubernetesRandomServiceEndpointProjectID = uuid.New()
var kubernetesTestServiceEndpointProjectID = &kubernetesRandomServiceEndpointProjectID

var kubernetesTestServiceEndpoint = serviceendpoint.ServiceEndpoint{
	Authorization: &serviceendpoint.EndpointAuthorization{},
	Id:            &kubernetesTestServiceEndpointID,
	Name:          converter.String("UNIT_TEST_CONN_NAME"),
	Owner:         converter.String("library"), // Supported values are "library", "agentcloud"
	Type:          converter.String("kubernetes"),
	Url:           converter.String("https://kubernetes.apiserver.com/"),
	Description:   converter.String("description"),
	ServiceEndpointProjectReferences: &[]serviceendpoint.ServiceEndpointProjectReference{
		{
			ProjectReference: &serviceendpoint.ProjectReference{
				Id: kubernetesTestServiceEndpointProjectID,
			},
			Name:        converter.String("UNIT_TEST_CONN_NAME"),
			Description: converter.String("description"),
		},
	},
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
		"clusterAdmin":          "false",
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
		"acceptUntrustedCerts": "false",
		"authorizationType":    "ServiceAccount",
	}

	return &serviceEndpoint
}

// verifies that the flatten/expand round trip yields the same service endpoint for autorization type "AzureSubscription"
func TestServiceEndpointKubernetesForAzureSubscriptionExpandFlattenRoundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointKubernetes().Schema, nil)
	kubernetesTestServiceEndpointForAzureSubscription := createkubernetesTestServiceEndpointForAzureSubscription()
	resourceData.Set("project_id", (*kubernetesTestServiceEndpointForAzureSubscription.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	doBaseFlattening(resourceData, kubernetesTestServiceEndpointForAzureSubscription)
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForAzureSubscription)

	serviceEndpointAfterRoundTrip, err := expandServiceEndpointKubernetes(resourceData)

	require.Nil(t, err)
	require.Equal(t, *kubernetesTestServiceEndpointForAzureSubscription, *serviceEndpointAfterRoundTrip)
	require.Equal(t, kubernetesTestServiceEndpointProjectID, (*serviceEndpointAfterRoundTrip.ServiceEndpointProjectReferences)[0].ProjectReference.Id)
}

// verifies that if an error is produced on create, the error is not swallowed
func TestServiceEndpointKubernetesForAzureSubscriptionCreateDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	kubernetesTestServiceEndpointForAzureSubscription := createkubernetesTestServiceEndpointForAzureSubscription()
	resourceData.Set("project_id", (*kubernetesTestServiceEndpointForAzureSubscription.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	doBaseFlattening(resourceData, kubernetesTestServiceEndpointForAzureSubscription)
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForAzureSubscription)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: kubernetesTestServiceEndpointForAzureSubscription}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New(errMsgCreateServiceEndpoint)).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), errMsgCreateServiceEndpoint)
}

// verifies that if an error is produced on a read, it is not swallowed
func TestServiceEndpointKubernetesForAzureSubscriptionReadDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	kubernetesTestServiceEndpointForAzureSubscription := createkubernetesTestServiceEndpointForAzureSubscription()
	resourceData.Set("project_id", (*kubernetesTestServiceEndpointForAzureSubscription.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	doBaseFlattening(resourceData, kubernetesTestServiceEndpointForAzureSubscription)
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForAzureSubscription)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{
		EndpointId: kubernetesTestServiceEndpointForAzureSubscription.Id,
		Project:    converter.String(kubernetesTestServiceEndpointProjectID.String()),
	}
	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, expectedArgs).
		Return(nil, errors.New(errMsgGetServiceEndpoint)).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), errMsgGetServiceEndpoint)
}

// verifies that if an error is produced on a delete, it is not swallowed
func TestServiceEndpointKubernetesForAzureSubscriptionDeleteDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	kubernetesTestServiceEndpointForAzureSubscription := createkubernetesTestServiceEndpointForAzureSubscription()
	resourceData.Set("project_id", (*kubernetesTestServiceEndpointForAzureSubscription.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	doBaseFlattening(resourceData, kubernetesTestServiceEndpointForAzureSubscription)
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForAzureSubscription)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{
		EndpointId: kubernetesTestServiceEndpointForAzureSubscription.Id,
		ProjectIds: &[]string{
			kubernetesTestServiceEndpointProjectID.String(),
		},
	}

	buildClient.
		EXPECT().
		DeleteServiceEndpoint(clients.Ctx, expectedArgs).
		Return(errors.New(errMsgDeleteServiceEndpoint)).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Contains(t, err.Error(), errMsgDeleteServiceEndpoint)
}

// verifies that if an error is produced on an update, it is not swallowed
func TestServiceEndpointKubernetesForAzureSubscriptionUpdateDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	kubernetesTestServiceEndpointForAzureSubscription := createkubernetesTestServiceEndpointForAzureSubscription()
	resourceData.Set("project_id", (*kubernetesTestServiceEndpointForAzureSubscription.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	doBaseFlattening(resourceData, kubernetesTestServiceEndpointForAzureSubscription)
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForAzureSubscription)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   kubernetesTestServiceEndpointForAzureSubscription,
		EndpointId: kubernetesTestServiceEndpointForAzureSubscription.Id,
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
func TestServiceEndpointKubernetesForKubeconfigExpandFlattenRoundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointKubernetes().Schema, nil)
	configureKubeconfig(resourceData)
	kubernetesTestServiceEndpointForKubeconfig := createkubernetesTestServiceEndpointForKubeconfig()
	resourceData.Set("project_id", (*kubernetesTestServiceEndpointForKubeconfig.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	doBaseFlattening(resourceData, kubernetesTestServiceEndpointForKubeconfig)
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForKubeconfig)

	serviceEndpointAfterRoundTrip, err := expandServiceEndpointKubernetes(resourceData)
	require.Nil(t, err)
	require.Equal(t, *kubernetesTestServiceEndpointForKubeconfig, *serviceEndpointAfterRoundTrip)
	require.Equal(t, kubernetesTestServiceEndpointProjectID, (*serviceEndpointAfterRoundTrip.ServiceEndpointProjectReferences)[0].ProjectReference.Id)
}

// verifies that if an error is produced on a read, it is not swallowed
func TestServiceEndpointKubernetesForKubeconfigCreateDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	configureKubeconfig(resourceData)
	kubernetesTestServiceEndpointForKubeconfig := createkubernetesTestServiceEndpointForKubeconfig()
	resourceData.Set("project_id", (*kubernetesTestServiceEndpointForKubeconfig.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	doBaseFlattening(resourceData, kubernetesTestServiceEndpointForKubeconfig)
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForKubeconfig)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: kubernetesTestServiceEndpointForKubeconfig}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New(errMsgCreateServiceEndpoint)).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), errMsgCreateServiceEndpoint)
}

// verifies that if an error is produced on a read, it is not swallowed
func TestServiceEndpointKubernetesForKubeconfigReadDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	configureKubeconfig(resourceData)
	kubernetesTestServiceEndpointForKubeconfig := createkubernetesTestServiceEndpointForKubeconfig()
	resourceData.Set("project_id", (*kubernetesTestServiceEndpointForKubeconfig.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	doBaseFlattening(resourceData, kubernetesTestServiceEndpointForKubeconfig)
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForKubeconfig)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{
		EndpointId: kubernetesTestServiceEndpointForKubeconfig.Id,
		Project:    converter.String(kubernetesTestServiceEndpointProjectID.String()),
	}
	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, expectedArgs).
		Return(nil, errors.New(errMsgGetServiceEndpoint)).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), errMsgGetServiceEndpoint)
}

// verifies that if an error is produced on a delete, it is not swallowed
func TestServiceEndpointKubernetesForKubeconfigDeleteDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	configureKubeconfig(resourceData)
	kubernetesTestServiceEndpointForKubeconfig := createkubernetesTestServiceEndpointForKubeconfig()
	resourceData.Set("project_id", (*kubernetesTestServiceEndpointForKubeconfig.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	doBaseFlattening(resourceData, kubernetesTestServiceEndpointForKubeconfig)
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForKubeconfig)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{
		EndpointId: kubernetesTestServiceEndpointForKubeconfig.Id,
		ProjectIds: &[]string{
			kubernetesTestServiceEndpointProjectID.String(),
		},
	}
	buildClient.
		EXPECT().
		DeleteServiceEndpoint(clients.Ctx, expectedArgs).
		Return(errors.New(errMsgDeleteServiceEndpoint)).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Contains(t, err.Error(), errMsgDeleteServiceEndpoint)
}

// verifies that if an error is produced on an update, it is not swallowed
func TestServiceEndpointKubernetesForKubeconfigUpdateDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	configureKubeconfig(resourceData)
	kubernetesTestServiceEndpointForKubeconfig := createkubernetesTestServiceEndpointForKubeconfig()
	resourceData.Set("project_id", (*kubernetesTestServiceEndpointForKubeconfig.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	doBaseFlattening(resourceData, kubernetesTestServiceEndpointForKubeconfig)
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForKubeconfig)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   kubernetesTestServiceEndpointForKubeconfig,
		EndpointId: kubernetesTestServiceEndpointForKubeconfig.Id,
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
func TestServiceEndpointKubernetesForServiceAccountExpandFlattenRoundtrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceServiceEndpointKubernetes().Schema, nil)
	configureServiceAccount(resourceData)
	kubernetesTestServiceEndpointForServiceAccount := createkubernetesTestServiceEndpointForServiceAccount()
	resourceData.Set("project_id", (*kubernetesTestServiceEndpointForServiceAccount.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	doBaseFlattening(resourceData, kubernetesTestServiceEndpointForServiceAccount)
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForServiceAccount)

	serviceEndpointAfterRoundTrip, err := expandServiceEndpointKubernetes(resourceData)

	require.Nil(t, err)
	require.Equal(t, *kubernetesTestServiceEndpointForServiceAccount, *serviceEndpointAfterRoundTrip)
	require.Equal(t, kubernetesTestServiceEndpointProjectID, (*serviceEndpointAfterRoundTrip.ServiceEndpointProjectReferences)[0].ProjectReference.Id)
}

// verifies that if an error is produced on a read, it is not swallowed
func TestServiceEndpointKubernetesForServiceAccountCreateDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	configureServiceAccount(resourceData)
	kubernetesTestServiceEndpointForServiceAccount := createkubernetesTestServiceEndpointForServiceAccount()
	resourceData.Set("project_id", (*kubernetesTestServiceEndpointForServiceAccount.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	doBaseFlattening(resourceData, kubernetesTestServiceEndpointForServiceAccount)
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForServiceAccount)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.CreateServiceEndpointArgs{Endpoint: kubernetesTestServiceEndpointForServiceAccount}
	buildClient.
		EXPECT().
		CreateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New(errMsgCreateServiceEndpoint)).
		Times(1)

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), errMsgCreateServiceEndpoint)
}

// verifies that if an error is produced on a read, it is not swallowed
func TestServiceEndpointKubernetesForServiceAccountReadDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	configureServiceAccount(resourceData)
	kubernetesTestServiceEndpointForServiceAccount := createkubernetesTestServiceEndpointForServiceAccount()
	resourceData.Set("project_id", (*kubernetesTestServiceEndpointForServiceAccount.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	doBaseFlattening(resourceData, kubernetesTestServiceEndpointForServiceAccount)
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForServiceAccount)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.GetServiceEndpointDetailsArgs{
		EndpointId: kubernetesTestServiceEndpointForServiceAccount.Id,
		Project:    converter.String(kubernetesTestServiceEndpointProjectID.String()),
	}
	buildClient.
		EXPECT().
		GetServiceEndpointDetails(clients.Ctx, expectedArgs).
		Return(nil, errors.New(errMsgGetServiceEndpoint)).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), errMsgGetServiceEndpoint)
}

// verifies that if an error is produced on a delete, it is not swallowed
func TestServiceEndpointKubernetesForServiceAccountDeleteDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	configureServiceAccount(resourceData)
	kubernetesTestServiceEndpointForServiceAccount := createkubernetesTestServiceEndpointForServiceAccount()
	resourceData.Set("project_id", (*kubernetesTestServiceEndpointForServiceAccount.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	doBaseFlattening(resourceData, kubernetesTestServiceEndpointForServiceAccount)
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForServiceAccount)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.DeleteServiceEndpointArgs{
		EndpointId: kubernetesTestServiceEndpointForServiceAccount.Id,
		ProjectIds: &[]string{
			kubernetesTestServiceEndpointProjectID.String(),
		},
	}
	buildClient.
		EXPECT().
		DeleteServiceEndpoint(clients.Ctx, expectedArgs).
		Return(errors.New(errMsgDeleteServiceEndpoint)).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Contains(t, err.Error(), errMsgDeleteServiceEndpoint)
}

// verifies that if an error is produced on an update, it is not swallowed
func TestServiceEndpointKubernetesForServiceAccountUpdateDoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceServiceEndpointKubernetes()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	configureServiceAccount(resourceData)
	kubernetesTestServiceEndpointForServiceAccount := createkubernetesTestServiceEndpointForServiceAccount()
	resourceData.Set("project_id", (*kubernetesTestServiceEndpointForServiceAccount.ServiceEndpointProjectReferences)[0].ProjectReference.Id.String())
	doBaseFlattening(resourceData, kubernetesTestServiceEndpointForServiceAccount)
	flattenServiceEndpointKubernetes(resourceData, kubernetesTestServiceEndpointForServiceAccount)

	buildClient := azdosdkmocks.NewMockServiceendpointClient(ctrl)
	clients := &client.AggregatedClient{ServiceEndpointClient: buildClient, Ctx: context.Background()}

	expectedArgs := serviceendpoint.UpdateServiceEndpointArgs{
		Endpoint:   kubernetesTestServiceEndpointForServiceAccount,
		EndpointId: kubernetesTestServiceEndpointForServiceAccount.Id,
	}

	buildClient.
		EXPECT().
		UpdateServiceEndpoint(clients.Ctx, expectedArgs).
		Return(nil, errors.New(errMsgUpdateServiceEndpoint)).
		Times(1)

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), errMsgUpdateServiceEndpoint)
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
