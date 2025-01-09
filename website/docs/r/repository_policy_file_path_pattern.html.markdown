---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_repository_policy_file_path_pattern"
description: |- Manages a file path pattern repository policy within Azure DevOps project.
---

# azuredevops_repository_policy_file_path_pattern

Manage a file path pattern repository policy within Azure DevOps project.

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

resource "azuredevops_repository_policy_file_path_pattern" "example" {
  project_id = azuredevops_project.example.id

  enabled           = true
  blocking          = true
  filepath_patterns = ["*.go", "/home/test/*.ts"]
  repository_ids    = [azuredevops_git_repository.example.id]
}
```

# Set project level repository policy
```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_repository_policy_file_path_pattern" "examplep" {
  project_id        = azuredevops_project.example.id
  enabled           = true
  blocking          = true
  filepath_patterns = ["*.go", "/home/test/*.ts"]
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The ID of the project in which the policy will be created.
- `enabled` - (Optional) A flag indicating if the policy should be enabled. Defaults to `true`.
- `blocking` - (Optional) A flag indicating if the policy should be blocking. Defaults to `true`.
- `filepath_patterns` - (Required) Block pushes from introducing file paths that match the following patterns. Exact paths begin with "/". You can specify exact paths and wildcards. You can also specify multiple paths using ";" as a separator. Paths prefixed with "!" are excluded. Order is important.
- `repository_ids` (Optional) Control whether the policy is enabled for the repository or the project. If `repository_ids` not configured, the policy will be set to the project.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the repository policy.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Policy Configurations](https://docs.microsoft.com/en-us/rest/api/azure/devops/policy/configurations?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the File Path Pattern Repository Policy.
* `read` - (Defaults to 5 minute) Used when retrieving the File Path Pattern Repository Policy.
* `update` - (Defaults to 30 minutes) Used when updating the File Path Pattern Repository Policy.
* `delete` - (Defaults to 30 minutes) Used when deleting the File Path Pattern Repository Policy.

## Import

Azure DevOps repository policies can be imported using the projectID/policyID or projectName/policyID:

```sh
terraform import azuredevops_repository_policy_file_path_pattern.example 00000000-0000-0000-0000-000000000000/0
```
