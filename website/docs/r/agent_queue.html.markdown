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
resource "azuredevops_project" "project" {
  project_name = "Sample Project"
}

data "azuredevops_agent_pool" "pool" {
  name = "contoso-pool"
}

resource "azuredevops_agent_queue" "queue" {
  project_id    = azuredevops_project.project.id
  agent_pool_id = data.azuredevops_agent_pool.pool.id
}

# Grant acccess to queue to all pipelines in the project
resource "azuredevops_resource_authorization" "auth" {
  project_id  = azuredevops_project.project.id
  resource_id = azuredevops_agent_queue.queue.id
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

- [Azure DevOps Service REST API 5.1 - Agent Queues](https://docs.microsoft.com/en-us/rest/api/azure/devops/distributedtask/queues?view=azure-devops-rest-5.1)

## Import

Azure DevOps Agent Pools can be imported using the project ID and agent queue ID, e.g.

```sh
terraform import azuredevops_agent_queue.q 44cbf614-4dfd-4032-9fae-87b0da3bec30/1381
```
