---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_group"
description: |-
  Manages a group within Azure DevOps organization.
---

# azuredevops_group

Manages a group within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name = "Example Project"
}

data "azuredevops_group" "example-readers" {
  project_id = azuredevops_project.example.id
  name       = "Readers"
}

data "azuredevops_group" "example-contributors" {
  project_id = azuredevops_project.example.id
  name       = "Contributors"
}

resource "azuredevops_group" "example" {
  scope        = azuredevops_project.example.id
  display_name = "Example group"
  description  = "Example description"

  members = [
    data.azuredevops_group.example-readers.descriptor,
    data.azuredevops_group.example-contributors.descriptor
  ]
}
```

## Argument Reference

The following arguments are supported:

* `scope` - (Optional) The scope of the group. A descriptor referencing the scope (collection, project) in which the group should be created. If omitted, will be created in the scope of the enclosing account or organization.x

* `origin_id` - (Optional) The OriginID as a reference to a group from an external AD or AAD backed provider. The `scope`, `mail` and `display_name` arguments cannot be used simultaneously with `origin_id`.

* `mail` - (Optional) The mail address as a reference to an existing group from an external AD or AAD backed provider. The `scope`, `origin_id` and `display_name` arguments cannot be used simultaneously with `mail`.

* `display_name` - (Optional) The name of a new Azure DevOps group that is not backed by an external provider. The `origin_id` and `mail` arguments cannot be used simultaneously with `display_name`.

* `description` - (Optional) The Description of the Project.

* `members` - (Optional) The member of the Group.

  ~> **NOTE:** It's possible to define group members both within the `azuredevops_group` resource via the members block and by using the `azuredevops_group_membership` resource. However it's not possible to use both methods to manage group members, since there'll be conflicts.

* `skip_destroy` - (Optional) Meta argument to skip group destroy API call. Might be useful when deployment identity does not have
sufficient organization permissions to do so. Only works in conjunction with `origin_id`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Group.
* `url` - This url is the full route to the source resource of this graph subject.
* `origin` - The type of source provider for the origin identifier (ex:AD, AAD, MSA)
* `subject_kind` - This field identifies the type of the graph subject (ex: Group, Scope, User).
* `domain` - This represents the name of the container of origin for a graph member.
* `principal_name` - This is the PrincipalName of this graph member from the source provider.
* `descriptor` - The identity (subject) descriptor of the Group.
* `group_id` - The ID of the Group.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Groups](https://docs.microsoft.com/en-us/rest/api/azure/devops/graph/groups?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Group.
* `read` - (Defaults to 5 minute) Used when retrieving the Group.
* `update` - (Defaults to 30 minutes) Used when updating the Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Group.

## Import

Azure DevOps groups can be imported using the group identity descriptor, e.g.

```sh
terraform import azuredevops_group.example aadgp.Uy0xLTktMTU1MTM3NDI0NS0xMjA0NDAwOTY5LTI0MDI5ODY0MTMtMjE3OTQwODYxNi0zLTIxNjc2NjQyNTMtMzI1Nzg0NDI4OS0yMjU4MjcwOTc0LTI2MDYxODY2NDU
```

## PAT Permissions Required

- **Project & Team**: Read, Write, & Manage
