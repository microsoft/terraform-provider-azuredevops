# azuredevops_serviceendpoint_kubernetes
Manages a Kubernetes service endpoint within Azure DevOps.

## Example Usage

```hcl
data "azuredevops_project" "p" {
  project_name = "contoso"
}

resource "azuredevops_serviceendpoint_kubernetes" "se_azure_sub" {
  project_id            = data.azuredevops_project.p
  service_endpoint_name = "Sample Kubernetes"
  apiserver_url         = "https://sample-kubernetes-cluster.hcp.westeurope.azmk8s.io"
  authorization_type    = "AzureSubscription"

  azure_subscription {
    subscription_id   = "001ac454-bb17-475a-8648-82c4234545be" # fake value
    subscription_name = "Microsoft Azure DEMO"
    tenant_id         = "8c46c3eb-ca1f-4a0b-8dfa-7c3baaf69d45" # fake value
    resourcegroup_id  = "sample-rg"
    namespace         = "default"
    cluster_name      = "sample-aks"
  }
}

resource "azuredevops_serviceendpoint_kubernetes" "se_kubeconfig" {
  project_id            = data.azuredevops_project.p
  service_endpoint_name = "Sample Kubernetes"
  apiserver_url         = "https://sample-kubernetes-cluster.hcp.westeurope.azmk8s.io"
  authorization_type    = "Kubeconfig"

  kubeconfig {
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
    accept_untrusted_certs = true
    cluster_context        = "dev-frontend"
  }
}

resource "azuredevops_serviceendpoint_kubernetes" "se_service_account" {
  project_id            = data.azuredevops_project.p
  service_endpoint_name = "Sample Kubernetes"
  apiserver_url         = "https://sample-kubernetes-cluster.hcp.westeurope.azmk8s.io"
  authorization_type    = "ServiceAccount"

  service_account {
    token   = "bXktYXBw[...]K8bPxc2uQ=="
    ca_cert = "Mzk1MjgkdmRnN0pi[...]mHHRUH14gw4Q=="
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The project ID or project name.
* `service_endpoint_name` - (Required) The Service Endpoint name.
* `apiserver_url` - (Required) The Service Endpoint description.
* `authorization_type` - (Required) The authentication method used to authenticate on the Kubernetes cluster. The value should be one of AzureSubscription, Kubeconfig, ServiceAccount.
* `azure_subscription` - (Optional) The configuration for authorization_type="AzureSubscription".
  * `azure_environment` - (Optional) Azure environment refers to whether the public cloud offering or domestic (government) clouds are being used. Currently, only the public cloud is supported. The value must be AzureCloud. This is also the default-value.
  * `cluster_name` - (Required) The name of the Kubernetes cluster.
  * `subscription_id` - (Required) The id of the Azure subscription.
  * `subscription_name` - (Required) The name of the Azure subscription.
  * `tenant_id` - (Required) The id of the tenant used by the subscription.
  * `resourcegroup_id` - (Required) The resource group id, to which the Kubernetes cluster is deployed.
  * `namespace` - (Optional) The Kubernetes namespace. Default value is "default".
* `kubeconfig` - (Optional) The configuration for authorization_type="Kubeconfig".
  * `kube_config` - (Required) The content of the kubeconfig in yaml notation to be used to communicate with the API-Server of Kubernetes.
  * `accept_untrusted_certs` - (Optional) Set this option to allow clients to accept a self-signed certificate.
  * `cluster_context` - (Optional) Context within the kubeconfig file that is to be used for identifying the cluster. Default value is the current-context set in kubeconfig.
* `service_account` - (Optional) The configuration for authorization_type="ServiceAccount". This type uses the credentials of a service account currently deployed to the cluster.
  * `token` - (Required) The token from a Kubernetes secret object.
  * `ca_cert` - (Required) The certificate from a Kubernetes secret object.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The project ID or project name.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links
* [Azure DevOps Service REST API 5.1 - Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-5.1)
