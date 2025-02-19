---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_client_config"
description: |-
  Use this data source to access information about the Azure DevOps organization configured for the provider.
---

# Data Source: azuredevops_client_config

Use this data source to access information about the Azure DevOps organization configured for the provider.

## Example Usage

```hcl
data "azuredevops_client_config" "example" {}

output "org_url" {
  value = data.azuredevops_client_config.example.organization_url
}
```

## Argument Reference

This data source has no arguments

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the organization.

* `name` - The name of the organization.

* `organization_url` - The URL of the organization.

* `owner_id` - The owner ID of the organization.

* `status` - The status of the organization.

* `tenant_id` - The Tenant ID of the connected Azure Directory.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minute) Used when retrieving the Client Config.