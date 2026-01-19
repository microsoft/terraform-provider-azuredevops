---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitemtrackingprocess_list"
description: |-
  Manages an organization-scoped list for work item tracking processes.
---

# azuredevops_workitemtrackingprocess_list

Manages an organization-scoped list for work item tracking processes.

~> **Note** Lists are organization-scoped resources, not process-specific. They can be referenced by custom fields across multiple processes.

## Example Usage

### Basic List

```hcl
resource "azuredevops_workitemtrackingprocess_list" "example" {
  name  = "Priority Levels"
  items = ["Low", "Medium", "High", "Critical"]
}
```

### List with Suggestions

```hcl
resource "azuredevops_workitemtrackingprocess_list" "example" {
  name         = "Environment"
  items        = ["Development", "Staging", "Production"]
  is_suggested = true
}
```

### Integer List

```hcl
resource "azuredevops_workitemtrackingprocess_list" "example" {
  name  = "Story Points"
  type  = "integer"
  items = ["1", "2", "3", "5", "8", "13", "21"]
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Name of the list.

* `items` - (Required) A list of items.

---

* `type` - (Optional) Data type of the list. Valid values: `string`, `integer`. Defaults to `string`. Changing this forces a new resource to be created.

* `is_suggested` - (Optional) Indicates whether items outside of the suggested list are allowed. Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the list.

* `url` - URL of the list.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Lists](https://learn.microsoft.com/en-us/rest/api/azure/devops/processes/lists?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the list.
* `read` - (Defaults to 5 minutes) Used when retrieving the list.
* `update` - (Defaults to 10 minutes) Used when updating the list.
* `delete` - (Defaults to 10 minutes) Used when deleting the list.

## Import

Lists can be imported using the list ID, e.g.

```shell
terraform import azuredevops_workitemtrackingprocess_list.example 00000000-0000-0000-0000-000000000000
```
