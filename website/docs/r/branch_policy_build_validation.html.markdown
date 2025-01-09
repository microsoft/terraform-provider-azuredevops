---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_branch_policy_build_validation"
description: |-
  Manages a build validation branch policy within Azure DevOps project.
---

# azuredevops_branch_policy_build_validation

Manages a build validation branch policy within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name = "Example Project"
}

resource "azuredevops_git_repository" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Repository"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_build_definition" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Build Definition"

  repository {
    repo_type = "TfsGit"
    repo_id   = azuredevops_git_repository.example.id
    yml_path  = "azure-pipelines.yml"
  }
}

resource "azuredevops_branch_policy_build_validation" "example" {
  project_id = azuredevops_project.example.id

  enabled  = true
  blocking = true

  settings {
    display_name                = "Example build validation policy"
    build_definition_id         = azuredevops_build_definition.example.id
    queue_on_source_update_only = true
    valid_duration              = 720
    filename_patterns = [
      "/WebApp/*",
      "!/WebApp/Tests/*",
      "*.cs"
    ]

    scope {
      repository_id  = azuredevops_git_repository.example.id
      repository_ref = azuredevops_git_repository.example.default_branch
      match_type     = "Exact"
    }

    scope {
      repository_id  = azuredevops_git_repository.example.id
      repository_ref = "refs/heads/releases"
      match_type     = "Prefix"
    }

    scope {
      match_type = "DefaultBranch"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The ID of the project in which the policy will be created.
- `enabled` - (Optional) A flag indicating if the policy should be enabled. Defaults to `true`.
- `blocking` - (Optional) A flag indicating if the policy should be blocking. Defaults to `true`.
- `settings` - (Required) Configuration for the policy. This block must be defined exactly once.

A `settings` block supports the following:

- `build_definition_id` - (Required) The ID of the build to monitor for the policy.
- `display_name` - (Required) The display name for the policy.
- `manual_queue_only` - (Optional) If set to true, the build will need to be manually queued. Defaults to `false`
- `queue_on_source_update_only` - (Optional) True if the build should queue on source updates only. Defaults to `true`.
- `valid_duration` - (Optional) The number of minutes for which the build is valid. If `0`, the build will not expire. Defaults to `720` (12 hours).

~> **Note** Combine `valid_duration` and `queue_on_source_update_only` to set the build expiration.   
    1.  Expire immediately when branch is updated: `valid_duration=0` and `queue_on_source_update_only=false`   
    2.  Expire after a period of time : `valid_duration=360` and `queue_on_source_update_only=true`   
    3.  Never expire: `valid_duration=0` and `queue_on_source_update_only=true`

- `filename_patterns` - (Optional) If a path filter is set, the policy will only apply when files which match the filter are changes. Not setting this field means that the policy will always apply. You can specify absolute paths and wildcards. Example: `["/WebApp/Models/Data.cs", "/WebApp/*", "*.cs"]`. Paths prefixed with "!" are excluded. Example: `["/WebApp/*", "!/WebApp/Tests/*"]`. Order is significant.
- `scope` (Required) Controls which repositories and branches the policy will be enabled for. This block must be defined at least once.

A `settings` `scope` block supports the following:

- `repository_id` - (Optional) The repository ID. Needed only if the scope of the policy will be limited to a single repository. If `match_type` is `DefaultBranch`, this should not be defined.
- `repository_ref` - (Optional) The ref pattern to use for the match when `match_type` other than `DefaultBranch`. If `match_type` is `Exact`, this should be a qualified ref such as `refs/heads/master`. If `match_type` is `Prefix`, this should be a ref path such as `refs/heads/releases`.
- `match_type` (Optional) The match type to use when applying the policy. Supported values are `Exact` (default), `Prefix` or `DefaultBranch`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of branch policy configuration.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Policy Configurations](https://docs.microsoft.com/en-us/rest/api/azure/devops/policy/configurations/create?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Build Validation Branch Policy.
* `read` - (Defaults to 2 minute) Used when retrieving the Build Validation Branch Policy.
* `update` - (Defaults to 5 minutes) Used when updating the Build Validation Branch Policy.
* `delete` - (Defaults to 5 minutes) Used when deleting the Build Validation Branch Policy.

## Import

Azure DevOps Branch Policies can be imported using the project ID and policy configuration ID:

```sh
terraform import azuredevops_branch_policy_build_validation.example 00000000-0000-0000-0000-000000000000/0
```
