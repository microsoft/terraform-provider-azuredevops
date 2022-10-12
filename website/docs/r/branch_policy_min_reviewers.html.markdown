---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_branch_policy_min_reviewers"
description: |-
  Manages a minimum reviewer branch policy within Azure DevOps project.
---

# azuredevops_branch_policy_min_reviewers

Branch policy for reviewers on pull requests. Includes the minimum number of reviewers and other conditions.

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

resource "azuredevops_branch_policy_min_reviewers" "example" {
  project_id = azuredevops_project.example.id

  enabled  = true
  blocking = true

  settings {
    reviewer_count                         = 7
    submitter_can_vote                     = false
    last_pusher_cannot_approve             = true
    allow_completion_with_rejects_or_waits = false
    on_push_reset_approved_votes           = true # OR on_push_reset_all_votes = true
    on_last_iteration_require_vote         = false

    scope {
      repository_id  = azuredevops_git_repository.example.id
      repository_ref = azuredevops_git_repository.example.default_branch
      match_type     = "Exact"
    }

    scope {
      repository_id  = null # All repositories in the project
      repository_ref = "refs/heads/releases"
      match_type     = "Prefix"
    }
    
    scope {
      match_type     = "DefaultBranch"
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

- `reviewer_count` - (Required) The number of reviewers needed to approve.
- `submitter_can_vote` - (Optional) Allow requesters to approve their own changes. Defaults to `false`.
- `last_pusher_cannot_approve`(Optional) Prohibit the most recent pusher from approving their own changes. Defaults to `false`.
- `allow_completion_with_rejects_or_waits` (Optional) Allow completion even if some reviewers vote to wait or reject. Defaults to `false`.
- `on_push_reset_approved_votes` (Optional) When new changes are pushed reset all approval votes (does not reset votes to reject or wait). Defaults to `false`.
- `on_push_reset_all_votes` (Optional) When new changes are pushed reset all code reviewer votes. Defaults to `false`.
- `on_last_iteration_require_vote` (Optional) On last iteration require vote. Defaults to `false`.

If `on_push_reset_all_votes` is `true`, then `on_push_reset_approved_votes` also must be `true`.

- `scope` (Required) Controls which repositories and branches the policy will be enabled for. This block must be defined at least once.

A `settings` `scope` block supports the following:

- `repository_id` - (Optional) The repository ID. Needed only if the scope of the policy will be limited to a single repository. If `match_type` is `DefaultBranch`, this should not be defined.
- `repository_ref` - (Optional) The ref pattern to use for the match when `match_type` other than `DefaultBranch`. If `match_type` is `Exact`, this should be a qualified ref such as `refs/heads/master`. If `match_type` is `Prefix`, this should be a ref path such as `refs/heads/releases`.
- `match_type` (Optional) The match type to use when applying the policy. Supported values are `Exact` (default), `Prefix` or `DefaultBranch`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of branch policy configuration.

## Relevant Links

- [Azure DevOps Service REST API 6.0 - Policy Configurations](https://docs.microsoft.com/en-us/rest/api/azure/devops/policy/configurations/create?view=azure-devops-rest-6.0)

## Import

Azure DevOps Branch Policies can be imported using the project ID and policy configuration ID:

```sh
terraform import azuredevops_branch_policy_min_reviewers.example 00000000-0000-0000-0000-000000000000/0
```
