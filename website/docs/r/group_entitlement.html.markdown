---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_group_entitlement"
description: |-
Manages a group entitlement within Azure DevOps organization.
---

# azuredevops_user_entitlement

Manages a group entitlement within Azure DevOps.

## Example Usage

### With group principal name
```hcl
resource "azuredevops_group_entitlement" "example" {
  principal_name = "[contoso]\\Group"
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

- `principal_name` - (Optional) The principal name is the PrincipalName of a graph member from the source provider. Usually, e-mail address.
- `origin_id` - (Optional) The unique identifier from the system of origin. Typically, a sid, object id or Guid. e.g. Used for member of other tenant on Azure Active Directory.
- `origin` - (Optional) The type of source provider for the origin identifier.
- `display_name` - (Optional) The display name is the name used in Azure DevOps UI. Cannot be used together with `principal_name`.
- `account_license_type` - (Optional) Type of Account License. Valid values: `advanced`, `earlyAdopter`, `express`, `none`, `professional`, or `stakeholder`. Defaults to `express`. In addition, the value `basic` is allowed which is an alias for `express` and reflects the name of the `express` license used in the Azure DevOps web interface.
- `licensing_source` - (Optional) The source of the licensing (e.g. Account. MSDN etc.) Valid values: `account` (Default), `auto`, `msdn`, `none`, `profile`, `trial`

> **NOTE:** A group can only be referenced by it's `principal_name` or by the combination of `origin_id` and `origin`.

## Attributes Reference

The following attributes are exported:

- `id` - The id of the entitlement.
- `descriptor` - The descriptor is the primary way to reference the graph subject while the system is running. This field will uniquely identify the group graph subject.

## Relevant Links

- [Azure DevOps Service REST API 6.0 - Group Entitlements - Add](https://learn.microsoft.com/en-us/rest/api/azure/devops/memberentitlementmanagement/group-entitlements/add?view=azure-devops-rest-6.0&tabs=HTTP)
- [Programmatic mapping of access levels](https://docs.microsoft.com/en-us/azure/devops/organizations/security/access-levels?view=azure-devops#programmatic-mapping-of-access-levels)

## Import

The resources allow the import via the UUID of a group entitlement.

## PAT Permissions Required

- **Member Entitlement Management**: Read & Write
