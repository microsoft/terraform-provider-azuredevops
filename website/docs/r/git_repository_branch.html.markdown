---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_git_repository_branch"
description: |-
  Manages a Git Repository Branch.
---

# azuredevops_git_repository_branch

Manages a Git Repository Branch.

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

resource "azuredevops_git_repository_branch" "example" {
  repository_id = azuredevops_git_repository.example.id
  name          = "example-branch-name"
  ref_branch    = azuredevops_git_repository.example.default_branch
}

resource "azuredevops_git_repository_branch" "example_from_commit_id" {
  repository_id = azuredevops_git_repository.example.id
  name          = "example-from-commit-id"
  ref_commit_id = azuredevops_git_repository_branch.example.last_commit_id
}
```

## Arguments Reference

The following arguments are supported:

- `name` - (Required) The name of the branch in short format not prefixed with `refs/heads/`.

- `repository_id` - (Required) The ID of the repository the branch is created against.

- `ref_branch` - (Optional) The reference to the source branch to create the branch from, in `<name>` or `refs/heads/<name>` format. Conflict with `ref_tag`, `ref_commit_id`.

- `ref_tag` - (Optional) The reference to the tag to create the branch from, in `<name>` or `refs/tags/<name>` format. Conflict with `ref_branch`, `ref_commit_id`.

- `ref_commit_id` - (Optional) The commit object ID to create the branch from. Conflict with `ref_branch`, `ref_tag`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

- `id` - The ID of the Git Repository Branch, in the format `<repository_id>:<name>`.

- `last_commit_id` - The commit object ID of last commit on the branch.
