---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_repository_policy_searchable_branches"
description: |- Manages searchable branches repository policy within Azure DevOps project.
---

# azuredevops_repository_policy_searchable_branches

Manage searchable branches repository policy within Azure DevOps project.

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

resource "azuredevops_repository_policy_searchable_branches" "example" {
  project_id          = data.azuredevops_project.example.id
  searchable_branches = ["examplebranch"]
  repository_ids      = [data.azuredevops_git_repository.example.id]

}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The ID of the project in which the policy will be created.
- `enabled` - (Computed) This is set to false by the provider and is not used.
- `blocking` - (Computed) This is set to false by the provider and is not used.
- `searchable_branches` - (Required) A list of branch names to be added as searchable branches for the repository. Branches do not have to exist.
- `repository_ids` (Required) ID of repository for which the policy is enabled. Note: Due to the API implementation of this policy, this only accepts 1 repository id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the repository policy.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Policy Configurations](https://docs.microsoft.com/en-us/rest/api/azure/devops/policy/configurations?view=azure-devops-rest-7.0)

## Import

Azure DevOps repository policies can be imported using the projectID/policyID or projectName/policyID:

```sh
terraform import azuredevops_repository_policy_searchable_branches.example 00000000-0000-0000-0000-000000000000/0
```
