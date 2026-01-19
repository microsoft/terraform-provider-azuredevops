---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitemtrackingprocess_inherited_state"
description: |-
  Manages an inherited state for a work item type.
---

# azuredevops_workitemtrackingprocess_inherited_state

Manages inherited states for a work item type.

Inherited states are predefined states that exist in all work item types (e.g., "New", "Active", "Closed"). This resource allows you to hide inherited states.

~> **Note:** When the resource is deleted, the state remains in Azure DevOps. Inherited states cannot be deleted from Azure DevOps.

~> **Note:** Only states with `customizationType` of `system` or `inherited` can be managed by this resource. Use `azuredevops_workitemtrackingprocess_state` to manage custom states.

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

# Hide an inherited state
resource "azuredevops_workitemtrackingprocess_inherited_state" "example" {
  process_id                    = azuredevops_workitemtrackingprocess_process.example.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.example.reference_name
  name                          = "Removed"
  hidden                        = true
}
```

## Arguments Reference

The following arguments are supported:

* `process_id` - (Required) The ID of the process. Changing this forces a new resource to be created.

* `work_item_type_reference_name` - (Required) The reference name of the work item type. Changing this forces a new resource to be created.

* `name` - (Required) Name of the inherited state to manage. This is used to look up the state and must match an existing inherited state name. Changing this forces a new resource to be created.

* `hidden` - (Optional) Whether the state is hidden.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the state.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - States](https://learn.microsoft.com/en-us/rest/api/azure/devops/processes/states?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when adopting the inherited state.
* `read` - (Defaults to 5 minutes) Used when retrieving the state.
* `update` - (Defaults to 10 minutes) Used when updating the state.
* `delete` - (Defaults to 10 minutes) Used when removing the resource from Terraform state.

## Import

Inherited states can be imported using the complete resource id `process_id/work_item_type_reference_name/state_name`, e.g.

```shell
terraform import azuredevops_workitemtrackingprocess_inherited_state.example 00000000-0000-0000-0000-000000000000/MyProcess.CustomWorkItemType/New
```
