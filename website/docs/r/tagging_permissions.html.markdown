---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_tagging_permissions"
description: |-
  Manages permissions for AzureDevOps Tagging
---

# azuredevops_tagging_permissions

Manages permissions for tagging

## Permission levels

Permissions for tagging within Azure DevOps can be applied only on Project level.
The project level is reflected by specifying the argument `project_id`.

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

resource "azuredevops_tagging_permissions" "root-permissions" {
  project_id  = azuredevops_project.project.id
  principal   = data.azuredevops_group.project-readers.id
  permissions = {
    Enumerate    = "allow"
    Create       = "allow"
    Update       = "allow"
    Delete       = "allow"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project to assign the permissions.
* `principal` - (Required) The **group or user** principal to assign the permissions.
* `permissions` - (Required) the permissions to assign. The following permissions are available.
* `replace` - (Optional) Replace (`true`) or merge (`false`) the permissions. Default: `true`

| Name               | Permission Description     |
| ------------------ | -------------------------- |
| Enumerate          | Enumerate tag definitions  |
| Create             | Create tag definition      | 
| Update             | Update tag definition      | 
| Delete             | Delete tag definition      |  

## Relevant Links

* [Azure DevOps Service REST API 6.0 - Security](https://docs.microsoft.com/en-us/rest/api/azure/devops/security/?view=azure-devops-rest-6.0)

## Import

The resource does not support import.

## PAT Permissions Required

- **Project & Team**: vso.security_manage - Grants the ability to read, write, and manage security permissions.
