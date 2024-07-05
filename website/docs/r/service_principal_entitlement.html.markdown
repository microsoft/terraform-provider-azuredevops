---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_service_principal_entitlement"
description: |-
  Manages a service principal entitlement within Azure DevOps organization.
---

# azuredevops_service_principal_entitlement

Manages a service principal entitlement within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_service_principal_entitlement" "example" {
  origin_id = "00000000-0000-0000-0000-000000000001"
}
```

## Argument Reference

- `origin_id` - (Required) The object ID of the enterprise application.
- `origin` - (Optional) The type of source provider for the origin identifier. Defaults to `aad`.
- `account_license_type` - (Optional) Type of Account License. Valid values: `advanced`, `earlyAdopter`, `express`, `none`, `professional`, or `stakeholder`. Defaults to `express`.

  ~> **Note**
  The value `basic` is allowed which is an alias for `express` and reflects the name of the `express` license used in the Azure DevOps web interface.


- `licensing_source` - (Optional) The source of the licensing (e.g. Account. MSDN etc.) Valid values: `account` (Default), `auto`, `msdn`, `none`, `profile`, `trial`

## Attributes Reference

The following attributes are exported:

- `id` - The id of the entitlement.
- `descriptor` - The descriptor is the primary way to reference the graph subject while the system is running. This field will uniquely identify the service principal graph subject.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - User Entitlements - Add](https://learn.microsoft.com/en-us/rest/api/azure/devops/memberentitlementmanagement/service-principal-entitlements/add?view=azure-devops-rest-7.1)
- [Programmatic mapping of access levels](https://docs.microsoft.com/en-us/azure/devops/organizations/security/access-levels?view=azure-devops#programmatic-mapping-of-access-levels)

## Import

The resources allows the import via the UUID of a service principal entitlement or by using the principal name of a service principal owning an entitlement.

## PAT Permissions Required

- **Member Entitlement Management**: Read & Write
