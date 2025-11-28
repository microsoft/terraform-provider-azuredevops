---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitemtrackingprocess_workitemtype"
description: |-
  Manages a work item type for a process.
---

# azuredevops_workitemtrackingprocess_workitemtype

Manages a work item type for a process.

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
  color       = "#FF5733"
  icon        = "icon_clipboard"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required)  Name of work item type. Changing this forces a new work item type to be created.

* `process_id` - (Required)  The ID of the process the work item type belongs to. Changing this forces a new work item type to be created.

---

* `color` - (Optional)  Color hexadecimal code to represent the work item type. Default: "#009ccc"

* `description` - (Optional)  Description of the work item type.

* `icon` - (Optional)  Icon to represent the work item type. Default: "icon_clipboard"

* `inherits_from` - (Optional)  Parent work item type for work item type. Changing this forces a new work item type to be created.

* `is_disabled` - (Optional)  True if the work item type need to be disabled. Default: false

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the work item type.

* `reference_name` -  Reference name of the work item type.

* `url` -  Url of the work item type.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Work Item Types](https://learn.microsoft.com/en-us/rest/api/azure/devops/processes/work-item-types?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the work item type.
* `read` - (Defaults to 5 minutes) Used when retrieving the work item type.
* `update` - (Defaults to 10 minutes) Used when updating the work item type.
* `delete` - (Defaults to 10 minutes) Used when deleting the work item type.

## Import

work item types can be imported using the complete resource id `process_id/reference_name`, e.g.

```shell
terraform import azuredevops_workitemtrackingprocess_workitemtype.example 00000000-0000-0000-0000-000000000000/MyProcess.CustomWorkItemType
```
