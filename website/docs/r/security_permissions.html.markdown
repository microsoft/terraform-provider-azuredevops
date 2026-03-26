---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_security_permissions"
description: |-
  Manages permissions for Azure DevOps security namespaces
---

# azuredevops_security_permissions

Manages permissions for Azure DevOps security namespaces. This is a generic permissions resource that can be used to manage permissions for any security namespace in Azure DevOps.

~> **Note** This is a low-level generic permissions resource. For specific resource types, consider using the dedicated permission resources such as `azuredevops_project_permissions`, `azuredevops_git_permissions`, `azuredevops_build_definition_permissions`, etc.

## Example Usage

### Collection-level Permissions

```hcl
data "azuredevops_security_namespace" "collection" {
  name = "Collection"
}

data "azuredevops_security_namespace_token" "collection" {
  namespace_name = "Collection"
}

data "azuredevops_group" "example" {
  name = "Project Collection Administrators"
}

resource "azuredevops_security_permissions" "collection_perms" {
  namespace_id = data.azuredevops_security_namespace.collection.id
  token        = data.azuredevops_security_namespace_token.collection.token
  principal    = data.azuredevops_group.example.descriptor
  permissions = {
    "GENERIC_READ"  = "allow"
    "GENERIC_WRITE" = "allow"
  }
}
```

### Project-level Permissions

```hcl
data "azuredevops_project" "example" {
  name = "Example Project"
}

data "azuredevops_security_namespace_token" "project" {
  namespace_name = "Project"
  identifiers = {
    project_id = data.azuredevops_project.example.id
  }
}

data "azuredevops_group" "example_readers" {
  project_id = data.azuredevops_project.example.id
  name       = "Readers"
}

data "azuredevops_namespace" "project" {
  name = "Project"
}

resource "azuredevops_security_permissions" "project_perms" {
  namespace_id = data.azuredevops_namespace.project.id
  token        = data.azuredevops_security_namespace_token.project.token
  principal    = data.azuredevops_group.example_readers.descriptor
  permissions = {
    "GENERIC_READ"  = "allow"
    "GENERIC_WRITE" = "deny"
    "DELETE"        = "deny"
  }
}
```

### Git Repository Permissions

```hcl
data "azuredevops_project" "example" {
  name = "Example Project"
}

data "azuredevops_git_repository" "example" {
  project_id = data.azuredevops_project.example.id
  name       = "Example Repository"
}

data "azuredevops_security_namespace" "git_repos" {
  name = "Git Repositories"
}

data "azuredevops_security_namespace_token" "git_repo" {
  namespace_name = "Git Repositories"
  identifiers = {
    project_id    = data.azuredevops_project.example.id
    repository_id = data.azuredevops_git_repository.example.id
  }
}

data "azuredevops_group" "example_contributors" {
  project_id = data.azuredevops_project.example.id
  name       = "Contributors"
}

resource "azuredevops_security_permissions" "git_perms" {
  namespace_id = data.azuredevops_security_namespace.git_repos.id
  token        = data.azuredevops_security_namespace_token.git_repo.token
  principal    = data.azuredevops_group.example_contributors.descriptor
  permissions = {
    "GenericRead"       = "allow"
    "GenericContribute" = "allow"
    "ForcePush"         = "deny"
    "ManagePermissions" = "deny"
  }
  replace = false
}
```

### Git Branch Permissions

```hcl
data "azuredevops_project" "example" {
  name = "Example Project"
}

data "azuredevops_git_repository" "example" {
  project_id = data.azuredevops_project.example.id
  name       = "Example Repository"
}

data "azuredevops_security_namespace" "git_repos" {
  name = "Git Repositories"
}

data "azuredevops_security_namespace_token" "main_branch" {
  namespace_name = "Git Repositories"
  identifiers = {
    project_id    = data.azuredevops_project.example.id
    repository_id = data.azuredevops_git_repository.example.id
    ref_name      = "refs/heads/main"
  }
}

data "azuredevops_group" "example_contributors" {
  project_id = data.azuredevops_project.example.id
  name       = "Contributors"
}

resource "azuredevops_security_permissions" "main_branch_perms" {
  namespace_id = "2e9eb7ed-3c0a-47d4-87c1-0ffdd275fd87"
  token        = data.azuredevops_security_namespace_token.main_branch.token
  principal    = data.azuredevops_group.example_contributors.descriptor
  permissions = {
    "ForcePush"     = "Deny"
    "RemoveOthersLocks" = "Deny"
  }
  replace = false
}
```

## Argument Reference

The following arguments are supported:

* `namespace_id` - (Required) The ID of the security namespace. Use the `azuredevops_security_namespaces` data source to discover available namespaces. Changing this forces a new resource to be created.

