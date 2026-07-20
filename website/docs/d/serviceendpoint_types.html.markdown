---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_serviceendpoint_types"
description: |-
  Gets information about all available service endpoint types in Azure DevOps.
---

# Data Source: azuredevops_serviceendpoint_types

Use this data source to query all available service endpoint types in Azure DevOps.

~> **NOTE:** If you need to find a specific service endpoint type by name, consider using the [azuredevops_serviceendpoint_type](serviceendpoint_type.html) data source.

## Example Usage

```hcl
data "azuredevops_serviceendpoint_types" "all" {
}

output "all_types" {
  value = data.azuredevops_serviceendpoint_types.all.types
}
```

### Filter service endpoint types

```hcl
data "azuredevops_serviceendpoint_types" "all" {
}

locals {
  git_types = [
    for type in data.azuredevops_serviceendpoint_types.all.types :
    type if length(regexall("git", lower(type.name))) > 0
  ]
}

output "git_related_types" {
  value = local.git_types
}
```

## Argument Reference

This data source has no arguments.

## Attributes Reference

The following attributes are exported:

* `types` - A list of service endpoint types. Each type has the following attributes:
  * `id` - The ID of the service endpoint type (typically same as name).
  * `name` - The name of the service endpoint type.
  * `display_name` - The display name of the service endpoint type.
  * `description` - The description of the service endpoint type.
  * `ui_contribution_id` - The UI contribution ID for this service endpoint type.
  * `authentication_schemes` - A list of available authentication schemes for this service endpoint type.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Service Endpoint Types](https://learn.microsoft.com/en-us/rest/api/azure/devops/serviceendpoint/types/list)

