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

* `environment_id` - (Optional) The ID of the Environment.

* `name` - (Optional) Name of the Environment.

~> **NOTE:** One of either `environment_id` or `name` must be specified.

## Attributes Reference

In addition to the Arguments list above - the following Attributes are exported:

* `name` - The name of the Environment.

* `description` - A description for the Environment.

## Relevant Links

* [Azure DevOps Service REST API 7.0 - Environments](https://docs.microsoft.com/en-us/rest/api/azure/devops/distributedtask/environments?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `read` - (Defaults to 5 minute) Used when retrieving the Environment.