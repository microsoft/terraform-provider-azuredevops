---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_git_repository"
description: |-
  Use this data source to access information about an existing Git Repository within Azure DevOps.
---

# Data Source: azuredevops_git_repository

Use this data source to access information about a **single** (existing) Git Repository within Azure DevOps.
To read informations about **multiple** Git Repositories use the data source [`azuredevops_git_repositories`](data_git_repositories.html)

## Example Usage

```hcl
# Load all projects of an organization, that are accessible by the current user
data "azuredevops_project" "project" {
  project_name = "contoso-project"
}

# Load a specific Git repository by name
data "azuredevops_git_repository" "single_repo" {
  project_id = data.azuredevops_project.project.id
  name       = "contoso-repo"
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

## Relevant Links

- [Azure DevOps Service REST API 5.1 - Git API](https://docs.microsoft.com/en-us/rest/api/azure/devops/git/?view=azure-devops-rest-5.1)
