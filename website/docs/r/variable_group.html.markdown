---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_variable_group"
description: |-
  Manages variable groups within Azure DevOps project.
---

# azuredevops_variable_group

Manages variable groups within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  project_name = "Test Project"
}

resource "azuredevops_variable_group" "variablegroup" {
  project_id   = azuredevops_project.project.id
  name         = "Test Variable Group"
  description  = "Test Variable Group Description"
  allow_access = true

  variable {
    name  = "key"
    value = "value"
  }

  variable {
    name      = "Account Password"
    value     = "p@ssword123"
    is_secret = true
  }
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The project ID or project name.
- `name` - (Required) The name of the Variable Group.
- `description` - (Optional) The description of the Variable Group.
- `allow_access` - (Required) Boolean that indicate if this variable group is shared by all pipelines of this project.
- `variable` - (Optional) One or more `variable` blocks as documented below.

A `variable` block supports the following:

- `name` - (Required) The key value used for the variable. Must be unique within the Variable Group.
- `value` - (Optional) The value of the variable. If omitted, it will default to empty string.
- `secret_value` - (Optional) The secret value of the variable. If omitted, it will default to empty string. Used when `is_secret` set to `true`.
- `is_secret` - (Optional) A boolean flag describing if the variable value is sensitive. Defaults to `false`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the Variable Group returned after creation in Azure DevOps.

## Relevant Links

- [Azure DevOps Service REST API 5.1 - Variable Groups](https://docs.microsoft.com/en-us/rest/api/azure/devops/distributedtask/variablegroups?view=azure-devops-rest-5.1)
- [Azure DevOps Service REST API 5.1 - Authorized Resources](https://docs.microsoft.com/en-us/rest/api/azure/devops/build/authorizedresources?view=azure-devops-rest-5.1)

## Import

Azure DevOps Variable groups can be imported using the project name/variable group ID or by the project Guid/variable group ID, e.g.

```sh
terraform import azuredevops_variable_group.variablegroup "Test Project/10"
```

or

```sh
terraform import azuredevops_variable_group.variablegroup 782a8123-1019-xxxx-xxxx-xxxxxxxx/10
```

_Note that for secret variables, the import command retrieve blank value in the tfstate._

## PAT Permissions Required

- **Variable Groups**: Read, Create, & Manage
