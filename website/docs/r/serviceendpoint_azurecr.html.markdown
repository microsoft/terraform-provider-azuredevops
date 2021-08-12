---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_azurecr"
description: |-
  Manages a Azure Container Registry service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_azurecr

Manages a Azure Container Registry service endpoint within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  name       = "Sample Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

# azure container registry service connection
resource "azuredevops_serviceendpoint_azurecr" "azurecr" {
  project_id             = azuredevops_project.project.id
  service_endpoint_name  = "Sample AzureCR"
  resource_group            = "sample-rg"
  azurecr_spn_tenantid      = "00000000-0000-0000-0000-000000000000"
  azurecr_name              = "sampleAcr"
  azurecr_subscription_id   = "00000000-0000-0000-0000-000000000000"
  azurecr_subscription_name = "sampleSub"
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The project ID or project name.
- `service_endpoint_name` - (Required) The name you will use to refer to this service connection in task inputs.
- `resource_group` - (Required) The resource group to which the container registry belongs.
- `azurecr_spn_tenantid` - (Required) The tenant id of the service principal.
- `azurecr_name` - (Required) The Azure container registry name.
- `azurecr_subscription_id` - (Required) The subscription id of the Azure targets.
- `azurecr_subscription_name` - (Required) The subscription name of the Azure targets.
- `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the service endpoint.
- `project_id` - The project ID or project name.
- `service_endpoint_name` - The Service Endpoint name.
- `service_principal_id` - The service principal ID.

## Relevant Links

- [Azure DevOps Service REST API 5.1 - Service Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-5.1)
- [Azure Container Registry REST API](https://docs.microsoft.com/en-us/rest/api/containerregistry/)

## Import

Azure DevOps Service Endpoint Azure Container Registry can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
$ terraform import azuredevops_serviceendpoint_azurecr.serviceendpoint 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```