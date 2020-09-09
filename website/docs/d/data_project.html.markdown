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
data "azuredevops_project" "p" {
  project_identifier = "Sample Project"
}

output "id" {
  value = data.azuredevops_project.p.id
}

output "project_name" {
  value = data.azuredevops_project.p.project_name
}

output "visibility" {
  value = data.azuredevops_project.p.visibility
}

output "version_control" {
  value = data.azuredevops_project.p.version_control
}

output "work_item_template" {
  value = data.azuredevops_project.p.work_item_template
}

output "process_template_id" {
  value = data.azuredevops_project.p.process_template_id
}
```

## Argument Reference

The following arguments are supported:

- `project_identifier` - (Required) Name or ID of the Project.

## Attributes Reference

The following attributes are exported:

`project_name` - The name of the referenced project
`description` - The description of the referenced project
`visibility` - The visibility of the referenced project
`version_control` - The version control of the referenced project
`work_item_template` - The work item template for the referenced project
`process_template_id` - The process template ID for the referenced project

## Relevant Links

- [Azure DevOps Service REST API 5.1 - Projects - Get](https://docs.microsoft.com/en-us/rest/api/azure/devops/core/projects/get?view=azure-devops-rest-5.1)
