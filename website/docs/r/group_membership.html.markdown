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
resource "azuredevops_project" "example" {
  name = "Example Project"
}

resource "azuredevops_user_entitlement" "example" {
  principal_name = "foo@contoso.com"
}

data "azuredevops_group" "example" {
  project_id = azuredevops_project.example.id
  name       = "Build Administrators"
}

resource "azuredevops_group_membership" "example" {
  group = data.azuredevops_group.example.descriptor
  members = [
    azuredevops_user_entitlement.example.descriptor
  ]
}
```

## Argument Reference

The following arguments are supported:

- `group` - (Required) The descriptor of the group being managed.
- `members` - (Required) A list of user or group descriptors that will become members of the group.

  ~> **NOTE** It's possible to define group members both within the `azuredevops_group_membership resource` via the members block and by using the `azuredevops_group` resource. However it's not possible to use both methods to manage group members, since there'll be conflicts.

  ~> **NOTE**  The `members` uses `descriptor` as the identifier not Resource ID or others.

- `mode` - (Optional) The mode how the resource manages group members.
  - `mode == add`: the resource will ensure that all specified members will be part of the referenced group
  - `mode == overwrite`: the resource will replace all existing members with the members specified within the `members` block
  
    ~> **NOTE** To clear all members from a group, specify an empty list of descriptors in the `members` attribute and set the `mode` member to `overwrite`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - A random ID for this resource. There is no "natural" ID, so a random one is assigned.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Memberships](https://docs.microsoft.com/en-us/rest/api/azure/devops/graph/memberships?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Group membership.
* `read` - (Defaults to 5 minute) Used when retrieving the Group membership.
* `update` - (Defaults to 10 minutes) Used when updating the Group membership.
* `delete` - (Defaults to 10 minutes) Used when deleting the Group membership.

## PAT Permissions Required

- **Deployment Groups**: Read & Manage
