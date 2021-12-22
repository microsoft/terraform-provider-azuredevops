---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_servicehook_permissions"
description: |-
  Manages ServiceHook Permissions.
---

# azuredevops_servicehook_permissions

Manages ServiceHook Permissions.

~> **Note** Permissions can be assigned to group principals and not to single user principals.

## Permission levels

Permissions for ServiceHooks within Azure DevOps can be applied on the
organizational level or, if specified using the `project_id` attribute, on a
single project.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Sample Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

data "azuredevops_group" "example" {
  project_id = azuredevops_project.example.id
  name       = "Readers"
}

resource "azuredevops_servicehook_permissions" "example" {
  project_id = azuredevops_project.example.id
  principal  = data.azuredevops_group.example.id

  permissions = {
    ViewSubscriptions   = "allow"
    EditSubscriptions   = "allow"
    DeleteSubscriptions = "allow"
    PublishEvents       = "allow"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `permissions` - (Required) The permissions to assign. The following permissions are available:

    | Permission          | Description                            |
    | ------------------- | -------------------------------------- |
    | ViewSubscriptions   | View view/list ServiceHook definitions |
    | EditSubscriptions   | Edit/Create ServiceHook definitions    |
    | DeleteSubscriptions | Delete ServiceHook definitions         |
    | PublishEvents       | Publish events to a ServiceHook        |

* `principal` - (Required) The **group** principal to assign the permissions. Changing this forces a new ServiceHook Permissions to be created.

---

* `project_id` - (Optional) The ID of the project to assign the permissions. Changing this forces a new ServiceHook Permissions to be created. If no `project_id` the permissions are assigned for the complete Azure DevOps Organization.

* `replace` - (Optional) Replace (`true`) or merge (`false`) the permissions. Default: `true`.

## Import

The resource does not support import.

## PAT Permissions Required

- **Project & Team**: vso.security_manage - Grants the ability to read, write, and manage security permissions.
