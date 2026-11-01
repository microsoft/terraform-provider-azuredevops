---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_git_repository_advanced_security"
description: |-
  Manages GitHub Advanced Security for Azure DevOps on a Git repository.
---

# azuredevops_git_repository_advanced_security

Manages GitHub Advanced Security for Azure DevOps (Secret Protection and Code Security) on a Git repository.

~> **NOTE:** Advanced Security is billed per active committer. Removing this resource disables Secret Protection and Code Security on the repository.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_git_repository" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Repository"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_git_repository_advanced_security" "example" {
  project_id    = azuredevops_project.example.id
  repository_id = azuredevops_git_repository.example.id

  secret_protection_enabled             = true
  block_pushes                          = true
  code_security_enabled                 = true
  dependency_scanning_injection_enabled = true
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project. Changing this forces a new resource to be created.

* `repository_id` - (Required) The ID of the Git repository. Changing this forces a new resource to be created.

---

* `secret_protection_enabled` - (Optional) Enable Secret Protection (secret scanning).

* `block_pushes` - (Optional) Block pushes that contain secrets.

* `code_security_enabled` - (Optional) Enable Code Security (dependency and code scanning).

* `dependency_scanning_injection_enabled` - (Optional) Inject dependency scanning into pipelines.

* `codeql_enabled` - (Optional) Use the CodeQL default setup.

* `autofix_enabled` - (Optional) Enable Copilot Autofix.

~> **NOTE:** Settings that are not configured are left unchanged and read from Azure DevOps.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID in the format `<project_id>/<repository_id>`.

## Relevant Links

- [Azure DevOps Service REST API 7.2 - Repo Enablement](https://learn.microsoft.com/en-us/rest/api/azure/devops/advancedsecurity/repo-enablement)
- [Configure GitHub Advanced Security for Azure DevOps](https://learn.microsoft.com/en-us/azure/devops/repos/security/configure-github-advanced-security-features)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when enabling Advanced Security.
* `read` - (Defaults to 5 minute) Used when retrieving the Advanced Security settings.
* `update` - (Defaults to 10 minutes) Used when updating the Advanced Security settings.
* `delete` - (Defaults to 10 minutes) Used when disabling Advanced Security.

## Import

Advanced Security settings of a Git repository can be imported using the project ID and repository ID, e.g.

```sh
terraform import azuredevops_git_repository_advanced_security.example 00000000-0000-0000-0000-000000000000/00000000-0000-0000-0000-000000000000
```

## PAT Permissions Required

- **Advanced Security**: Read, write, & manage
