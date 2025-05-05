---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_branch_policy_comment_resolution"
description: |-
  Configure a comment resolution policy for your branch within Azure DevOps project.
---

# azuredevops_branch_policy_comment_resolution

Configure a comment resolution policy for your branch within Azure DevOps project.

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

resource "azuredevops_branch_policy_comment_resolution" "example" {
  project_id = azuredevops_project.example.id

  enabled  = true
  blocking = true

  settings {

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

* `project_id` - (Required) The ID of the project in which the policy will be created.

* `settings` - (Required) A `settings` block as defined below. Configuration for the policy. This block must be defined exactly once.

---

* `enabled` - (Optional) A flag indicating if the policy should be enabled. Defaults to `true`.

* `blocking` - (Optional) A flag indicating if the policy should be blocking. Defaults to `true`.

---

A `settings` block supports the following:

* `scope` (Required) A `scope` block as defined below. Controls which repositories and branches the policy will be enabled for. This block must be defined at least once.

---

A `scope` block supports the following:

* `repository_id` - (Optional) The repository ID. Needed only if the scope of the policy will be limited to a single repository. If `match_type` is `DefaultBranch`, this should not be defined.

* `repository_ref` - (Optional) The ref pattern to use for the match when `match_type` other than `DefaultBranch`. If `match_type` is `Exact`, this should be a qualified ref such as `refs/heads/master`. If `match_type` is `Prefix`, this should be a ref path such as `refs/heads/releases`.

* `match_type` (Optional) The match type to use when applying the policy. Supported values are `Exact` (default), `Prefix` or `DefaultBranch`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of branch policy configuration.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Policy Configurations](https://docs.microsoft.com/en-us/rest/api/azure/devops/policy/configurations/create?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Comment Resolution Branch Policy.
* `read` - (Defaults to 2 minute) Used when retrieving the Comment Resolution Branch Policy.
* `update` - (Defaults to 5 minutes) Used when updating the Comment Resolution Branch Policy.
* `delete` - (Defaults to 5 minutes) Used when deleting the Comment Resolution Branch Policy.

## Import

Azure DevOps Branch Policies can be imported using the project ID and policy configuration ID:

```sh
terraform import azuredevops_branch_policy_comment_resolution.example 00000000-0000-0000-0000-000000000000/0
```
