---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitemquery_folder"
description: |-
  Manages a Work Item Query Folder in Azure DevOps.
---

# azuredevops_workitemquery_folder

Manages a Work Item Query Folder in Azure DevOps.

Folders allow you to organize queries in a hierarchy beneath either the `Shared Queries` or `My Queries` root folder (area).
You must provide exactly one of `area` (either `Shared Queries` or `My Queries`) or `parent_id` (an existing folder's ID) when creating a folder.

## Example Usage

### Basic folder under Shared Queries

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

resource "azuredevops_workitemquery_folder" "team_folder" {
  project_id = azuredevops_project.example.id
  name       = "Team"
  area       = "Shared Queries"
}
```

### Nested folder

```hcl
resource "azuredevops_workitemquery_folder" "parent" {
  project_id = azuredevops_project.example.id
  name       = "Parent"
  area       = "Shared Queries"
}

resource "azuredevops_workitemquery_folder" "child" {
  project_id = azuredevops_project.example.id
  name       = "Child"
  parent_id  = azuredevops_workitemquery_folder.parent.id
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the Project containing the folder.

* `name` - (Required) The display name of the folder.

---

And one of the following must be specified:

* `area` - Root folder. Must be one of `Shared Queries` or `My Queries`.

* `parent_id` - The ID of the parent query folder.

## Attributes Reference

In addition to the arguments above, the following attribute is exported:

* `id` - The ID of the Query Folder. Can be used as the `parent_id` argument when creating queries or sub-folders.

## Relevant Links

* [Azure DevOps REST API - Work Item Query (Queries)](https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/queries?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Work Item Query Folder.
* `read` - (Defaults to 2 minutes) Used when retrieving the Work Item Query Folder.
* `delete` - (Defaults to 5 minutes) Used when deleting the Work Item Query Folder.

## Import

The resource does not support import.

## PAT Permissions Required

* **Work Items**: Read & Write
