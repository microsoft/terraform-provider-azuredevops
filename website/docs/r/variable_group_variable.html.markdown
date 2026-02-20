---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_variable_group_variable"
description: |-
  Manages variable group variables within a variable group.
---

# azuredevops_variable_group_variable

Manages variable group variables within a variable group.

~> **Note** Variable group variables can also be managed inlined in the `variable` blocks in `azuredevops_variable_group`.

## Example Usage

### Basic usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

resource "azuredevops_variable_group" "example" {
  project_id   = azuredevops_project.example.id
  name         = "Example Variable Group"
  description  = "Example Variable Group Description"
  allow_access = true

  variable {
    name  = "key1"
    value = "val1"
  }

  lifecycle {
    ignore_changes = [variable]
  }
}


resource "azuredevops_variable_group_variable" "example" {
  project_id        = azuredevops_project.example.id
  variable_group_id = azuredevops_variable_group.example.id
  name              = "key2"
  value             = "val2"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `variable_group_id` - (Required) The ID of the variable group.

* `name` - (Required) The name of the variable. Must be unique within the Variable Group.

* `value` - (Optional) The value of the variable.

* `secret_value` - (Optional) The value of the secret variable.

-> **NOTE** Exactly one of `value` and `secret_value` must be specified.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Variable Group returned after creation in Azure DevOps.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Variable Groups](https://docs.microsoft.com/en-us/rest/api/azure/devops/distributedtask/variablegroups?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Variable Group Variable.
* `read` - (Defaults to 5 minute) Used when retrieving the Variable Group Variable.
* `update` - (Defaults to 10 minutes) Used when updating the Variable Group Variable.
* `delete` - (Defaults to 10 minutes) Used when deleting the Variable Group Variable.

## Import

**Secret variable cannot be imported.**

Azure DevOps Variable group variables can be imported using the `project ID/variable group ID/variable name`, e.g.


```sh
terraform import azuredevops_variable_group_variable.example 00000000-0000-0000-0000-000000000000/0/key1
```

## PAT Permissions Required

- **Variable Groups**: Read, Create, & Manage
