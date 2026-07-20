---
layout: "azuredevops"
page_title: "AzureDevops: Data Source: azuredevops_service_principal"
description: |-
  Gets information about an existing Service Principal.
---

# Data Source: azuredevops_service_principal

Use this data source to access information about an existing Service Principal.

## Example Usage

### By Display Name

```hcl
data "azuredevops_service_principal" "example" {
  display_name = "existing"
}

output "id" {
  value = data.azuredevops_service_principal.example.id
}
```

### By Origin ID

```hcl
data "azuredevops_service_principal" "example" {
  origin_id = "00000000-0000-0000-0000-000000000000"
}

output "id" {
  value = data.azuredevops_service_principal.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `display_name` - (Optional) The Display Name of the Service Principal. Changing this forces a new Service Principal to be created.

* `origin_id` - (Optional) The origin ID of the Service Principal.

~> **NOTE:** Exactly one of `display_name` or `origin_id` must be specified.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Service Principal.

* `descriptor` - The descriptor of the Service Principal.

* `origin` - The origin of the Service Principal.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 30 minutes) Used when retrieving the Service Principal.
