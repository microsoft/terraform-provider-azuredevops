---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_servicehook_storage_queue"
description: |-
  Manages a Service Hook Storage Queue.
---

# azuredevops_servicehook_storage_queue

Manages a Service Hook Storage Queue.

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name = "example-project" 
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe" 
}

resource "azurerm_storage_account" "example" {
  name                     = "servicehookexamplestacc"
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_queue" "example" {
  name                  = "examplequeue"
  storage_account_name  = azurerm_storage_account.example.name
}

resource "azuredevops_servicehook_storage_queue" "example" {
  project_id   = azuredevops_project.example.id
  account_name = azurerm_storage_account.example.name
  account_key  = azurerm_storage_account.example.primary_access_key 
  queue_name   = azurerm_storage_queue.example.name
  visi_timeout = 30
  publisher {
    name = "pipelines"
    stage_state_changed {
      state_filter  = "Completed"
      result_filter = "Succeeded"
    }
  }
}
```

## Arguments Reference

The following arguments are supported:

* `account_key` - (Required) A valid account key from the queue's storage account.

* `account_name` - (Required) The queue's storage account name.

* `project_id` - (Required) The ID of the associated project. Changing this forces a new Service Hook Storage Queue to be created.

* `publisher` - (Required) A `publisher` block as defined below.

* `queue_name` - (Required) The name of the queue that will store the events.

---

* `ttl` - (Optional) event time-to-live - the duration a message can remain in the queue before it's automatically removed.

* `visi_timeout` - (Optional) event visibility timout - how long a message is invisible to other consumers after it's been dequeued.

---

A `publisher` block supports the following:

* `name` - (Required) The name of the publisher.

* `run_state_changed` - (Optional) A `run_state_changed` block as defined below.

* `stage_state_changed` - (Optional) A `stage_state_changed` block as defined below.

---

A `run_state_changed` block supports the following:

* `pipeline_id` - (Optional) The ID of the pipeline that will generate the event.

* `result_filter` - (Optional) Which result should generate an event.

* `state_filter` - (Optional) Which final state should generate an event.

---

A `stage_state_changed` block supports the following:

* `pipeline_id` - (Optional) The ID of the pipeline that will generate the event.

* `result_filter` - (Optional) Which result should generate an event.

* `stage_name` - (Optional) The name of the stage that, in case of state change, will generate an event.

* `state_filter` - (Optional) Which final state should generate an event.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Service Hook Storage Queue.



## Import

Service Hook Storage Queues can be imported using the `resource id`, e.g.

```shell
terraform import azuredevops_servicehook_storage_queue.example 00000000-0000-0000-0000-000000000000
```
