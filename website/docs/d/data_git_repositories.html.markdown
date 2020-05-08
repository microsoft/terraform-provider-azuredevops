---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_git_repositories"
description: |-
  Use this data source to access information about an existing Projects within Azure DevOps.
---

# Data Source: azuredevops_git_repositories

Use this data source to access information about an existing Projects within Azure DevOps.

## Example Usage

```hcl

# Make sure to set the following environment variables:
#   AZDO_PERSONAL_ACCESS_TOKEN
#   AZDO_ORG_SERVICE_URL
provider "azuredevops" {
  version = ">= 0.0.1"
}

# Load all projects of an organization,
# that are accessible by the current user
data "azuredevops_projects" "tf-projects" {
}

# Build a local map, to access projects by name
locals {
  project_map = {
    for project in data.azuredevops_projects.tf-projects.projects : project["name"] => project
  }
}

# Load all Git repositories of an organization,
# which are accessible for the current user
data "azuredevops_git_repositories" "tf-git-repos-all" {
}

output "out-tf-git-repos-all" {
  value = data.azuredevops_git_repositories.tf-git-repos-all.repositories
}

# Build a local map, to access Git repositories by name
locals {
  repo_map = {
    for repo in data.azuredevops_git_repositories.tf-git-repos-all.repositories : repo["name"] => repo
  }
}

# Load all Git repositories of a project,
# which are accessible for the current user
data "azuredevops_git_repositories" "tf-git-repos-project" {
  project_id = local.project_map[var.project_name].project_id
}

# Load a specific Git repository by name
data "azuredevops_git_repositories" "tf-git-repos-project-reponame" {
  project_id = local.project_map[var.project_name].project_id
  name       = var.git_repo_name
}

```

## Argument Reference

The following arguments are supported:

- `project_id` - (Optional) ID of project to list Git repositories
- `name` - (Optional) Name of the Git repository to retrieve; requires `project_id` to be specified as well
- `include_hidden` - (Optional, default: false)

DataSource without specifying any arguments will return all Git repositories of an organization.

## Attributes Reference

The following attributes are exported:

- `repositories` - A list of existing projects in your Azure DevOps Organization with details about every project which includes:

  - `id` - Git repository identifier.
  - `name` - Git repository name.
  - `url` - Details REST API endpoint for the Git Repository.
  - `ssh_url` - SSH Url to clone the Git repository
  - `web_url` - Url of the Git repository web view
  - `remote_url` - HTTPS Url to clone the Git repository
  - `project_id` - Project identifier to which the Git repository belongs.
  - `size` - Compressed size (bytes) of the repository.
  - `default_branch` - The name of the default branch.

## Relevant Links

- [Azure DevOps Service REST API 5.1 - Git API](https://docs.microsoft.com/en-us/rest/api/azure/devops/git/?view=azure-devops-rest-5.1)
