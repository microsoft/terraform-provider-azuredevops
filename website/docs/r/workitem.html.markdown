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
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}
resource "azuredevops_workitem" "example" {
  project_id    = data.azuredevops_project.example.id
  title         = "Testing Terraform Item"
  type          = "Issue"
  state = "Active"
  tags  = ["Tag"]
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the Project.

* `title` - (Required) Title of the Work Item.

* `type` - (Required) Type of the Work Item

---

* `custom_fields` - (Optional) Specifies a list with Custom Fields for the Work Item.

* `state` - (Optional) Initial State of the Work Item.

* `tags` - (Optional) Specifies a list of Tags.
  
## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Work Item.



## Import

Work Items  can be imported using the project name/variable group ID or by the project Guid/variable group ID, e.g.

```sh
terraform import azuredevops_workitem.example "Example Project/10"
```
