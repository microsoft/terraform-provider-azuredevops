---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_repository_policy_author_email_pattern"
description: |- Manages author email pattern repository policy within Azure DevOps project.
---

# azuredevops_repository_policy_author_email_pattern

Manage author email pattern repository policy within Azure DevOps project.

## Example Usage

```hcl
resource "azuredevops_project" "p" {
  name               = "Sample Project"
  description        = "Managed by Terraform"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
}

resource "azuredevops_git_repository" "r" {
  project_id = azuredevops_project.p.id
  name       = "Sample Repo"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_repository_policy_author_email_pattern" "p" {
  project_id = azuredevops_project.p.id

  enabled  = true
  blocking = true

  settings {
    author_email_patterns = ["user1@test.com", "user2@test.com"]
    scope {
      repository_id = azuredevops_git_repository.r.id
    }
  }
}
```

## Set project level repository policy
```hcl
resource "azuredevops_repository_policy_author_email_pattern" "p" {
  project_id = azuredevops_project.p.id

  enabled  = true
  blocking = true

  settings {
    author_email_patterns = ["user1@test.com", "user2@test.com"]
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

- `author_email_patterns` - (Required) Block pushes with a commit author email that does not match the patterns. You can specify exact emails or use wildcards. 
  Email patterns prefixed with "!" are excluded. Order is important.
- `scope` (Optional) Control whether the policy is enabled for the repository or the project. If `scope` not configured, the policy will be set to the project.
  
  `scope` block supports the following:
    - `repository_id` - (Required) The repository ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of repository policy configuration.

## Relevant Links

- [Azure DevOps Service REST API 5.1 - Policy Configurations](https://docs.microsoft.com/en-us/rest/api/azure/devops/policy/configurations/create?view=azure-devops-rest-5.1)

## Import

Azure DevOps Branch Policies can be imported using the project ID and policy configuration ID:

```sh
$ terraform import azuredevops_repository_policy_author_email_pattern.p 00000000-0000-0000-0000-000000000000/0
```
