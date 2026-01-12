---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitemtracking_field"
description: |-
  Manages a work item field in Azure DevOps.
---

# azuredevops_workitemtracking_field

Manages a work item field in Azure DevOps.

~> **Note:** Custom fields are created at the organization level, not the project level. The Azure DevOps API does not support project-scoped custom fields.

## Example Usage

### Basic Field

```hcl
resource "azuredevops_workitemtracking_field" "example" {
  name           = "My Custom Field"
  reference_name = "Custom.MyCustomField"
  type           = "string"
}
```

### Restore a Deleted Field

```hcl
resource "azuredevops_workitemtracking_field" "restored" {
  name           = "Restored Field"
  reference_name = "Custom.RestoredField"
  type           = "string"
  restore        = true
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The friendly name of the field. Changing this forces a new field to be created.

* `reference_name` - (Required) The reference name of the field (e.g., `Custom.MyField`). Changing this forces a new field to be created.

* `type` - (Required) The type of the field. Possible values: `string`, `integer`, `dateTime`, `plainText`, `html`, `treePath`, `history`, `double`, `guid`, `boolean`, `identity`. Changing this forces a new field to be created.

---

* `description` - (Optional) The description of the field. Changing this forces a new field to be created.

* `usage` - (Optional) The usage of the field. Possible values: `none`, `workItem`, `workItemLink`, `tree`, `workItemTypeExtension`. Default: `workItem`. Changing this forces a new field to be created.

* `read_only` - (Optional) Indicates whether the field is read-only. Default: `false`. Changing this forces a new field to be created.

* `can_sort_by` - (Optional) Indicates whether the field can be sorted in server queries. Default: `true`. Changing this forces a new field to be created.

* `is_queryable` - (Optional) Indicates whether the field can be queried in the server. Default: `true`. Changing this forces a new field to be created.

* `is_identity` - (Optional) Indicates whether this field is an identity field. Default: `false`. Changing this forces a new field to be created.

* `is_picklist` - (Optional) Indicates whether this field is a picklist. Default: `false`. Changing this forces a new field to be created.

* `is_picklist_suggested` - (Optional) Indicates whether this field is a suggested picklist. Default: `false`. Changing this forces a new field to be created.

* `picklist_id` - (Optional) The identifier of the picklist associated with this field, if applicable. Changing this forces a new field to be created.

* `is_locked` - (Optional) Indicates whether this field is locked for editing. Default: `false`.

* `restore` - (Optional) Set to `true` to restore a previously deleted field instead of creating a new one. When set to `true`, the resource will attempt to restore the field with the specified `reference_name`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the field.

* `url` - The URL of the field resource.

* `supported_operations` - The supported operations on this field. A `supported_operations` block as defined below.

---

A `supported_operations` block exports:

* `name` - The friendly name of the operation.

* `reference_name` - The reference name of the operation.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Fields](https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/fields?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the field.
* `read` - (Defaults to 5 minutes) Used when retrieving the field.
* `update` - (Defaults to 10 minutes) Used when updating the field.
* `delete` - (Defaults to 10 minutes) Used when deleting the field.

## Import

Fields can be imported using the reference name:

```shell
terraform import azuredevops_workitemtracking_field.example Custom.MyCustomField
```
