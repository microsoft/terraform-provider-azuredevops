---
layout: "azuredevops"
page_title: "AzureDevops: Data Source: azuredevops_git_repository_file"
description: |-
  Gets an existing Git Repository File.
---

# Data Source: azuredevops_git_repository_file

Use this data source to get an existing Git Repository File.

## Example Usage

```hcl
data "azuredevops_project" "example" {
  name = "Example Project"
}

# Load a specific Git repository by name
data "azuredevops_git_repository" "example" {
  project_id = data.azuredevops_project.example.id
  name       = "Example Repository"
}

data "azuredevops_git_repository_file" "example" {
  repository_id = data.azuredevops_git_repository.example.id
  branch = "refs/heads/main"
  file   = "MyFile.txt"
}
```

## Arguments Reference

The following arguments are supported:

* `file` - (Required) The path of the file to get.

* `repository_id` - (Required) The ID of the Git repository.

---

* `branch` - (Optional) The git branch to use. Conflicts with `tag`; one or the other must be specified.

* `tag` - (Optional) The tag to use.Conflicts with `branch`; one or the other must be specified.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Git Repository File. Note this is different from that of the corresponding resource being a combination of the `repository ID` and `file`, followed by either ':branch:' or ':tag:' and then the `branch` or `tag` used.

* `last_commit_message` - The commit message for the file.

* `content` - The file content.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the Git Repository File.
