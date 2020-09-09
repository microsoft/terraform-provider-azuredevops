---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_group"
description: |-
  Use this data source to access information about an existing Group within Azure DevOps.
---

# Data Source: azuredevops_group

Use this data source to access information about an existing Group within Azure DevOps

## Example Usage

```hcl
data "azuredevops_project" "p" {
  project_name = "contoso-project"
}

data "azuredevops_group" "test" {
  project_id = data.azuredevops_project.p.id
  name       = "Test Group"
}

output "group_id" {
  value = data.azuredevops_group.test.id
}

output "group_descriptor" {
  value = data.azuredevops_group.test.descriptor
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The Project Id.
- `name` - (Required) The Group Name.

## Attributes Reference

The following attributes are exported:

- `id` - The ID for this resource is the group descriptor. See below.
- `descriptor` - The Descriptor is the primary way to reference the graph subject. This field will uniquely identify the same graph subject across both Accounts and Organizations.
- `origin` - The type of source provider for the origin identifier (ex:AD, AAD, MSA)
- `origin_id` - The unique identifier from the system of origin. Typically a sid, object id or Guid. Linking and unlinking operations can cause this value to change for a user because the user is not backed by a different provider and has a different unique id in the new provider.

## Relevant Links

- [Azure DevOps Service REST API 5.1 - Groups - Get](https://docs.microsoft.com/en-us/rest/api/azure/devops/graph/groups/get?view=azure-devops-rest-5.1)
