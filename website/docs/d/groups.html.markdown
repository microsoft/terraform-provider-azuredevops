---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_groups"
description: |-
  Use this data source to access information about existing Groups within Azure DevOps
---

# Data Source: azuredevops_groups

Use this data source to access information about existing Groups within Azure DevOps

## Example Usage

```hcl
# load all existing groups inside an organization
data "azuredevops_groups" "all-groups" {
}

# load all existing groups inside a specific project
data "azuredevops_groups" "project-groups" {
  project_id = data.azuredevops_project.p.id
}

```

## Argument Reference

The following arguments are supported:

- `project_id` - (Optional) The Project ID. If no project ID is specified all groups of an organization will be returned

## Attributes Reference

The following attributes are exported:

- `groups` - A set of existing groups in your Azure DevOps Organization or project with details about every single group which includes:

  - `description` - A short phrase to help human readers disambiguate groups with similar names
  - `descriptor` - The descriptor is the primary way to reference the graph subject while the system is running. This field will uniquely identify the same graph subject across both Accounts and Organizations.
  - `display_name` - This is the non-unique display name of the graph subject. To change this field, you must alter its value in the source provider.
  - `domain` - This represents the name of the container of origin for a graph member. (For MSA this is "Windows Live ID", for AD the name of the domain, for AAD the tenantID of the directory, for VSTS groups the ScopeId, etc)
  - `mail_address` - The email address of record for a given graph member. This may be different than the principal name.
  - `origin` - The type of source provider for the origin identifier (ex:AD, AAD, MSA)
  - `origin_id` - The unique identifier from the system of origin. Typically a sid, object id or Guid. Linking and unlinking operations can cause this value to change for a user because the user is not backed by a different provider and has a different unique id in the new provider.
  - `principal_name` - This is the PrincipalName of this graph member from the source provider. The source provider may change this field over time and it is not guaranteed to be immutable for the life of the graph member by VSTS.
  - `url` - This url is the full route to the source resource of this graph subject.

## Relevant Links

- [Azure DevOps Service REST API 6.0 - Groups - List](https://docs.microsoft.com/en-us/rest/api/azure/devops/graph/groups/list?view=azure-devops-rest-6.0)
