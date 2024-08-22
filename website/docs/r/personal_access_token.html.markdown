---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_personal_access_token"
description: |-
  Manages a User's Personal Access Tokens for an Azure DevOps Organization.
---

# azuredevops_personal_access_token

Manages a User's Personal Access Tokens for an Azure DevOps Organization.

## Example Usage

```hcl
resource "azuredevops_personal_access_token" "example" {
  name            = "Example Token"
  all_orgs        = false
  scopes          = ["vso.dashboards", "vso.taskgroups_manage"]
  valid_to        = "2025-01-01 00:00:00Z"
}

data "azuredevops_personal_access_token" "example-pat" {
  authorization_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The Token Name.
- `all_orgs` - (Optional) True, if this personal access token (PAT) is for all of the user's accessible organizations. False, if otherwise (e.g. if the token is for a specific organization).
- `scope` - (Optional) The token scopes for accessing Azure DevOps resources.
- `valid_to` - (Optional) The token expiration date.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:
- `authorization_id` - Unique guid identifier.
- `target_accounts` - The organizations for which the token is valid; null if the token applies to all of the user's accessible organizations.
- `token` - The Personal Access Token (Sensitive)
- `valid_from` - The token creation date

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Token Lifecycle Management](https://learn.microsoft.com/en-us/rest/api/azure/devops/tokens/?view=azure-devops-rest-7.0)

## Import

Azure DevOps Personal Access Tokens can be imported using the authorization ID:

```sh
terraform import azuredevops_personal_access_token.example 00000000-0000-0000-0000-000000000000
```
