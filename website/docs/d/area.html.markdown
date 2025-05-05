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
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

data "azuredevops_area" "example" {
  project_id     = azuredevops_project.example.id
  path           = "/"
  fetch_children = "false"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The project ID.

---

* `path` - (Optional) The path to the Area; _Format_: URL relative; if omitted, or value `"/"` is used, the root Area will be returned

* `fetch_children` - (Optional) Read children nodes, _Depth_: 1, _Default_: `true`

## Attributes Reference

The following attributes are exported:

* `id` - The id of the Area node

* `name` - The name of the Area node

* `has_children` - Indicator if an Area node has child nodes

* `children` - A list of `children` blocks as defined below, empty if `has_children == false`

---

A `children` block supports the following:

* `id` - The ID of the child Area node

* `name` - The name of the child Area node

* `project_id` - The ID of project.

* `path` - The complete path (in relative URL format) of the child Area

* `has_children` - Indicator if the child Area node has child nodes

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Classification Nodes - Get Classification Nodes](https://docs.microsoft.com/en-us/rest/api/azure/devops/wit/classification-nodes/create-or-update?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minute) Used when retrieving the Area.

## PAT Permissions Required

- **Project & Team**: vso.work - Grants the ability to read work items, queries, boards, area and iterations paths, and other work item tracking related metadata. Also grants the ability to execute queries, search work items and to receive notifications about work item events via service hooks. 
