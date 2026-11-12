---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_analytics_permissions"
description: |-
  Manages Analytics permissions for a AzureDevOps project
---

# azuredevops_analytics_permissions

Manages Analytics permissions for a AzureDevOps project.

These permissions are shown in the **Analytics** section of the project **Permissions** page.

~> **Note** Permissions can be assigned to group principals and not to single user principals.

~> **Note** Permissions for shared Analytics views are managed by the `azuredevops_analytics_views_permissions` resource.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

data "azuredevops_group" "example-readers" {
  project_id = azuredevops_project.example.id
  name       = "Readers"
}

resource "azuredevops_analytics_permissions" "example-permissions" {
  project_id = azuredevops_project.example.id
  principal  = data.azuredevops_group.example-readers.id
  permissions = {
    Read = "Allow"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project to assign the permissions.

* `principal` - (Required) The `group` principal to assign the permissions.

* `permissions` - (Required) the permissions to assign. The following permissions are available

    | Permission | Description    |
    |------------|----------------|
    | Read       | View analytics |

~> **Note** The `Analytics` security namespace also defines the `Administer`, `Stage`, `ExecuteUnrestrictedQuery` and `ReadEuii` actions, but those bits are reserved by Azure DevOps and cannot be assigned.

---

* `replace` - (Optional) Replace (`true`) or merge (`false`) the permissions. Default: `true`

## Relevant Links

* [Azure DevOps Service REST API 7.1 - Security](https://learn.microsoft.com/en-us/rest/api/azure/devops/security/?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Analytics Permission.
* `read` - (Defaults to 5 minute) Used when retrieving the Analytics Permission.
* `update` - (Defaults to 10 minutes) Used when updating the Analytics Permission.
* `delete` - (Defaults to 10 minutes) Used when deleting the Analytics Permission.

## Import

The resource does not support import.

## PAT Permissions Required

- **Project & Team**: vso.security_manage - Grants the ability to read, write, and manage security permissions.
