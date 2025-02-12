---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_user"
description: |-
  Use this data source to access information about an existing user within Azure DevOps.
---

# Data Source: azuredevops_user

Use this data source to access information about an existing user within Azure DevOps.

~>**NOTE:** If you only have the Storage Key(UUID) of the user, you can use `azuredevops_descriptor` to resolve the Storage Key(UUID) to a `descriptor`.

## Example Usage

```hcl
data "azuredevops_user" "example" {
  principal_name = "example@example.com"
}
```

## Argument Reference

The following arguments are supported:

* `descriptor` - (Required) The descriptor of the user.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the User.

* `display_name` - This is the non-unique display name of the graph subject.

* `domain` - This is the non-unique display name of the graph subject.

* `mail_address` - The email address of record for a given graph member. This may be different than the principal name.

* `origin` - The type of source provider for the origin identifier (ex:`AD`, `AAD`, `MSA`).

* `origin_id` - The unique identifier from the system of origin. Typically a sid, object id or Guid. Linking and unlinking operations can cause this value to change for a user because the user is not backed by a different provider and has a different unique id in the new provider.

* `principal_name` - This is the PrincipalName of this graph member from the source provider. The source provider may change this field over time and it is not guaranteed to be immutable for the life of the graph member by VSTS.

* `subject_kind` - The subject kind of the user (ex: `Group`, `Scope`, `User`).

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Graph Users API](https://learn.microsoft.com/en-us/rest/api/azure/devops/graph/users/get?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 2 minute) Used when retrieving the User.