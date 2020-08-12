---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_iteration_permissions"
description: |-
  Manages permissions for a AzureDevOps Iteration (Sprint)
---

# azuredevops_iteration_permissions

Manages permissions for an Iteration (Sprint)

~> **Note** Permissions can be assigned to group principals and not to single user principals.

## Permission levels

Permission for Iterations within Azure DevOps can be applied on two different levels.
Those levels are reflected by specifying (or omitting) values for the arguments `project_id` and `path`.

## Example Usage

```hcl

resource "azuredevops_project" "project" {
  project_name       = "Sample Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

data "azuredevops_group" "project-readers" {
	project_id = azuredevops_project.project.id
	name       = "Readers"
}

resource "azuredevops_iteration_permissions" "root-permissions" {
	project_id  = azuredevops_project.project.id
	principal   = azuredevops_group.project-readers.id
	permissions = {
	  CREATE_CHILDREN = "Deny"
	  GENERIC_READ    = "NotSet"
	  DELETE          = "Deny"
	}
}

resource "azuredevops_iteration_permissions" "iteration-permissions" {
	project_id  = azuredevops_project.project.id
	principal   = azuredevops_group.project-readers.id
	path        = "Iteration 1"
	permissions = {
	  CREATE_CHILDREN = "Allow"
	  GENERIC_READ    = "NotSet"
	  DELETE          = "Allow"
	}
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project to assign the permissions.
* `principal` - (Required) The **group** principal to assign the permissions.
* `permissions` - (Required) the permissions to assign. The following permissions are available.
* `path` - (Optional) The name of the branch to assign the permissions. 
* `replace` - (Optional) Replace (`true`) or merge (`false`) the permissions. Default: `true`

| Permission      | Description                    |
|-----------------|--------------------------------|
| GENERIC_READ    | View permissions for this node |
| GENERIC_WRITE   | Edit this node                 |
| CREATE_CHILDREN | Create child nodes             |
| DELETE          | Delete this node               |

## Relevant Links

* [Azure DevOps Service REST API 5.1 - Security](https://docs.microsoft.com/en-us/rest/api/azure/devops/security/?view=azure-devops-rest-5.1)

## Import

The resource does not support import.

## PAT Permissions Required

- **Project & Team**: vso.security_manage - Grants the ability to read, write, and manage security permissions.
