---
layout: "azuredevops"
page_title: "AzureDevops: Data Source: azuredevops_storage_key"
description: |-
  Resolve a descriptor to a storage key.
---

# Data Source: azuredevops_storage_key

Use this data source to access information about an existing Storage Key.

## Example Usage

```hcl
data "azuredevops_storage_key" "example" {
  descriptor = "aad.000000000000000000000000000000000000"
}

output "id" {
  value = data.azuredevops_storage_key.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `descriptor` - (Required) The descriptor that will be resolved to a storage key.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Storage Key.

* `storage_key` - The Storage Key of the descriptor.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Storage Key - Get](https://learn.microsoft.com/en-us/rest/api/azure/devops/graph/storage-keys/get?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 2 minutes) Used when retrieving the Storage Key.
