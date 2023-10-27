---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_area_permissions"
description: |-
  Manages permissions for a AzureDevOps Area (Component)
---

# azuredevops_area_permissions

Manages permissions for an Area (Component)

~> **Note** Permissions can be assigned to group principals and not to single user principals.

## Permission levels

Permission for Areas within Azure DevOps can be applied on two different levels.
Those levels are reflected by specifying (or omitting) values for the arguments `project_id` and `path`.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

data "azuredevops_group" "example-project-readers" {
  project_id = azuredevops_project.example.id
  name       = "Readers"
}

resource "azuredevops_area_permissions" "example-root-permissions" {
  project_id = azuredevops_project.example.id
  principal  = data.azuredevops_group.example-project-readers.id
  path       = "/"
  permissions = {
    CREATE_CHILDREN = "Deny"
    GENERIC_READ    = "Allow"
    DELETE          = "Deny"
    WORK_ITEM_READ  = "Allow"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project to assign the permissions.
* `principal` - (Required) The **group** principal to assign the permissions.
* `permissions` - (Required) the permissions to assign. The following permissions are available.
* `path` - (Optional) The name of the branch to assign the permissions. 
* `replace` - (Optional) Replace (`true`) or merge (`false`) the permissions. Default: `true`.

| Permission             | Description                          |
|------------------------|--------------------------------------|
| GENERIC_READ           | View permissions for this node       |
| GENERIC_WRITE          | Edit this node                       |
| CREATE_CHILDREN        | Create child nodes                   |
| DELETE                 | Delete this node                     |
| WORK_ITEM_READ         | View work items in this node         |
| WORK_ITEM_WRITE        | Edit work items in this node         |
| MANAGE_TEST_PLANS      | Manage test plans                    |
| MANAGE_TEST_SUITES     | Manage test suites                   |
| WORK_ITEM_SAVE_COMMENT | Edit work item comments in this node |

## Relevant Links

* [Azure DevOps Service REST API 7.0 - Security](https://docs.microsoft.com/en-us/rest/api/azure/devops/security/?view=azure-devops-rest-7.0)

## Import

The resource does not support import.

## PAT Permissions Required

- **Project & Team**: vso.security_manage - Grants the ability to read, write, and manage security permissions.
