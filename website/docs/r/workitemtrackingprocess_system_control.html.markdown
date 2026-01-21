---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitemtrackingprocess_system_control"
description: |-
  Manages a system control customization for a work item type.
---

# azuredevops_workitemtrackingprocess_system_control

Manages a system control customization for a work item type.

System controls are built-in controls like Area Path, Iteration Path, and Reason that can have their visibility and label customized. Unlike regular controls, system controls cannot be removed - only their display properties can be modified.

~> **Note:** This resource modifies system controls. When the resource is deleted, the system control reverts to its default state rather than being removed.

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

resource "azuredevops_workitemtrackingprocess_system_control" "example" {
  process_id                    = azuredevops_workitemtrackingprocess_process.example.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.example.reference_name
  control_id                    = "System.AreaPath"
  visible                       = false
}
```

## Arguments Reference

The following arguments are supported:

* `process_id` - (Required) The ID of the process. Changing this forces a new resource to be created.

* `work_item_type_reference_name` - (Required) The reference name of the work item type. Changing this forces a new resource to be created.

* `control_id` - (Required) The ID of the system control (e.g., `System.AreaPath`, `System.IterationPath`, `System.Reason`). Changing this forces a new resource to be created.

---

* `label` - (Optional) Label for the control.

* `visible` - (Optional) Whether the control should be visible. Defaults to `true`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the system control.

* `control_type` - Type of the control.

* `read_only` - Whether the control is read-only.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - System Controls](https://learn.microsoft.com/en-us/rest/api/azure/devops/processes/system-controls?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the system control customization.
* `read` - (Defaults to 5 minutes) Used when retrieving the system control.
* `update` - (Defaults to 10 minutes) Used when updating the system control.
* `delete` - (Defaults to 10 minutes) Used when reverting the system control to its default state.

## Import

System control customizations can be imported using the complete resource id `process_id/work_item_type_reference_name/control_id`, e.g.

```shell
terraform import azuredevops_workitemtrackingprocess_system_control.example 00000000-0000-0000-0000-000000000000/MyProcess.CustomWorkItemType/System.AreaPath
```
