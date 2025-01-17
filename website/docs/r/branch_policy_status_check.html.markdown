---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_branch_policy_status_check"
description: |- Manages status check branch policy within Azure DevOps project.
---

# azuredevops_branch_policy_status_check

Manages a status check branch policy within Azure DevOps.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  features = {
    testplans = "disabled"
    artifacts = "disabled"
  }
  description = "Managed by Terraform"
}

resource "azuredevops_git_repository" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Repository"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_user_entitlement" "example" {
  principal_name       = "mail@email.com"
  account_license_type = "basic"
}

resource "azuredevops_branch_policy_status_check" "example" {
  project_id = azuredevops_project.example.id

  enabled  = true
  blocking = true

  settings {
    name                 = "Release"
    author_id            = azuredevops_user_entitlement.example.id
    invalidate_on_update = true
    applicability        = "conditional"
    display_name         = "PreCheck"

    scope {
      repository_id  = azuredevops_git_repository.example.id
      repository_ref = azuredevops_git_repository.example.default_branch
      match_type     = "Exact"
    }
    
    scope {
      match_type     = "DefaultBranch"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project in which the policy will be created.

* `settings` - (Required) A `settings` block as defined below. Configuration for the policy. This block must be defined exactly once. 

* `scope` (Required) Controls which repositories and branches the policy will be enabled for. This block must be defined
  at least once.

---

* `blocking` - (Optional) A flag indicating if the policy should be blocking. Defaults to `true`.

* `enabled` - (Optional) A flag indicating if the policy should be enabled. Defaults to `true`.

---

A `settings` block supports the following:

* `name` - (Required) The status name to check.

* `scope` - (Required) A `scope` block as defined below.

* `genre` - (Optional) The genre of the status to check (see [Microsoft Documentation](https://docs.microsoft.com/en-us/azure/devops/repos/git/pull-request-status?view=azure-devops#status-policy))

* `author_id` - (Optional) The authorized user can post the status.

* `invalidate_on_update` - (Optional) Reset status whenever there are new changes.

* `applicability` - (Optional) Policy applicability. If policy `applicability=default`, apply unless "Not Applicable"
  status is posted to the pull request. If policy `applicability=conditional`, policy is applied only after a status 
  is posted to the pull request. Possible values `default`, `conditional`. Defaults to `default`.

* `filename_patterns` - (Optional) If a path filter is set, the policy will only apply when files which match the filter are changed. Not setting this field means that the policy is always applied.
  
  ~>**NOTE** 1. Specify absolute paths and wildcards. Example: `["/WebApp/Models/Data.cs", "/WebApp/*", "*.cs"]`. 
  <br> 2. Paths prefixed with "!" are excluded. Example: `["/WebApp/*", "!/WebApp/Tests/*"]`. Order is significant.

* `display_name` - (Optional) The display name.

---

A `scope` block supports the following:

* `repository_id` - (Optional) The repository ID. Needed only if the scope of the policy will be limited to a single repository. If `match_type=DefaultBranch`, this should not be defined.

* `repository_ref` - (Optional) The ref pattern to use for the match when `match_type` other than `DefaultBranch`. If `match_type=Exact`, this should be a qualified ref such as `refs/heads/master`. If `match_type=Prefix`, this should be a ref path such as `refs/heads/releases`.

* `match_type` (Optional) The match type to use when applying the policy. Supported values are `Exact` (default), `Prefix` or `DefaultBranch`.
    

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of branch policy configuration.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Policy Configurations](https://docs.microsoft.com/en-us/rest/api/azure/devops/policy/configurations/create?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Status Check Branch Policy.
* `read` - (Defaults to 2 minute) Used when retrieving the Status Check Branch Policy.
* `update` - (Defaults to 5 minutes) Used when updating the Status Check Branch Policy.
* `delete` - (Defaults to 5 minutes) Used when deleting the Status Check Branch Policy.

## Import

Azure DevOps Branch Policies can be imported using the project ID and policy configuration ID:

```sh
terraform import azuredevops_branch_policy_status_check.example 00000000-0000-0000-0000-000000000000/0
```
