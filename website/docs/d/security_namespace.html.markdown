---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_security_namespace"
description: |-
  Use this data source to access information about a specific security namespace within Azure DevOps.
---

# Data Source: azuredevops_security_namespace

Use this data source to access information about a specific security namespace within Azure DevOps. Security namespaces define the security model for different resources and operations in Azure DevOps.

## Example Usage

### Find a Specific Namespace by Name

```hcl
data "azuredevops_security_namespace" "git" {
  name = "Git Repositories"
}

output "git_id" {
  value = data.azuredevops_security_namespace.git.id
}

output "git_permissions" {
  value = data.azuredevops_security_namespace.git.actions
}
```

### Find a Specific Namespace by ID

```hcl
data "azuredevops_security_namespace" "project" {
  id = "52d39943-cb85-4d7f-8fa8-c6baac873819"
}

output "project_namespace_name" {
  value = data.azuredevops_security_namespace.project.name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the security namespace.
* `id` - (Optional) The ID of the security namespace.

~> **NOTE:** One of `name` or `id` must be specified.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the data source (same as `id`).
* `id` - The unique identifier (UUID) of the security namespace.
* `name` - The name of the security namespace.
* `display_name` - The display name of the security namespace.
* `actions` - A set of available actions (permissions) in this namespace. Each `action` block exports the following:
  * `name` - The name of the action/permission.
  * `display_name` - The display name of the action/permission.
  * `bit` - The bit value for this permission (used in permission calculations).

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Security Namespaces - Query](https://docs.microsoft.com/en-us/rest/api/azure/devops/security/security-namespaces/query?view=azure-devops-rest-7.0)
- [Security Namespaces Documentation](https://docs.microsoft.com/en-us/azure/devops/organizations/security/security-glossary?view=azure-devops)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving security namespaces.

