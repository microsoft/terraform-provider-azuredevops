---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_personal_access_token"
description: |-
  Use this data source to access information about an existing Personal Access Token within Azure DevOps
---

# Data Source: azuredevops_personal_access_token

Use this data source to access information about an existing Personal Access Token within Azure DevOps.

## Example Usage

```hcl
data "azuredevops_personal_access_token" "example-pat" {
  authorization_id = "00000000-0000-0000-0000-000000000000"
}

output "accounts" {
  value = data.azuredevops_personal_access_token.example.target_accounts
}

output "accounts" {
  value = data.azuredevops_personal_access_token.example.scope
}

output "creation_date" {
  value = data.azuredevops_personal_access_token.example.valid_from
}

output "expiration_date" {
  value = data.azuredevops_personal_access_token.example.valid_to
}
```

## Argument Reference

The following arguments are supported:

- `authorization_id` - (Required) Unique guid identifier of each Personal Access Token.

## Attributes Reference

The following attributes are exported:

- `name` - (Required) The Token Name.
- `scope` - (Optional) The token scopes for accessing Azure DevOps resources.
- `target_accounts` - The organizations for which the token is valid; null if the token applies to all of the user's accessible organizations.
- `token` - The Personal Access Token (Sensitive)
- `valid_to` - (Optional) The token expiration date.
- `valid_from` - The token creation date

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Token Lifecycle Management](https://learn.microsoft.com/en-us/rest/api/azure/devops/tokens/?view=azure-devops-rest-7.0)
