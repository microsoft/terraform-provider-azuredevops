---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_git_ref"
description: |-
  Manages a Git Ref.
---

# azuredevops_git_ref

Manages a Git Ref.

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

resource "azuredevops_git_ref" "example" {
  repository_id = azuredevops_git_repository.example.id
  name          = "refs/heads/example"
  ref_branch    = azuredevops_git_repository.example.default_branch
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Git Ref. Changing this forces a new Git Ref to be created.

* `repository_id` - (Required) The ID of the Git Repository. Changing this forces a new Git Ref to be created.

---

* `ref_branch` - (Optional) The name of the branch to create the ref from. Conflicts with `ref_tag`,`ref_commit_id`. Changing this forces a new Git Ref to be created.

* `ref_commit_id` - (Optional) The commit ID to create the ref from. Conflicts with `ref_branch`,`ref_tag`. Changing this forces a new Git Ref to be created.

* `ref_tag` - (Optional) The name of the tag to create the ref from. Conflicts with `ref_branch`,`ref_commit_id`. Changing this forces a new Git Ref to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Git Ref.

* `object_id` - The commit ID the ref points to.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Git Ref.
* `read` - (Defaults to 5 minutes) Used when retrieving the Git Ref.
* `delete` - (Defaults to 10 minutes) Used when deleting the Git Ref.

## Import

Git Refs can be imported using the `resource id`, e.g.

```shell
terraform import azuredevops_git_ref.example 00000000-0000-0000-0000-000000000000:refs/heads/main
```
