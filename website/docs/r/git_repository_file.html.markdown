---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_git_repository_file"
description: |- Manage files within an Azure DevOps Git repository.
---

# azuredevops_git_repository_file

Manage files within an Azure DevOps Git repository.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_git_repository" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Git Repository"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_git_repository_file" "example" {
  repository_id       = azuredevops_git_repository.example.id
  file                = ".gitignore"
  content             = "**/*.tfstate"
  branch              = "refs/heads/master"
  commit_message      = "First commit"
  overwrite_on_create = false
}
```

## Argument Reference

The following arguments are supported:

- `repository_id` - (Required) The ID of the Git repository.
- `file` - (Required) The path of the file to manage.
- `content` - (Required) The file content.
- `branch` - (Optional) Git branch (defaults to `refs/heads/master`). The branch must already exist, it will not be created if it
  does not already exist.
- `commit_message` - (Optional) Commit message when adding or updating the managed file.
- `overwrite_on_create` - (Optional) Enable overwriting existing files (defaults to `false`).

## Import

Repository files can be imported using a combination of the `repository ID` and `file`, e.g.

```sh
terraform import azuredevops_git_repository_file.example 00000000-0000-0000-0000-000000000000/.gitignore
```

To import a file from a branch other than `master`, append `:` and the branch name, e.g.

```sh
terraform import azuredevops_git_repository_file.example 00000000-0000-0000-0000-000000000000/.gitignore:refs/heads/master
```

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Git API](https://docs.microsoft.com/en-us/rest/api/azure/devops/git/?view=azure-devops-rest-7.0)
