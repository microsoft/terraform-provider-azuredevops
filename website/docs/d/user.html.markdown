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

* `display_name` - The display name of the User.

* `domain` - The domain of the user.

* `mail_address` - The email address of the user.

* `origin` - The type of source provider for the origin identifier (ex:`AD`, `AAD`, `MSA`).

* `origin_id` - The origin ID of the user.

* `principal_name` - The principal name of the user.

* `subject_kind` - The subject kind of the user (ex: `Group`, `Scope`, `User`).

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Graph Users API](https://learn.microsoft.com/en-us/rest/api/azure/devops/graph/users/get?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 2 minute) Used when retrieving the User.
