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

}
```

## Arguments Reference

The following arguments are supported:

* `account_license_type` - (Optional) TODO.

* `application_id` - (Optional) The ID of the TODO.Conflicts with `origin_id`,`origin` Changing this forces a new Service Principal Entitlement to be created.

* `licensing_source` - (Optional) TODO.

* `origin` - (Optional) TODO.Conflicts with `application_id` Changing this forces a new Service Principal Entitlement to be created.

* `origin_id` - (Optional) The ID of the TODO.Conflicts with `application_id` Changing this forces a new Service Principal Entitlement to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Service Principal Entitlement.

* `descriptor` - TODO.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Service Principal Entitlement.
* `read` - (Defaults to 5 minutes) Used when retrieving the Service Principal Entitlement.
* `update` - (Defaults to 30 minutes) Used when updating the Service Principal Entitlement.
* `delete` - (Defaults to 30 minutes) Used when deleting the Service Principal Entitlement.

## Import

Service Principal Entitlements can be imported using the `resource id`, e.g.

```shell
terraform import azuredevops_service_principal_entitlement.example 00000000-0000-0000-0000-000000000000
```