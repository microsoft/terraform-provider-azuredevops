---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitemtrackingprocess_rule"
description: |-
  Manages a rule for a work item type in a process.
---

# azuredevops_workitemtrackingprocess_rule

Manages a rule for a work item type in a process. Rules define conditions and actions that are triggered during work item lifecycle events.

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

resource "azuredevops_workitemtrackingprocess_rule" "example" {
  process_id        = azuredevops_workitemtrackingprocess_process.example.id
  work_item_type_id = azuredevops_workitemtrackingprocess_workitemtype.example.reference_name
  name              = "Require Title When New"

  condition {
    condition_type = "when"
    field          = "System.State"
    value          = "New"
  }

  action {
    action_type  = "makeRequired"
    target_field = "System.Title"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `process_id` - (Required) The ID of the process. Changing this forces a new rule to be created.

* `work_item_type_id` - (Required) The ID (reference name) of the work item type. Changing this forces a new rule to be created.

* `name` - (Required) Name of the rule.

* `condition` - (Required) One or more `condition` blocks as defined below.

* `action` - (Required) One or more `action` blocks as defined below.

---

* `is_disabled` - (Optional) Indicates if the rule is disabled. Default: `false`

---

A `condition` block supports the following:

* `condition_type` - (Required) Type of condition. Valid values: `when`, `whenNot`, `whenChanged`, `whenNotChanged`, `whenWas`, `whenCurrentUserIsMemberOfGroup`, `whenCurrentUserIsNotMemberOfGroup`.

* `field` - (Optional) Field reference name for the condition. Required for most condition types.

* `value` - (Optional) Value to match for the condition.

---

A `action` block supports the following:

* `action_type` - (Required) Type of action. Valid values: `makeRequired`, `makeReadOnly`, `setDefaultValue`, `setDefaultFromClock`, `setDefaultFromCurrentUser`, `setDefaultFromField`, `copyValue`, `copyFromClock`, `copyFromCurrentUser`, `copyFromField`, `setValueToEmpty`, `copyFromServerClock`, `copyFromServerCurrentUser`, `hideTargetField`, `disallowValue`.

* `target_field` - (Required) Field to act on.

* `value` - (Optional) Value to set on the target field.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the rule.

* `url` - URL of the rule resource.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Rules](https://learn.microsoft.com/en-us/rest/api/azure/devops/processes/rules?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the rule.
* `read` - (Defaults to 5 minutes) Used when retrieving the rule.
* `update` - (Defaults to 10 minutes) Used when updating the rule.
* `delete` - (Defaults to 10 minutes) Used when deleting the rule.

## Import

Rules can be imported using the complete resource id `process_id/work_item_type_id/rule_id`, e.g.

```shell
terraform import azuredevops_workitemtrackingprocess_rule.example 00000000-0000-0000-0000-000000000000/MyProcess.CustomWorkItemType/11111111-1111-1111-1111-111111111111
```
