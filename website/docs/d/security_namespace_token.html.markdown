---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_security_namespace_token"
description: |-
  Use this data source to generate security tokens for Azure DevOps security namespaces.
---

# Data Source: azuredevops_security_namespace_token

Use this data source to generate security tokens for Azure DevOps security namespaces. Security tokens are required when managing permissions with the `azuredevops_security_permissions` resource.

## Example Usage

### Discovering Required Identifiers for a Namespace

```hcl
data "azuredevops_security_namespace_token" "git_info" {
  namespace_name         = "Git Repositories"
  return_identifier_info = true
}

output "git_required_identifiers" {
  value = data.azuredevops_security_namespace_token.git_info.required_identifiers
}

output "git_optional_identifiers" {
  value = data.azuredevops_security_namespace_token.git_info.optional_identifiers
}
```

### Collection-level Token

```hcl
data "azuredevops_security_namespace_token" "collection" {
  namespace_name = "Collection"
}

output "collection_token" {
  value = data.azuredevops_security_namespace_token.collection.token
}
```

### Project-level Token

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

output "project_token" {
  value = data.azuredevops_security_namespace_token.project.token
}
```

### Git Repository Token

```hcl
data "azuredevops_project" "example" {
  name = "Example Project"
}

data "azuredevops_git_repository" "example" {
  project_id = data.azuredevops_project.example.id
  name       = "Example Repository"
}

data "azuredevops_security_namespace_token" "git_repo" {
  namespace_name = "Git Repositories"
  identifiers = {
    project_id    = data.azuredevops_project.example.id
    repository_id = data.azuredevops_git_repository.example.id
  }
}

output "git_repo_token" {
  value = data.azuredevops_security_namespace_token.git_repo.token
}
```

### Git Repository Branch Token

```hcl
data "azuredevops_project" "example" {
  name = "Example Project"
}

data "azuredevops_git_repository" "example" {
  project_id = data.azuredevops_project.example.id
  name       = "Example Repository"
}

data "azuredevops_security_namespace_token" "git_branch" {
  namespace_name = "Git Repositories"
  identifiers = {
    project_id    = data.azuredevops_project.example.id
    repository_id = data.azuredevops_git_repository.example.id
    ref_name      = "refs/heads/main"
  }
}

output "git_branch_token" {
  value = data.azuredevops_security_namespace_token.git_branch.token
}
```

### Using Namespace ID

```hcl
data "azuredevops_security_namespace_token" "analytics" {
  namespace_id = "58450c49-b02d-465a-ab12-59ae512d6531"
  identifiers = {
    project_id = data.azuredevops_project.example.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `namespace_id` - (Optional) The ID of the security namespace. Conflicts with `namespace_name`.

* `namespace_name` - (Optional) The name of the security namespace. Conflicts with `namespace_id`. Common values include:
  - `Collection` - Organization/collection-level permissions
  - `Project` - Project-level permissions
  - `Git Repositories` - Git repository permissions
  - `Analytics` - Analytics permissions
  - `AnalyticsViews` - Analytics Views permissions
  - `Process` - Process permissions
  - `AuditLog` - Audit log permissions
  - `BuildAdministration` - Build administration permissions
  - `Server` - Server-level permissions
  - `VersionControlPrivileges` - Version control privileges

* `identifiers` - (Optional) A map of identifiers required for token generation. The required identifiers depend on the namespace. Not used when `return_identifier_info` is `true`.

* `return_identifier_info` - (Optional) When set to `true`, the data source will return the lists of required and optional identifiers for the namespace instead of generating a token. This is useful for discovering what identifiers are needed for a particular namespace. Default: `false`.

~> **NOTE:** One of either `namespace_id` or `namespace_name` must be specified.

### Namespace-Specific Identifiers

Different namespaces require different identifiers:

| Namespace | Required Identifiers | Optional Identifiers | Example |
|-----------|---------------------|---------------------|---------|
| **Collection** | None | None | `{}` |
| **Project** | `project_id` | None | `{project_id = "..."}` |
| **Git Repositories** | `project_id` | `repository_id`, `ref_name` | `{project_id = "...", repository_id = "..."}` |
| **Analytics** | `project_id` | None | `{project_id = "..."}` |
| **AnalyticsViews** | `project_id` | None | `{project_id = "..."}` |
| **Process** | None | None | `{}` |
| **AuditLog** | None | None | `{}` |
| **BuildAdministration** | None | None | `{}` |
| **Server** | None | None | `{}` |
| **VersionControlPrivileges** | None | None | `{}` |

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the data source in the format `ns-token-{namespace_id}-{token}` or `ns-info-{namespace_id}` when `return_identifier_info` is `true`.

* `token` - The generated security token for the namespace. This token can be used with the `azuredevops_security_permissions` resource. Only populated when `return_identifier_info` is `false`.

* `required_identifiers` - A list of required identifier names for this namespace. Only populated when `return_identifier_info` is `true`.

* `optional_identifiers` - A list of optional identifier names for this namespace. Only populated when `return_identifier_info` is `true`.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Security](https://docs.microsoft.com/en-us/rest/api/azure/devops/security/?view=azure-devops-rest-7.0)
- [Security Namespaces Documentation](https://docs.microsoft.com/en-us/azure/devops/organizations/security/security-glossary?view=azure-devops)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when generating the security token.

## Notes

- Security tokens are namespace-specific string identifiers that represent resources within Azure DevOps
- Tokens are used in conjunction with security permissions to control access to various resources
- The format of the token varies depending on the namespace and the resource being targeted
- For Git repositories, you can specify tokens at the repository level or branch level by providing the `ref_name` identifier
- Branch reference names must follow the Git reference format (e.g., `refs/heads/main`, `refs/tags/v1.0.0`)

