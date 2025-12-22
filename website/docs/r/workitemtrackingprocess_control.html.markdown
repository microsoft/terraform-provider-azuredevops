---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitemtrackingprocess_control"
description: |-
  Manages a control within a group for a work item type.
---

# azuredevops_workitemtrackingprocess_control

Manages a control within a group for a work item type. Controls can be field controls or contribution controls (extensions).

## Example Usage

### Basic Field Control

```hcl
resource "azuredevops_workitemtrackingprocess_process" "example" {
  name                   = "example-process"
  parent_process_type_id = "adcc42ab-9882-485e-a3ed-7678f01f66bc"
}

resource "azuredevops_workitemtrackingprocess_workitemtype" "example" {
  process_id  = azuredevops_workitemtrackingprocess_process.example.id
  name        = "example"
}

resource "azuredevops_workitemtrackingprocess_group" "example" {
  process_id                    = azuredevops_workitemtrackingprocess_process.example.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.example.reference_name
  page_id                       = azuredevops_workitemtrackingprocess_workitemtype.example.pages[0].id
  section_id                    = azuredevops_workitemtrackingprocess_workitemtype.example.pages[0].sections[0].id
  label                         = "Custom Group"
}

resource "azuredevops_workitemtrackingprocess_control" "example" {
  process_id                    = azuredevops_workitemtrackingprocess_process.example.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.example.reference_name
  group_id                      = azuredevops_workitemtrackingprocess_group.example.id
  control_id                    = "System.Title"
  label                         = "Title"
}
```

### Contribution Control (Extension)

```hcl
resource "azuredevops_workitemtrackingprocess_process" "example" {
  name                   = "example-process"
  parent_process_type_id = "adcc42ab-9882-485e-a3ed-7678f01f66bc"
}

resource "azuredevops_workitemtrackingprocess_workitemtype" "example" {
  process_id  = azuredevops_workitemtrackingprocess_process.example.id
  name        = "example"
}

resource "azuredevops_workitemtrackingprocess_group" "example" {
  process_id                    = azuredevops_workitemtrackingprocess_process.example.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.example.reference_name
  page_id                       = azuredevops_workitemtrackingprocess_workitemtype.example.pages[0].id
  section_id                    = azuredevops_workitemtrackingprocess_workitemtype.example.pages[0].sections[0].id
  label                         = "Custom Group"
}

resource "azuredevops_workitemtrackingprocess_control" "example" {
  process_id                    = azuredevops_workitemtrackingprocess_process.example.id
  work_item_type_reference_name = azuredevops_workitemtrackingprocess_workitemtype.example.reference_name
  group_id                      = azuredevops_workitemtrackingprocess_group.example.id
  control_id                    = "MultiValueControl"
  is_contribution               = true

  contribution {
    contribution_id = "ms-devlabs.vsts-extensions-multivalue-control.multivalue-form-control"
    height          = 50
    inputs = {
      FieldName = "System.Tags"
      Values    = "Option1;Option2;Option3"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `process_id` - (Required) The ID of the process. Changing this forces a new control to be created.

* `work_item_type_reference_name` - (Required) The reference name of the work item type. Changing this forces a new control to be created.

* `group_id` - (Required) The ID of the group to add the control to. Changing this moves the control to the new group.

* `control_id` - (Required) The ID for the control. For field controls, this is the field reference name. Changing this forces a new control to be created.

---

* `label` - (Optional) Label for the control.

* `order` - (Optional) Order in which the control should appear in its group.

* `visible` - (Optional) A value indicating if the control should be visible or not. Default: `true`

* `read_only` - (Optional) A value indicating if the control is readonly. Default: `false`

* `metadata` - (Optional) Inner text of the control.

* `watermark` - (Optional) Watermark text for the textbox.

* `height` - (Optional) Height of the control, for HTML controls.

* `control_type` - (Optional) Type of the control.

* `inherited` - (Optional) A value indicating whether this layout node has been inherited from a parent layout.

* `overridden` - (Optional) A value indicating whether this layout node has been overridden by a child layout.

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

* `id` - The ID of the control.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Controls](https://learn.microsoft.com/en-us/rest/api/azure/devops/processes/controls?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the control.
* `read` - (Defaults to 5 minutes) Used when retrieving the control.
* `update` - (Defaults to 10 minutes) Used when updating the control.
* `delete` - (Defaults to 10 minutes) Used when deleting the control.

## Import

Controls can be imported using the complete resource id `process_id/work_item_type_reference_name/group_id/control_id`, e.g.

```shell
terraform import azuredevops_workitemtrackingprocess_control.example 00000000-0000-0000-0000-000000000000/MyProcess.CustomWorkItemType/group-id/System.Title
```
