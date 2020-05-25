---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_branch_policy_min_reviewers"
description: |-
  Manages a minimum reviewer branch policy within Azure DevOps project.
---

# azuredevops_branch_policy_min_reviewers
Manages a minimum reviewer branch policy within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "p" {
  project_name = "Sample Project"
}

resource "azuredevops_git_repository" "r" {
  project_id = azuredevops_project.p.id
  name       = "Sample Repo"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_branch_policy_min_reviewers" "p" {
  project_id = azuredevops_project.p.id

  enabled  = true
  blocking = true

  settings {
    reviewer_count     = 2
    submitter_can_vote = false

    scope {
      repository_id  = azuredevops_git_repository.r.id
      repository_ref = azuredevops_git_repository.r.default_branch
      match_type     = "Exact"
    }

    scope {
      repository_id  = azuredevops_git_repository.r.id
      repository_ref = "refs/heads/releases"
      match_type     = "Prefix"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project in which the policy will be created.
* `enabled` - (Optional) A flag indicating if the policy should be enabled. Defaults to `true`.
* `blocking` - (Optional) A flag indicating if the policy should be blocking. Defaults to `true`.
* `settings` - (Required) Configuration for the policy. This block must be defined exactly once.

A `settings` block supports the following:

* `reviewer_count` - (Required) The number of reviewrs needed to approve.
* `submitter_can_vote` - (Optional) Controls whether or not the submitter's vote counts. Defaults to `false`.
* `scope` (Required) Controls which repositories and branches the policy will be enabled for. This block must be defined at least once.

A `settings` `scope` block supports the following:
* `repository_id` - (Optional) The repository ID. Needed only if the scope of the policy will be limited to a single repository.
* `repository_ref` - (Optional) The ref pattern to use for the match. If `match_type` is `Exact`, this should be a qualified ref such as `refs/heads/master`. If `match_type` is `Prefix`, this should be a ref path such as `refs/heads/releases`.
* `match_type` (Optional) The match type to use when applying the policy. Supported values are `Exact` (default) or `Prefix`.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of branch policy configuration.


## Relevant Links
* [Azure DevOps Service REST API 5.1 - Policy Configurations](https://docs.microsoft.com/en-us/rest/api/azure/devops/policy/configurations/create?view=azure-devops-rest-5.1)

## Import
Azure DevOps Branch Policies can be imported using the project ID and policy configuration ID:

```
terraform import azuredevops_branch_policy_min_reviewers.p aa4a9756-8a86-4588-86d7-b3ee2d88b033/60```
