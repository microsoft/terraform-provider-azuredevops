---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_team"
description: |-
  Manages a team within a project in a Azure DevOps organization.
---


# azuredevops_team

Manages a team within a project in a Azure DevOps organization.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

data "azuredevops_group" "example-project-contributors" {
  project_id = azuredevops_project.example.id
  name       = "Contributors"
}

data "azuredevops_group" "example-project-readers" {
  project_id = azuredevops_project.example.id
  name       = "Readers"
}

resource "azuredevops_team" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Team"
  administrators = [
    data.azuredevops_group.example-project-contributors.descriptor
  ]
  members = [
    data.azuredevops_group.example-project-readers.descriptor
  ]
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The Project ID.

* `name` - (Required) The name of the Team.

---

* `description`- (Optional) The description of the Team.

* `administrators` - (Optional) List of subject descriptors to define administrators of the team.

  ~> **NOTE:** It's possible to define team administrators both within the
   `azuredevops_team` resource via the `administrators` block and by using the
   `azuredevops_team_administrators` resource. However it's not possible to use
   both methods to manage team administrators, since there'll be conflicts.

---

* `members` - (Optional) List of subject descriptors to define members of the team.

  ~> **NOTE:** It's possible to define team members both within the
   `azuredevops_team` resource via the `members` block and by using the
   `azuredevops_team_members` resource. However it's not possible to use
   both methods to manage team members, since there'll be conflicts.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Team.
* `descriptor` - The descriptor of the Team.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Teams - Create](https://docs.microsoft.com/en-us/rest/api/azure/devops/core/teams/create?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Team.
* `read` - (Defaults to 5 minute) Used when retrieving the Team.
* `update` - (Defaults to 10 minutes) Used when updating the Team.
* `delete` - (Defaults to 10 minutes) Used when deleting the Team.

## Import

Azure DevOps teams can be imported using the complete resource id `<project_id>/<team_id>` e.g.

```sh
terraform import azuredevops_team.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```

## PAT Permissions Required

- **vso.project_manage**:	Grants the ability to create, read, update, and delete projects and teams.
