---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_servicehook_permissions"
description: |-
  Manages permissions for AzureDevOps service hooks
---

# azuredevops_servicehook_permissions

Manages permissions for service hooks

## Permission levels

Permissions for service hooks within Azure DevOps can be applied on the Organizational level or, if the optional attribute `project_id` is specified, on Project level.
Those levels are reflected by specifying (or omitting) values for the argument `project_id`.

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

resource "azuredevops_servicehook_permissions" "example-permissions" {
  project_id = azuredevops_project.example.id
  principal  = data.azuredevops_group.example-readers.id
  permissions = {
    ViewSubscriptions   = "allow"
    EditSubscriptions   = "allow"
    DeleteSubscriptions = "allow"
    PublishEvents       = "allow"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (optional) The ID of the project.
* `principal` - (Required) The **group** principal to assign the permissions.
* `permissions` - (Required) the permissions to assign. The following permissions are available.
* `replace` - (Optional) Replace (`true`) or merge (`false`) the permissions. Default: `true`

| Name               | Permission Description   |
| ------------------ | ------------------------ |
| ViewSubscriptions  | View Subscriptions       |
| EditSubscriptions  | Edit Subscription        | 
| DeleteSubscriptions| Delete Subscriptions     | 
| PublishEvents      | Publish Events           | 

## Relevant Links

* [Azure DevOps Service REST API 7.0 - Security](https://docs.microsoft.com/en-us/rest/api/azure/devops/security/?view=azure-devops-rest-7.0)

## Import

The resource does not support import.

## PAT Permissions Required

- **Project & Team**: vso.security_manage - Grants the ability to read, write, and manage security permissions.
