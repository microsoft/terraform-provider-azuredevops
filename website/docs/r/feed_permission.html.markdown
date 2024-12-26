---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_feed_permission"
description: |-
  Manages creation of the Feed Permission within Azure DevOps organization.
---

# azuredevops_feed_permission

Manages creation of the Feed Permission within Azure DevOps organization.

## Example Usage

### Create Feed Permission
```hcl
resource "azuredevops_project" "example" {
  name = "Example Project"
}

resource "azuredevops_group" "example" {
  scope        = azuredevops_project.example.id
  display_name = "Example group"
  description  = "Example description"
}

resource "azuredevops_feed" "example" {
  name = "releases"
}

resource "azuredevops_feed_permission" "permission" {
  feed_id             = azuredevops_feed.example.id
  role                = "reader"
  identity_descriptor = azuredevops_group.example.descriptor
}
```


## Argument Reference

The following arguments are supported:

- `feed_id` - (Required) The ID of the Feed.
- `identity_descriptor` - (Required) The Descriptor of identity you want to assign a role.
- `role` - (Required) The role to be assigned, possible values : `reader`, `contributor`, `collaborator`, `administrator`
- `project_id` - (Optional) The ID of the Project Feed is created in. If not specified, feed will be created at the organization level.
- `display_name` - (Optional) The display name of the assignment

## Attributes Reference

The following attributes are exported:

- `feed_id` - The ID of the Feed.
- `identity_descriptor` - The Descriptor of  the identity.
- `identity_id` - The ID of the identity.
- `role` - The assigned role
- `project_id` - The ID of the Project Feed is created in (if one exists).
- `display_name` - The display name of the assignment (if one exists).

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Feed Management](https://learn.microsoft.com/en-us/rest/api/azure/devops/artifacts/feed-management?view=azure-devops-rest-7.0)