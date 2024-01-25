---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_identity_group"
description: |-
  Use this data source to access information about existing Groups within Azure DevOps
---

# Data Source: azuredevops_identity_group

Use this data source to access information about existing Groups within Azure DevOps On-Premise(Azure DevOps Server).

## Example Usage

```hcl
data "azuredevops_project" "example" {
  name = "Example Project"
}

# load all existing groups inside an organization
data "azuredevops_identity_group" "example-all-group" {
  name = "Group-Name"
}

# load all existing groups inside a specific project
data "azuredevops_identity_group" "example-project-group" {
  project_id = data.azuredevops_project.example.id
  name = "[Project-Name]\\Group-Name"
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required) The name of the group.
- `project_id` - (Optional) The Project ID. If no project ID is specified all groups of an organization will be returned

## Attributes Reference

The following attributes are exported:

  - `id` - The ID is the primary way to reference the identity subject. This field will uniquely identify the same identity subject across both Accounts and Organizations.
  - `name` - This is the non-unique display name of the identity subject. To change this field, you must alter its value in the source provider.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Identities](https://docs.microsoft.com/en-us/rest/api/azure/devops/ims/?view=azure-devops-rest-7.2)