* `token` - (Required) The security token for the resource. Use the `azuredevops_security_namespace_token` data source to generate tokens for specific resources. Changing this forces a new resource to be created.

* `principal` - (Required) The descriptor or identity ID of the principal (user or group). Changing this forces a new resource to be created.

* `permissions` - (Required) A map of permission names to permission values. All permission names specified must be valid for the given namespace, or an error will be returned. Permission values must be one of:
  - `Allow` (or `allow`, `ALLOW`) - Grant the permission
  - `Deny` (or `deny`, `DENY`) - Explicitly deny the permission
  - `NotSet` (or `notset`, `NOTSET`) - Remove the permission (inherit from parent)

* `replace` - (Optional) Replace (`true`) or merge (`false`) the permissions with existing permissions. When `true`, all existing permissions for the principal on this token will be replaced with the specified permissions. When `false`, the specified permissions will be merged with existing permissions. Default: `true`.

### Permission Names by Namespace

Permission names vary by namespace. Use the `azuredevops_security_namespaces` data source to discover available permissions for each namespace. Common namespaces and their permissions:

#### Collection Namespace (`3e65f728-f8bc-4ecd-8764-7e378b19bfa7`)

* `GENERIC_READ` - View instance-level information
* `GENERIC_WRITE` - Edit instance-level information
* `DELETE_FIELD` - Delete field from organization
* `MANAGE_PROPERTIES` - Manage collection properties
* `MANAGE_TEST_CONTROLLERS` - Manage test controllers
* `TRIGGER_EVENT` - Trigger organization-level events

#### Project Namespace (`52d39943-cb85-4d7f-8fa8-c6baac873819`)

* `GENERIC_READ` - View project-level information
* `GENERIC_WRITE` - Edit project-level information
* `DELETE` - Delete team project
* `PUBLISH_TEST_RESULTS` - Create test runs
* `MANAGE_PROPERTIES` - Manage project properties
* `RENAME` - Rename team project
* `UPDATE_VISIBILITY` - Update project visibility
* And many more...

#### Git Repositories Namespace (`2e9eb7ed-3c0a-47d4-87c1-0ffdd275fd87`)

* `GenericRead` - Read repository
* `GenericContribute` - Contribute
* `ForcePush` - Force push (rewrite history, delete branches and tags)
* `CreateBranch` - Create branch
* `CreateTag` - Create tag
* `ManageNote` - Manage notes
* `PolicyExempt` - Bypass policies when pushing
* `CreateRepository` - Create repository
* `DeleteRepository` - Delete repository
* `RenameRepository` - Rename repository
* `EditPolicies` - Edit policies
* `RemoveOthersLocks` - Remove others' locks
* `ManagePermissions` - Manage permissions
* `PullRequestContribute` - Contribute to pull requests
* `PullRequestBypassPolicy` - Bypass policies when completing pull requests

~> **Note** Permission names are case-sensitive and must match exactly as defined in the namespace. Use the `azuredevops_security_namespaces` data source to discover the exact permission names.

## Attributes Reference

No additional attributes are exported.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Security Permission.
* `read` - (Defaults to 5 minutes) Used when retrieving the Security Permission.
* `update` - (Defaults to 10 minutes) Used when updating the Security Permission.
* `delete` - (Defaults to 10 minutes) Used when deleting the Security Permission.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Security](https://docs.microsoft.com/en-us/rest/api/azure/devops/security/?view=azure-devops-rest-7.0)
- [Azure DevOps Security Namespaces](https://docs.microsoft.com/en-us/azure/devops/organizations/security/security-glossary?view=azure-devops)

## Import

The resource does not support import.

## PAT Permissions Required

- **Project & Team**: vso.security_manage - Grants the ability to read, write, and manage security permissions.

## Notes

- This is a generic low-level resource for managing permissions across any Azure DevOps security namespace
- For better user experience and type safety, consider using dedicated permission resources when available (e.g., `azuredevops_project_permissions`, `azuredevops_git_permissions`)
- Permission names are namespace-specific and case-sensitive. All permission names in the `permissions` map are validated against the namespace - if any permission name is invalid, an error will be returned
- When `replace = true`, all existing permissions for the principal will be removed and replaced with the specified permissions
- When `replace = false`, the specified permissions will be merged with existing permissions, allowing you to manage only a subset of permissions
- when `replace = false`, deletion of the resource only removes the permissions specified in the `permissions` map, rather than all permissions for the principal
- The `principal` must be a group descriptor or identity ID. Individual user principals are not supported
- Use the `azuredevops_security_namespace_token` data source to generate correct tokens for different resource types
- Permissions are propagated asynchronously and may take a few moments to take effect

