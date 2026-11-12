---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_analytics_views_permissions"
description: |-
  Manages permissions for shared Analytics views of a AzureDevOps project
---

# azuredevops_analytics_views_permissions

Manages permissions for the shared Analytics views of a AzureDevOps project.

These permissions are shown as `Edit shared Analytics views` and `Delete shared Analytics views` in the **Analytics** section of the project **Permissions** page.

~> **Note** Permissions can be assigned to group principals and not to single user principals.

~> **Note** The `View analytics` permission is managed by the `azuredevops_analytics_permissions` resource.

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

resource "azuredevops_analytics_views_permissions" "example-permissions" {
  project_id = azuredevops_project.example.id
  principal  = data.azuredevops_group.example-readers.id
  permissions = {
    Read   = "Allow"
    Edit   = "Deny"
    Delete = "Deny"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project to assign the permissions.

* `principal` - (Required) The `group` principal to assign the permissions.

* `view_id` - (Optional) The ID of a shared Analytics view to assign the permissions to. If not set, the permissions are assigned to all shared Analytics views of the project.

* `permissions` - (Required) the permissions to assign. The following permissions are available

    | Permission | Description                   |
    |------------|-------------------------------|
    | Read       | View shared Analytics views   |
    | Edit       | Edit shared Analytics views   |
    | Delete     | Delete shared Analytics views |

~> **Note** The `AnalyticsViews` security namespace also defines the `Execute` and `ManagePermissions` actions. `Execute` is reserved by Azure DevOps and cannot be assigned, and `ManagePermissions` is not shown in the Azure DevOps UI.

~> **Note** Permissions assigned without `view_id` are inherited by every shared Analytics view of the project. A resource that sets `view_id` overrides that inherited value for the given view. Avoid managing the same permission and principal both with and without `view_id`.

---

* `replace` - (Optional) Replace (`true`) or merge (`false`) the permissions. Default: `true`

## Relevant Links

* [Azure DevOps Service REST API 7.1 - Security](https://learn.microsoft.com/en-us/rest/api/azure/devops/security/?view=azure-devops-rest-7.1)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Analytics Views Permission.
* `read` - (Defaults to 5 minute) Used when retrieving the Analytics Views Permission.
* `update` - (Defaults to 10 minutes) Used when updating the Analytics Views Permission.
* `delete` - (Defaults to 10 minutes) Used when deleting the Analytics Views Permission.

## Import

The resource does not support import.

## PAT Permissions Required

- **Project & Team**: vso.security_manage - Grants the ability to read, write, and manage security permissions.
