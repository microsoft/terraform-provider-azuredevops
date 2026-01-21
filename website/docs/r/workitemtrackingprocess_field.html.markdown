---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitemtrackingprocess_field"
description: |-
  Manages a field for a work item type in a process.
---

# azuredevops_workitemtrackingprocess_field

Manages a field for a work item type in a process. This resource adds an existing field to a work item type and allows configuring field-specific settings like default value, required, and read-only status.

## Example Usage

```hcl
resource "azuredevops_workitemtrackingprocess_process" "example" {
  name                   = "example-process"
  parent_process_type_id = "adcc42ab-9882-485e-a3ed-7678f01f66bc"
}

resource "azuredevops_workitemtrackingprocess_workitemtype" "example" {
  process_id  = azuredevops_workitemtrackingprocess_process.example.id
  name        = "example"
  description = "Example work item type"
}

resource "azuredevops_workitemtracking_field" "example" {
  name           = "Priority Level"
  reference_name = "Custom.PriorityLevel"
  type           = "string"
}

resource "azuredevops_workitemtrackingprocess_field" "example" {
  process_id              = azuredevops_workitemtrackingprocess_process.example.id
  work_item_type_ref_name = azuredevops_workitemtrackingprocess_workitemtype.example.reference_name
  reference_name          = azuredevops_workitemtracking_field.example.reference_name
  required                = true
  default_value           = "Medium"
}
```

## Arguments Reference

The following arguments are supported:

* `process_id` - (Required) The ID of the process. Changing this forces a new field to be created.

* `work_item_type_ref_name` - (Required) The reference name of the work item type. Changing this forces a new field to be created.

* `reference_name` - (Required) The reference name of the field. Changing this forces a new field to be created.

---

* `default_value` - (Optional) The default value of the field.

* `read_only` - (Optional) If true, the field cannot be edited. Default: `false`.

* `required` - (Optional) If true, the field cannot be empty. Default: `false`.

* `allow_groups` - (Optional) Allow setting field value to a group identity. Only applies to identity fields.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The reference name of the field.

* `name` - The name of the field.

* `type` - The type of the field.

* `description` - The description of the field.

* `customization` - Indicates the type of customization on this work item. Possible values are `system`, `inherited`, or `custom`.

* `is_locked` - Indicates whether the field definition is locked for editing.

* `url` - The URL of the field resource.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Fields](https://learn.microsoft.com/en-us/rest/api/azure/devops/processes/fields?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the field.
* `read` - (Defaults to 5 minutes) Used when retrieving the field.
* `update` - (Defaults to 10 minutes) Used when updating the field.
* `delete` - (Defaults to 10 minutes) Used when deleting the field.

## Import

Fields can be imported using the complete resource id `process_id/work_item_type_ref_name/field_ref_name`, e.g.

```shell
terraform import azuredevops_workitemtrackingprocess_field.example 00000000-0000-0000-0000-000000000000/MyProcess.CustomWorkItemType/Custom.MyField
```
