---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_group_membership"
description: |-
  Manages group membership within Azure DevOps organization.
---

# azuredevops_group_membership
Manages group membership within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "project" {
  project_name = "Test Project"
}

resource "azuredevops_user_entitlement" "user" {
  principal_name = "foo@contoso.com"
}

data "azuredevops_group" "group" {
  project_id = azuredevops_project.project.id
  name       = "Build Administrators"
}

resource "azuredevops_group_membership" "membership" {
  group = data.azuredevops_group.group.descriptor
  members = [
    azuredevops_user_entitlement.user.descriptor
  ]
}
```

## Argument Reference

The following arguments are supported:

* `group` - (Required) The descriptor of the group being managed.
* `members` - (Required) A list of user or group descriptors that will become members of the group.
> NOTE: It's possible to define group members both within the `azuredevops_group_membership resource` via the members block and by using the `azuredevops_group` resource. However it's not possible to use both methods to manage group members, since there'll be conflicts.
* `mode` - (Optional) The mode how the resource manages group members.
  * `mode == add`: the resource will ensure that all specified members will be part of the referenced group
  * `mode == overwrite`: the resource will replace all existing members with the members specified within the `members` block 
> NOTE: To clear all members from a group, specify an empty list of descriptors in the `members` attribute and set the `mode` member to `overwrite`.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A random ID for this resource. There is no "natural" ID, so a random one is assigned.

## Relevant Links
* [Azure DevOps Service REST API 5.1 - Memberships](https://docs.microsoft.com/en-us/rest/api/azure/devops/graph/memberships?view=azure-devops-rest-5.0)

## Import

Not supported.

## PAT Permissions Required

- **Deployment Groups**: Read & Manage
