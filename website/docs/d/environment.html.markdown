---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_environment"
description: |-
  Use this data source to access information about an Environment.
---

# Data Source: azuredevops_environment

Use this data source to access information about an Environment.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

resource "azuredevops_environment" "example" {
  project_id = azuredevops_project.example.id
  name       = "Example Environment"
  description = "Managed by Terraform"
}

data "azuredevops_environment" "example" {
  project_id= azuredevops_project.example.id
  environment_id = azuredevops_environment.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project.

* `environment_id` - (Required) The ID of the Environment.

## Attributes Reference

In addition to the Arguments list above - the following Attributes are exported:

* `name` - The name of the Environment.

* `description` - A description for the Environment.

## Relevant Links

* [Azure DevOps Service REST API 7.0 - Environments](https://docs.microsoft.com/en-us/rest/api/azure/devops/distributedtask/environments?view=azure-devops-rest-7.0)

