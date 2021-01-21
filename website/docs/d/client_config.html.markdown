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
data "azuredevops_client_config" "c" {}

output "org_url" {
  value = data.azuredevops_client_config.c.organization_url
}
```

## Argument Reference

This data source has no arguments

## Attributes Reference

The following attributes are exported:

`organization_url` - The organization configured for the provider
