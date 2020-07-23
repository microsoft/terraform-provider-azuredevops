---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitemquery_permissions"
description: |-
  Manages permissions for Work Item Queries
---

# azuredevops_git_permissions

Manages permissions for Work Item Queries. 

~> **Note** Permissions can be assigned to group principals and not to single user principals.

## Permission levels

Permission for Work Item Queries within Azure DevOps can be applied on two different levels.
Those levels are reflected by specifying (or omitting) values for the arguments `project_id` and `path`.

### Project level

Permissions for all Work Item Queries inside a project (existing or newly created ones) are specified, if only the argument `project_id` has a value.

#### Example usage

```hcl
resource "azuredevops_workitemquery_permissions" "project-wiq-root-permissions" {
  project_id  = azuredevops_project.project.id  
  principal   = data.azuredevops_group.project-readers.id
  permissions = {
    CreateRepository = "Deny"
    DeleteRepository = "Deny"
    RenameRepository = "NotSet"
  }
}
```

### Shared Queries folder level

Permissions for a specific folder inside Shared Queries are specified if the arguments `project_id` and `path` are set.

~> **Note** To set permissions for the Shared Queries folder itself use `/` as path value

#### Example usage

```hcl
resource "azuredevops_workitemquery_permissions" "wiq-folder-permissions" {
  project_id = azuredevops_project.project.id
  path = "/Team"
  principal   = data.azuredevops_group.project-readers.id
  permissions = {
    Contribute = "Allow"
    Delete = "Deny"
    Read = "NotSet"
  }
}
```

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  project_name       = "Test Project"
  description        = "Test Project Description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

data "azuredevops_group" "project-readers" {
  project_id = azuredevops_project.project.id
  name       = "Readers"
}

data "azuredevops_group" "project-contributors" {
  project_id = azuredevops_project.project.id
  name       = "Contributors"
}

data "azuredevops_group" "project-administrators" {
  project_id = azuredevops_project.project.id
  name       = "Project administrators"
}

resource "azuredevops_workitemquery_permissions" "wiq-project-permissions" {
  project_id  = azuredevops_project.project.id
  principal   = data.azuredevops_group.project-administrators.id
  permissions = {
    Delete            = "Deny"
    Contribute        = "Allow"
    ManagePermissions = "NotSet"
  }
}

resource "azuredevops_workitemquery_permissions" "wiq-sharedqueries-permissions" {
  project_id = azuredevops_project.project.id
  path = "/"
  principal   = data.azuredevops_group.project-contributors.id
  permissions = {
    FullControl              = "Allow"
  }
}

```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project to assign the permissions.
* `path` - (Optional) Path to a query or folder beneath `Shared Queries`
* `principal` - (Required) The **group** principal to assign the permissions.
* `replace` - (Optional) Replace (`true`) or merge (`false`) the permissions. Default: `true`
* `permissions` - (Required) the permissions to assign. The follwing permissions are available


| Permissions              | Description                        |
|--------------------------|------------------------------------|
| Read                     | Read                               |
| Contribute               | Contribute                         |
| Delete                   | Delete                             |
| ManagePermissions        | Manage Permissions                 |
| FullControl              | Full Control                       |

## Relevant Links

* [Azure DevOps Service REST API 5.1 - Security](https://docs.microsoft.com/en-us/rest/api/azure/devops/security/?view=azure-devops-rest-5.1)

## Import

The resource does not support import.

## PAT Permissions Required

- **Project & Team**: vso.security_manage - Grants the ability to read, write, and manage security permissions.
