---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitemtracking_area"
description: |-
  Manages a Work Item Tracking Area.
---

# azuredevops_workitemtracking_area

Manages a Work Item Tracking Area.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

resource "azuredevops_workitemtracking_area" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Area"
  path       = "/"
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the Project. Changing this forces a new Area to be created.
* `name` - (Required) The name of the Area.
* `path` - (Optional) The path of the Area. Changing this forces a new Area to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The Node ID of the Area.
* `identifier` - The ID (UUID) of the Area.
* `has_children` - Indicator if the child Area node exists.

## Import

Project Areas can be imported using the `project_id/id` or `project_id/path`, e.g.

```shell
terraform import azuredevops_workitemtracking_area.example 00000000-0000-0000-0000-000000000000/12345
```
