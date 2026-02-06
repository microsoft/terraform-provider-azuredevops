---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_work_team_settings"
description: |-
  Manages a Project Team Settings.
---

# azuredevops_work_team_settings

Manages a Project Team Settings.

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

resource "azuredevops_work_team_settings" "example" {
  project_id = azuredevops_project.example.id
  team_id    = azuredevops_team.example.id
  backlog_iteration_id = "<some guid>" #need iteration data lookup by name
  default_iteration_id = "<some guid>" #need iteration data lookup by name
  
  backlog_visibilities = [
    "Microsoft.EpicCategory",
  ]

  working_days =  [
    "monday",
    "tuesday",
    "wednesday",
    "thursday",
    "friday",
  ]

  "bugs_behavior" = "asRequirements"
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the Project. Changing this forces a new Project Team Settings to be created.

* `team_id` - (Required) The ID of the Team. Changing this forces a new Project Team Settings to be created.

---

* `backlog_iteration_id` - (Optional) Determines which work items appear on your team’s backlog and board.

* `bugs_behavior` - (Optional) Set your team's preference for how they manage bugs. Your selection determines where bugs appear in the hierarchy and on backlogs and boards. [Learn more about the bug management setting|https://learn.microsoft.com/en-us/azure/devops/organizations/settings/show-bugs-on-backlog?view=azure-devops].

* `default_iteration_id` - (Optional) Assigns a default Iteration value to work items created from your team context.Conflicts with `default_iteration_macro`

* `default_iteration_macro` - (Optional) default Iteration Macro.Conflicts with `default_iteration_id`

* `working_days` - (Optional) Specifies a list of Capacity and burndown are based on the days your team works..

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Project Team Settings.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Project Team Settings.
* `read` - (Defaults to 5 minutes) Used when retrieving the Project Team Settings.
* `update` - (Defaults to 10 minutes) Used when updating the Project Team Settings.
* `delete` - (Defaults to 10 minutes) Used when deleting the Project Team Settings.

## Import

Project Team Settingss can be imported using the `resource id`, e.g.

```shell
terraform import azuredevops_work_team_settings.example 00000000-0000-0000-0000-000000000000
```
