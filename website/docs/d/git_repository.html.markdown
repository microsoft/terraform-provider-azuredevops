---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_git_repository"
description: |-
  Use this data source to access information about an existing Git Repository within Azure DevOps.
---

# Data Source: azuredevops_git_repository

Use this data source to access information about a **single** (existing) Git Repository within Azure DevOps.
To read information about **multiple** Git Repositories use the data source [`azuredevops_git_repositories`](data_git_repositories.html)

## Example Usage

```hcl
data "azuredevops_project" "example" {
  name = "Example Project"
}

# Load a specific Git repository by name
data "azuredevops_git_repository" "example-single-repo" {
  project_id = data.azuredevops_project.example.id
  name       = "Example Repository"
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) ID of project to list Git repositories
- `name` - (Required) Name of the Git repository to retrieve

## Attributes Reference

The following attributes are exported:

- `id` - Git repository identifier.
- `name` - Git repository name.
- `url` - Details REST API endpoint for the Git Repository.
- `ssh_url` - SSH Url to clone the Git repository
- `web_url` - Url of the Git repository web view
- `remote_url` - HTTPS Url to clone the Git repository
- `project_id` - Project identifier to which the Git repository belongs.
- `size` - Compressed size (bytes) of the repository.
- `default_branch` - The ref of the default branch.
- `disabled` - Is the repository disabled?

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Git API](https://docs.microsoft.com/en-us/rest/api/azure/devops/git/?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 30 minute) Used when retrieving the Git Repository.