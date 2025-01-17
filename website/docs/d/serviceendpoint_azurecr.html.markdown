---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_azurecr"
description: |-
  Gets information about an existing Azure Container Registry Service Endpoint. 
---

# Data Source : azuredevops_serviceendpoint_azurecr

Use this data source to access information about an existing Azure Container Registry Service Endpoint.

## Example Usage

```hcl
data "azuredevops_serviceendpoint_azurecr" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example Azure Container Registry"
}

output "service_endpoint_id" {
  value = data.azuredevops_serviceendpoint_azurecr.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `service_endpoint_id` - (Optional) the ID of the Service Endpoint.

* `service_endpoint_name` - (Optional) the Name of the Service Endpoint.

~> **NOTE:** One of either `service_endpoint_id` or `service_endpoint_name` must be specified.

## Attributes Reference

In addition to the Arguments list above - the following Attributes are exported:

* `authorization` - Specifies the Authorization Scheme Map.

* `description` - The Service Endpoint description.

* `resource_group` - The Resource Group to which the Container Registry belongs.

* `azurecr_spn_tenantid` - The Tenant ID of the service principal.

* `azurecr_name` - The Azure Container Registry name.

* `azurecr_subscription_id` - The Subscription ID of the Azure targets.

* `azurecr_subscription_name` - The Subscription Name of the Azure targets.

* `app_object_id` - The Object ID of the Service Principal.

* `spn_object_id` - The ID of the Service Principal.

* `az_spn_role_assignment_id` - The ID of Service Principal Role Assignment.

* `az_spn_role_permissions` - The Service Principal Role Permissions.

* `service_principal_id` - The Application(Client) ID of the Service Principal.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minute) Used when retrieving the Azure Container Registry Service Endpoint.