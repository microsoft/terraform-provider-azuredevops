---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_project_tags"
description: |-
  Manages Project Tags within Azure DevOps organization.
---

# azuredevops_project_tags

Manages Project Tags within Azure DevOps organization.

## Example Usage
```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  visibility         = "private"
  version_control    = "Git"
  work_item_template = "Agile"
  description        = "Managed by Terraform"
}

resource "azuredevops_project_tags" "example" {
  project_id = azuredevops_project.example.id
  tags       = ["tag1", "tag2"]
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Optional) The ID of the Project. Changing this forces a new resource to be created.

* `tags` - A mapping of tags assigned to the Project.

## Attributes Reference

The following attributes are exported:

* `project_id` - The ID of the Project.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Project Properties](https://learn.microsoft.com/en-us/rest/api/azure/devops/core/projects/get-project-properties?view=azure-devops-rest-7.1&tabs=HTTP)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the Project Tags.
* `read` - (Defaults to 2 minute) Used when retrieving the Project Tags.
* `update` - (Defaults to 5 minutes) Used when updating the Project Tags.
* `delete` - (Defaults to 5 minutes) Used when deleting the Project Tags.

## Import

Azure DevOps Project Tags can be imported using the Project ID e.g.:

```sh
terraform import azuredevops_project_tags.example 00000000-0000-0000-0000-000000000000
```