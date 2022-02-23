---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_repository_policy_reserved_names"
description: |- Manage a reserved names repository policy within Azure DevOps project.
---

# azuredevops_repository_policy_reserved_names

Manage a reserved names repository policy within Azure DevOps project. Block pushes that introduce files, folders, or branch names that include platform reserved names or incompatible characters.

~> If both project and project policy are enabled, the project policy has high priority.

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

resource "azuredevops_repository_policy_reserved_names" "example" {
  project_id     = azuredevops_project.example.id
  enabled        = true
  blocking       = true
  repository_ids = [azuredevops_git_repository.example.id]
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

resource "azuredevops_repository_policy_reserved_names" "example" {
  project_id = azuredevops_project.example.id
  enabled    = true
  blocking   = true
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The ID of the project in which the policy will be created.
- `enabled` - (Optional) A flag indicating if the policy should be enabled. Defaults to `true`. 
- `blocking` - (Optional) A flag indicating if the policy should be blocking. Defaults to `true`.
- `repository_ids` (Optional) Control whether the policy is enabled for the repository or the project. If `repository_ids` not configured, the policy will be set to the project.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the repository policy.

## Relevant Links

- [Azure DevOps Service REST API 6.0 - Policy Configurations](https://docs.microsoft.com/en-us/rest/api/azure/devops/policy/configurations?view=azure-devops-rest-6.0)

## Import

Azure DevOps repository policies can be imported using the projectID/policyID or projectName/policyID:

```sh
terraform import azuredevops_repository_policy_reserved_names.p 00000000-0000-0000-0000-000000000000/0
```
