---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_identity_user"
description: |-
  Use this data source to access information about an existing users within Azure DevOps.
---

# Data Source: azuredevops_identity_user

Use this data source to access information about an existing users within Azure DevOps On-Premise(Azure DevOps Server).

## Example Usage

```hcl
# Load single user by using it's principal name
data "azuredevops_identity_user" "contoso-user" {
  name = "contoso-user"
}

# Use MailAddress instead of principal name.
data "azuredevops_user" "contoso-user-upn" {
  name = "contoso-user@contoso.onmicrosoft.com"
  search_filter = "MailAddress"
}


# Use MailAddress instead of principal name.
data "azuredevops_user" "contoso-user-upn" {
  name = "contoso-user@contoso.onmicrosoft.com"
  search_filter = "MailAddress"
}

# Use DisplayName instead of principal name.
data "azuredevops_user" "contoso-user-upn" {
  name = "Contoso User"
  search_filter = "DisplayName"
}

```

## Argument Reference

The following arguments are supported:

- `name` - (required) The PrincipalName of this identity member from the source provider.
- `search_filter` - (Optional) Default is General, but other options are AccountName, DisplayName, and MailAddress.


## Attributes Reference

The following attributes are exported:

- `users` - A set of existing users in your Azure DevOps Organization with details about every single user which includes:

  - `id` - The ID is the primary way to reference the identity subject while the system is running. This field will uniquely identify the same identity subject across both Accounts and Organizations.
  - `name` - This is the PrincipalName of this identity member from the source provider. The source provider may change this field over time and it is not guaranteed to be immutable for the life of the identity member.


## Relevant Links

- [Azure DevOps Service REST API 7.0 - Identities](https://docs.microsoft.com/en-us/rest/api/azure/devops/ims/?view=azure-devops-rest-7.2)
