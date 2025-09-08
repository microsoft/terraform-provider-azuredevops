---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_project_features"
description: |-
  Use this data source to access information about the features of an existing Project within Azure DevOps.
---

# Data Source: azuredevops_project

Use this data source to access information about the features of an existing Project within Azure DevOps.

## Example Usage

```hcl
data "azuredevops_project" "example" {
  name = "Example Project"
}

data "azuredevops_project_features" "example" {
  project_id = data.azuredevops_project.example.id
}

output "project" {
  value = data.azuredevops_project_features.features
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) ID of the Project.

## Attributes Reference

The following attributes are exported:

* `project_id` - The ID of the project.

* `features` - The features for projects.

## Relevant Links

- [Azure DevOps Service REST API 7.1 - Queries - Get](https://learn.microsoft.com/en-us/rest/api/azure/devops/wit/queries/get?view=azure-devops-rest-7.1&tabs=HTTP)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minute) Used when retrieving the Project.

## PAT Permissions Required

- **Project & Team**: Read
