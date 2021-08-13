---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_teams"
description: |-
  Use this data source to access information about existing Teams in a Project or globally within an Azure DevOps organization
---

# Data Source: azuredevops_team

Use this data source to access information about existing Teams in a Project or globally within an Azure DevOps organization

## Example Usage

```hcl
data "azuredevops_teams" "test" {
}

output "project_id" {
  value = data.azuredevops_teams.test.teams.*.project_id
}

output "name" {
  value = data.azuredevops_teams.test.teams.*.name
}

output "administrators" {
  value = data.azuredevops_teams.test.teams.*.administrators
}

output "administrators" {
  value = data.azuredevops_teams.test.teams.*.members
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Optional) The Project ID. If no project ID all teams of the organization will be returned.

## Attributes Reference

The following attributes are exported:

- `teams` - A list of existing projects in your Azure DevOps Organization with details about every project which includes:

  - `project_id` - Project identifier.
  - `id - Team identifier
  - `name` - Team name.
  - `description` - Team description.
  - `administrators` - List of subject descriptors for `administrators` of the team.
  - `members` - List of subject descriptors for `members` of the team.

## Relevant Links

- [Azure DevOps Service REST API 5.1 - Teams - Get](https://docs.microsoft.com/en-us/rest/api/azure/devops/core/teams/get?view=azure-devops-rest-5.1)

## PAT Permissions Required

- **vso.project**:	Grants the ability to read projects and teams.
