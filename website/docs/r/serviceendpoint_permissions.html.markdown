---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_project_permissions"
description: |-
  Manages project permissions for a AzureDevOps Service Endpoint
---

# azuredevops_serviceendpoint_project_permissions

Manages project permissions for a Service Endpoint, allowing sharing a service connection with multiple projects including optional service_endpoint_name and description.

~> **Note** Permissions can be assigned to group principals and not to single user principals.

## Permission levels

Permission for Service Endpoints within Azure DevOps can be applied on two different levels.
Those levels are reflected by specifying (or omitting) values for the arguments `project_id` and `serviceendpoint_id`.

## Example Usage

```hcl
resource "azuredevops_serviceendpoint_project_permissions" "example-share" {
  serviceendpoint_id = azuredevops_serviceendpoint_azurerm.example.id

  project_reference {
    project_id            = azuredevops_project.example_one.id
    service_endpoint_name = "service-connection-shared"
    description           = "Service Connection Shared by Terraform - Cluster One"
  }

  project_reference {
    project_id            = azuredevops_project.example_two.id
    service_endpoint_name = "service-connection-shared"
    description           = "Service Connection Shared by Terraform - Cluster Two"
  }
}
```

## Argument Reference

The following arguments are supported:

* `serviceendpoint_id` - (Required) The ID of the service endpoint to share.
* `project_reference` - (Required) A list of `project_reference` block as defined below. Objects describing which projects the service connection will be shared.

An `project_reference` block supports the following:

* `project_id` - (Required) Project id which service endpoint will be shared.
* `service_endpoint_name` - (Optional) Name for service connection in the shared project. Default keep the same name.
* `description` - (Optional) Description for service connection in the shared project. Default keep the same description.

## Relevant Links

* [Azure DevOps Service REST API 7.0 - Security](https://docs.microsoft.com/en-us/rest/api/azure/devops/security/?view=azure-devops-rest-7.0)

## Import

The resource does not support import.

## PAT Permissions Required

- **Project & Team**: vso.security_manage - Grants the ability to read, write, and manage security permissions.
