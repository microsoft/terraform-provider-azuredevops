---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_identity_group"
description: |-
  Use this data source to access information about existing Groups within Azure DevOps
---

# Data Source: azuredevops_identity_group

Use this data source to access information about an existing Group within Azure DevOps On-Premise(Azure DevOps Server).

## Example Usage

```hcl
# load existing group with specific name
data "azuredevops_identity_group" "example-project-group" {
  project_id = data.azuredevops_project.example.id
  name       = "[Project-Name]\\Group-Name"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the group.

* `project_id` - (Required) The Project ID.

## Attributes Reference

The following attributes are exported:

* `id` - The ID is the primary way to reference the identity subject. This field will uniquely identify the same identity subject across both Accounts and Organizations.

* `name` - This is the non-unique display name of the identity subject. To change this field, you must alter its value in the source provider.

* `descriptor` - The descriptor of the identity group.

* `subject_descriptor` - The subject descriptor of the identity group.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Identities](https://docs.microsoft.com/en-us/rest/api/azure/devops/ims/?view=azure-devops-rest-7.2)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minute) Used when retrieving the Identity Group.
