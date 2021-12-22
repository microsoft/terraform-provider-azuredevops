---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_servicehooks_permissions"
description: |-
  Manages permissions for AzureDevOps service hooks
---

# azuredevops_servicehooks_permissions

Manages permissions for service hooks

## Permission levels

Permissions for service hooks within Azure DevOps can be applied only on Project level.
Those levels are reflected by specifying (or omitting) values for the argument `project_id`.

## Example Usage

```hcl

resource "azuredevops_project" "project" {
  name               = "Sample Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

data "azuredevops_group" "project-readers" {
  project_id = azuredevops_project.project.id
  name       = "Readers"
}

resource "azuredevops_servicehooks_permissions" "root-permissions" {
  project_id  = azuredevops_project.project.id
  principal   = data.azuredevops_group.project-readers.id
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

* `project_id` - (optional) The ID of the project to assign the permissions.
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

* [Azure DevOps Service REST API 5.1 - Security](https://docs.microsoft.com/en-us/rest/api/azure/devops/security/?view=azure-devops-rest-5.1)

## Import

The resource does not support import.

## PAT Permissions Required

- **Project & Team**: vso.security_manage - Grants the ability to read, write, and manage security permissions.
