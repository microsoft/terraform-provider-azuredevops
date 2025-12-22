---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitemtrackingprocess_group"
description: |-
  Manages a group within a page and section for a work item type.
---

# azuredevops_workitemtrackingprocess_group

Manages a group within a page and section for a work item type.

## Example Usage

### Basic Group

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

resource "azuredevops_workitemtrackingprocess_group" "example" {
  process_id                    = azuredevops_workitemtrackingprocess_process.example.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.example.reference_name
  page_id                       = azuredevops_workitemtrackingprocess_workitemtype.example.pages[0].id
  section_id                    = azuredevops_workitemtrackingprocess_workitemtype.example.pages[0].sections[0].id
  label                         = "Custom Group"
}
```

### Group with Controls

This example creates a group with multiple controls, including an HTML field control (System.Description) which can only be added when creating the group.

```hcl
resource "azuredevops_workitemtrackingprocess_group" "example" {
  process_id                    = azuredevops_workitemtrackingprocess_process.example.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.example.reference_name
  page_id                       = azuredevops_workitemtrackingprocess_workitemtype.example.pages[0].id
  section_id                    = azuredevops_workitemtrackingprocess_workitemtype.example.pages[0].sections[0].id
  label                         = "Details Group"

  control {
    id    = "System.Description"
    label = "Description"
  }

  control {
    id    = "System.Title"
    label = "Title"
  }
}
```

### Group with Contribution Control (Extension)

```hcl
resource "azuredevops_extension" "multivalue" {
  publisher_id = "ms-devlabs"
  extension_id = "vsts-extensions-multivalue-control"
}

resource "azuredevops_workitemtrackingprocess_group" "example" {
  depends_on                    = [azuredevops_extension.multivalue]
  process_id                    = azuredevops_workitemtrackingprocess_process.example.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.example.reference_name
  page_id                       = azuredevops_workitemtrackingprocess_workitemtype.example.pages[0].id
  section_id                    = azuredevops_workitemtrackingprocess_workitemtype.example.pages[0].sections[0].id
  label                         = "Extension Group"

  control {
    id              = "ms-devlabs.vsts-extensions-multivalue-control.multivalue-form-control"
    is_contribution = true
    contribution {
      contribution_id = "ms-devlabs.vsts-extensions-multivalue-control.multivalue-form-control"
      height          = 50
      inputs = {
        FieldName = "System.Tags"
        Values    = "Option1;Option2;Option3"
      }
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `process_id` - (Required) The ID of the process. Changing this forces a new group to be created.

* `work_item_type_reference_name` - (Required) The reference name of the work item type. Changing this forces a new group to be created.

* `page_id` - (Required) The ID of the page to add the group to. Changing this moves the group to the new page.

* `section_id` - (Required) The ID of the section to add the group to. Changing this moves the group to the new section.

* `label` - (Required) Label for the group.

---

* `order` - (Optional) Order in which the group should appear in the section.

* `visible` - (Optional) A value indicating if the group should be visible or not. Default: `true`

* `control` - (Optional) Controls to be created with the group. Required for HTML controls which cannot be added to existing groups. A `control` block as defined below.

---

A `control` block supports the following:

* `id` - (Required) The ID of the control. This is the field reference name (e.g., System.Description) or the contribution ID for extension controls.

* `label` - (Optional) Label for the control.

* `visible` - (Optional) A value indicating if the control should be visible or not. Default: `true`

* `read_only` - (Optional) A value indicating if the control is read only. Default: `false`

* `order` - (Optional) Order in which the control should appear in the group.

* `metadata` - (Optional) Inner text of the control.

* `watermark` - (Optional) Watermark text for the textbox.

* `is_contribution` - (Optional) A value indicating if the control is a contribution (extension) control. Default: `false`

* `contribution` - (Optional) Contribution configuration for extension controls. A `contribution` block as defined below.

---

A `contribution` block supports the following:

* `contribution_id` - (Required) The ID of the contribution (extension).

* `height` - (Optional) The height for the contribution.

* `inputs` - (Optional) A dictionary holding key value pairs for contribution inputs.

* `show_on_deleted_work_item` - (Optional) A value indicating if the contribution should be shown on deleted work items. Default: `false`

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the group.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Groups](https://learn.microsoft.com/en-us/rest/api/azure/devops/processes/groups?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the group.
* `read` - (Defaults to 5 minutes) Used when retrieving the group.
* `update` - (Defaults to 10 minutes) Used when updating the group.
* `delete` - (Defaults to 10 minutes) Used when deleting the group.

## Import

Groups can be imported using the complete resource id `process_id/work_item_type_reference_name/page_id/section_id/group_id`, e.g.

```shell
terraform import azuredevops_workitemtrackingprocess_group.example 00000000-0000-0000-0000-000000000000/MyProcess.CustomWorkItemType/page-id/section-id/group-id
```
