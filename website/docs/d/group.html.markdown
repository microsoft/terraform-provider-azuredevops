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
data "azuredevops_project" "example" {
  name = "Example Project"
}

data "azuredevops_group" "example" {
  project_id = data.azuredevops_project.example.id
  name       = "Example Group"
}

output "group_id" {
  value = data.azuredevops_group.example.id
}

output "group_descriptor" {
  value = data.azuredevops_group.example.descriptor
}

data "azuredevops_group" "example-collection-group" {
  name = "Project Collection Administrators"
}

output "collection_group_id" {
  value = data.azuredevops_group.example.id
}

output "collection_group_descriptor" {
  value = data.azuredevops_group.example.descriptor
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The Group Name.
- `project_id` - (Optional) The Project ID. If no project ID is specified the project collection groups will be searched.

## Attributes Reference

The following attributes are exported:

- `id` - The ID for this resource is the group descriptor. See below.
- `descriptor` - The Descriptor is the primary way to reference the graph subject. This field will uniquely identify the same graph subject across both Accounts and Organizations.
- `origin` - The type of source provider for the origin identifier (ex:AD, AAD, MSA)
- `origin_id` - The unique identifier from the system of origin. Typically a sid, object id or Guid. Linking and unlinking operations can cause this value to change for a user because the user is not backed by a different provider and has a different unique id in the new provider.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Groups - Get](https://docs.microsoft.com/en-us/rest/api/azure/devops/graph/groups/get?view=azure-devops-rest-7.0)
