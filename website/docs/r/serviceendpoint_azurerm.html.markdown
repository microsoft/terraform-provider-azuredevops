# azuredevops_serviceendpoint_azurerm
Manages a AzureRM service endpoint within Azure DevOps.

## Requirements
Before to create a service end point in Azure DevOps, you need to create a Service Principal in your Azure subscription.

For detailled steps to create a service principal with Azure cli see the [documentation](https://docs.microsoft.com/en-us/cli/azure/create-an-azure-service-principal-azure-cli?view=azure-cli-latest)

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  project_name       = "Sample Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_serviceendpoint_azurerm" "endpointazure" {
  project_id                = azuredevops_project.project.id
  service_endpoint_name     = "TestServiceRM"
  azurerm_spn_clientid      = "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxx"
  azurerm_spn_clientsecret  = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  azurerm_spn_tenantid      = "xxxxxxx-xxxx-xxx-xxxxx-xxxxxxxx"
  azurerm_subscription_id   = "xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxx"
  azurerm_subscription_name = "Microsoft Azure DEMO"
  azurerm_scope             = "/subscriptions/xxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The project ID or project name.
* `service_endpoint_name` - (Required) The Service Endpoint name.
* `azurerm_spn_clientid` - (Required) The service principal application Id
* `azurerm_spn_clientsecret` - (Required) The service principal secret.
* `azurerm_spn_tenantid` - (Required) The tenant id if the service principal.
* `azurerm_subscription_id` - (Required) The subscription Id of the Azure targets.
* `azurerm_subscription_name` - (Required) The subscription Name of the targets.
* `azurerm_scope` - (Required) The Azure scope of the end point (ID of the subscription or resource group).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.
* `project_id` - The project ID or project name.
* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links
* [Azure DevOps Service REST API 5.1 - Service End points](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-5.1)
