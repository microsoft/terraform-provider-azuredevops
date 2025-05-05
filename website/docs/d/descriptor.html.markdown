---
layout: "azuredevops"
page_title: "AzureDevops: Data Source: azuredevops_descriptor"
description: |-
  Resolve a storage key(`user`, `group`, `scope`, etc.) to a descriptor.
---

# Data Source: azuredevops_descriptor

Use this data source to access information about an existing Descriptor.

## Example Usage

```hcl
data "azuredevops_descriptor" "example" {
  storage_key = "00000000-0000-0000-0000-000000000000"
}

output "id" {
  value = data.azuredevops_descriptor.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `storage_key` - (Required) The ID of the resource(`user`, `group`, `scope`, etc.) that will be resolved to a descriptor.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Descriptor.

* `descriptor` - The descriptor of the storage key.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Descriptors - Get](https://learn.microsoft.com/en-us/rest/api/azure/devops/graph/descriptors/get?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 2 minutes) Used when retrieving the Descriptor.
