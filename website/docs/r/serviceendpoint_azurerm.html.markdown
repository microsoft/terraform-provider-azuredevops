---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_azurerm"
description: |-
  Manages a AzureRM service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_azurerm

Manages Manual or Automatic AzureRM service endpoint within Azure DevOps.

## Requirements (Manual AzureRM Service Endpoint)

Before to create a service end point in Azure DevOps, you need to create a Service Principal in your Azure subscription.

For detailed steps to create a service principal with Azure cli see the [documentation](https://docs.microsoft.com/en-us/cli/azure/create-an-azure-service-principal-azure-cli?view=azure-cli-latest)

## Example Usage

### Manual AzureRM Service Endpoint

```hcl
resource "azuredevops_project" "project" {
  name       = "Sample Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_azurerm" "endpointazure" {
  project_id            = azuredevops_project.project.id
  service_endpoint_name = "Sample AzureRM"
  description = "Managed by Terraform" 
  credentials {
    serviceprincipalid  = "00000000-0000-0000-0000-000000000000"
    serviceprincipalkey = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  }
  azurerm_spn_tenantid      = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_id   = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_name = "Sample Subscription"
}
```

### Automatic AzureRM Service Endpoint

```hcl
resource "azuredevops_project" "project" {
  name       = "Sample Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_azurerm" "endpointazure" {
  project_id                = azuredevops_project.project.id
  service_endpoint_name     = "Sample AzureRM"
  description = "Managed by Terraform" 
  azurerm_spn_tenantid      = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_id   = "00000000-0000-0000-0000-000000000000"
  azurerm_subscription_name = "Microsoft Azure DEMO"
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The project ID or project name.
- `service_endpoint_name` - (Required) The Service Endpoint name.
- `azurerm_spn_tenantid` - (Required) The tenant id if the service principal.
- `azurerm_subscription_id` - (Required) The subscription Id of the Azure targets.
- `azurerm_subscription_name` - (Required) The subscription Name of the targets.
- `description` - (Optional) Service connection description.
- `credentials` - (Optional) A `credentials` block.
- `resource_group` - (Optional) The resource group used for scope of automatic service endpoint.

---

A `credentials` block supports the following:

- `serviceprincipalid` - (Required) The service principal application Id
- `serviceprincipalkey` - (Required) The service principal secret.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the service endpoint.
- `project_id` - The project ID or project name.
- `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 6.0 - Service End points](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-6.0)

## Import

Azure DevOps Service Endpoint Azure Resource Manage can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
terraform import azuredevops_serviceendpoint_azurerm.serviceendpoint 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
