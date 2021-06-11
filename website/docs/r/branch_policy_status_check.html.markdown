---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_branch_policy_status_check"
description: |- Manages status check branch policy within Azure DevOps project.
---

# azuredevops_branch_policy_status_check

Manages a status check branch policy within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "p" {
  name               = "Sample Project"
  description        = "Managed by Terraform"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  features = {
    "testplans" = "disabled"
    "artifacts" = "disabled"
  }
}

resource "azuredevops_git_repository" "r" {
  project_id = azuredevops_project.p.id
  name       = "Sample Repo"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_user_entitlement" "user" {
  principal_name       = "mail@email.com"
  account_license_type = "basic"
}

resource "azuredevops_branch_policy_status_check" "p" {
  project_id = azuredevops_project.p.id

  enabled  = true
  blocking = true

  settings {
    name                 = "Release"
    author_id            = azuredevops_user_entitlement.user.id
    invalidate_on_update = true
    applicability        = "conditional"
    display_name         = "PreCheck"

    scope {
      repository_id  = azuredevops_git_repository.r.id
      repository_ref = azuredevops_git_repository.r.default_branch
      match_type     = "Exact"
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

`settings` block supports the following:

- `status_name` - (Required) The status name to check.
- `author_id` - (Optional) The authorized user can post the status.
- `invalidate_on_update` - (Optional) Reset status whenever there are new changes.
- `applicability` - (Optional) Policy applicability. If policy `applicability` is `default`, apply unless "Not Applicable" 
  status is posted to the pull request. If policy `applicability` is `conditional`, policy is applied only after a status 
  is posted to the pull request.
- `filename_patterns` - (Optional) If a path filter is set, the policy will only apply when files which match the filter are changes. Not setting this field means that the policy will always apply. You can specify absolute paths and wildcards. Example: `["/WebApp/Models/Data.cs", "/WebApp/*", "*.cs"]`. Paths prefixed with "!" are excluded. Example: `["/WebApp/*", "!/WebApp/Tests/*"]`. Order is significant.
- `display_name` - (Optional) The display name.
- `scope` (Required) Controls which repositories and branches the policy will be enabled for. This block must be defined
  at least once.

  `scope` block supports the following:

    - `repository_id` - (Optional) The repository ID. Needed only if the scope of the policy will be limited to a single
      repository.
    - `repository_ref` - (Optional) The ref pattern to use for the match. If `match_type` is `Exact`, this should be a
      qualified ref such as `refs/heads/master`. If `match_type` is `Prefix`, this should be a ref path such
      as `refs/heads/releases`.
    - `match_type` (Optional) The match type to use when applying the policy. Supported values are `Exact` (default)
      or `Prefix`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of branch policy configuration.

## Relevant Links

- [Azure DevOps Service REST API 5.1 - Policy Configurations](https://docs.microsoft.com/en-us/rest/api/azure/devops/policy/configurations/create?view=azure-devops-rest-5.1)

## Import

Azure DevOps Branch Policies can be imported using the project ID and policy configuration ID:

```sh
$ terraform import azuredevops_branch_policy_status_check.p 00000000-0000-0000-0000-000000000000/0
```
