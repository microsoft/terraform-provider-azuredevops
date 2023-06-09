---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_projects"
description: |-
  Use this data source to access information about a existing Projects within Azure DevOps.
---

# Data Source: azuredevops_projects

Use this data source to access information about existing Projects within Azure DevOps.

## Example Usage

```hcl
data "azuredevops_projects" "example" {
  name  = "Example Project"
  state = "wellFormed"
}

output "project_id" {
  value = data.azuredevops_projects.example.projects.*.project_id
}

output "name" {
  value = data.azuredevops_projects.example.projects.*.name
}

output "project_url" {
  value = data.azuredevops_projects.example.projects.*.project_url
}

output "state" {
  value = data.azuredevops_projects.example.projects.*.state
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Optional) Name of the Project, if not specified all projects will be returned.

- `state` - (Optional) State of the Project, if not specified all projects will be returned. Valid values are `all`, `deleting`, `new`, `wellFormed`, `createPending`, `unchanged`,`deleted`.

DataSource without specifying any arguments will return all projects.

## Attributes Reference

The following attributes are exported:

- `projects` - A list of existing projects in your Azure DevOps Organization with details about every project which includes:

  - `project_id` - The ID of the Project.
  - `name` - The name of the Project.
  - `project_url` - Url to the full version of the object.
  - `state` - Project state.

## Relevant Links

- [Azure DevOps Service REST API 6.0 - Projects - Get](https://docs.microsoft.com/en-us/rest/api/azure/devops/core/projects/get?view=azure-devops-rest-6.0)
