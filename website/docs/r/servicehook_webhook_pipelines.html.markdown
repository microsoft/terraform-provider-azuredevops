---
layout: "azuredevops"
page_title: "AzureDevops: azuredevops_servicehook_webhook_pipelines"
description: |-
  Manages a Pipelines publisher webhook (HTTP POST) service hook subscription within Azure DevOps.
---

# azuredevops_servicehook_webhook_pipelines

Manages a service hook subscription that uses the `pipelines` publisher and posts
events to a webhook (`webHooks` consumer / `httpRequest` action).

This is the webhook counterpart to
[`azuredevops_servicehook_storage_queue_pipelines`](servicehook_storage_queue_pipelines.html.markdown)
and the pipelines counterpart to
[`azuredevops_servicehook_webhook_tfs`](servicehook_webhook_tfs.html.markdown).

## Example Usage

```hcl
resource "azuredevops_project" "example" {
  name = "Example Project"
}

resource "azuredevops_servicehook_webhook_pipelines" "build_failed_canceled" {
  project_id = azuredevops_project.example.id
  url        = "https://example.azurewebsites.net/api/BuildFailed"

  http_headers = {
    "x-functions-key" = var.function_token
  }

  stage_state_changed_event {
    stage_state_filter  = "Completed"
    stage_result_filter = "Canceled"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the project. Forces a new resource.
* `url` - (Required) The URL to which the HTTP POST will be sent.
* `accept_untrusted_certs` - (Optional) Accept untrusted SSL certificates. Defaults to `false`.
* `basic_auth_username` - (Optional) Basic authentication username.
* `basic_auth_password` - (Optional) Basic authentication password.
* `http_headers` - (Optional) Map of HTTP headers to send with the POST request.
* `resource_details_to_send` - (Optional) `all`, `minimal`, or `none`. Defaults to `all`.
* `messages_to_send` - (Optional) `all`, `text`, `html`, `markdown`, or `none`. Defaults to `all`.
* `detailed_messages_to_send` - (Optional) `all`, `text`, `html`, `markdown`, or `none`. Defaults to `all`.
* `resource_version` - (Optional) Resource version of the subscription. Defaults to `5.1-preview.1` (required for the pipelines publisher).

Exactly one of the following event blocks must be specified:

* `stage_state_changed_event` - (Optional)
  * `pipeline_id` - (Optional) Pipeline ID to filter on.
  * `stage_name` - (Optional) Stage name to filter on.
  * `stage_state_filter` - (Optional) One of `NotStarted`, `Waiting`, `Running`, `Completed`.
  * `stage_result_filter` - (Optional) One of `Canceled`, `Failed`, `Rejected`, `Skipped`, `Succeeded`.

* `run_state_changed_event` - (Optional)
  * `pipeline_id` - (Optional) Pipeline ID to filter on.
  * `run_state_filter` - (Optional) One of `InProgress`, `Canceling`, `Completed`.
  * `run_result_filter` - (Optional) One of `Canceled`, `Failed`, `Succeeded`.

## Attributes Reference

* `id` - The ID (UUID) of the service hook subscription.

## Import

Service hook subscriptions can be imported using the subscription ID:

```sh
terraform import azuredevops_servicehook_webhook_pipelines.example 00000000-0000-0000-0000-000000000000
```
