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
  project_id  = data.azuredevops_project.example.id
  title       = "Example Work Item"
  description = "Managed by Terraform"
  type        = "Issue"
  state       = "Active"
  tags        = ["Tag"]
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
### With Parent Work Item

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
}

resource "azuredevops_workitem" "epic" {
  project_id = azuredevops_project.example.id
  title      = "Example EPIC Title"
  type       = "Epic"
  state      = "New"
}

resource "azuredevops_workitem" "example" {
  project_id = azuredevops_project.example.id
  title      = "Example Work Item"
  type       = "Issue"
  state      = "Active"
  tags       = ["Tag"]
  parent_id  = azuredevops_workitem.epic.id
}
```

### With Additional Fields

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
  type       = "User Story"
  state      = "New"
  tags       = ["Tag"]
  additional_fields_json = jsonencode({
    "Microsoft.VSTS.Scheduling.StoryPoints"    = 5
    "Microsoft.VSTS.Common.AcceptanceCriteria" = "This is our definition of done"
    "Microsoft.VSTS.Common.Priority"           = 2
    "Microsoft.VSTS.Common.ValueArea"          = "Business"
  })
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the Project.

* `title` - (Required) The Title of the Work Item.

* `type` - (Required) The Type of the Work Item. The work item type varies depending on the process used when creating the project(`Agile`, `Basic`, `Scrum`, `Scrum`). See [Work Item Types](https://learn.microsoft.com/en-us/azure/devops/boards/work-items/about-work-items?view=azure-devops) for more details.

---

* `additional_fields_json` - (Optional) A JSON-formatted string of extra fields. **Note**: Removing this attribute from your configuration will not clear existing fields in the API. To remove all fields, set this value to an empty JSON string (`"{}"`).

* `area_path` - (Optional) Specifies the area where the Work Item is used.

* `custom_fields` - (Optional, **Deprecated** use `additional_fields_json` argument instead) Specifies a list with Custom Fields for the Work Item.

* `description` - (Optional) A description for the Work Item. **Note**: Due to current lifecycle behavior, omitting this field or setting it to an empty string will not clear the description in Azure DevOps; the provider will instead read the existing value. To avoid a breaking change, the ability to clear this field will be introduced in a future major release.

* `iteration_path` - (Optional) Specifies the iteration in which the Work Item is used.

* `parent_id` - (Optional) The parent work item.

* `state` - (Optional) The state of the Work Item. The four main states that are defined for the User Story (`Agile`) are `New`, `Active`, `Resolved`, and `Closed`. See [Workflow states](https://learn.microsoft.com/en-us/azure/devops/boards/work-items/workflow-and-state-categories?view=azure-devops&tabs=agile-process#workflow-states) for more details.

* `tags` - (Optional) Specifies a list of Tags.
  
## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Work Item.

* `url` - The URL of the Work Item.

* `relations` - A `relations` blocks as documented below.


---

An `relations` block supports the following:

* `rel` - The type of relationship. For example: `System.LinkTypes.Hierarchy-Reverse` is a parent relationship. More details [item link type](https://learn.microsoft.com/en-us/azure/devops/boards/queries/link-type-reference?view=azure-devops#example).

* `url` - The URL of the Work Item.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Work Item.
* `read` - (Defaults to 5 minute) Used when retrieving the Work Item.
* `update` - (Defaults to 10 minutes) Used when updating the Work Item.
* `delete` - (Defaults to 10 minutes) Used when deleting the Work Item.

## Import

Azure DevOps Work Item can be imported using the Project ID and Work Item ID, e.g.

```sh
terraform import azuredevops_workitem.example 00000000-0000-0000-0000-000000000000/0
```
