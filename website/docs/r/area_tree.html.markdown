---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_area_tree"
description: |-
  Manages an entire Area Path hierarchy for a project in Azure DevOps.
---

# azuredevops_area_tree

Manages an entire Area Path (classification node) hierarchy for a project as a single resource.

Unlike [`azuredevops_area`](area.html.markdown), which manages one area path node per resource instance, `azuredevops_area_tree` takes the complete, desired area path hierarchy - expressed as an infinitely-deep nested object - and owns its full lifecycle: creating every node (including implied ancestors), pruning nodes that are removed on update, and tearing down the whole tree on delete.

Because there is only ever one resource instance per project (no `for_each` over individual nodes), there is no possibility of the "Cycle" error that can arise when Terraform resources of the same type reference sibling instances of themselves, and orphaned ancestor nodes are correctly pruned on update since this resource has full visibility into both the old and new desired state.

## Example Usage

### Basic tree

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

resource "azuredevops_area_tree" "example" {
  project_id = azuredevops_project.example.id

  paths = jsonencode({
    "Team A" = {}
    "Team B" = {}
  })
}
```

### Arbitrarily deep, nested tree

Every key in the object is an area name, and its value is the (possibly empty) object describing that area's children. Nesting can go as deep as needed; every intermediate/ancestor node (e.g. `Team A` and `Team A/Sub Area` below) is created and removed automatically, so it does not need to be listed as its own top-level entry.

```hcl
resource "azuredevops_area_tree" "example" {
  project_id = azuredevops_project.example.id

  paths = jsonencode({
    "Team A" = {
      "Sub Area" = {
        "Grandchild" = {}
      }
    }
    "Team B" = {}
  })
}
```

Using this tree with a team's default area:

```hcl
resource "azuredevops_team" "team_a" {
  project_id = azuredevops_project.example.id
  name       = "Team A"

  area = [
    azuredevops_area_tree.example.area_paths["Team A"],
  ]
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project. Changing this forces a new resource to be created.
* `paths` - (Required) The full area path hierarchy to manage, expressed as a JSON-encoded, infinitely-deep object tree (typically built with `jsonencode`). Every key is an area name and its value is the (possibly empty) object describing that area's children.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the project (same as `project_id`).
* `area_ids` - Map of every managed area path (including implied ancestors, keyed by slash-separated path relative to the project root, e.g. `Team A/Sub Area`) to its integer node ID.
* `area_paths` - Map of every managed area path (keyed the same way as `area_ids`) to its full canonical Azure DevOps path (e.g. `Example Project\Area\Team A\Sub Area`), suitable for direct use as an entry in an `azuredevops_team` resource's `area` attribute.

## Import

Area trees can be imported using the project ID, e.g.:

```shell
terraform import azuredevops_area_tree.example 00000000-0000-0000-0000-000000000000
```
