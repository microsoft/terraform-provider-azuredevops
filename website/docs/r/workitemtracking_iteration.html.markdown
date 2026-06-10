---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitemtracking_iteration"
description: |-
  Manages a Work Item Tracking Iteration.
---

# azuredevops_workitemtracking_iteration

Manages a Work Item Tracking Iteration.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

resource "azuredevops_workitemtracking_iteration" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Iteration"
  path       = "/"
  attributes {
    start_date  = "2023-01-01T00:00:00Z"
    finish_date = "2023-01-31T00:00:00Z"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the Project. Changing this forces a new Iteration to be created.
* `name` - (Required) The name of the Iteration.
* `path` - (Optional) The path of the Iteration. Changing this forces a new Iteration to be created.
* `attributes` - (Optional) A `attributes` block as defined below.

---

A `attributes` block supports the following:

* `start_date` - (Optional) The start date of the Iteration.
* `finish_date` - (Optional) The finish date of the Iteration.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Node ID of the Iteration.
* `identifier` - The ID (UUID) of the Iteration.
* `has_children` - Indicator if the child Iteration node exists.

## Import

Project Iterations can be imported using the `project_id/id` or `project_id/path`, e.g.

```shell
terraform import azuredevops_workitemtracking_iteration.example 00000000-0000-0000-0000-000000000000/12345
```
