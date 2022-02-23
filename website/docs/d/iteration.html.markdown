---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_iteration"
description: |-
  Use this data source to access information about an existing Iteration (Sprint) within Azure DevOps.
---

# Data Source: azuredevops_iteration

Use this data source to access information about an existing Iteration (Sprint) within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

data "azuredevops_iteration" "example-root-iteration" {
  project_id     = azuredevops_project.example.id
  path           = "/"
  fetch_children = true
}

data "azuredevops_iteration" "example-child-iteration" {
  project_id     = azuredevops_project.example.id
  path           = "/Iteration 1"
  fetch_children = true
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The project ID.
- `path` - (Optional) The path to the Iteration, _Format_: URL relative; if omitted, or value `"/"` is used, the root Iteration will be returned
- `fetch_children` - (Optional) Read children nodes, _Depth_: 1, _Default_: `true`

## Attributes Reference

The following attributes are exported:

- `id` - The id of the Iteration node
- `name` - The name of the Iteration node
- `has_children` - Indicator if a Iteration node has child nodes
- `children` - A list of `children` blocks as defined below, empty if `has_children == false`

A `children` block supports the following:

- `id` - The id of the child Iteration node
- `name` - The name of the child Iteration node
- `project_id` - The project ID of the child Iteration node
- `path` - The complete path (in relative URL format) of the child Iteration
- `has_children` - Indicator if the child Iteration node has child nodes

## Relevant Links

- [Azure DevOps Service REST API 6.0 - Classification Nodes - Get Classification Nodes](https://docs.microsoft.com/en-us/rest/api/azure/devops/wit/classification-nodes/get-classification-nodes?view=azure-devops-rest-6.0)

## PAT Permissions Required

- **Project & Team**: vso.work - Grants the ability to read work items, queries, boards, area and iterations paths, and other work item tracking related metadata. Also grants the ability to execute queries, search work items and to receive notifications about work item events via service hooks. 
