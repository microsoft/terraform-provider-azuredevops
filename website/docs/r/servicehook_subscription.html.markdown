---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_servicehook_subscription"
description: |-
  Manages a Service Hook subscription.
---

# azuredevops_servicehook_subscription

Manages a Service Hook subscription in Azure DevOps. Service Hooks provide a way to subscribe to events that occur in Azure DevOps and have those events invoke a service endpoint.

## Example Usage

### WebHook Subscription for Git Push Events

```hcl
resource "azuredevops_project" "example" {
  name               = "example-project"
  work_item_template = "Agile"
  version_control    = "Git"
  visibility         = "private"
  description        = "Managed by Terraform"
}

resource "azuredevops_servicehook_subscription" "webhook_git_push" {
  project_id          = azuredevops_project.example.id
  publisher_id        = "tfs"
  event_type          = "git.push"
  consumer_id         = "webHooks"
  consumer_action_id  = "httpRequest"

  publisher_inputs = {
    repository = "MyRepository"
    branch     = "refs/heads/main"
  }

  consumer_inputs = {
    url = "https://example.com/webhook"
  }

  resource_version = "1.0"
  status          = "enabled"
}
```

### Azure Service Bus Subscription for Build Complete Events

```hcl
resource "azuredevops_servicehook_subscription" "servicebus_build" {
  project_id          = azuredevops_project.example.id
  publisher_id        = "tfs"
  event_type          = "build.complete"
  consumer_id         = "azureServiceBus"
  consumer_action_id  = "serviceBusQueueMessage"

  publisher_inputs = {
    buildStatus = "Succeeded"
  }

  consumer_inputs = {
    connectionString = "Endpoint=sb://example.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=your-key"
    queueName       = "build-notifications"
  }

  resource_version = "1.0"
  status          = "enabled"
}
```

### Azure Storage Queue Subscription for Pipeline Events

```hcl
resource "azuredevops_servicehook_subscription" "storage_pipeline" {
  project_id          = azuredevops_project.example.id
  publisher_id        = "pipelines"
  event_type          = "ms.vss-pipelines.run-state-changed-event"
  consumer_id         = "azureStorageQueue"
  consumer_action_id  = "enqueue"

  publisher_inputs = {
    pipelineId  = "123"
    runStateId  = "Completed"
    runResultId = "Succeeded"
  }

  consumer_inputs = {
    accountName = "mystorageaccount"
    accountKey  = "your-storage-account-key"
    queueName   = "pipeline-notifications"
  }

  resource_version = "5.1-preview.1"
  status          = "enabled"
}
```

### Organization-level Subscription

```hcl
resource "azuredevops_servicehook_subscription" "org_webhook" {
  # No project_id specified - organization-level subscription
  publisher_id       = "tfs"
  event_type         = "workitem.created"
  consumer_id        = "webHooks"
  consumer_action_id = "httpRequest"

  publisher_inputs = {
    workItemType = "Bug"
  }

  consumer_inputs = {
    url = "https://example.com/org-webhook"
  }

  resource_version = "1.0"
  status          = "enabled"
}
```

## Argument Reference

The following arguments are supported:

* `publisher_id` - (Required) The publisher ID that identifies the event source. Common values include:
  * `tfs` - Team Foundation Server (for Git, Work Items, Builds, etc.)
  * `pipelines` - Azure Pipelines  
  * `boards` - Azure Boards

* `event_type` - (Required) The event type to subscribe to. Examples:
  * `git.push` - Git push events
  * `git.pullrequest.created` - Pull request created
  * `build.complete` - Build completed
  * `workitem.created` - Work item created
  * `ms.vss-pipelines.run-state-changed-event` - Pipeline run state changed

* `consumer_id` - (Required) The consumer ID that identifies the target service. Common values:
  * `webHooks` - HTTP webhooks
  * `azureServiceBus` - Azure Service Bus
  * `azureStorageQueue` - Azure Storage Queue

* `consumer_action_id` - (Required) The action ID for the consumer. Examples:
  * `httpRequest` - For webhooks
  * `serviceBusQueueMessage` - For Service Bus queues
  * `enqueue` - For Storage queues

* `consumer_inputs` - (Required) A map of consumer-specific configuration inputs. This field is sensitive as it may contain secrets like connection strings or API keys.

---

* `project_id` - (Optional) The ID of the project. If not provided, the subscription will be created at the organization level.

* `publisher_inputs` - (Optional) A map of publisher-specific configuration inputs for filtering events.

* `resource_version` - (Optional) The resource version for the subscription. Default: `1.0`

* `status` - (Optional) The status of the subscription. Possible values: `enabled`, `disabled`, `disabledByUser`, `disabledBySystem`, `onProbation`. Default: `enabled`

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the service hook subscription

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 minutes) Used when creating the Service Hook subscription.
* `read` - (Defaults to 5 minutes) Used when retrieving the Service Hook subscription.
* `update` - (Defaults to 10 minutes) Used when updating the Service Hook subscription.
* `delete` - (Defaults to 10 minutes) Used when deleting the Service Hook subscription.

## Import

Service Hook subscriptions can be imported using the subscription ID:

```shell
terraform import azuredevops_servicehook_subscription.example 12345678-1234-5678-9012-123456789012
```

## PAT Permissions Required

- **Service Hooks**: Read & Write - Grants the ability to create and manage service hook subscriptions.

## Relevant Links

* [Azure DevOps Service Hooks](https://docs.microsoft.com/en-us/azure/devops/service-hooks/)
* [Azure DevOps Service REST API 7.0 - Service Hooks](https://docs.microsoft.com/en-us/rest/api/azure/devops/hooks/?view=azure-devops-rest-7.0)