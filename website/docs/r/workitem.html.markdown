---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_workitem"
description: |-
  Manages a Work Item in Azure Devops.
---

# azuredevops_workitem

Manages a Work Item in Azure Devops.

## Example Usage

```hcl
resource "azuredevops_workitem" "example" {
  title = "Testing Terraform with"
  project_id = "9a6ec7fd-a679-4b11-af60-0b1fc8cae540"
  type = "Issue"
  custom_fields = {
    foo = "SomeCustomField"
  } 
 state = "To Do"

}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) UUID of the Project.

* `title` - (Required) Title of the Work Item.

* `type` - (Required) Type of the Work Item

---

* `custom_fields` - (Optional) Specifies a list with Custom Fields for the Work Item.

* `state` - (Optional) Initial State of the Work Item.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Work Item.



## Import

Work Item can be imported using the `resource id`, e.g.

```shell
terraform import azuredevops_workitem.example 00000000-0000-0000-0000-000000000000
```
