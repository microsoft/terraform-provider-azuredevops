---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitemtrackingprocess_state"
description: |-
  Manages a state for a work item type.
---

# azuredevops_workitemtrackingprocess_state

Manages a state for a work item type.

## Example Usage

```hcl
resource "azuredevops_workitemtrackingprocess_process" "example" {
  name                   = "example-process"
  parent_process_type_id = "adcc42ab-9882-485e-a3ed-7678f01f66bc"
}

resource "azuredevops_workitemtrackingprocess_workitemtype" "example" {
  process_id = azuredevops_workitemtrackingprocess_process.example.id
  name       = "example"
}

resource "azuredevops_workitemtrackingprocess_state" "example" {
  process_id                    = azuredevops_workitemtrackingprocess_process.example.id
  work_item_type_id = azuredevops_workitemtrackingprocess_workitemtype.example.reference_name
  name                          = "Ready"
  color                         = "#5688E0"
  state_category                = "Proposed"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Name of the state. Changing this forces a new state to be created.

* `process_id` - (Required) The ID of the process. Changing this forces a new state to be created.

* `work_item_type_id` - (Required) The ID (reference name) of the work item type. Changing this forces a new state to be created.

* `state_category` - (Required) Category of the state. Valid values: `Proposed`, `InProgress`, `Resolved`, `Completed`, `Removed`.

* `color` - (Required) Color hexadecimal code to represent the state, e.g. `#b2b2b2`.

---

* `order` - (Optional) Order within the category where the state should appear.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the state.

* `url` - URL of the state.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - States](https://learn.microsoft.com/en-us/rest/api/azure/devops/processes/states?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the state.
* `read` - (Defaults to 5 minutes) Used when retrieving the state.
* `update` - (Defaults to 10 minutes) Used when updating the state.
* `delete` - (Defaults to 10 minutes) Used when deleting the state.

## Import

States can be imported using the complete resource id `process_id/work_item_type_id/state_id`, e.g.

```shell
terraform import azuredevops_workitemtrackingprocess_state.example 00000000-0000-0000-0000-000000000000/MyProcess.CustomWorkItemType/00000000-0000-0000-0000-000000000000
```
