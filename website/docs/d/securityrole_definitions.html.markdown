---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_securityrole_definitions"
description: |-
  Use this data source to access information about existing Security Role Definitions within a given scope in Azure DevOps.
---

# Data Source: azuredevops_securityrole_definitions

Use this data source to access information about existing Security Role Definitions within a given scope in Azure DevOps.

## Example Usage

```hcl
data "azuredevops_securityrole_definitions" "example" {
  scope = "distributedtask.environmentreferencerole"
}

output "securityrole_definitions" {
  value = data.aazuredevops_securityrole_definitions.example.definitions
}

```

## Argument Reference

The following arguments are supported:

- `scope` - (Required) Name of the Scope for which Security Role Definitions will be returned.

DataSource without specifying any arguments will return all projects.

## Attributes Reference

The following attributes are exported:

- `definitions` - A list of existing Security Role Definitions in a Scope in your Azure DevOps Organization with details about every definition which includes:

- `name` - The name of the Security Role Definition.
- `display_name` - The display name of the Security Role Definition.
- `allow_permissions` - The mask of allowed permissions of the Security Role Definition.
- `deny_permissions` - The mask of the denied permissions of the Security Role Definition.
- `identifier` - The identifier of the Security Role Definition.
- `description` - The description of the Security Role Definition.
- `scope` - The scope of the Security Role Definition.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Projects - Get](https://docs.microsoft.com/en-us/rest/api/azure/devops/core/projects/get?view=azure-devops-rest-7.0)
