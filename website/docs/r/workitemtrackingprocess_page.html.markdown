---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitemtrackingprocess_page"
description: |-
  Manages a page in the work item form layout for a work item type.
---

# azuredevops_workitemtrackingprocess_page

Manages a page in the work item form layout for a work item type.

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

resource "azuredevops_workitemtrackingprocess_page" "example" {
  process_id                    = azuredevops_workitemtrackingprocess_process.example.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.example.reference_name
  label                         = "Custom Page"
}
```

## Arguments Reference

The following arguments are supported:

* `process_id` - (Required) The ID of the process. Changing this forces a new page to be created.

* `work_item_type_reference_name` - (Required) The reference name of the work item type. Changing this forces a new page to be created.

* `label` - (Required) The label for the page.

---

* `order` - (Optional) Order in which the page should appear in the layout.

* `visible` - (Optional) A value indicating if the page should be visible or not. Default: `true`

* `locked` - (Optional) A value indicating whether any user operations are permitted on this page. Default: `false`

* `page_type` - (Optional) The type of the page. Possible values are `custom`, `history`, `links`, `attachments`. Default: `custom`

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the page.

* `inherited` - A value indicating whether this page has been inherited from a parent layout.

* `overridden` - A value indicating whether this page has been overridden by a child layout.

* `section` - The sections of the page. A `section` block as defined below.

---

A `section` block exports the following:

* `id` - The ID of the section.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Pages](https://learn.microsoft.com/en-us/rest/api/azure/devops/processes/pages?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the page.
* `read` - (Defaults to 5 minutes) Used when retrieving the page.
* `update` - (Defaults to 10 minutes) Used when updating the page.
* `delete` - (Defaults to 10 minutes) Used when deleting the page.

## Import

Pages can be imported using the complete resource id `process_id/work_item_type_reference_name/page_id`, e.g.

```shell
terraform import azuredevops_workitemtrackingprocess_page.example 00000000-0000-0000-0000-000000000000/MyProcess.CustomWorkItemType/page-id
```
