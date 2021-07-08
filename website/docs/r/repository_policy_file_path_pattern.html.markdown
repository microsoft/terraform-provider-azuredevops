---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_repository_policy_file_path_pattern"
description: |- Manages a file path pattern repository policy within Azure DevOps project.
---

# azuredevops_repository_policy_file_path_pattern

Manage a file path pattern repository policy within Azure DevOps project.

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

resource "azuredevops_repository_policy_file_path_pattern" "p" {
  project_id = azuredevops_project.p.id

  enabled  = true
  blocking = true

  settings {
    filepath_patterns = ["*.go", "/home/test/*.ts"]
    scope {
      repository_id = azuredevops_git_repository.r.id
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

- `filepath_patterns` - (Required) Block pushes from introducing file paths that match the following patterns. Exact paths begin with "/". You can specify exact paths and wildcards. You can also specify multiple paths using ";" as a separator. Paths prefixed with "!" are excluded. Order is important.
- `scope` (Optional) Controls which repositories and branches the policy will be enabled for. This block must be defined
  at least once.   
  
  `scope` block supports the following:
    - `repository_id` - (Required) The repository ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the repository policy.

## Relevant Links

- [Azure DevOps Service REST API 5.1 - Policy Configurations](https://docs.microsoft.com/en-us/rest/api/azure/devops/policy/configurations/create?view=azure-devops-rest-5.1)

## Import

Azure DevOps repository policies can be imported using the projectID/policyID or projectName/policyID:

```sh
$ terraform import azuredevops_repository_policy_file_path_pattern.p 00000000-0000-0000-0000-000000000000/0
```
