---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_team_area_path"
description: |-
  Manages a single area path assignment for a team within a project in Azure DevOps.
---

# azuredevops_team_area_path

Manages a single area path assignment for a team within a project in Azure DevOps.

~> **Note** This resource only assigns an existing area path to a team. It does not create the area path node itself. Area path nodes must be created separately (e.g., via the Azure DevOps UI or API).

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
}

resource "azuredevops_team" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Team"
}

resource "azuredevops_team_area_path" "example" {
  project_id       = azuredevops_project.example.id
  team_id          = azuredevops_team.example.id
  area_path        = "Example Project"
  include_children = true
}
```

### Nested Area Path

~> **Note** The nested area path node must already exist before it can be assigned to a team.

```hcl
resource "azuredevops_team_area_path" "nested" {
  project_id       = azuredevops_project.example.id
  team_id          = azuredevops_team.example.id
  area_path        = "Example Project\\Example Team"
  include_children = true
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The Project ID. Changing this forces a new resource to be created.

* `team_id` - (Required) The ID of the Team. Changing this forces a new resource to be created.

* `area_path` - (Required) The area path to assign to the team (e.g., `"ProjectName"` for the root or `"ProjectName\\AreaName"` for a nested path). The area path node must already exist. Changing this forces a new resource to be created.

* `include_children` - (Optional) Whether child area paths are included. Defaults to `true`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the resource in the format `<project_id>/<team_id>/<url-encoded-area_path>`.

## Import

Azure DevOps Team Area Paths can be imported using the resource ID, e.g.

```shell
terraform import azuredevops_team_area_path.example 00000000-0000-0000-0000-000000000000/11111111-1111-1111-1111-111111111111/Example+Project
```

The format is `<project_id>/<team_id>/<url-encoded-area_path>`.
