# azuredevops_group
Manages a group within Azure DevOps.

## Example Usage

```hcl
provider "azuredevops" {
  version = ">= 0.0.1"
}

resource "azuredevops_project" "tf-project-test-001" {
  project_name = "Test Project"
}

data "azuredevops_group" "tf-project-readers" {
  project_id = azuredevops_project.tf-project-test-001.id
  name = "Readers"
}

data "azuredevops_group" "tf-project-contributors" {
  project_id = azuredevops_project.tf-project-test-001.id
  name = "Contributors"
}

resource "azuredevops_group" "tf-project-group-001" {
  scope        = azuredevops_project.tf-project-test-001.id
  display_name = "Test group"
  description  = "Test description"

  members = [
      data.azuredevops_group.tf-project-readers.descriptor,data.azuredevops_group.tf-project-contributors.descriptor
  ]
}
```

## Argument Reference

The following arguments are supported:

* `scope` - (Optional) The scope of the group. A descriptor referencing the scope (collection, project) in which the group should be created. If omitted, will be created in the scope of the enclosing account or organization.
* `origin_id` - (Optional) The OriginID as a reference to a group from an external AD or AAD backed provider.
* `mail` - (Optional) The mail address as a reference to an existing group from an external AD or AAD backed provider.
* `display_name` - (Optional) The name of a new Azure DevOps group that is not backed by an external provider.
* `description` - (Optional) The Description of the Project.
* `members` - (Optional)
> NOTE: It's possible to define group members both within the azuredevops_group resource via the members block and by using the azuredevops_group_membership resource. However it's not possible to use both methods to manage group members, since there'll be conflicts.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the Group.
* `url` - This url is the full route to the source resource of this graph subject.
* `origin` - The type of source provider for the origin identifier (ex:AD, AAD, MSA)
* `subject_kind` - This field identifies the type of the graph subject (ex: Group, Scope, User).
* `domain` - This represents the name of the container of origin for a graph member.
* `principal_name` - This is the PrincipalName of this graph member from the source provider. 
* `descriptor` - The identity (subject) descriptor of the Group.

## Relevant Links
* [Azure DevOps Service REST API 5.1 - Groups](https://docs.microsoft.com/en-us/rest/api/azure/devops/graph/groups?view=azure-devops-rest-5.1)

## Import
Azure DevOps Projects can be imported using the group identity descriptor, e.g.

```
terraform import azuredevops_project.id aadgp.Uy0xLTktMTU1MTM3NDI0NS0xMjA0NDAwOTY5LTI0MDI5ODY0MTMtMjE3OTQwODYxNi0zLTIxNjc2NjQyNTMtMzI1Nzg0NDI4OS0yMjU4MjcwOTc0LTI2MDYxODY2NDU
```

## PAT Permissions Required

- **Project & Team**: Read, Write, & Manage