---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_project_pipeline_retention_settings"
description: |-
  Manages Pipeline Retention Settings for Azure DevOps projects.
---

# azuredevops_project_pipeline_retention_settings

Manages Pipeline Retention Settings for Azure DevOps projects

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_project_pipeline_retention_settings" "example" {
  project_id = azuredevops_project.example.id

  run_retention                    = 30
  artifact_retention               = 20
  pull_request_run_retention       = 15
  retain_runs_per_protected_branch = 10
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project for which the project pipeline retention settings will be managed.

---

* `run_retention` - (Optional) The number of days to retain pipeline runs.

* `artifact_retention` - (Optional) The number of days to retain artifacts. Artifacts can not live longer than a run, so will be overridden by a shorter run retention setting.

* `pull_request_run_retention` - (Optional) The number of days to retain pull request pipeline runs.

* `retain_runs_per_protected_branch` - (Optional) The number of runs to retain per protected branch.

~> **NOTE:** The allowed range for each setting is determined by the organization's retention policy and is not managed by this resource. If a value outside the allowed range is specified, the Azure DevOps API will reject the change.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the project.

## Relevant Links

No official documentation available

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Project Pipeline Retention Settings.
* `read` - (Defaults to 5 minute) Used when retrieving the Project Pipeline Retention Settings.
* `update` - (Defaults to 10 minutes) Used when updating the Project Pipeline Retention Settings.
* `delete` - (Defaults to 10 minutes) Used when deleting the Project Pipeline Retention Settings.

## Import

Azure DevOps project pipeline retention settings can be imported using the project id, e.g.

```sh
terraform import azuredevops_project_pipeline_retention_settings.example 00000000-0000-0000-0000-000000000000
```

## PAT Permissions Required

- Full Access
