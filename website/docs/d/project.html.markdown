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

* `name` - (Required) Name of the Project.

* `project_id` - (Required) ID of the Project.

~> **NOTE:** One of either `project_id` or `name` must be specified.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the project.

* `name` - The name of the project.

* `description` - The description of the project.

* `visibility` - The visibility of the project.

* `version_control` - The version control of the project.

* `work_item_template` - The work item template for the project.

* `process_template_id` - The process template ID for the project.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Projects - Get](https://docs.microsoft.com/en-us/rest/api/azure/devops/core/projects/get?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minute) Used when retrieving the Project.

## PAT Permissions Required

- **Project & Team**: Read
- **Work Items**: Read