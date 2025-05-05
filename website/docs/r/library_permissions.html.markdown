---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_library_permissions"
description: |-
  Manages permissions for an Azure DevOps Library
---

# azuredevops_library_permissions

Manages permissions for a Library

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  name               = "Testing"
  description        = "Testing-description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

data "azuredevops_group" "tf-project-readers" {
  project_id = azuredevops_project.project.id
  name       = "Readers"
}

resource "azuredevops_library_permissions" "permissions" {
  project_id = azuredevops_project.project.id
  principal  = data.azuredevops_group.tf-project-readers.id
  permissions = {
    "View" : "allow",
    "Administer" : "allow",
    "Use" : "allow",
  }
}
```

## Roles

The Azure DevOps UI uses roles to assign permissions for the Library.

| Role          | Allowed Permissions   |
|---------------|-----------------------|
| Reader        | View                  |
| Creator       | View, Create          |
| User          | View, Use             |
| Administrator | View, Use, Administer |

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `principal` - (Required) The **group** principal to assign the permissions.

* `variable_group_id` - (Required) The id of the variable group to assign the permissions.

* `permissions` - (Required) the permissions to assign. The following permissions are available.

  | Permission  | Description               |
  |-------------|---------------------------|
  | View        | View library item         |
  | Administer  | Administer library item   |
  | Create      | Create library item       |
  | ViewSecrets | View library item secrets |
  | Use         | Use library item          |
  | Owner       | Owner library item        |

---

* `replace` - (Optional) Replace (`true`) or merge (`false`) the permissions. Default: `true`

## Relevant Links

* [Azure DevOps Service REST API 6.0 - Security](https://docs.microsoft.com/en-us/rest/api/azure/devops/security/?view=azure-devops-rest-6.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Library Permission.
* `read` - (Defaults to 5 minute) Used when retrieving the Library Permission.
* `update` - (Defaults to 10 minutes) Used when updating the Library Permission.
* `delete` - (Defaults to 10 minutes) Used when deleting the Library Permission.

## Import

The resource does not support import.

## PAT Permissions Required

- **Project & Team**: vso.security_manage - Grants the ability to read, write, and manage security permissions.
