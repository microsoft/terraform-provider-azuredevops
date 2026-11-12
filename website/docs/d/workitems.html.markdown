---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitems"
description: |-
  Use this data source to access information about existing Work Items within Azure DevOps.
---

# Data Source: azuredevops_workitems

Use this data source to access information about existing Work Items within Azure DevOps by running either a
[Work Item Query Language (WIQL)](https://learn.microsoft.com/en-us/azure/devops/boards/queries/wiql-syntax) query
or a saved query.

## Example Usage

```hcl
data "azuredevops_project" "example" {
  name = "Example Project"
}

# Flat query
data "azuredevops_workitems" "active_bugs" {
  project_id = data.azuredevops_project.example.id
  wiql       = "SELECT [System.Id] FROM WorkItems WHERE [System.WorkItemType] = 'Bug' AND [System.State] = 'Active'"
}

# Recursive tree query, returns a Work Item and all of its transitive children
data "azuredevops_workitems" "tree" {
  project_id = data.azuredevops_project.example.id
  wiql       = <<-EOT
    SELECT [System.Id]
    FROM WorkItemLinks
    WHERE [Source].[System.Id] = 1234
      AND [System.Links.LinkType] = 'System.LinkTypes.Hierarchy-Forward'
    MODE (Recursive)
  EOT
}

# Saved query
resource "azuredevops_workitemquery" "active_tasks" {
  project_id = data.azuredevops_project.example.id
  area       = "Shared Queries"
  name       = "Active Tasks"
  wiql       = "SELECT [System.Id] FROM WorkItems WHERE [System.WorkItemType] = 'Task' AND [System.State] = 'Active'"
}

data "azuredevops_workitems" "active_tasks" {
  project_id = data.azuredevops_project.example.id
  query_id   = azuredevops_workitemquery.active_tasks.id
}

output "titles" {
  value = [for item in data.azuredevops_workitems.tree.work_items : item.title]
}
```

## Argument Reference

The following arguments are supported:

* `wiql` - (Optional) The WIQL statement used to look up the Work Items. Both flat (`FROM WorkItems`) and link (`FROM WorkItemLinks`) queries are supported.

* `query_id` - (Optional) The ID of a saved query used to look up the Work Items, for example from an `azuredevops_workitemquery` resource. Query folders are not supported.

~> **NOTE:** Exactly one of `wiql` or `query_id` must be specified.

---

* `project_id` - (Optional) The ID of the Project the query is scoped to. Required when the query uses the `@project` macro.

## Attributes Reference

The following attributes are exported:

* `work_items` - A list of `work_items` blocks as defined below.

---

A `work_items` block exports the following:

* `id` - The ID of the Work Item.

* `title` - The title of the Work Item.

* `description` - The description of the Work Item.

* `type` - The type of the Work Item, for example `Issue`, `Bug` or `Task`.

* `state` - The state of the Work Item.

* `area_path` - The area path of the Work Item.

* `iteration_path` - The iteration path of the Work Item.

* `parent_id` - The ID of the parent Work Item. `0` when the Work Item has no parent.

* `tags` - A set of tags assigned to the Work Item.

* `url` - The REST API URL of the Work Item.

* `fields` - A map of all Work Item fields, keyed by their reference name, for example `System.Title` or `Custom.MyField`. Values which are not strings are JSON encoded.

* `relations` - A list of `relations` blocks as defined below.

---

A `relations` block exports the following:

* `rel` - The reference name of the relation, for example `System.LinkTypes.Hierarchy-Forward`.

* `url` - The URL of the related object.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Wiql - Query By Wiql](https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/wiql/query-by-wiql?view=azure-devops-rest-7.1)
- [Azure DevOps Service REST API 7.1 - Wiql - Query By Id](https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/wiql/query-by-id?view=azure-devops-rest-7.1)
- [Azure DevOps Service REST API 7.1 - Work Items - List](https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/work-items/list?view=azure-devops-rest-7.1)
- [Work Item Query Language (WIQL) syntax reference](https://learn.microsoft.com/en-us/azure/devops/boards/queries/wiql-syntax)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 10 minute) Used when retrieving the Work Items.

## PAT Permissions Required

- **Project & Team**: vso.work - Grants the ability to read work items, queries, boards, area and iterations paths, and other work item tracking related metadata. Also grants the ability to execute queries.
