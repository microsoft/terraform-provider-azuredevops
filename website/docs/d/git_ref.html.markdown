---
layout: "azuredevops"
page_title: "AzureDevops: Data Source: azuredevops_git_ref"
description: |-
  Gets information about an existing Git Ref.
---

# Data Source: azuredevops_git_ref

Use this data source to access information about an existing Git Ref.

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

data "azuredevops_git_ref" "example" {
  repository_id = azuredevops_git_repository.example.id
  name          = "refs/heads/master"
}
```

## Arguments Reference

The following arguments are supported:

* `repository_id` - (Required) The ID of the Git Repository.
* `name` - (Required) The full name of the Git Ref (e.g., `refs/heads/master`).
* `project_id` - (Optional) The ID of the Project.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Git Ref.
* `object_id` - The commit ID the ref points to.
* `peeled_object_id` - The peeled object ID of the ref (for annotated tags).
* `creator` - The ID of the creator of the ref.
* `url` - The URL of the ref.
* `is_locked` - Whether the ref is locked.
* `is_locked_by` - The ID of the user who locked the ref.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Git Ref.
