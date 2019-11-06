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

## Arugument Reference

The following arguments are supported:

* `group` - (Required) The descriptor of the group being managed.
* `members` - (Required) A list of entity user or group descriptors that will become members of the group.

## Attributes Reference

The following attributes are exported:

* `id` - A random ID for this resource. There is no "natural" ID, so a random one is assigned.

## Relevant Links
* [Azure DevOps Service REST API 5.1 - Memberships](https://docs.microsoft.com/en-us/rest/api/azure/devops/graph/memberships?view=azure-devops-rest-5.0)

## Import

Not supported.

## PAT Permissions Required

- **Deployment Groups**: Read & Manage