---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_agent_queue"
description: |-
  Manages an agent queue within Azure DevOps project.
---

# azuredevops_agent_queue

Manages an agent queue within Azure DevOps. In the UI, this is equivalent to adding an
Organization defined pool to a project.

The created queue is not authorized for use by all pipelines in the project. However,
the `azuredevops_resource_authorization` resource can be used to grant authorization.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name = "Example Project"
}

data "azuredevops_agent_pool" "example" {
  name = "example-pool"
}

resource "azuredevops_agent_queue" "example" {
  project_id    = azuredevops_project.example.id
  agent_pool_id = data.azuredevops_agent_pool.example.id
}

# Grant access to queue to all pipelines in the project
resource "azuredevops_resource_authorization" "example" {
  project_id  = azuredevops_project.example.id
  resource_id = azuredevops_agent_queue.example.id
  type        = "queue"
  authorized  = true
}
```

## Argument Reference

The following arguments are supported:

- `project_id` - (Required) The ID of the project in which to create the resource.
- `agent_pool_id` - (Required) The ID of the organization agent pool.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the agent queue reference.

## Relevant Links

- [Azure DevOps Service REST API 7.0 - Agent Queues](https://docs.microsoft.com/en-us/rest/api/azure/devops/distributedtask/queues?view=azure-devops-rest-7.0)

## Import

Azure DevOps Agent Pools can be imported using the project ID and agent queue ID, e.g.

```sh
terraform import azuredevops_agent_queue.example 00000000-0000-0000-0000-000000000000/0
```
