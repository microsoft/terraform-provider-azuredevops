---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_iteration"
description: |-
  Manages a Project Iteration.
---

# azuredevops_iteration

Manages a Project Iteration.

## Example Usage

```hcl
resource "azuredevops_iteration" "example" {
  project_id = "Iteration"
  name       = "example"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Project Iteration.

* `project_id` - (Required) The ID of the Iteration. Changing this forces a new Project Iteration to be created.

---

* `attributes` - (Optional) A `attributes` block as defined below.

* `path` - (Optional) Iteration. Changing this forces a new Project Iteration to be created.

---

A `attributes` block supports the following:

* `finish_date` - (Optional) Iteration.

* `start_date` - (Optional) Iteration.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Project Iteration.

* `has_children` - Iteration.

* `node_id` - The ID of the Iteration.



## Import

Project Iterations can be imported using the `resource id`, e.g.

```shell
terraform import azuredevops_iteration.example 00000000-0000-0000-0000-000000000000
```
