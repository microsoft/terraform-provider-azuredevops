---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitemtrackingprocess_process_permissions"
description: |-
  Manages permissions for an Azure DevOps Process
---

# azuredevops_workitemtrackingprocess_process_permissions

Manages permissions for an Azure DevOps Process

## Example Usage

### Permissions on an inherited process

```hcl
resource "azuredevops_workitemtrackingprocess_process" "example" {
  name                   = "Example Process"
  parent_process_type_id = "adcc42ab-9882-485e-a3ed-7678f01f66bc" # Agile
}

data "azuredevops_group" "example-group" {
  name = "Example Group"
}

resource "azuredevops_workitemtrackingprocess_process_permissions" "example" {
  process_id = azuredevops_workitemtrackingprocess_process.example.id
  principal  = data.azuredevops_group.example-group.id
  permissions = {
    Edit                         = "Allow"
    Delete                       = "Deny"
    AdministerProcessPermissions = "Allow"
  }
}
```

### Permissions on a system process

```hcl
data "azuredevops_group" "example-group" {
  name = "Example Group"
}

resource "azuredevops_workitemtrackingprocess_process_permissions" "example" {
  process_id = "adcc42ab-9882-485e-a3ed-7678f01f66bc"  # Agile system process
  principal  = data.azuredevops_group.example-group.id
  permissions = {
    Create = "Deny"  # Prevent creating inherited processes from Agile
  }
}
```

## Argument Reference

The following arguments are supported:

* `process_id` - (Required) The ID of the process to assign the permissions.

* `principal` - (Required) The principal to assign the permissions.

* `permissions` - (Required) the permissions to assign. The following permissions are available

    **Inherited process permissions:**

    | Permission                     | Description                    |
    |--------------------------------|--------------------------------|
    | Edit                           | Edit process                   |
    | Delete                         | Delete process                 |
    | AdministerProcessPermissions   | Administer process permissions |

    **System process permissions:**

    | Permission                     | Description                    |
    |--------------------------------|--------------------------------|
    | Create                         | Create inherited process       |

---

* `replace` - (Optional) Replace (`true`) or merge (`false`) the permissions. Default: `true`

## Relevant Links

* [Azure DevOps Service REST API 7.1 - Security](https://docs.microsoft.com/en-us/rest/api/azure/devops/security/?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Process Permission.
* `read` - (Defaults to 5 minute) Used when retrieving the Process Permission.
* `update` - (Defaults to 10 minutes) Used when updating the Process Permission.
* `delete` - (Defaults to 10 minutes) Used when deleting the Process Permission.

## Import

The resource does not support import.

## PAT Permissions Required

- **Security**: Manage
- **Identity**: Read
