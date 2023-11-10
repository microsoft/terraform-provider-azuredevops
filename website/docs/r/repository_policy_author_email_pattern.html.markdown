---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_repository_policy_author_email_pattern"
description: |- Manages author email pattern repository policy within Azure DevOps project.
---

# azuredevops_repository_policy_author_email_pattern

Manage author email pattern repository policy within Azure DevOps project.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_git_repository" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Repository"
  initialization {
    init_type = "Clean"
  }
}

resource "azuredevops_repository_policy_author_email_pattern" "example" {
  project_id            = azuredevops_project.example.id
  enabled               = true
  blocking              = true
  author_email_patterns = ["user1@test.com", "user2@test.com"]
  repository_ids        = [azuredevops_git_repository.example.id]
}
```

## Set project level repository policy
```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_repository_policy_author_email_pattern" "example" {
  project_id            = azuredevops_project.example.id
  enabled               = true
  blocking              = true
  author_email_patterns = ["user1@test.com", "user2@test.com"]
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The ID of the project in which the policy will be created.
- `enabled` - (Optional) A flag indicating if the policy should be enabled. Defaults to `true`.
- `blocking` - (Optional) A flag indicating if the policy should be blocking. Defaults to `true`.
- `author_email_patterns` - (Required) Block pushes with a commit author email that does not match the patterns. You can specify exact emails or use wildcards. 
  Email patterns prefixed with "!" are excluded. Order is important.
- `repository_ids` (Optional) Control whether the policy is enabled for the repository or the project. If `repository_ids` not configured, the policy will be set to the project.   
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of repository policy configuration.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Policy Configurations](https://docs.microsoft.com/en-us/rest/api/azure/devops/policy/configurations?view=azure-devops-rest-7.0)

## Import

Azure DevOps Branch Policies can be imported using the project ID and policy configuration ID:

```sh
terraform import azuredevops_repository_policy_author_email_pattern.example 00000000-0000-0000-0000-000000000000/0
```
