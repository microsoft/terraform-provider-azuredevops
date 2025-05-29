---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_secure_file_permissions"
description: |-
  Manages permissions for an Azure DevOps Secure File
---

# azuredevops_secure_file_permissions

Manages permissions for a Secure File in Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  name               = "Testing"
  description        = "Testing-description"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_secure_file" "example" {
  project_id = azuredevops_project.project.id
  name       = "my-secure-file.txt"
  content    = file("./my-secure-file.txt")
}

data "azuredevops_group" "tf-project-readers" {
  project_id = azuredevops_project.project.id
  name       = "Readers"
}

resource "azuredevops_secure_file_permissions" "permissions" {
  project_id      = azuredevops_project.project.id
  secure_file_id  = azuredevops_secure_file.example.id
  principal       = data.azuredevops_group.tf-project-readers.id
  permissions = {
    "View" : "allow",
    "Manage" : "allow",
    "Use" : "allow",
  }
}
```

## Roles

The Azure DevOps UI uses roles to assign permissions for secure files.

| Role          | Allow Permissions     |
|---------------|-----------------------|
| Reader        | View                  |
| User          | View, Use             |
| Administrator | View, Use, Administer |

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.
* `secure_file_id` - (Required) The id of the secure file to assign the permissions.
* `principal` - (Required) The **group** principal to assign the permissions.
* `permissions` - (Required) The permissions to assign. The following permissions are available:

  | Permission  | Description               |
  |-------------|---------------------------|
  | View        | View library item         |
  | Manage      | Administer library item   |
  | Create      | Create library item       |
  | Use         | Use library item          |
  | Owner       | Owner library item        |

---

* `replace` - (Optional) Replace (`true`) or merge (`false`) the permissions. Default: `true`

## Relevant Links

* [Azure DevOps Service REST API 7.1 - Security](https://docs.microsoft.com/en-us/rest/api/azure/devops/security/?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Secure File Permissions.
* `read` - (Defaults to 5 minute) Used when retrieving the Secure File Permissions.
* `update` - (Defaults to 10 minutes) Used when updating the Secure File Permissions.
* `delete` - (Defaults to 10 minutes) Used when deleting the Secure File Permissions.

## Import

The resource does not support import.

## PAT Permissions Required

- **Project & Team**: vso.security_manage - Grants the ability to read, write, and manage security permissions.

