---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitemtrackingprocess_inherited_page"
description: |-
  Manages an inherited page customization for a work item type.
---

# azuredevops_workitemtrackingprocess_inherited_page

Manages inherited page customizations for a work item type.

Inherited pages are pages that exist in the parent process template and are inherited by derived processes.

~> **Note:** When the resource is deleted, the page reverts to its inherited state rather than being removed.

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

# Relabel the Details page
resource "azuredevops_workitemtrackingprocess_inherited_page" "example" {
  process_id        = azuredevops_workitemtrackingprocess_process.example.id
  work_item_type_id = azuredevops_workitemtrackingprocess_workitemtype.example.reference_name
  page_id           = azuredevops_workitemtrackingprocess_workitemtype.example.pages[0].id
  label             = "Custom Details"
}
```

## Arguments Reference

The following arguments are supported:

* `process_id` - (Required) The ID of the process. Changing this forces a new resource to be created.

* `work_item_type_id` - (Required) The ID (reference name) of the work item type. Changing this forces a new resource to be created.

* `page_id` - (Required) The ID of the inherited page to customize. Changing this forces a new resource to be created.

* `label` - (Required) Label for the page.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the page.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Pages](https://learn.microsoft.com/en-us/rest/api/azure/devops/processes/pages?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when customizing the inherited page.
* `read` - (Defaults to 5 minutes) Used when retrieving the page.
* `update` - (Defaults to 10 minutes) Used when updating the page.
* `delete` - (Defaults to 10 minutes) Used when reverting the page to its inherited state.

## Import

Inherited page customizations can be imported using the complete resource id `process_id/work_item_type_id/page_id`, e.g.

```shell
terraform import azuredevops_workitemtrackingprocess_inherited_page.example 00000000-0000-0000-0000-000000000000/MyProcess.CustomWorkItemType/page-id
```
