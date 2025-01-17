---
layout: "azuredevops"
page_title: "AzureDevops: Data Source: azuredevops_serviceendpoint_azurerm"
description: |-
  Gets information about an existing AzureRM Service Endpoint.
---

# Data Source : azuredevops_serviceendpoint_azurerm

Use this data source to access information about an existing AzureRM service Endpoint.

## Example Usage

### By Service Endpoint ID

```hcl
data "azuredevops_project" "sample" {
  name = "Sample Project"
}

data "azuredevops_serviceendpoint_azurerm" "serviceendpoint" {
  project_id          = data.azuredevops_project.sample.id
  service_endpoint_id = "00000000-0000-0000-0000-000000000000"
}

output "service_endpoint_name" {
  value = data.azuredevops_serviceendpoint_azurerm.serviceendpoint.service_endpoint_name
}
```

### By Service Endpoint Name

```hcl
data "azuredevops_project" "sample" {
  name = "Sample Project"
}

data "azuredevops_serviceendpoint_azurerm" "serviceendpoint" {
  project_id            = data.azuredevops_project.sample.id
  service_endpoint_name = "Example-Service-Endpoint"
}

output "service_endpoint_id" {
  value = data.azuredevops_serviceendpoint_azurerm.serviceendpoint.id
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `service_endpoint_id` - (Optional) the ID of the Service Endpoint.

* `service_endpoint_name` - (Optional) the Name of the Service Endpoint.

~> **NOTE:** 1. One of either `service_endpoint_id` or `service_endpoint_name` must be specified.
    <br>2. When supplying `service_endpoint_name`, take care to ensure that this is a unique name.

## Attributes Reference

In addition to the Arguments list above - the following Attributes are exported:

* `id` - The ID of the Azure Resource Manager Service Endpoint.

* `authorization` - The Authorization scheme.

* `azurerm_management_group_id` - The Management Group ID of the Service Endpoint is target, if available.

* `azurerm_management_group_name` - The Management Group Name of the Service Endpoint target, if available.

* `azurerm_subscription_id` - The Subscription ID of the Service Endpoint target, if available.

* `azurerm_subscription_name` - The Subscription Name of the Service Endpoint target, if available.

* `resource_group` - The Resource Group of the Service Endpoint target, if available.

* `azurerm_spn_tenantid` - The Tenant ID of the Azure targets.

* `service_principal_id` - The Application(Client) ID of the Service Principal.

* `description` - The description of the Service Endpoint.

* `server_url` - The server URL of the service Endpoint.

* `environment` - The Cloud Environment.

* `service_endpoint_authentication_scheme` - The authentication scheme of Azure Resource Management Endpoint

* `workload_identity_federation_issuer` - The issuer if `of the Workload Identity Federation Subject

* `workload_identity_federation_subject` - The subject of the Workload Identity Federation Subject.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minute) Used when retrieving the Azure Resource Manager Service Endpoint.