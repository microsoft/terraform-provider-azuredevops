---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_servicehook_storage_queue_pipelines"
description: |-
  Manages a Service Hook Storage Queue Pipelines.
---

# azuredevops_servicehook_storage_queue_pipelines

Manages a Service Hook Storage Queue Pipelines.

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

resource "azuredevops_servicehook_storage_queue_pipelines" "example" {
  project_id   = azuredevops_project.example.id
  account_name = azurerm_storage_account.example.name
  account_key  = azurerm_storage_account.example.primary_access_key 
  queue_name   = azurerm_storage_queue.example.name
  visi_timeout = 30
  published_event = "RunStateChanged"
  event_config {
    state_filter = "Completed"
    result_filter = "Succeeded"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `account_key` - (Required) A valid account key from the queue's storage account.

* `account_name` - (Required) The queue's storage account name.

* `project_id` - (Required) The ID of the associated project. Changing this forces a new Service Hook Storage Queue Pipelines to be created.

* `published_event` - (Required) The trigger event. Possible options are `RunStateChanged`, and `StageStateChanged`.

* `queue_name` - (Required) The name of the queue that will store the events.

---

* `event_config` - (Optional) A `event_config` block as defined below.

* `ttl` - (Optional) TODO.

* `visi_timeout` - (Optional) TODO.

---

A `event_config` block supports the following:

* `pipeline_id` - (Optional) The ID of the TODO.

* `result_filter` - (Optional) TODO.

* `stage_name` - (Optional) TODO.

* `state_filter` - (Optional) TODO.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Service Hook Storage Queue Pipelines.



## Import

Service Hook Storage Queue Pipeliness can be imported using the `resource id`, e.g.

```shell
terraform import azuredevops_servicehook_storage_queue_pipelines.example 00000000-0000-0000-0000-000000000000
```
