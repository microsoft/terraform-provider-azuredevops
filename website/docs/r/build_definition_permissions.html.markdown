---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_build_definition_permissions"
description: |-
  Manages permissions for a AzureDevOps Build Definition
---

# azuredevops_build_definition_permissions

Manages permissions for a Build Definition

~> **Note** Permissions can be assigned to group principals and not to single user principals.

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

resource "azuredevops_git_repository" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Repository"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_build_definition" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Build Definition"
  path       = "\\ExampleFolder"

  ci_trigger {
    use_yaml = true
  }

  repository {
    repo_type   = "TfsGit"
    repo_id     = azuredevops_git_repository.example.id
    branch_name = azuredevops_git_repository.example.default_branch
    yml_path    = "azure-pipelines.yml"
  }
}

resource "azuredevops_build_definition_permissions" "example" {
  project_id = azuredevops_project.example.id
  principal  = data.azuredevops_group.example-readers.id

  build_definition_id = azuredevops_build_definition.example.id

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

* [Azure DevOps Service REST API 7.0 - Security](https://docs.microsoft.com/en-us/rest/api/azure/devops/security/?view=azure-devops-rest-7.0)

## Import

The resource does not support import.

## PAT Permissions Required

- **Project & Team**: vso.security_manage - Grants the ability to read, write, and manage security permissions.
