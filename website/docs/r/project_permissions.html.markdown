---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_project_permissions"
description: |-
  Manages permissions for a AzureDevOps project
---

# azuredevops_project_permissions

Manages permissions for a AzureDevOps project

~> **Note** Permissions can be assigned to group principals and not to single user principals.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

data "azuredevops_group" "example-readers" {
  project_id = azuredevops_project.example.id
  name       = "Readers"
}

resource "azuredevops_project_permissions" "example-permission" {
  project_id = azuredevops_project.example.id
  principal  = data.azuredevops_group.example-readers.id
  permissions = {
    DELETE              = "Deny"
    EDIT_BUILD_STATUS   = "NotSet"
    WORK_ITEM_MOVE      = "Allow"
    DELETE_TEST_RESULTS = "Deny"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project to assign the permissions.
* `principal` - (Required) The **group** principal to assign the permissions.
* `replace` - (Optional) Replace (`true`) or merge (`false`) the permissions. Default: `true`
* `permissions` - (Required) the permissions to assign. The following permissions are available

| Permission                   | Description                                  |
|------------------------------|----------------------------------------------|
| GENERIC_READ                 | View project-level information               |
| GENERIC_WRITE                | Edit project-level information               |
| DELETE                       | Delete team project                          |
| PUBLISH_TEST_RESULTS         | Create test runs                             |
| ADMINISTER_BUILD             | Administer a build                           |
| START_BUILD                  | Start a build                                |
| EDIT_BUILD_STATUS            | Edit build quality                           |
| UPDATE_BUILD                 | Write to build operational store             |
| DELETE_TEST_RESULTS          | Delete test runs                             |
| VIEW_TEST_RESULTS            | View test runs                               |
| MANAGE_TEST_ENVIRONMENTS     | Manage test environments                     |
| MANAGE_TEST_CONFIGURATIONS   | Manage test configurations                   |
| WORK_ITEM_DELETE             | Delete and restore work items                |
| WORK_ITEM_MOVE               | Move work items out of this project          |
| WORK_ITEM_PERMANENTLY_DELETE | Permanently delete work items                |
| RENAME                       | Rename team project                          |
| MANAGE_PROPERTIES            | Manage project properties                    |
| MANAGE_SYSTEM_PROPERTIES     | Manage system project properties             |
| BYPASS_PROPERTY_CACHE        | Bypass project property cache                |
| BYPASS_RULES                 | Bypass rules on work item updates            |
| SUPPRESS_NOTIFICATIONS       | Suppress notifications for work item updates |
| UPDATE_VISIBILITY            | Update project visibility                    |
| CHANGE_PROCESS               | Change process of team project.              |
| AGILETOOLS_BACKLOG           | Agile backlog management.                    |
| AGILETOOLS_PLANS             | Agile plans.                                 |

## Relevant Links

* [Azure DevOps Service REST API 7.0 - Security](https://docs.microsoft.com/en-us/rest/api/azure/devops/security/?view=azure-devops-rest-7.0)

## Import

The resource does not support import.

## PAT Permissions Required

- **Project & Team**: vso.security_manage - Grants the ability to read, write, and manage security permissions.
