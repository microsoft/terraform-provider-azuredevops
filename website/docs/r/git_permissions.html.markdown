---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_git_permissions"
description: |-
  Manages permissions for Git repositories
---

# azuredevops_git_permissions

Manages permissions for Git repositories. 

~> **Note** Permissions can be assigned to group principals and not to single user principals.

## Permission levels

Permission for Git Repositories within Azure DevOps can be applied on three different levels.
Those levels are reflected by specifying (or omitting) values for the arguments `project_id`, `repository_id` and `branch_name`.

### Project level

Permissions for all Git Repositories inside a project (existing or newly created ones) are specified, if only the argument `project_id` has a value.

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

resource "azuredevops_git_permissions" "example-permissions" {
  project_id = azuredevops_project.example.id
  principal  = data.azuredevops_group.example-readers.id
  permissions = {
    CreateRepository = "Deny"
    DeleteRepository = "Deny"
    RenameRepository = "NotSet"
  }
}
```

### Repository level

Permissions for a specific Git Repository and all existing or newly created branches are specified if the arguments `project_id` and `repository_id` are set.

#### Example usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

data "azuredevops_group" "example-group" {
  name = "Project Collection Administrators"
}

resource "azuredevops_git_repository" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Empty Git Repository"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_git_permissions" "example-permissions" {
  project_id    = azuredevops_git_repository.example.project_id
  repository_id = azuredevops_git_repository.example.id
  principal     = data.azuredevops_group.example-group.id
  permissions = {
    RemoveOthersLocks = "Allow"
    ManagePermissions = "Deny"
    CreateTag         = "Deny"
    CreateBranch      = "NotSet"
  }
}
```

### Branch level

Permissions for a specific branch inside a Git Repository are specified if all above mentioned the arguments are set.

#### Example usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

resource "azuredevops_git_repository" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Empty Git Repository"
  initialization {
    init_type = "Clean"
  }
}

data "azuredevops_group" "example-group" {
  name = "Project Collection Administrators"
}

resource "azuredevops_git_permissions" "example-permissions" {
  project_id    = azuredevops_git_repository.example.project_id
  repository_id = azuredevops_git_repository.example.id
  branch_name   = "refs/heads/master"
  principal     = data.azuredevops_group.example-group.id
  permissions = {
    RemoveOthersLocks = "Allow"
    ForcePush         = "Deny"
  }
}
```

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

data "azuredevops_group" "example-project-readers" {
  project_id = azuredevops_project.example.id
  name       = "Readers"
}

data "azuredevops_group" "example-project-contributors" {
  project_id = azuredevops_project.example.id
  name       = "Contributors"
}

data "azuredevops_group" "example-project-administrators" {
  project_id = azuredevops_project.example.id
  name       = "Project administrators"
}

resource "azuredevops_git_permissions" "example-permissions" {
  project_id = azuredevops_project.example.id
  principal  = data.azuredevops_group.example-project-readers.id
  permissions = {
    CreateRepository = "Deny"
    DeleteRepository = "Deny"
    RenameRepository = "NotSet"
  }
}

resource "azuredevops_git_repository" "example" {
  project_id     = azuredevops_project.example.id
  name           = "TestRepo"
  default_branch = "refs/heads/master"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_git_permissions" "example-repo-permissions" {
  project_id    = azuredevops_git_repository.example.project_id
  repository_id = azuredevops_git_repository.example.id
  principal     = data.azuredevops_group.example-project-administrators.id
  permissions = {
    RemoveOthersLocks = "Allow"
    ManagePermissions = "Deny"
    CreateTag         = "Deny"
    CreateBranch      = "NotSet"
  }
}

resource "azuredevops_git_permissions" "example-branch-permissions" {
  project_id    = azuredevops_git_repository.example.project_id
  repository_id = azuredevops_git_repository.example.id
  branch_name   = "master"
  principal     = data.azuredevops_group.example-project-contributors.id
  permissions = {
    RemoveOthersLocks = "Allow"
    ForcePush         = "Deny"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project to assign the permissions.
* `repository_id` - (Optional) The ID of the GIT repository to assign the permissions
* `branch_name` - (Optional) The name of the branch to assign the permissions. 

   ~> **Note** to assign permissions to a branch, the `repository_id` must be set as well.

* `principal` - (Required) The **group** principal to assign the permissions.
* `replace` - (Optional) Replace (`true`) or merge (`false`) the permissions. Default: `true`
* `permissions` - (Required) the permissions to assign. The following permissions are available


| Permissions             | Description                                            |
|-------------------------|--------------------------------------------------------|
| Administer              | Administer                                             |
| GenericRead             | Read                                                   |
| GenericContribute       | Contribute                                             |
| ForcePush               | Force push (rewrite history, delete branches and tags) |
| CreateBranch            | Create branch                                          |
| CreateTag               | Create tag                                             |
| ManageNote              | Manage notes                                           |
| PolicyExempt            | Bypass policies when pushing                           |
| CreateRepository        | Create repository                                      |
| DeleteRepository        | Delete repository                                      |
| RenameRepository        | Rename repository                                      |
| EditPolicies            | Edit policies                                          |
| RemoveOthersLocks       | Remove others' locks                                   |
| ManagePermissions       | Manage permissions                                     |
| PullRequestContribute   | Contribute to pull requests                            |
| PullRequestBypassPolicy | Bypass policies when completing pull requests          |

## Relevant Links

* [Azure DevOps Service REST API 7.0 - Security](https://docs.microsoft.com/en-us/rest/api/azure/devops/security/?view=azure-devops-rest-7.0)

## Import

The resource does not support import.

## PAT Permissions Required

- **Project & Team**: vso.security_manage - Grants the ability to read, write, and manage security permissions.
