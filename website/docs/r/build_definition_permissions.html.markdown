---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_build_definition"
description: |-
  Manages permissions for a AzureDevOps Build Definition
---

# azuredevops_build_definition

Manages permissions for a Build Definition

~> **Note** Permissions can be assigned to group principals and not to single user principals.

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  name       = "Sample Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

data "azuredevops_group" "project-readers" {
	project_id = azuredevops_project.project.id
	name       = "Readers"
}

resource "azuredevops_git_repository" "repository" {
  project_id = azuredevops_project.project.id
  name       = "Sample Repository"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_build_definition" "build" {
  project_id = azuredevops_project.project.id
  name       = "Sample Build Definition"
  path       = "\\ExampleFolder"

  ci_trigger {
    use_yaml = true
  }

  repository {
    repo_type   = "TfsGit"
    repo_id     = azuredevops_git_repository.repository.id
    branch_name = azuredevops_git_repository.repository.default_branch
    yml_path    = "azure-pipelines.yml"
  }
}

resource "azuredevops_build_definition_permissions" "permissions" {
	project_id  = azuredevops_project.project.id
	principal   = data.azuredevops_group.project-readers.id

	build_definition_id   = azuredevops_build_definition.build.id
	build_definition_path = azuredevops_build_definition.build.path

	permissions = {
	  ViewBuilds       = "Allow"
	  EditBuildQuality = "Deny"
	  DeleteBuilds     = "Deny"
	  StopBuilds       = "Allow"
	}
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project to assign the permissions.
* `principal` - (Required) The **group** principal to assign the permissions.
* `build_definition_id` - (Required) The id of the build definition to assign the permissions. 
* `build_definition_path` - (Required) The path of the build definition to assign the permissions. 
* `replace` - (Optional) Replace (`true`) or merge (`false`) the permissions. Default: `true`.
* `permissions` - (Required) the permissions to assign. The following permissions are available.

| Permission                     | Description                           |
|--------------------------------|---------------------------------------|
| ViewBuilds                     | View builds                           |
| EditBuildQuality               | Edit build quality                    |
| RetainIndefinitely             | Retain indefinitely                   |
| DeleteBuilds                   | Delete builds                         |
| ManageBuildQualities           | Manage build qualities                |
| DestroyBuilds                  | Destroy builds                        |
| UpdateBuildInformation         | Update build information              |
| QueueBuilds                    | Queue builds                          |
| ManageBuildQueue               | Manage build queue                    |
| StopBuilds                     | Stop builds                           |
| ViewBuildDefinition            | View build pipeline                   |
| EditBuildDefinition            | Edit build pipeline                   |
| DeleteBuildDefinition          | Delete build pipeline                 |
| OverrideBuildCheckInValidation | Override check-in validation by build |
| AdministerBuildPermissions     | Administer build permissions          |

## Relevant Links

* [Azure DevOps Service REST API 5.1 - Security](https://docs.microsoft.com/en-us/rest/api/azure/devops/security/?view=azure-devops-rest-5.1)

## Import

The resource does not support import.

## PAT Permissions Required

- **Project & Team**: vso.security_manage - Grants the ability to read, write, and manage security permissions.
