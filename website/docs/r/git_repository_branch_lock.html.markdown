---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_git_repository_branch_lock"
description: |-
  Manages a Git Repository Branch Lock.
---

# azuredevops_git_repository_branch_lock

Manages a Git Repository Branch Lock.

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

resource "azuredevops_git_repository_branch_lock" "example" {
  is_locked = true
  repository_id = azuredevops_git_repository.example.id
  branch = azuredevops_git_repository_branch.example.name
}
```

## Arguments Reference

The following arguments are supported:

* `branch` - (Required) The name of the branch to lock. Changing this forces a new Git Repository Branch Lock to be created.

* `is_locked` - (Required) Whether the branch is locked. Changing this forces a new Git Repository Branch Lock to be created.

* `repository_id` - (Required) The ID of the Git Repository. Changing this forces a new Git Repository Branch Lock to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Git Repository Branch Lock.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Git Repository Branch Lock.
* `read` - (Defaults to 5 minutes) Used when retrieving the Git Repository Branch Lock.
* `delete` - (Defaults to 5 minutes) Used when deleting the Git Repository Branch Lock.

## Import

Git Repository Branch Locks can be imported using the `resource id`, e.g.

```shell
terraform import azuredevops_git_repository_branch_lock.example 00000000-0000-0000-0000-000000000000
```
