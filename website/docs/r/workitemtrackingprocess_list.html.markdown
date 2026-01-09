---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitemtrackingprocess_list"
description: |-
  Manages an organization-scoped picklist for work item tracking processes.
---

# azuredevops_workitemtrackingprocess_list

Manages an organization-scoped picklist for work item tracking processes.

~> **Note** Picklists are organization-scoped resources, not process-specific. They can be referenced by custom fields across multiple processes.

## Example Usage

### Basic Picklist

```hcl
resource "azuredevops_workitemtrackingprocess_list" "example" {
  name  = "Priority Levels"
  items = ["Low", "Medium", "High", "Critical"]
}
```

### Picklist with Suggestions

```hcl
resource "azuredevops_workitemtrackingprocess_list" "example" {
  name         = "Environment"
  items        = ["Development", "Staging", "Production"]
  is_suggested = true
}
```

### Integer Picklist

```hcl
resource "azuredevops_workitemtrackingprocess_list" "example" {
  name  = "Story Points"
  type  = "integer"
  items = ["1", "2", "3", "5", "8", "13", "21"]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Name of the picklist.

* `items` - (Required) A list of picklist items.

---

* `type` - (Optional) DataType of the picklist. Valid values: `string`, `integer`. Defaults to `string`. Changing this forces a new resource to be created.

* `is_suggested` - (Optional) Indicates whether items outside of the suggested list are allowed. Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the picklist.

* `url` - URL of the picklist.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Lists](https://learn.microsoft.com/en-us/rest/api/azure/devops/processes/lists?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the picklist.
* `read` - (Defaults to 5 minutes) Used when retrieving the picklist.
* `update` - (Defaults to 10 minutes) Used when updating the picklist.
* `delete` - (Defaults to 10 minutes) Used when deleting the picklist.

## Import

Picklists can be imported using the picklist ID, e.g.

```shell
terraform import azuredevops_workitemtrackingprocess_list.example 00000000-0000-0000-0000-000000000000
```
