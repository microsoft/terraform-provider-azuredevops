---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_argocd"
description: |-
  Manages a ArgoCD server endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_argocd
Manages a ArgoCD service endpoint within Azure DevOps. Using this service endpoint requires you to first install [Argo CD Extension](https://marketplace.visualstudio.com/items?itemName=scb-tomasmortensen.vsix-argocd).

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  name               = "Sample Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_argocd" "serviceendpoint" {

  project_id            = azuredevops_project.project.id
  service_endpoint_name = "Sample ArgoCD"
  description           = "Managed by Terraform"
  url                   = "https://argocd.my.com"
  authentication_token {
      token      = "0000000000000000000000000000000000000000"
  }
  description           = "Managed by Terraform"
}
```
Alternatively a username and password may be used.

```hcl
resource "azuredevops_serviceendpoint_argocd" "serviceendpoint" {

  project_id            = azuredevops_project.project.id
  service_endpoint_name = "Sample ArgoCD"
  description           = "Managed by Terraform"
  url                   = "https://argocd.my.com"
  authentication_basic {
      username              = "sampleuser"
      password              = "0000000000000000000000000000000000000000"
  }
}
```
## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The project ID or project name.
- `service_endpoint_name` - (Required) The Service Endpoint name.
- `url` - (Required) URL of the ArgoCD server to connect with.
- `description` - (Optional) The Service Endpoint description.
- `authentication_token` - An `authentication_token` block for the ArgoCD as documented below.
- `authentication_basic` - An `authentication_token` block for the ArgoCD as documented below.

~> **NOTE:** `authentication_basic` and `authentication_token` conflict with each other, only one is required.

---

A `authentication_token` block supports the following:

  - `token` - Authentication Token generated through ArgoCD.

A `authentication_basic` block supports the following:
  - `username` - ArgoCD Username. 
  - `password` - ArgoCD Password.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The project ID or project name.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links
* [Azure DevOps Service Connections](https://docs.microsoft.com/en-us/azure/devops/pipelines/library/service-endpoints?view=azure-devops&tabs=yaml)
* [ArgoCD Project/User Token](https://argo-cd.readthedocs.io/en/stable/user-guide/commands/argocd_account_generate-token/)
* [Argo CD Extension](https://marketplace.visualstudio.com/items?itemName=scb-tomasmortensen.vsix-argocd)

## Import
Azure DevOps Service Endpoint ArgoCD can be imported using the **projectID/serviceEndpointID**, e.g.


```shell
$ terraform import azuredevops_serviceendpoint_argocd.serviceendpoint 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
