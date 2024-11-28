---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_azure_service_bus"
description: |-
  Manages a Azure Service Bus service endpoint within Azure DevOps organization.
---

# azuredevops_serviceendpoint_azure_service_bus

Manages an Azure Service Bus endpoint within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_serviceendpoint_azure_service_bus" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example Azure Service Bus"
  queue_name            = "queue"
  connection_string     = "connection string"
  description           = "Managed by Terraform"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `service_endpoint_name` - (Required) The Service Endpoint name.

* `queue_name` - (Required) The Azure Service Bus Queue Name.

* `connection_string` - (Required) The  Azure Service Bus Connection string.

* `description` - (Optional) The Service Endpoint description. Defaults to `Managed by Terraform`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the service endpoint.

* `project_id` - The ID of the project.

* `service_endpoint_name` - The Service Endpoint name.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Endpoints](https://docs.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/endpoints?view=azure-devops-rest-7.0)

## Import

Azure DevOps Azure Service Bus Service Endpoint can be imported using **projectID/serviceEndpointID** or **projectName/serviceEndpointID**

```sh
terraform import azuredevops_serviceendpoint_azure_service_bus.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
