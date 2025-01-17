---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_identity_groups"
description: |-
  Use this data source to access information about existing Groups within Azure DevOps
---

# Data Source: azuredevops_identity_groups

Use this data source to access information about existing Groups within Azure DevOps On-Premise(Azure DevOps Server).

## Example Usage

```hcl
data "azuredevops_project" "example" {
  name = "Example Project"
}

# load all existing groups inside an organization
data "azuredevops_identity_groups" "example-all-groups" {
}

# load all existing groups inside a specific project
data "azuredevops_identity_groups" "example-project-groups" {
  project_id = data.azuredevops_project.example.id
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Optional) The Project ID. If no project ID is specified all groups of an organization will be returned

## Attributes Reference

The following attributes are exported:

* `groups` - A `groups` blocks as documented below. A set of existing groups in your Azure DevOps Organization or project with details about every single group.

---

A `groups` block supports the following:

* `id` - The ID of the Identity Group.

* `name` - This is the non-unique display name of the identity subject. To change this field, you must alter its value in the source provider.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Identities](https://docs.microsoft.com/en-us/rest/api/azure/devops/ims/?view=azure-devops-rest-7.2)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minute) Used when retrieving the Identity Groups.