---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_team_members"
description: |-
  Manages members of a team within a project in a Azure DevOps organization.
---

# azuredevops_team_members

Manages members of a team within a project in a Azure DevOps organization.

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  name               = "Test Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "My first project"
}

data "azuredevops_group" "builtin_project_readers" {
  project_id = azuredevops_project.project.id
  name       = "Readers"
}

resource "azuredevops_team" "team" {
  project_id = azuredevops_project.project.id
  name       = "${azuredevops_project.project.name} Team 2"
}

resource "azuredevops_team_members" "team_members" {
  project_id = azuredevops_team.team.project_id
  team_id    = azuredevops_team.team.id
  mode       = "overwrite"
  members    = [
    azuredevops_group.builtin_project_readers.descriptor
  ]
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The Project ID.
- `team_id` - (Required) The ID of the Team.
- `members` - (Required) List of subject descriptors to define members of the team.

  > NOTE: It's possible to define team members both within the
  > `azuredevops_team` resource via the `members` block and by using the
  > `azuredevops_team_members` resource. However it's not possible to use
  > both methods to manage team members, since there'll be conflicts.
- `mode` - (Optional) The mode how the resource manages team members.
  - `mode == add`: the resource will ensure that all specified members will be part of the referenced team
  - `mode == overwrite`: the resource will replace all existing members with the members specified within the `members` block

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - A random ID for this resource. There is no "natural" ID, so a random one is assigned.

## Relevant Links

- [Azure DevOps Service REST API 5.1 - Teams - Update](https://docs.microsoft.com/en-us/rest/api/azure/devops/core/teams/update?view=azure-devops-rest-5.1)

## Import

The resource does not support import.

## PAT Permissions Required

- **vso.project_write**:	Grants the ability to read and update projects and teams. 
