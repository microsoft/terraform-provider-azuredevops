---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitemtrackingprocess_inherited_control"
description: |-
  Manages an inherited control customization for a work item type.
---

# azuredevops_workitemtrackingprocess_inherited_control

Manages an inherited control customization for a work item type.

Inherited controls are controls that exist in the parent process template and are inherited by derived processes. This resource allows you to customize inherited controls.

~> **Note:** This resource customizes inherited controls. When the resource is deleted, the control reverts to its inherited state rather than being removed.

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

# Customize the first control in the first group
resource "azuredevops_workitemtrackingprocess_inherited_control" "example" {
  process_id                    = azuredevops_workitemtrackingprocess_process.example.id
  work_item_type_id = azuredevops_workitemtrackingprocess_workitemtype.example.reference_name
  group_id                      = azuredevops_workitemtrackingprocess_workitemtype.example.pages[0].sections[0].groups[0].id
  control_id                    = azuredevops_workitemtrackingprocess_workitemtype.example.pages[0].sections[0].groups[0].controls[0].id
  visible                       = false
}
```

## Arguments Reference

The following arguments are supported:

* `process_id` - (Required) The ID of the process. Changing this forces a new resource to be created.

* `work_item_type_id` - (Required) The ID (reference name) of the work item type. Changing this forces a new resource to be created.

* `group_id` - (Required) The ID of the group containing the control. Changing this forces a new resource to be created.

* `control_id` - (Required) The ID of the inherited control to customize. Changing this forces a new resource to be created.

---

* `label` - (Optional) Label for the control.

* `visible` - (Optional) Whether the control should be visible.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the control.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Controls](https://learn.microsoft.com/en-us/rest/api/azure/devops/processes/controls?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when customizing the inherited control.
* `read` - (Defaults to 5 minutes) Used when retrieving the control.
* `update` - (Defaults to 10 minutes) Used when updating the control.
* `delete` - (Defaults to 10 minutes) Used when reverting the control to its inherited state.

## Import

Inherited control customizations can be imported using the complete resource id `process_id/work_item_type_id/group_id/control_id`, e.g.

```shell
terraform import azuredevops_workitemtrackingprocess_inherited_control.example 00000000-0000-0000-0000-000000000000/MyProcess.CustomWorkItemType/group-id/System.Title
```
