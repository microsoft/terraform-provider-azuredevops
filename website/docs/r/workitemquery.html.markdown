---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitemquery"
description: |-
  Manages a Work Item Query in Azure DevOps.
---

# azuredevops_workitemquery

Manages a Work Item Query (WIQL based list / tree query) in Azure DevOps.

A query can live either directly under one of the root areas `Shared Queries` or `My Queries`, or inside another query folder. You must provide exactly one of `area` (either `Shared Queries` or `My Queries`) or `parent_id` (an existing folder's ID) when creating a query.

The WIQL (Work Item Query Language) statement is used to define the query logic. See the [WIQL Syntax Reference](https://learn.microsoft.com/en-us/azure/devops/boards/queries/wiql-syntax?view=azure-devops) for more information.

## Example Usage

### Basic query under Shared Queries

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
}

resource "azuredevops_workitemquery" "all_issues" {
  project_id = azuredevops_project.example.id
  name       = "All Active Issues"
  area       = "Shared Queries"
  wiql       = <<-WIQL
    SELECT [System.Id], [System.Title], [System.State]
    FROM WorkItems
    WHERE [System.WorkItemType] = 'Issue'
      AND [System.State] <> 'Closed'
    ORDER BY [System.ChangedDate] DESC
  WIQL
}
```

### Query inside a custom folder

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
}

resource "azuredevops_workitemquery_folder" "team_folder" {
  project_id = azuredevops_project.example.id
  name       = "Team"
  area       = "Shared Queries"
}

resource "azuredevops_workitemquery" "my_team_bugs" {
  project_id = azuredevops_project.example.id
  name       = "Team Bugs"
  parent_id  = azuredevops_workitemquery_folder.team_folder.id
  wiql       = <<-WIQL
    SELECT [System.Id], [System.Title], [System.State], [System.AssignedTo]
    FROM WorkItems
    WHERE [System.WorkItemType] = 'Bug'
      AND [System.State] <> 'Closed'
    ORDER BY [System.CreatedDate] DESC
  WIQL
}
```

### Applying permissions to a query

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
}

resource "azuredevops_workitemquery_folder" "team_folder" {
  project_id = azuredevops_project.example.id
  name       = "Team"
  area       = "Shared Queries"
}

resource "azuredevops_workitemquery" "my_team_bugs" {
  project_id = azuredevops_project.example.id
  name       = "Team Bugs"
  parent_id  = azuredevops_workitemquery_folder.team_folder.id
  wiql       = <<-WIQL
    SELECT [System.Id], [System.Title], [System.State], [System.AssignedTo]
    FROM WorkItems
    WHERE [System.WorkItemType] = 'Bug'
      AND [System.State] <> 'Closed'
    ORDER BY [System.CreatedDate] DESC
  WIQL
}

data "azuredevops_group" "example-readers" {
  project_id = azuredevops_project.example.id
  name       = "Readers"
}

resource "azuredevops_workitemquery_permissions" "query_permissions" {
  project_id = azuredevops_project.example.id
  # Permissions can only be set for folders and queries under 'Shared Queries'.
  # The path here is relative to the 'Shared Queries' folder.
  path = format(
    "%s/%s",
    azuredevops_workitemquery_folder.team_folder.name,
    azuredevops_workitemquery.my_team_bugs.name
  )
  principal  = data.azuredevops_group.example-readers.id
  permissions = {
    "Read"   = "Allow"
    "Contribute"   = "Deny"
    "ManagePermissions" = "Deny"
    "Delete" = "Deny"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the Project containing the query.

* `name` - (Required) The display name of the query.

* `wiql` - (Required) The WIQL (Work Item Query Language) statement. Length 1–32000 characters.

---

And one of the following must be specified:

* `area` - Root folder for the query. Must be one of `Shared Queries` or `My Queries`.

* `parent_id` - The ID of the parent query folder under which to create the query.

## Attributes Reference

In addition to the arguments above, the following attribute is exported:

* `id` - The ID of the Work Item Query.

## Relevant Links

* [Azure DevOps REST API - Work Item Query (Queries)](https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/queries?view=azure-devops-rest-7.1)
* [WIQL Syntax Reference](https://learn.microsoft.com/en-us/azure/devops/boards/queries/wiql-syntax?view=azure-devops)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Work Item Query.
* `read` - (Defaults to 2 minutes) Used when retrieving the Work Item Query.
* `update` - (Defaults to 5 minutes) Used when updating the Work Item Query.
* `delete` - (Defaults to 5 minutes) Used when deleting the Work Item Query.

## Import

A Work Item Query can be imported using the following format `projectId/queryId`.

For example:

```sh
terraform import azuredevops_workitemquery.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```

## PAT Permissions Required

* **Work Items**: Read & Write
