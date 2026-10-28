---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitem"
description: |-
  Use this data source to access information about an existing Work Item within Azure DevOps.
---

# Data Source: azuredevops_workitem

Use this data source to access information about an existing Work Item within Azure DevOps.

## Example Usage

```hcl
data "azuredevops_project" "example" {
  name = "Example Project"
}

data "azuredevops_workitem" "example" {
  id         = 1234
  project_id = data.azuredevops_project.example.id
}

output "title" {
  value = data.azuredevops_workitem.example.title
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the Work Item.

* `project_id` - (Required) The ID of the Project the Work Item belongs to.

## Attributes Reference

The following attributes are exported:

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

- [Azure DevOps Service REST API 7.1 - Work Items - Get Work Item](https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/work-items/get-work-item?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minute) Used when retrieving the Work Item.

## PAT Permissions Required

- **Project & Team**: vso.work - Grants the ability to read work items, queries, boards, area and iterations paths, and other work item tracking related metadata.
