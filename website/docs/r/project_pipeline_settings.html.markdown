---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_project_pipeline_settings"
description: |-
  Manages Pipeline Settings for Azure DevOps projects.
---

# azuredevops_project_pipeline_settings

Manages Pipeline Settings for Azure DevOps projects

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_project_pipeline_settings" "example" {
  project_id = azuredevops_project.example.id

  enforce_job_scope                    = true
  enforce_referenced_repo_scoped_token = false
  enforce_settable_var                 = true
  publish_pipeline_metadata            = false
  status_badges_are_private            = true
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project for which the project pipeline settings will be managed.

---

* `enforce_job_scope` - (Optional) Limit job authorization scope to current project for non-release pipelines.

* `enforce_referenced_repo_scoped_token` - (Optional) Protect access to repositories in YAML pipelines.

* `enforce_settable_var` - (Optional) Limit variables that can be set at queue time.

* `publish_pipeline_metadata` - (Optional) Publish metadata from pipelines.

* `status_badges_are_private` - (Optional) Disable anonymous access to badges.

* `enforce_job_scope_for_release` - (Optional) Limit job authorization scope to current project for release pipelines.

~> **NOTE:** The settings at the organization will override settings specified on the project. 
  For example, if `enforce_job_scope` is true at the organization, the `azuredevops_project_pipeline_settings` resource cannot set it to false. 
  In this scenario, the plan will always show that the resource is trying to change `enforce_job_scope` from `true` to `false`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the project.

## Relevant Links

No official documentation available

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Project Pipeline Settings.
* `read` - (Defaults to 5 minute) Used when retrieving the Project Pipeline Settings.
* `update` - (Defaults to 10 minutes) Used when updating the Project Pipeline Settings.
* `delete` - (Defaults to 10 minutes) Used when deleting the Project Pipeline Settings.

## Import

Azure DevOps feature settings can be imported using the project id, e.g.

```sh
terraform import azuredevops_project_pipeline_settings.example 00000000-0000-0000-0000-000000000000
```

## PAT Permissions Required

- Full Access
