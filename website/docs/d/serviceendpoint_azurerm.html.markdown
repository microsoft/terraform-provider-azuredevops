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
  project_id = data.azuredevops_project.sample.id
  service_endpoint_id         = "00000000-0000-0000-0000-000000000000"
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

~> **NOTE:** One of either `service_endpoint_id` or `service_endpoint_name` must be specified.
~> **NOTE:** When supplying `service_endpoint_name`, take care to ensure that this is a unique name.

## Attributes Reference

In addition to the Arguments list above - the following Attributes are exported:

* `authorization` - Specifies the Authorization Scheme Map.
* `azurerm_management_group_id` - Specified the Management Group ID of the Service Endpoint is target, if available.
* `azurerm_management_group_name` - Specified the Management Group Name of the Service Endpoint target, if available.
* `azurerm_subscription_id` - Specifies the Subscription ID of the Service Endpoint target, if available.
* `azurerm_subscription_name` - Specifies the Subscription Name of the Service Endpoint target, if available.
* `resource_group` - Specifies the Resource Group of the Service Endpoint target, if available.
* `azurerm_spn_tenantid` - Specifies the Tenant ID of the Azure targets.
* `description` - Specifies the description of the Service Endpoint.
