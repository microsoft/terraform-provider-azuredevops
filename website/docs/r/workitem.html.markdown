---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitem"
description: |-
  Manages a Work Item in Azure Devops.
---

# azuredevops_workitem

Manages a Work Item in Azure Devops.

## Example Usage

### Basic usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}
resource "azuredevops_workitem" "example" {
  project_id = data.azuredevops_project.example.id
  title      = "Example Work Item"
  type       = "Issue"
  state      = "Active"
  tags       = ["Tag"]
}
```

### With custom fields

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}
resource "azuredevops_workitem" "example" {
  project_id = data.azuredevops_project.example.id
  title      = "Example Work Item"
  type       = "Issue"
  state      = "Active"
  tags       = ["Tag"]
  custom_fields = {
    example : "example"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the Project.

* `title` - (Required) The Title of the Work Item.

* `type` - (Required) The Type of the Work Item. The work item type varies depending on the process used when creating the project(`Agile`, `Basic`, `Scrum`, `Scrum`). See [Work Item Types](https://learn.microsoft.com/en-us/azure/devops/boards/work-items/about-work-items?view=azure-devops) for more details.

---

* `area_path` - (Optional) Specifies the area where the Work Item is used.

* `custom_fields` - (Optional) Specifies a list with Custom Fields for the Work Item.

* `iteration_path` - (Optional) Specifies the iteration in which the Work Item is used.

* `state` - (Optional) The state of the Work Item. The four main states that are defined for the User Story (`Agile`) are `New`, `Active`, `Resolved`, and `Closed`. See [Workflow states](https://learn.microsoft.com/en-us/azure/devops/boards/work-items/workflow-and-state-categories?view=azure-devops&tabs=agile-process#workflow-states) for more details.

* `tags` - (Optional) Specifies a list of Tags.
  
## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Work Item.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Work Item.
* `read` - (Defaults to 5 minute) Used when retrieving the Work Item.
* `update` - (Defaults to 10 minutes) Used when updating the Work Item.
* `delete` - (Defaults to 10 minutes) Used when deleting the Work Item.

## Import

Work Item resource does not support import.
