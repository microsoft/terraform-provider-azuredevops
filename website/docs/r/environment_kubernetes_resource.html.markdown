---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_environment_resource_kubernetes"
description: |-
  Manages a Kubernetes Resource for an Environment.
---

# azuredevops_environment_resource_kubernetes

Manages a Kubernetes Resource for an Environment.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

resource "azuredevops_environment" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Environment"
}

resource "azuredevops_serviceendpoint_kubernetes" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example Kubernetes"
  apiserver_url         = "https://sample-kubernetes-cluster.hcp.westeurope.azmk8s.io"
  authorization_type    = "AzureSubscription"

  azure_subscription {
    subscription_id   = "00000000-0000-0000-0000-000000000000"
    subscription_name = "Example"
    tenant_id         = "00000000-0000-0000-0000-000000000000"
    resourcegroup_id  = "example-rg"
    namespace         = "default"
    cluster_name      = "example-aks"
  }
}

resource "azuredevops_environment_resource_kubernetes" "example" {
  project_id          = azuredevops_project.example.id
  environment_id      = azuredevops_environment.example.id
  service_endpoint_id = azuredevops_serviceendpoint_kubernetes.example.id

  name         = "Example"
  namespace    = "default"
  cluster_name = "example-aks"
  tags         = ["tag1", "tag2"]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name for the Kubernetes Resource.

* `namespace` - (Required) The namespace for the Kubernetes Resource.

* `project_id` - (Required) The ID of the project.

* `environment_id` - (Required) The ID of the environment under which to create the Kubernetes Resource.

* `service_endpoint_id` - (Required) The ID of the service endpoint to associate with the Kubernetes Resource.

---

* `cluster_name` - (Optional) A cluster name for the Kubernetes Resource.

* `tags` - (Optional) A set of tags for the Kubernetes Resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Kubernetes Resource.

## Relevant Links

* [Azure DevOps Service REST API 6.0 - Kubernetes](https://docs.microsoft.com/en-us/rest/api/azure/devops/distributedtask/kubernetes?view=azure-devops-rest-6.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Environment Kubernetes Resource.
* `read` - (Defaults to 5 minute) Used when retrieving the Environment Kubernetes Resource.
* `delete` - (Defaults to 30 minutes) Used when deleting the Environment Kubernetes Resource.

## Import

The resource does not support import.
