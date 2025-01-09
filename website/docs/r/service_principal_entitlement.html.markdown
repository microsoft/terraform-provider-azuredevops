---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_service_principal_entitlement"
description: |-
  Manages a Service Principal Entitlement.
---

# azuredevops_service_principal_entitlement

Manages a Service Principal Entitlement.

## Example Usage

```hcl
resource "azuredevops_service_principal_entitlement" "example" {
  origin_id = "00000000-0000-0000-0000-000000000000"
}
```

## Arguments Reference

The following arguments are supported:

* `origin_id` - (Required) The Object ID of the service principal in Entra ID. Changing this forces a new Service Principal Entitlement to be created.

---

* `account_license_type` - (Optional) Type of Account License. Valid values: `advanced`, `earlyAdopter`, `express`, `none`, `professional`, or `stakeholder`. Defaults to `express`. In addition the value `basic` is allowed which is an alias for `express` and reflects the name of the `express` license used in the Azure DevOps web interface.

* `licensing_source` - (Optional) The source of the licensing (e.g. Account. MSDN etc.) Valid values: `account` (Default), `auto`, `msdn`, `none`, `profile`, `trial`

* `origin` - (Optional) The type of source provider for the origin identifier.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Service Principal Entitlement.

* `descriptor` - The descriptor is the primary way to reference the graph subject while the system is running. This field will uniquely identify the user graph subject.

* `display_name` - The display name of service principal.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Service Principal Entitlement.
* `read` - (Defaults to 2 minutes) Used when retrieving the Service Principal Entitlement.
* `update` - (Defaults to 5 minutes) Used when updating the Service Principal Entitlement.
* `delete` - (Defaults to 5 minutes) Used when deleting the Service Principal Entitlement.

## Import

Service Principal Entitlements can be imported using the `resource id`.
The `resource id` can be found using DEV Tools in the `Users` section of the ADO organization.


```shell
terraform import azuredevops_service_principal_entitlement.example 8480c6eb-ce60-47e9-88df-eca3c801638b
```
