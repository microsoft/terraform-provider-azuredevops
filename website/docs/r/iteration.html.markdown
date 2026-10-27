---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_iteration"
description: |-
  Manages an Iteration Path in Azure DevOps.
---

# azuredevops_iteration

Manages an Iteration Path (classification node) in Azure DevOps.

Iteration paths allow you to group work items into sprints or other time-boxed intervals. They form a hierarchy under the project's root iteration node.

## Example Usage

### Basic iteration at root level

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

resource "azuredevops_iteration" "example" {
  project_id = azuredevops_project.example.id
  name       = "Release 1"
}
```

### Nested iteration paths with dates

```hcl
resource "azuredevops_iteration" "parent" {
  project_id = azuredevops_project.example.id
  name       = "Release 1"
}

resource "azuredevops_iteration" "child" {
  project_id          = azuredevops_project.example.id
  name                = "Sprint 1"
  parent_iteration_id = azuredevops_iteration.parent.id
  start_date          = "2026-01-05"
  finish_date         = "2026-01-16"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project. Changing this forces a new resource to be created.
* `name` - (Required) The name of the iteration path node. Must conform to [naming restrictions](https://learn.microsoft.com/en-us/azure/devops/organizations/settings/about-areas-iterations?view=azure-devops#naming-restrictions).
* `parent_iteration_id` - (Optional) The integer ID of the parent iteration node. If not specified, the iteration is created at the root level. Changing this forces a new resource to be created.
* `start_date` - (Optional) The start date of the iteration, in `YYYY-MM-DD` format. Must be set together with `finish_date`.
* `finish_date` - (Optional) The finish date of the iteration, in `YYYY-MM-DD` format. Must be set together with `start_date`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The integer ID of the iteration path node.
* `path` - The path of the iteration, in backslash-separated format including the project name (e.g., `Example Project\Release 1\Sprint 1`).

## Import

Iteration paths can be imported using the project ID and the iteration's integer node ID, e.g.:

```shell
terraform import azuredevops_iteration.example 00000000-0000-0000-0000-000000000000/42
```
