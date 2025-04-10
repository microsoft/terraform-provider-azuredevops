---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_servicehook_storage_queue_pipelines"
description: |-
  Manages a Storage Queue Pipelines Service Hook.
---

# azuredevops_servicehook_storage_queue_pipelines

Manages a Storage Queue Pipelines Service Hook .

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
  name                 = "examplequeue"
  storage_account_name = azurerm_storage_account.example.name
}

resource "azuredevops_servicehook_storage_queue_pipelines" "example" {
  project_id   = azuredevops_project.example.id
  account_name = azurerm_storage_account.example.name
  account_key  = azurerm_storage_account.example.primary_access_key
  queue_name   = azurerm_storage_queue.example.name
  visi_timeout = 30
  run_state_changed_event {
    run_state_filter  = "Completed"
    run_result_filter = "Succeeded"
  }
}
```

An empty configuration block will occur in all events triggering the associated action.

```hcl
resource "azuredevops_servicehook_storage_queue_pipelines" "example" {
  project_id   = azuredevops_project.example.id
  account_name = azurerm_storage_account.example.name
  account_key  = azurerm_storage_account.example.primary_access_key
  queue_name   = azurerm_storage_queue.example.name
  visi_timeout = 30
  run_state_changed_event {}
}
```


## Arguments Reference

The following arguments are supported:

* `account_key` - (Required)  A valid account key from the queue's storage account.

* `account_name` - (Required) The queue's storage account name.

* `project_id` - (Required) The ID of the associated project. Changing this forces a new Service Hook Storage Queue Pipelines to be created.

* `queue_name` - (Required) The name of the queue that will store the events.

---

* `run_state_changed_event` - (Optional) A `run_state_changed_event` block as defined below. Conflicts with `stage_state_changed_event`

* `stage_state_changed_event` - (Optional) A `stage_state_changed_event` block as defined below. Conflicts with `run_state_changed_event`

-> **Note** At least one of `run_state_changed_event` and `stage_state_changed_event` has to be set.

* `ttl` - (Optional) event time-to-live - the duration a message can remain in the queue before it's automatically removed. Defaults to `604800`.

* `visi_timeout` - (Optional) event visibility timout - how long a message is invisible to other consumers after it's been dequeued. Defaults to `0`.

---

A `run_state_changed_event` block supports the following:

* `pipeline_id` - (Optional) The pipeline ID that will generate an event. If not specified, all pipelines in the project will trigger the event.

* `run_result_filter` - (Optional) Which run result should generate an event. Only valid if published_event is `RunStateChanged`. If not specified, all results will trigger the event.

* `run_state_filter` - (Optional) Which run state should generate an event. Only valid if published_event is `RunStateChanged`. If not specified, all states will trigger the event.

---

A `stage_state_changed_event` block supports the following:

* `pipeline_id` - (Optional) The pipeline ID that will generate an event.

* `stage_name` - (Optional) Which stage should generate an event. Only valid if published_event is `StageStateChanged`. If not specified, all stages will trigger the event.

* `stage_result_filter` - (Optional) Which stage result should generate an event. Only valid if published_event is `StageStateChanged`. If not specified, all results will trigger the event.

* `stage_state_filter` - (Optional) Which stage state should generate an event. Only valid if published_event is `StageStateChanged`. If not specified, all states will trigger the event.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Service Hook Storage Queue Pipelines.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Storage Queue Pipelines Service Hook.
* `read` - (Defaults to 5 minute) Used when retrieving the Storage Queue Pipelines Service Hook.
* `update` - (Defaults to 10 minutes) Used when updating the Storage Queue Pipelines Service Hook.
* `delete` - (Defaults to 10 minutes) Used when deleting the Storage Queue Pipelines Service Hook.

## Import

Storage Queue Pipelines Service Hook can be imported using the `resource id`, e.g.

```shell
terraform import azuredevops_servicehook_storage_queue_pipelines.example 00000000-0000-0000-0000-000000000000
```
