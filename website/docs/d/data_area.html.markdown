---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_area"
description: |-
  Use this data source to access information about an existing Area (Component) within Azure DevOps.
---

# Data Source: azuredevops_area

Use this data source to access information about an existing Area (Component) within Azure DevOps.

## Example Usage

```hcl
resource "random_id" "rand_id" {
  keepers = {
    seed = var.random-anchor
  }

  byte_length = 6
}

locals {
  project_name = "test-acc-project-${random_id.rand_id.hex}"
}

#---------------------------------------------------------------------------
# Azure DevOps project
resource "azuredevops_project" "project" {
  project_name       = local.project_name
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "My first project"
}

data "azuredevops_area" "area" {
  project_id = azuredevops_project.project.id
  path = "/"
  fetch_children = "false"
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The project ID.
- `path` - (Optional) The path to the Area; _Format_: URL relative; if omitted, or value `"/"` is used, the root Area will be returned
- `fetch_children` - (Optional) Read children nodes, _Depth_: 1, _Default_: `true`

## Attributes Reference

The following attributes are exported:

- `id` - The id of the Area node
- `name` - The name of the Area node
- `has_children` - Indicator if a Area node has child nodes
- `children` - A list of `children` blocks as defined below, empty if `has_children == false`

A `children` block supports the following:

- `id` - The id of the child Area node
- `name` - The name of the child Area node
- `project_id` - The project ID of the child Area node
- `path` - The complete path (in relative URL format) of the child Area
- `has_children` - Indicator if the child Area node has child nodes

## Relevant Links

- [Azure DevOps Service REST API 5.1 - Classification Nodes - Get Classification Nodes](https://docs.microsoft.com/en-us/rest/api/azure/devops/wit/classification%20nodes/get%20classification%20nodes?view=azure-devops-rest-5.1)

## PAT Permissions Required

- **Project & Team**: vso.work - Grants the ability to read work items, queries, boards, area and iterations paths, and other work item tracking related metadata. Also grants the ability to execute queries, search work items and to receive notifications about work item events via service hooks. 
