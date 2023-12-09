---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_environment"
description: |-
  Manages an Environment.
---

# azuredevops_environment

Manages an Environment.

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
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Environment.

* `project_id` - (Required) The ID of the project. Changing this forces a new Environment to be created.

---

* `description` - (Optional) A description for the Environment.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Environment.

## Relevant Links

* [Azure DevOps Service REST API 7.0 - Environments](https://docs.microsoft.com/en-us/rest/api/azure/devops/distributedtask/environments?view=azure-devops-rest-7.0)


## Import

Azure DevOps Environments can be imported using the project ID and environment ID, e.g.:

```sh
terraform import azuredevops_environment.example 00000000-0000-0000-0000-000000000000/0
```
