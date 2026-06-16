---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_area"
description: |-
  Manages an Area Path in Azure DevOps.
---

# azuredevops_area

Manages an Area Path (classification node) in Azure DevOps.

Area paths allow you to group work items by team or product area. They form a hierarchy under the project's root area node.

## Example Usage

### Basic area at root level

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

resource "azuredevops_area" "example" {
  project_id = azuredevops_project.example.id
  name       = "Frontend"
}
```

### Nested area paths

```hcl
resource "azuredevops_area" "parent" {
  project_id = azuredevops_project.example.id
  name       = "Engineering"
}

resource "azuredevops_area" "child" {
  project_id     = azuredevops_project.example.id
  name           = "Frontend"
  parent_area_id = azuredevops_area.parent.area_id
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project. Changing this forces a new resource to be created.
* `name` - (Required) The name of the area path node. Must conform to [naming restrictions](https://learn.microsoft.com/en-us/azure/devops/organizations/settings/about-areas-iterations?view=azure-devops#naming-restrictions).
* `parent_area_id` - (Optional) The integer ID of the parent area node. If not specified, the area is created at the root level. Changing this forces a new resource to be created. Use the `area_id` attribute of another `azuredevops_area` resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The UUID identifier of the area path node.
* `area_id` - The integer ID of this area node. Use this to reference as `parent_area_id` in child areas.
* `full_path` - The full path of the area, relative to the project (e.g., `/Engineering/Frontend`).

## Import

Area paths can be imported using the project ID and the area's integer node ID, e.g.:

```shell
terraform import azuredevops_area.example 00000000-0000-0000-0000-000000000000/42
```
