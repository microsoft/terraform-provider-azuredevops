---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_group_entitlement"
description: |-
Manages a group entitlement within Azure DevOps organization.
---

# azuredevops_group_entitlement

Manages a group entitlement within Azure DevOps.

## Example Usage

### With an Azure DevOps local group managed by this resource
```hcl
resource "azuredevops_group_entitlement" "example" {
  display_name = "Group Name"
}
```

### With group origin ID
```hcl
resource "azuredevops_group_entitlement" "example" {
  origin    = "aad"
  origin_id = "00000000-0000-0000-0000-000000000000"
}
```

## Argument Reference

- `display_name` - (Optional) The display name is the name used in Azure DevOps UI. Cannot be set together with `origin_id` and `origin`.
- `origin_id` - (Optional) The unique identifier from the system of origin. Typically, a sid, object id or Guid. e.g. Used for member of other tenant on Azure Active Directory.
- `origin` - (Optional) The type of source provider for the origin identifier.
- `account_license_type` - (Optional) Type of Account License. Valid values: `advanced`, `earlyAdopter`, `express`, `none`, `professional`, or `stakeholder`. Defaults to `express`. In addition, the value `basic` is allowed which is an alias for `express` and reflects the name of the `express` license used in the Azure DevOps web interface.
- `licensing_source` - (Optional) The source of the licensing (e.g. Account. MSDN etc.) Valid values: `account` (Default), `auto`, `msdn`, `none`, `profile`, `trial`

> **NOTE:** A existing group in Azure AD can only be referenced by the combination of `origin_id` and `origin`.

## Attributes Reference

The following attributes are exported:

- `id` - The id of the entitlement.
- `principal_name` - The principal name of a graph member on Azure DevOps
- `descriptor` - The descriptor is the primary way to reference the graph subject while the system is running. This field will uniquely identify the group graph subject.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Group Entitlements](https://learn.microsoft.com/en-us/rest/api/azure/devops/memberentitlementmanagement/group-entitlements?view=azure-devops-rest-7.1)
- [Programmatic mapping of access levels](https://docs.microsoft.com/en-us/azure/devops/organizations/security/access-levels?view=azure-devops#programmatic-mapping-of-access-levels)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Group Entitlement.
* `read` - (Defaults to 5 minute) Used when retrieving the Group Entitlement.
* `update` - (Defaults to 30 minutes) Used when updating the Group Entitlement.
* `delete` - (Defaults to 30 minutes) Used when deleting the Group Entitlement.

## Import

The resource allows the import via the ID of a group entitlement, which is a UUID.

```
terraform import azuredevops_group_entitlement.example 00000000-0000-0000-0000-000000000000
```

## PAT Permissions Required

- **Member Entitlement Management**: Read & Write
