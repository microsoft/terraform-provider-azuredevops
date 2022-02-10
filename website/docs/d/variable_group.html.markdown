---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_variable_group"
description: |-
  Use this data source to access information about existing Variable Groups within Azure DevOps.
---

# Data Source: azuredevops_variable_group

Use this data source to access information about existing Variable Groups within Azure DevOps.

## Example Usage

```hcl
data "azuredevops_project" "project" {
  name = "Test Project"
}

data "azuredevops_variable_group" "vg" {
  project_id = data.azuredevops_project.project.id
  name       = "Sample Variable Group"
}

output "id" {
  value = data.azuredevops_variable_group.vg.id
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The project ID or project name.
- `name` - (Required) The name of the Variable Group to retrieve.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the Variable Group.
- `description` - The description of the Variable Group.
- `allow_access` - Boolean that indicate if this Variable Group is shared by all pipelines of this project.
- `variable` - One or more `variable` blocks as documented below.
- `key_vault` - A list of `key_vault` blocks as documented below.

A `variable` block supports the following:

- `name` - The key value used for the variable.
- `value` - The value of the variable.
- `secret_value` - The secret value of the variable.
- `is_secret` - A boolean flag describing if the variable value is sensitive.

A `key_vault` block supports the following:

- `name` - The name of the Azure key vault to link secrets from as variables.
- `service_endpoint_id` - The id of the Azure subscription endpoint to access the key vault.

## Relevant Links

- [Azure DevOps Service REST API 6.0 - Variable Groups](https://docs.microsoft.com/en-us/rest/api/azure/devops/distributedtask/variablegroups?view=azure-devops-rest-6.0)
