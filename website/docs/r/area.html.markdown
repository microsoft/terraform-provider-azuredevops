---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_area"
description: |-
  Manages a Project Area.
---

# azuredevops_area

Manages a Project Area.

## Example Usage

```hcl
resource "azuredevops_area" "example" {
  project_id = "Area"
  name = "example"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Project Area.

* `project_id` - (Required) The ID of the Area. Changing this forces a new Project Area to be created.

---

* `path` - (Optional) Area. Changing this forces a new Project Area to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Project Area.

* `has_children` - Area.

* `node_id` - The ID of the Area.



## Import

Project Areas can be imported using the `resource id`, e.g.

```shell
terraform import azuredevops_area.example 00000000-0000-0000-0000-000000000000
```
