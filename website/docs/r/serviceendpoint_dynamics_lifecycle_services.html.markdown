---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_dynamics_lifecycle_services"
description: |-
  Manages a Dynamics Lifecycle Services service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_dynamics_lifecycle_services

Manages a Dynamics Lifecycle Services service endpoint within Azure DevOps. Using this service endpoint requires you to install: [Dynamics Lifecycle Services](https://marketplace.visualstudio.com/items?itemName=Dyn365FinOps.dynamics365-finops-tools)

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_dynamics_lifecycle_services" "example" {
  project_id                      = azuredevops_project.example.id
  service_endpoint_name           = "Example Service connection"
  authorization_endpoint          = "https://login.microsoftonline.com/organization"
  lifecycle_services_api_endpoint = "https://lcsapi.lcs.dynamics.com"
  client_id                       = "00000000-0000-0000-0000-000000000000"
  username                        = "username"
  password                        = "password"
  description                     = "Managed by Terraform"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `service_endpoint_name` - (Required) The Service Endpoint name.

* `authorization_endpoint` - (Required) The URL of the Authentication Endpoint.

* `lifecycle_services_api_endpoint` - (Required) The URL of the Lifecycle Services API Endpoint.

* `client_id` - (Required) The client ID for a native application registration in Azure Active Directory with API permissions for Dynamics Lifecycle Services.
 
* `username` - (Required) The E-mail address of user with sufficient permissions to interact with LCS asset library and environments.

* `password` - (Required) The Password for the Azure Active Directory account.

* `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

## Attributes Reference

The following attributes are exported:

- `id` - The ID of the service endpoint.
- `project_id` - The ID of the project.
- `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 minutes) Used when creating the Dynamic Life Cycle Service Endpoint.
* `read` - (Defaults to 1 minute) Used when retrieving the Dynamic Life Cycle Service Endpoint.
* `update` - (Defaults to 2 minutes) Used when updating the Dynamic Life Cycle Service Endpoint.
* `delete` - (Defaults to 2 minutes) Used when deleting the Dynamic Life Cycle Service Endpoint.

## Import

Azure DevOps Dynamics Life Cycle Service Endpoint can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
terraform import azuredevops_serviceendpoint_dynamics_lifecycle_services.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
