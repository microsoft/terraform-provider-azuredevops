---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_project"
description: |-
  Use this data source to access information about an existing Projects within Azure DevOps.
---

# Data Source: azuredevops_project

Use this data source to access information about an existing Projects within Azure DevOps.

## Example Usage

```hcl

data "azuredevops_projects" "test" {
    project_name = "contoso"
    state       = "wellFormed"
}

output "project_id" {
  value = "${data.azuredevops_projects.test.projects.*.project_id}"
}

output "project_name" {
  value = "${data.azuredevops_projects.test.projects.*.name}"
}

output "project_url" {
  value = "${data.azuredevops_projects.test.projects.*.project_url}"
}

output "state" {
  value = "${data.azuredevops_projects.test.projects.*.state}"
}


output "project_url" {

  value = { for project in data.azuredevops_projects.test.projects :
    project.name => project.project_url
  }

}
```

## Argument Reference

The following arguments are supported:

- `project_name` - (Optional) Name of the Project, if not specified all projects will be returned.

- `state` - (Optional) State of the Project, if not specified all projects will be returned. Valid values are `all`, `deleting`, `new`, `wellFormed`, `createPending`, `unchanged`,`deleted`.

DataSource without specifying any arguments will return all projects.

## Attributes Reference

The following attributes are exported:

- `projects` - A list of existing projects in your Azure DevOps Organization with details about every project which includes:

  - `id` - Project identifier.
  - `name` - Project name.
  - `project_url` - Url to the full version of the object.
  - `state` - Project state.

## Relevant Links

- [Azure DevOps Service REST API 5.1 - Projects - Get](https://docs.microsoft.com/en-us/rest/api/azure/devops/core/projects/get?view=azure-devops-rest-5.1)
