---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_security_namespaces"
description: |-
  Use this data source to access information about security namespaces within Azure DevOps.
---

# Data Source: azuredevops_security_namespaces

Use this data source to access information about security namespaces within Azure DevOps. Security namespaces define the security model for different resources and operations in Azure DevOps.

## Example Usage

### List All Security Namespaces

```hcl
data "azuredevops_security_namespaces" "all" {
}

output "namespaces" {
  value = data.azuredevops_security_namespaces.all.namespaces
}
```

### Find a Specific Namespace by Name

```hcl
data "azuredevops_security_namespaces" "all" {
}

locals {
  git_namespace = [
    for ns in data.azuredevops_security_namespaces.all.namespaces :
    ns if ns.name == "Git Repositories"
  ][0]
}

output "git_namespace_id" {
  value = local.git_namespace.namespace_id
}

output "git_permissions" {
  value = local.git_namespace.actions
}
```

### Discover Available Permissions for a Namespace

```hcl
data "azuredevops_security_namespaces" "all" {
}

locals {
  project_namespace = [
    for ns in data.azuredevops_security_namespaces.all.namespaces :
    ns if ns.name == "Project"
  ][0]
}

output "project_permissions" {
  value = {
    for action in local.project_namespace.actions :
    action.name => {
      display_name = action.display_name
      bit          = action.bit
    }
  }
}
```

## Argument Reference

This data source does not require any arguments.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the data source in the format `security-namespaces-{uuid}`.

* `namespaces` - A set of security namespaces. Each `namespace` block exports the following:

---

A `namespace` block exports the following:

* `namespace_id` - The unique identifier (UUID) of the security namespace.

* `name` - The name of the security namespace (e.g., "Git Repositories", "Project").

* `display_name` - The display name of the security namespace.

* `description` - The description of the security namespace (if available).

* `actions` - A set of available actions (permissions) in this namespace. Each `action` block exports the following:
  * `name` - The name of the action/permission.
  * `display_name` - The display name of the action/permission.
  * `bit` - The bit value for this permission (used in permission calculations).
  * `namespace_id` - The namespace ID this action belongs to.

## Common Security Namespaces

The following are common security namespaces available in Azure DevOps:

| Namespace Name | Namespace ID | Description |
|---------------|--------------|-------------|
| **Collection** | `3e65f728-f8bc-4ecd-8764-7e378b19bfa7` | Organization/collection-level security |
| **Project** | `52d39943-cb85-4d7f-8fa8-c6baac873819` | Project-level security |
| **Git Repositories** | `2e9eb7ed-3c0a-47d4-87c1-0ffdd275fd87` | Git repository security |
| **Analytics** | `58450c49-b02d-465a-ab12-59ae512d6531` | Analytics security |
| **AnalyticsViews** | `d34d3680-dfe5-4cc6-a949-7d9c68f73cba` | Analytics Views security |
| **Process** | `2dab47f9-bd70-49ed-9bd5-8eb051e59c02` | Process template security |
| **AuditLog** | `a6cc6381-a1ca-4b36-b3c1-4e65211e82b6` | Audit log security |
| **BuildAdministration** | `302acaca-b667-436d-a946-87133492041c` | Build administration security |
| **Server** | `1f4179b3-6bac-4d01-b421-71ea09171400` | Server-level security |
| **VersionControlPrivileges** | `66312704-deb5-43f9-b51c-ab4ff5e351c3` | Version control privileges |

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Security Namespaces - Query](https://docs.microsoft.com/en-us/rest/api/azure/devops/security/security-namespaces/query?view=azure-devops-rest-7.0)
- [Security Namespaces Documentation](https://docs.microsoft.com/en-us/azure/devops/organizations/security/security-glossary?view=azure-devops)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving security namespaces.

## PAT Permissions Required

- **Project & Team**: Read

## Notes

- Security namespaces define the security model for different resources and operations in Azure DevOps
- Each namespace has a unique identifier (UUID) that doesn't change across organizations
- Namespaces contain actions (permissions) that can be granted or denied to users and groups
- Permission bits are used to calculate effective permissions when multiple permissions are set
- This data source is useful for discovering available permissions and namespace IDs for use with `azuredevops_security_permissions` resources

