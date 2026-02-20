---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_deployment_group"
description: |-
  Manages a Deployment Group.
---

# azuredevops_deployment_group

Manages a Deployment Group used by classic release pipelines.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name               = "Example Project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

resource "azuredevops_deployment_group" "example" {
  project_id  = azuredevops_project.example.id
  name        = "Example Deployment Group"
  description = "Managed by Terraform"
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name of the Deployment Group.

* `project_id` - (Required) The ID of the project. Changing this forces a new Deployment Group to be created.

---

* `description` - (Optional) A description for the Deployment Group. Defaults to `""`.

* `pool_id` - (Optional) The ID of the deployment pool in which deployment agents are registered. If not specified, a new pool will be created. Changing this forces a new Deployment Group to be created.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Deployment Group.

* `machine_count` - The number of deployment targets in the Deployment Group.

## Relevant Links

* [Azure DevOps Service REST API 7.0 - Deployment Groups](https://docs.microsoft.com/en-us/rest/api/azure/devops/distributedtask/deployment-groups?view=azure-devops-rest-7.0)

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Deployment Group.
* `read` - (Defaults to 5 minute) Used when retrieving the Deployment Group.
* `update` - (Defaults to 10 minutes) Used when updating the Deployment Group.
* `delete` - (Defaults to 10 minutes) Used when deleting the Deployment Group.

## Import

Azure DevOps Deployment Groups can be imported using the project ID and deployment group ID, e.g.:

```sh
terraform import azuredevops_deployment_group.example 00000000-0000-0000-0000-000000000000/0
```
