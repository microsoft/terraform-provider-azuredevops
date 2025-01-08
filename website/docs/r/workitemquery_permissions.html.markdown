---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitemquery_permissions"
description: |-
  Manages permissions for Work Item Queries
---

# azuredevops_workitemquery_permissions

Manages permissions for Work Item Queries. 

~> **Note** Permissions can be assigned to group principals and not to single user principals.

## Permission levels

Permission for Work Item Queries within Azure DevOps can be applied on two different levels.
Those levels are reflected by specifying (or omitting) values for the arguments `project_id` and `path`.

### Project level

Permissions for all Work Item Queries inside a project (existing or newly created ones) are specified, if only the argument `project_id` has a value.

#### Example usage

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

resource "azuredevops_workitemquery_permissions" "project-wiq-root-permissions" {
  project_id = azuredevops_project.example.id
  principal  = data.azuredevops_group.example-readers.id
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

resource "azuredevops_workitemquery_permissions" "example-permissions" {
  project_id = azuredevops_project.example.id
  path       = "/Team"
  principal  = data.azuredevops_group.example-readers.id
  permissions = {
    Contribute = "Allow"
    Delete     = "Deny"
    Read       = "NotSet"
  }
}
```

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

data "azuredevops_group" "example-contributors" {
  project_id = azuredevops_project.example.id
  name       = "Contributors"
}

resource "azuredevops_workitemquery_permissions" "example-project-permissions" {
  project_id = azuredevops_project.example.id
  principal  = data.azuredevops_group.example-readers.id
  permissions = {
    Read              = "Allow"
    Delete            = "Deny"
    Contribute        = "Deny"
    ManagePermissions = "Deny"
  }
}

resource "azuredevops_workitemquery_permissions" "example-sharedqueries-permissions" {
  project_id = azuredevops_project.example.id
  path       = "/"
  principal  = data.azuredevops_group.example-contributors.id
  permissions = {
    Read   = "Allow"
    Delete = "Deny"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project to assign the permissions.
* `path` - (Optional) Path to a query or folder beneath `Shared Queries`
* `principal` - (Required) The **group** principal to assign the permissions.
* `replace` - (Optional) Replace (`true`) or merge (`false`) the permissions. Default: `true`
* `permissions` - (Required) the permissions to assign. The following permissions are available

| Permissions              | Description                        |
|--------------------------|------------------------------------|
| Read                     | Read                               |
| Contribute               | Contribute                         |
| Delete                   | Delete                             |
| ManagePermissions        | Manage Permissions                 |

## Relevant Links

* [Azure DevOps Service REST API 7.0 - Security](https://docs.microsoft.com/en-us/rest/api/azure/devops/security/?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Work Item Query Permissions.
* `read` - (Defaults to 5 minute) Used when retrieving the Work Item Query Permissions.
* `update` - (Defaults to 10 minutes) Used when updating the Work Item Query Permissions.
* `delete` - (Defaults to 10 minutes) Used when deleting the Work Item Query Permissions.

## Import

The resource does not support import.

## PAT Permissions Required

- **Project & Team**: vso.security_manage - Grants the ability to read, write, and manage security permissions.
