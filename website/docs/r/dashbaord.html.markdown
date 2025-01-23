---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_dashboard"
description: |-
  Manages Dashboard within Azure DevOps project.
---

# azuredevops_dashboard

Manages Dashboard within Azure DevOps project.

~> **NOTE:** Project level Dashboard allows to be created with the same name. Dashboard held by a team must have a different name.

## Example Usage


### Manage Project dashboard

```hcl
resource "azuredevops_project" "example" {
  name        = "Example Project"
  description = "Managed by Terraform"
}

resource "azuredevops_dashboard" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example dashboard"
}
```

### Manage Team dashboard

```hcl
resource "azuredevops_project" "example" {
  name        = "Example Project"
  description = "Managed by Terraform"
}

resource "azuredevops_team" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example team"
}

resource "azuredevops_dashboard" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example dashboard"
  team_id    = azuredevops_team.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Wiki.

* `project_id` - (Required) The ID of the Project. Changing this forces a new resource to be created.

---

* `description` - (Optional) The description of the dashboard.
 
* `team_id` - (Optional) The ID of the Team.

* `refresh_interval` - (Optional) The interval for client to automatically refresh the dashboard. Expressed in minutes. Possible values are: `0`, `5`.Defaults to `0`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Dashboard.
* `owner_id` - The owner of the Dashboard, could be the project or a team.

## Relevant Links

- [Azure DevOps dashboards REST API 7.1 - Wiki ](https://learn.microsoft.com/en-us/rest/api/azure/devops/dashboard/dashboards?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Dashboard.
* `read` - (Defaults to 2 minute) Used when retrieving the Dashboard.
* `update` - (Defaults to 5 minutes) Used when updating the Dashboard.
* `delete` - (Defaults to 5 minutes) Used when deleting the Dashboard.

## Import

Azure DevOps Dashboard can be imported using the `projectId/dasboardId` or `projectId/teamId/dasboardId`

```shell
terraform import azuredevops_dashboard.dashboard 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```

or 

```shell
terraform import azuredevops_dashboard.dashboard 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```
