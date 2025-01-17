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

# Use MailAddress
data "azuredevops_identity_user" "contoso-user-upn" {
  name = "contoso-user@contoso.onmicrosoft.com"
  search_filter = "MailAddress"
}

# Use AccountName
data "azuredevops_identity_user" "contoso-user-upn" {
  name = "contoso-user@contoso.onmicrosoft.com"
  search_filter = "AccountName"
}

# Use DisplayName
data "azuredevops_identity_user" "contoso-user-upn" {
  name = "Contoso User"
  search_filter = "DisplayName"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (required) The PrincipalName of this identity member from the source provider.

* `search_filter` - (Optional) The type of search to perform. Possible values are: `AccountName`, `DisplayName`, and `MailAddress`. Default is `General`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the user.

* `descriptor` - The descriptor of the user.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Identities](https://docs.microsoft.com/en-us/rest/api/azure/devops/ims/?view=azure-devops-rest-7.2)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minute) Used when retrieving the Identity Users.