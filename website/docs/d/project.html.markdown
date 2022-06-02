---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_project"
description: |-
  Use this data source to access information about an existing Project within Azure DevOps.
---

# Data Source: azuredevops_project

Use this data source to access information about an existing Project within Azure DevOps.

## Example Usage

```hcl
data "azuredevops_project" "example" {
  name = "Example Project"
}

output "project" {
  value = data.azuredevops_project.example
}
```

## Argument Reference

The following arguments are supported:

- `name` - (Required if `project_id` not set) Name of the Project.
- `project_id` - (Required if `name` not set) ID of the Project.

## Attributes Reference

The following attributes are exported:

`name` - The name of the referenced project
`description` - The description of the referenced project
`visibility` - The visibility of the referenced project
`version_control` - The version control of the referenced project
`work_item_template` - The work item template for the referenced project
`process_template_id` - The process template ID for the referenced project

## Relevant Links

- [Azure DevOps Service REST API 6.0 - Projects - Get](https://docs.microsoft.com/en-us/rest/api/azure/devops/core/projects/get?view=azure-devops-rest-6.0)
