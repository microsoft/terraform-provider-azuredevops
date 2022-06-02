---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_permissions"
description: |-
  Manages permissions for a AzureDevOps Service Endpoint
---

# azuredevops_serviceendpoint_permissions

Manages permissions for a Service Endpoint

~> **Note** Permissions can be assigned to group principals and not to single user principals.

## Permission levels

Permission for Service Endpoints within Azure DevOps can be applied on two different levels.
Those levels are reflected by specifying (or omitting) values for the arguments `project_id` and `serviceendpoint_id`.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

data "azuredevops_group" "example-readers" {
  project_id = azuredevops_project.example.id
  name       = "Readers"
}

resource "azuredevops_serviceendpoint_permissions" "example-root-permissions" {
  project_id = azuredevops_project.example.id
  principal  = data.azuredevops_group.example-readers.id
  permissions = {
    Use               = "allow"
    Administer        = "allow"
    Create            = "allow"
    ViewAuthorization = "allow"
    ViewEndpoint      = "allow"
  }
}

resource "azuredevops_serviceendpoint_dockerregistry" "example" {
  project_id            = azuredevops_project.example.id
  service_endpoint_name = "Example Docker Hub"
  docker_username       = "username"
  docker_email          = "email@example.com"
  docker_password       = "password"
  registry_type         = "DockerHub"
}

resource "azuredevops_serviceendpoint_permissions" "example-permissions" {
  project_id         = azuredevops_project.example.id
  principal          = data.azuredevops_group.example-readers.id
  serviceendpoint_id = azuredevops_serviceendpoint_dockerregistry.example.id
  permissions = {
    Use               = "allow"
    Administer        = "deny"
    Create            = "deny"
    ViewAuthorization = "allow"
    ViewEndpoint      = "allow"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.
* `principal` - (Required) The **group** principal to assign the permissions.
* `permissions` - (Required) the permissions to assign. The following permissions are available.
* `serviceendpoint_id` - (Optional) The id of the service endpoint to assign the permissions.
* `replace` - (Optional) Replace (`true`) or merge (`false`) the permissions. Default: `true`

| Permission        | Description                         |
| ----------------- | ----------------------------------- |
| Use               | Use service endpoint                |
| Administer        | Full control over service endpoints |
| Create            | Create service endpoints            |
| ViewAuthorization | View authorizations                 |
| ViewEndpoint      | View service endpoint properties    |

## Relevant Links

* [Azure DevOps Service REST API 6.0 - Security](https://docs.microsoft.com/en-us/rest/api/azure/devops/security/?view=azure-devops-rest-6.0)

## Import

The resource does not support import.

## PAT Permissions Required

- **Project & Team**: vso.security_manage - Grants the ability to read, write, and manage security permissions.
